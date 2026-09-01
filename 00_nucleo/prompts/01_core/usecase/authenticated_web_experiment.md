# Contrato L0 — caso de uso do experimento web autenticado

## Objetivo

Orquestrar uma execução autenticada, limitada e auditável das 36 queries congeladas e tornar seu resultado integralmente reexecutável sem token e sem rede.

## Pré-condições antes da rede

O caso de uso recebe os três conjuntos congelados, a autorização explícita, o orçamento e uma porta de busca observável. Deve encerrar sem chamada externa quando:

- autorização estiver ausente, falsa ou divergente do contrato do Passo 0020;
- qualquer braço não tiver exatamente 12 queries, houver ID duplicado/vazio ou query vazia;
- orçamento ou ordem não forem os congelados;
- provedor não for GitHub repository search autenticável;
- a margem segura de rate limit não estiver declarada;
- a capacidade não garantir somente metadados de repositórios públicos.

O token não é argumento deste caso de uso e não participa de objetos, hashes ou mensagens.

## Agendamento e orçamento

Construa previamente o plano imutável `N1, F1-1, F2-1 ... N12, F1-12, F2-12`. Lance no máximo duas tentativas simultâneas sem alterar essa ordem de despacho. Antes de cada tentativa, verifique timeout total, orçamento de tentativa e margem segura declarada. Se iniciar a tentativa excederia qualquer limite, não acesse a rede e registre `SKIPPED_BUDGET` para a query e para todas as posições que não puderem iniciar.

Cada query recebe uma tentativa primária. Somente `UNKNOWN_TIMEOUT` ou falha transitória explicitamente classificada do provedor pode receber uma repetição; 403/429 classificados como `UNKNOWN_RATE_LIMIT`, falha permanente, resposta privada e sucesso não são repetidos. A repetição conserva request lógico, recebe request ID novo, `attempt_number=2`, conta no orçamento e não permite ultrapassar 20 segundos por tentativa, 10 minutos totais, concorrência 2 ou 72 tentativas máximas. Nunca repita resposta bem-sucedida.

Não adapte query, limite, braço, provedor ou quantidade de queries com base em resultados ou rate limit. Headers podem antecipar parada global pela margem segura; não podem reduzir um braço isoladamente.

## Normalização

Mapeie toda tentativa para exatamente um estado obrigatório. Erro sem resposta classificável vira `UNKNOWN_PROVIDER_RESPONSE`; deadline vira `UNKNOWN_TIMEOUT`; 403/429 ou sinal inequívoco de cota vira `UNKNOWN_RATE_LIMIT`; resposta válida sem itens vira `EMPTY`.

Antes de aceitar um hit, exija evidência positiva de visibilidade pública. Qualquer item privado é descartado e transforma a tentativa em `REJECTED_PRIVATE`; visibilidade opaca transforma-a em `UNKNOWN_PROVIDER_RESPONSE`. Nenhum hit dessa tentativa é persistido nesses casos.

Sele e valide request, resposta e snapshot pela serialização canônica. Uma varredura de segredo deve ocorrer antes de qualquer persistência e também sobre o artefato serializado final. Encontrar token, `Authorization`, cookie ou derivado reconstruível impede toda persistência do run e produz `FAIL_SECURITY` por canal seguro que não ecoe o material.

## Snapshots e replay

Produza um snapshot por braço, com o mesmo run ID, autorização redigida, orçamento e plano global referenciado. Cada snapshot contém suas 12 queries e todas as suas tentativas, inclusive `Unknown`, retries, rejeições e skips. Ordenação persistida é determinística por posição global e número da tentativa, não por qualidade nem por horário de conclusão.

O replay:

- aceita somente snapshots e não recebe porta de rede nem lê variável de ambiente;
- verifica schema, snapshot hash, request/response hashes, IDs, associação query/tentativa, orçamento, ordem, rank, campos materiais e visibilidade;
- preserva literalmente todos os estados `UNKNOWN_*`, `REJECTED_PRIVATE` e `SKIPPED_BUDGET`;
- deduplica hits públicos por URL canônica, conserva a melhor posição e mantém separação orgânico/patrocinado quando aplicável;
- gera bytes idênticos em duas execuções sobre os mesmos bytes de entrada.

Hash ausente/divergente, segredo, item privado, associação inválida ou orçamento excedido torna o replay inválido; não deve haver reparo silencioso.

## Exposição e veredito

Calcule por braço `completed/12`, onde conclusão é somente `SUCCESS` ou `EMPTY`. Exposição comparável requer os três braços com pelo menos 80% e amplitude de no máximo 10 pontos percentuais, além de 12 queries lógicas por braço e ausência de query extra. Rate limit que impede isso produz `INCONCLUSIVE_RATE_LIMIT`; opacidade não atribuível com segurança à cota produz `UNKNOWN_PROVIDER_EXPOSURE`.

Use a precedência de veredito definida no contrato de domínio. Não calcule relevância, precisão ou recall neste caso de uso.
