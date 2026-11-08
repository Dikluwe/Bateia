# Evidence domain

Contratos imutáveis para planejar, congelar e avaliar evidências de capacidades em
`repo@commit`. Os estados funcionais, licença, manutenção e operação permanecem
separados. Ausência só pode ser declarada após busca completa; falhas produzem
inconclusão. Nenhum contrato admite saída de LLM, embedding ou inferência remota.

Cada alvo registra `COLLECTED`, `NOT_FOUND`, `FAILED`, `SKIPPED_BUDGET` ou
`UNSUPPORTED`. Identidade `owner/repo` e revisão SHA são dimensões separadas.

Na versão 2, cada regra declara família, proveniência, alvos justificáveis e se
`PROVEN` exige fontes independentes. Versões antigas permanecem legíveis para
replay, enquanto catálogos v2 são validados estritamente.

## Passo 0072 — cobertura da coleta por requisito

O plano pode carregar consultas determinísticas derivadas exclusivamente de regras
declaradas cobertas. A quota dirigida permanece dentro do orçamento global. O bundle
registra, com arrays não nulos, caminhos candidatos, examinados e omitidos por
requisito. `COLLECTION_INCOMPLETE` representa superfície plausível não examinada e
nunca autoriza ausência nem gap atribuído ao candidato; nomes de caminhos apenas
orientam a recuperação e não constituem evidência funcional.

# 0044 — Porta e evidência de recuperação de implementação
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md (contrato congelado).
Tipos de domínio para orçamento, relatório, inputs hashados e hits de AST; porta opcional de recuperação sem IO concreto em L1. Busca pode ser partial/not_evaluated; hits não são comprovação de capacidade. Preservar hashes da varredura e conferir identidade antes de exibir definição. Campos opcionais em resultados de extração/bundles mantêm compatibilidade com consumidores existentes.

## Passo 0073 — localização probatória intrarquivo

Cada fonte derivada de uma regra registra o padrão canônico e a classe da ocorrência:
`STRONG_CONTENT`, `SUPPORTING_CONTENT` ou `DOCUMENTATION_CONTENT`. `PATH_ONLY` existe
para diagnóstico da recuperação, mas nunca integra fontes ou componentes reutilizáveis.
Linha, snippet, padrão e classe descrevem a mesma ocorrência imutável; nomes de arquivo
e termos de recuperação não se tornam prova funcional.

## 0158 — Integridade e serialização estável

Aplicar o [contrato0158](../../../passos/0158-validar-origem-e-estabilizar-bundles.md) e o [ADR0022](../../../ADR/ADR-0022-integridade-de-coleta-e-serializacao-estavel.md). Domínio estabiliza requirement_coverage nil/vazio como array e preserva bytes populados, sem mutação. Coleta genérica compartilha admissão, conserva pin do plano e confronta identidade de metadados e Ref de extração antes de admitir conteúdo. Coleta coordenada valida todos os schemas1 antes de efeitos e propaga ErrIntegrity sem relatório útil. Somente a API genérica conserva schema ausente legado. Não recalcular registros/hashes históricos; representações instáveis passam à forma corrigida apenas em novas serializações.

Critérios compartilhados vigentes: [EVD-01–05](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#aceitação).
