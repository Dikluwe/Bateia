# Prompt espelhado — oráculos e ataques do schema `EngineeringInvestigation` R2

## Autoridade e regime

- regime: protocolo Tekt completo, pré-candidato;
- papéis autorizados pelo coordenador: `AUTOR_DE_ORACULOS` e `ADVERSARIO`;
- executor: `/root/oraculo_0029`;
- teste espelhado: `04_wiring/step0029_experiment/schema_r2_oracle_test.go`;
- pacote de teste: `step0029_experiment_test`;
- alegação máxima: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

Esta autoria conjunta produz apenas controles positivos e ataques negativos do
schema de transporte. Ela não lê nem edita `harness.go`, `types.go`, saída do
implementador, receipts, selos, certificados ou vereditos. O workspace é
compartilhado e, portanto, não oferece atestação material de isolamento.

## Entradas congeladas

| Entrada | SHA-256 |
|---|---|
| `00_nucleo/benchmark/0029/engineering-investigation.schema-r2.json` | `7b16306c0188d627d06392edf4438f087c26264e801980f3649570bbe16ab7d8` |
| `00_nucleo/benchmark/0029/schema-transport-refinement-r1.json` | `54cae1534b1cd46ee7a88206b1df28867a5cdeb7546112a7fb59b24ac7c5a345` |
| `00_nucleo/benchmark/0029/prompt-contract-r2.json` | `204fcd66d878657a8a61a8933707dd15912be6fab70af14b9d859759aada179c` |
| `00_nucleo/benchmark/0029/contract.md` | `48f0f3ff313caf2cd1c17f0dacb3a0143bf7908928e6922d3838d25e4f1118e7` |
| `00_nucleo/passos/0029-medir-ganho-da-ia-usando-bateia.md` | `135f8de345e2b7dd2f53f45d7ec30611ea836d360b29111f9cc751df64b4877a` |

Entradas operacionais adicionais são a skill
`tekt-materializacao-segregada` e suas referências de papéis e gates. Nenhum
artefato candidato é uma entrada causal desta suíte.

## Hipótese do refinamento

Todos os cinco objetos aninhados precisam possuir forma fechada e inteiramente
obrigatória para que o schema comum seja capturável sem aceitar projeções
parciais ou propriedades não contratadas. O controle positivo deve demonstrar
que uma investigação contendo simultaneamente os cinco objetos válidos continua
aceita.

## Controles positivos

1. Uma `EngineeringInvestigation` completa mistura instâncias válidas de
   `alternatives`, `mechanisms`, `improvement_hypotheses`,
   `evidence_references` e `validation_plan`.
2. O schema standalone R2 é estruturalmente idêntico a `output_schema` no
   `prompt-contract-r2.json`.
3. O caminhamento mecânico encontra seis schemas de objeto — raiz e cinco
   aninhados — todos com `properties`, `required` exatamente igual ao conjunto
   de propriedades e `additionalProperties=false`.

## Ataques negativos

Para cada um dos cinco objetos aninhados, aplicar isoladamente:

- remoção de um campo obrigatório;
- substituição de um campo por valor de tipo incompatível;
- inserção de propriedade extra.

Aplicar ainda enum inválido em todos os campos aninhados que possuem enum no
schema R2:

- `alternatives[].decision`;
- `evidence_references[].source_type`;
- `evidence_references[].verification_status`.

Os outros três tipos aninhados não possuem enum e não recebem restrição
inventada. Cada mutação deve ser rejeitada sem alterar o controle positivo.

## Método mecânico independente

O teste implementa um avaliador mínimo e determinístico somente para as
palavras-chave presentes no schema congelado: `type` (`object`, `array` e
`string`), `properties`, `required`, `additionalProperties`, `items` e `enum`.
Ele lê apenas os dois artefatos JSON públicos necessários em tempo de teste; não
importa nem chama o candidato e não realiza rede, processo ou ação externa.

## Limites

Esta suíte demonstra a forma observável do schema R2 e seu transporte no prompt;
não demonstra comportamento do harness, ganho do Bateia, mutation score global,
refinamento de produto, selo ou veredito. A verificação desta autoria permanece
segregada e sob autoridade posterior.
