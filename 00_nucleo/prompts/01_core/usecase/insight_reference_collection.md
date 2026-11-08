# Coleta limitada de referências

**Camada:** L1. **Contrato:** passo 0137.

Pela porta estreita, ler árvore e primeiro README regular de raiz conforme prioridade. Só usar SHA válido sem conflito, conferir resposta ref/path/cardinalidade/UTF-8/32 KiB, conservar hash completo e bundle limitado explicitamente incompleto. Erros seguros sem conteúdo arbitrário. Não truncar nem buscar HEAD/fallbacks na política. Sem I/O concreto.

Conservar cobertura explícita somente do documento root README examinado, conforme passo: não é avaliação de capacidade. Não emitir estrutura de cobertura vazia que alterne entre null/[] no marshal legado; preservar o replay de contexto após publicação sem modificar o formato antigo.

## Evolução 0139 (precedência sobre limites de fontes do recorte0138)

Aplicar o passo0139: RepositoryReader opcional exclusivo Gitee/AtomGit, resolução imutável e contextos de leitura com identidade por sítio. Reader legado permanece só GitHub. Preservar resultados antigos quando o novo leitor não está vinculado. Erros do leitor classificados sem mensagens remotas; budgets compartilhados; tarefas agregam IDs dos cartões com origem.

## 0154 — Descoberta contextual
Se a referência não tem pin e não há conflito, tentar porta InsightSiteRepositoryReader quando disponível. Exigir identidade, URL canônica e commit imutável antes de ler; preservar metadados originais e inserir apenas pin resolvido. Falha explícita, sem bypass. Reader sem resolver mantém OPAQUE_REVISION. Nenhuma resolução por referência com conflito.
