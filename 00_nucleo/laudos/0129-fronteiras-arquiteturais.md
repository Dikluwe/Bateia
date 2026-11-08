# 0129 — Fronteiras arquiteturais materializadas

**Data:** 2026-10-28  
**Contrato:** [passo 0129](../passos/0129-concluir-fronteiras-arquiteturais.md), [ADR-0012](../ADR/ADR-0012-portas-de-execucao-e-arquivos.md).  
**Evidência:** [bancada 0129](../../lab/bancada/arquitetura/0129/README.md).

## Resultado

O recorte encerra os efeitos temporais concretos remanescentes em L1 e os efeitos concretos de arquivos, ambiente, sinais e captura local na CLI. Portas pertencem a L1, sua implementação a `03_infra/host` e sua composição a L4. A CLI conserva entrada, apresentação e validação de transporte. A captura e o provedor de snapshots não possuem implementação física alternativa em L2.

O runtime durável conserva store, fila e lifecycle único. Containers internos possuem configuração de descoberta, extração e síntese; a capacidade de extração agrega/ordena mecanismos e lacunas, e o serviço publica essas saídas. A resolução de fontes de 0127 continua separada da análise e do protocolo. Estado operacional não promove hipótese.

Testes de protocolo/parse de Play Store, F-Droid, forjas, iTunes e replay foram transferidos para L3, preservando seus oráculos. Testes de composição continuam em L4. O limite MCP é interpretado em L2. Datas, durações e deadlines usam porta explícita; a data do relatório é reutilizada pelo extrator do Lab.

## Preservação e reparos de verificação

Leituras limitadas, retenção de descritores, recusa de aliases/hardlinks, modos, quotas, hash/cópia na mesma leitura, publicação atômica e cleanup foram mantidos nas extrações por capacidade e nos testes ativos de segurança. As fachadas públicas Configure permanecem. A extensão temporal da recomendação está num arquivo aditivo, conservando o hash histórico de `recommendation.go`.

A suíte encontrou ordem física variável em `Examined` do adaptador local. A projeção passou a ordenar por caminho; varredura limitada, hashes e oráculo de paridade local/remoto foram preservados. A falha inicial e a correção possuem logs na bancada.

## Verificação e limites

O auditor AST verifica os 346 arquivos Go produtivos, a gravidade, os efeitos inspecionados, o isolamento do Lab e a ausência de material executável em L0. Seus controles incluem casos que devem falhar; não constitui prova semântica universal de arquitetura.

O lint instalado confundia funções Go livres com impls e não distinguia imports de teste. Um patch reproduzível do parser, aplicado somente a uma cópia temporária da revisão fixada, preserva regras/severidades e verifica tipos reais e métodos produtivos. A ferramenta passa nos 671 testes de biblioteca e nos 11 controles do parser, antes do lint da árvore. A instalação global permanece intacta; usar o script documentado para reproduzir o resultado.

Os códigos de saída finais de lint, suíte raiz, vet e build constam de [verification.json](../../lab/bancada/arquitetura/0129/verification.json). Gates selados não foram promovidos e o acervo histórico não foi reescrito. As alterações anteriores 0126–0128 e arquivos não relacionados da worktree foram preservados.

A CLI consumida como biblioteca requer host explicitamente composto no startup. Runtime temporal ausente deixa metadados desconhecidos e recusa deadline solicitado. Concorrência e estado em memória do funil permanecem: não se afirma pureza referencial de toda a orquestração.

A tentativa de substituir `.agents/AGENTS.md` por referência foi recusada pela revisão automática, que a classificou como alteração não autorizada de instruções protegidas. Esse arquivo não foi alterado; a conclusão deste recorte cobre as fronteiras produtivas e seus contratos.
