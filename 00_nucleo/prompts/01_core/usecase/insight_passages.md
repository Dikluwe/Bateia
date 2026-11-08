# Leitura automática e seleção de trechos

**Camada:** L1. **Contrato:** passo 0136.

Extrair candidatos por arquivo elegível conforme janelas/limites do passo, sem interpretar comportamento ou executar. Registrar termos literais, rótulos familiares lexicais, prefixos/itens ignorados. Escolher janela por matches/menos familiaridade/linha e arquivos por novidade estrutural de repositório/tipo, matches e menos familiaridade, ordem canônica; desduplicar trechos idênticos, no máximo um por arquivo. Preservar bytes/hash/linhas e IDs causais, com razões/diagnósticos explícitos. Respeitar cancelamento e bounds. Sem I/O ou catálogo fechado.

## 0153 — Passagens na revisão opt-in
No modo ReviewAlternatives, conservar a seleção histórica como primeira passagem por arquivo; então preencher vagas sob MaxPassages com até duas janelas adicionais não sobrepostas às já selecionadas do mesmo arquivo. Escolher por matches/menos familiaridade/riqueza lexical, ordem canônica em empate, desduplicar conteúdos globalmente. Preservar quotas de prefixo, bytes, linhas, cancelamento e teto global. Registrar ADDITIONAL_FILE_PASSAGE e limite de seleção quando houver candidatos não representados. Nenhuma janela adicional é interpretação ou garantia de alternativa. Fora do modo, comportamento idêntico.
