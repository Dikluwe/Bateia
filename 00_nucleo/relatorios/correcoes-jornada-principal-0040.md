# Correções da jornada principal — passo 0040

Data: 2026-09-07. Implementação e validação local concluídas. O foco é a qualidade da evidência para decisões de reúso, entre projetos e linguagens.

## Mudanças implementadas

Em `01_core/usecase/evidence.go`, a avaliação lexical deixa de emitir PROVEN por nomes de símbolos, exemplos e presença de testes. Esses indícios permanecem nas fontes, com força limitada e explicação “comportamento não verificado”. A versão das regras passa a `2026-09-07.0040`. Não foi acrescentada execução automática de código nem um verificador comportamental geral.

O reconhecimento de licença exige documento com nome apropriado na raiz, mesmo quando um bundle antigo já classifica o item como licença. Arquivos como `license_checker_test.go` e licenças de diretórios internos não estabelecem autoridade sobre todo o projeto. Licenças com escopo complexo permanecem inconclusivas.

A CLI deixa de construir a descrição do candidato local com a ideia. A triagem em `probabilistic_screening.go` também deixa de usar as consultas como conteúdo do candidato; cobertura preliminar usa descrição e snippets retornados. Consultas continuam registradas como procedência da busca. Pontuação de triagem continua não calibrada e não comprova solução.

Menções a não objetivos, inclusive frases negadas, permanecem INCONCLUSIVE até que haja evidência de comportamento. Não se infere conflito apenas de palavras. A linha de uma fonte é calculada a partir do trecho; localização desconhecida é omitida, em vez de receber linha 1.

A seleção prioriza código de produção e penaliza diretórios temporários, gerados e de validação; alvos explícitos podem ser examinados. Arquivos elegíveis omitidos por quantidade ou profundidade recebem SKIPPED_BUDGET, e a coleta fica incompleta. Isso evita converter um orçamento esgotado em ausência funcional.

Em `reuse_decision.go` e na projeção CLI, o resultado INCONCLUSIVE informa a condição desconhecida e a próxima observação necessária por candidato. Capacidades já comprovadas por avaliações fornecidas ao decisor não são reclassificadas como lacunas porque falta manutenção ou licença.

## Antes e depois observados

| Caso do 0039 | Antes | Depois |
| --- | --- | --- |
| Eco Go rotulado como assíncrono | PROVEN; uma parte reutilizável | PARTIAL; nenhuma parte promovida |
| Licença da autoanálise | MIT PROVEN via teste do detector | INCONCLUSIVE |
| Dois candidatos locais sem descrição observada | Alinhamento 1, cobertura 1, ranking 0,6033 | Alinhamento 0, cobertura 0, ranking 0,2633 |
| “This tool does not upload code.” | CONFLICT; linha 1 | INCONCLUSIVE; linha 3 |
| Orçamento da autoanálise | Arquivos centrais não coletados sem registro individual do descarte | 805 arquivos elegíveis registrados como SKIPPED_BUDGET; coleta incompleta |
| Manutenção desconhecida | Razão genérica no dossiê principal | Instrução para obter metadados datados de atividade/arquivamento |
| Exportador Python fora do catálogo | INCONCLUSIVE sem regra | Permanece INCONCLUSIVE sem regra |

A autoanálise ainda não examinou `01_core/usecase/evidence.go`; agora esse arquivo está explicitamente registrado como encontrado e não examinado. A seleção melhorou a preferência por produção, mas a cobertura semântica continua insuficiente nesse caso. O resultado não autoriza dizer que o Bateia compreende sua implementação inteira.

## Validação

Cinco execuções locais e um replay concluíram com código 0. O replay preservou byte a byte os quatro JSONs de decisão/avaliação/alternativas/expansão comparados. Descoberta vazia e com erro continuou INCONCLUSIVE no serviço e MCP reais com adaptador controlado. As requisições ao provedor foram registradas, evitando um falso teste que parasse antes da descoberta.

Novas regressões verificam eco Go, autoridade de licença inclusive em bundles mal classificados, linha real, menções negadas/positivas, prioridade C++/Python, metadados locais independentes da ideia, consultas sem cobertura fabricada, omissão de orçamento e separação entre manutenção desconhecida e lacuna funcional.

`go test ./...` passou; `crystalline-lint .` retornou “No violations found”; `go build -o bateia ./04_wiring` passou. O comando antigo com somente `main.go`, indicado nas instruções locais, falha porque omite arquivos auxiliares do pacote; a compilação válida usa o pacote completo.

Expectativas públicas que promoviam tipos ou testes vazios a PROVEN foram substituídas pelo contrato 0040, inclusive expectativas de componentes e top choice derivados dessas fontes. O golden de identidade de execução mudou com a versão das regras; hashes de evidência congelada foram preservados. Os artefatos históricos do 0039 não foram reescritos. Não há alegação de auditoria independente.

## Limites e continuidade

Este passo reduz conclusões indevidas; não demonstra ganho de recall, descoberta pública ou compreensão comportamental geral. O avaliador lexical não produz comprovação funcional automática. Isso reduz recomendações fortes que dependiam dos falsos positivos anteriores, preservando os indícios para revisão.

O próximo ganho de qualidade exige associar capacidades a implementações e verificações específicas, com cobertura declarada, sem depender de adicionar uma regra para cada exemplo. O catálogo fora dos domínios cobertos e a seleção de implementações em repositórios grandes continuam limitados.

[Artefatos, hashes e comandos](../benchmark/0040/README.md).
