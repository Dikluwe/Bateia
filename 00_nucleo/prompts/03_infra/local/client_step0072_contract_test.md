## Supersessão explícita — passo 0095

O requisito abaixo de sucesso acima de 20.000 entradas é histórico. A revisão
0095 mantém não truncar silenciosamente, mas agora exige erro explícito e árvore
nil no teto, não uma coleta completa ilimitada. Certificados não são reemitidos.

# Contrato independente — enumeração local e header C++, passo 0072

Derivar testes somente do passo 0072 e do baseline público congelado. Não ler código
candidato.

Em diretório temporário confinado, verificar que `FetchTree` não corta silenciosamente a
superfície no antigo limite de 20.000 entradas e que inclui um arquivo posterior ao
limite. Verificar também que a coleta pública aceita `span` extensionless quando dirigido
por consulta, mas não classifica arquivo extensionless irrelevante como código.

Não seguir aliases externos, não relaxar limites de leitura e não depender da ordem do
filesystem.

## Revisão de segurança 2026-09-26

Aplicar `security.go`: árvore em lotes (máximo 20.000 entradas), sem alocar diretório inteiro antes de controlar orçamento. Exceder teto produz erro explícito, não sucesso truncado. Inventário expõe symlink/hardlink/nonregular; leitura rooted mantém alias interno mas recusa hardlinks por nlink do descritor. Extração até 100 arquivos/8 MiB (2 MiB por arquivo). Não executar Git ou herdar ambiente/configuração. BaseCommit consultivo vem de HEAD/refs/packed-refs confinados; gitdir externo fica indisponível explicitamente. Metadata não autentica bytes da worktree e certificados históricos não são atualizados.
