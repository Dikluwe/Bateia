# Laudo 0016 — Superfícies públicas de catálogo

**Estado**: `EM VALIDAÇÃO`  
**Data**: 2026-08-26  
**Regime**: protocolo completo Tekt, executado sem atestação de isolamento

## Materializado

- contrato nominal para `REPOSITORY`, `README`, `TOPICS`, `PACKAGE_REGISTRY` e
  `EXTERNAL_PROJECT_CATALOG`;
- planejador que preserva discriminator completo por superfície;
- busca GitHub restrita a README raiz/documentação inicial;
- busca estruturada por tópico com igualdade do tópico normalizado;
- snippets, campos, queries, superfícies, owner/repo e commits rastreáveis;
- união canônica que conta cada superfície uma vez;
- identidade, snippet ou source opacos permanecem `Unknown`;
- avanço exige duas superfícies independentes;
- CLI `idea catalog discover-surfaces` e `screen-multiview`;
- manifesto, selo, snapshots, avaliação e certificado reproduzíveis.

O arquivo do Passo 0016 permanece com estado literal `PRONTO PARA
IMPLEMENTAÇÃO` porque foi congelado antes do contrato. O índice e este laudo
registram o estado operacional sem invalidar o hash do prompt.

## Gates e ataques

O contrato rejeitou 3/3 mutações: substituição nominal de superfície, promoção
de pacote opaco e contagem de queries repetidas como visões independentes.
Mutation score: 1,0.

Um ataque adicional comprovou que owner/repo e commit não bastam quando o
snippet documental está ausente. O núcleo passou a manter esse caso em
`Unknown`; o snapshot público não contém match com snippet opaco, portanto não
precisou ser repetido.

Suíte completa, linter V15 e `git diff --check` passaram. Manifesto, contrato e
oráculos mantiveram os hashes selados.

## Experimento público

| Superfície | Planejadas | Sucesso | Unknown | Matches |
|---|---:|---:|---:|---:|
| README | 12 | 10 | 2 rate-limited | 16 |
| tópicos | 12 | 12 | 0 | 20 |
| pacotes Go | — | — | indisponível | — |
| ecosyste.ms | — | — | capacidade de busca não comprovada | — |

A superfície ecosyste.ms foi avaliada por sua documentação pública como fonte
de metadados de repositórios e vínculo com pacotes. Sem evidência de busca
textual de produtos pela intenção, ela não foi chamada nem classificada como
resultado vazio.

Resultados consolidados:

| Observável | Resultado |
|---|---:|
| candidatos README+tópicos | 35 |
| união com braço C | 69 |
| sementes repo-backed | 0/4 |
| recall competitivo | 0% |
| `ADVANCE` | 0 |
| `REVIEW` | 35 |
| queries de prova | 0 |

Nenhum repositório apareceu em duas superfícies independentes. Assim, nenhum
candidato avançou e Sourcegraph não foi executado.

## Veredito

O fragmento multiview preserva provenance e opacidade, mas nome, descrição,
README e tópicos continuam com recall 0%. O resultado está abaixo do piso de
40% pré-registrado.

O próximo passo deve reformular a geração de queries por **decomposição
funcional determinística**, ainda sem LLM: pesquisar comportamentos observáveis
do produto — por exemplo recomendação de repositórios a partir de requisitos,
comparação de ferramentas e decisão build-versus-reuse — em vez de apenas os
rótulos abstratos da ontologia.

Não é permitido promover `REVIEW`, afrouxar discriminator ou inserir nomes das
sementes. Package registry e catálogo externo continuam `Unknown`, não falhas
de recall.

A cadeia foi segregada por manifesto, selo e hashes, mas compartilhou a mesma
autoridade conversacional; por isso não há alegação de isolamento forte.
