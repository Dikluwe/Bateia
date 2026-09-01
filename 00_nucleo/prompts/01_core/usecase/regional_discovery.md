# Orquestração da descoberta regional

## Linhagem proprietária

Esta é a projeção consumer-specific para `01_core/usecase` da intenção proprietária em
`00_nucleo/prompts/01_core/domain/regional_discovery.md` e do contrato candidato M21. Ela
não amplia a semântica desses artefatos: apenas atribui ao caso de uso a normalização, o
planejamento, a propagação e a vinculação de replay já exigidos.

## Responsabilidade do caso de uso

O caso de uso recebe zero ou mais valores repetíveis de locale natural e preserva
`--language` exclusivamente como filtro opcional de linguagem de programação.

Os únicos valores de locale reconhecidos são `pt-BR`, `en-US`, `es-ES`, `fr-FR`, `de-DE`,
`ja-JP`, `ko-KR`, `ru-RU`, `ar-SA`, `hi-IN` e o alias `all`. O alias representa os dez
locales. A normalização expande `all`, remove repetições e usa exatamente a ordem canônica
acima. Reordenação, repetição e combinação de `all` com valores explícitos não alteram o
plano. Locales desconhecidos são erro de contrato e nunca há mais de dez locales
normalizados.

Zero locales preservam exatamente o fluxo, os artefatos de domínio e a identidade de
execução anteriores a M21. O conjunto normalizado participa do hash da execução somente
quando não vazio.

## Planejamento pelo catálogo controlado

Para o conceito `recipe-manager`, o caso de uso deriva propostas somente do catálogo
`regional-recipes-v1`, com consulta-fonte `gerenciador de receitas` e estratégia
`controlled-catalog-exact-v1`:

| Locale | Consulta nativa |
| --- | --- |
| `pt-BR` | `gerenciador de receitas` |
| `en-US` | `recipe manager` |
| `es-ES` | `gestor de recetas` |
| `fr-FR` | `gestionnaire de recettes` |
| `de-DE` | `Rezeptmanager` |
| `ja-JP` | `レシピ管理` |
| `ko-KR` | `레시피 관리자` |
| `ru-RU` | `менеджер рецептов` |
| `ar-SA` | `مدير الوصفات` |
| `hi-IN` | `रेसिपी प्रबंधक` |

Cada proposta preserva locale, versão do catálogo, regra/conceito, consulta-fonte,
estratégia e consulta nativa. Tradução livre, tradução implícita, transliteração e
derivação sem origem são proibidas. A identidade de cada proposta, consulta e tentativa
localizada contém o locale canônico para impedir colisões.

O plano regional registra em ordem canônica:

- todos os locales solicitados;
- os locales realizados por propostas válidas do catálogo;
- uma lacuna por locale solicitado sem tradução controlada para conceito válido, com
  estado `Unknown` e motivo não vazio.

Uma lacuna não é sucesso nem erro. Locale inválido, tradução sem regra, origem sem
consulta-fonte e catálogo inconsistente são erros de contrato.

## Execução e procedência

Cada proposta localizada é uma consulta normal para orçamento e limites. Não existe busca
regional implícita. O locale acompanha a requisição ao provedor, o snapshot, a tentativa
efetiva e a origem do candidato.

Resultados são deduplicados pela identidade do repositório. Quando o mesmo repositório é
encontrado por propostas diferentes, permanece um candidato com todas as origens
distintas, deduplicadas e ordenadas primeiro pelo locale canônico e depois pela identidade
determinística da proposta/consulta.

Reordenar locales, propostas ou respostas do provedor não altera plano, snapshot,
candidatos nem relatórios de domínio.

## Avaliação e relatórios

O caso de uso propaga as origens regionais ao relatório de alternativas JSON e à projeção
Markdown. Ambos expõem os mesmos locales canônicos por alternativa.

Locale é somente procedência. Ele não altera capacidade, cobertura, alinhamento do
problema, licença, manutenção, operação, blockers, fit, desempate ou posição no ranking.
Texto localizado encontrado é apenas descoberta e não substitui evidência rastreável do
repositório. Ausência ou ambiguidade funcional pode permanecer `INCONCLUSIVE`.

## Replay, retomada e opacidade

Checkpoint, replay e retomada vinculam exatamente o conjunto canônico não vazio ao hash
da execução. Um conjunto diferente, inclusive vazio versus não vazio, é incompatível.
Reordenações equivalentes têm a mesma identidade.

O caso de uso nunca converte `Unknown` em cobertura, evidência ou sucesso presumido; nunca
substitui tradução ausente por locale padrão, inglês implícito ou tradução livre.

## Limite da projeção

Esta projeção cobre somente descoberta regional controlada do conceito congelado e sua
propagação determinística. Não afirma tradução geral, cobertura de lojas proprietárias,
compreensão automática de documentação multilíngue, paridade entre mercados ou adequação
funcional das alternativas encontradas.
