# Prompt espelhado — oráculos independentes do harness experimental 0029

## Linhagem e autoridade

- papel: `AUTOR_DOS_ORACULOS`;
- executor: `/root/oraculo_0029`;
- regime: protocolo Tekt completo, anterior à implementação candidata;
- contrato consumido: `00_nucleo/benchmark/0029/contract.md`, revisão candidata R1;
- teste espelhado: `04_wiring/step0029_experiment/harness_oracle_test.go`;
- pacote de teste: `step0029_experiment_test`;
- alegação máxima: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

O autor não lê implementação candidata, ataques, mutantes, casos `TEST`
protegidos, receipts, selos, certificados, laudos, resultados ou vereditos do
Passo 0029. O workspace compartilhado impede atestar isolamento material; os
hashes identificam entradas e ordem causal, mas não provam ausência de acesso
por outros processos.

## Objetivo observável

Exercitar exclusivamente a API pública do pacote
`bateia/04_wiring/step0029_experiment`, como consumidor externo, sem I/O de
rede, disco, processos, relógio ou aleatoriedade implícita. Os testes devem
compilar quando a implementação pública contratada existir e não incorporam nem
despacham agente.

## Obrigações cobertas

1. `BuildPlan` materializa exatamente 108 células e 54 pares, com partição
   18/36/54 para `CONTRACT`/`VALIDATION`/`TEST`, 18 células por caso e dois
   braços por identidade de par.
2. Os nove perfis são exatamente LUNA, TERRA e SOL cruzados com `low`, `medium`
   e `high`, sem fallback. Uma mutação de `ProfileID` é rejeitada por
   `ValidatePlan`.
3. Dentro do par, caso, perfil, prompt base e budget são iguais. O tratamento
   altera o prompt efetivo pela instrução literal e troca a allowlist para a
   superfície Bateia; A possui zero e B 25 chamadas Bateia.
4. Violações públicas são estáveis e ordenadas por `Code`, `Path`, `Message`.
5. `ClassifyAttempt` preserva os status enumerados. `EvaluateRetry` permite uma
   repetição somente para `INFRA_FAILURE` confirmada antes de qualquer saída;
   output observado ou opaco impede retry.
6. `BuildRunReceipt` preserva todas as tentativas e hashes causais, calcula hash
   determinístico e não publica prompt privado, chain of thought, credenciais
   ou bodies integrais de tool.
7. `Anonymize` remove braço, modelo, esforço, par e timestamps; aplica mapa
   global consistente de aliases; preserva texto substantivo, inclusive erros;
   e produz hashes e bytes idênticos em reexecução. O controle positivo O6 usa
   `EngineeringInvestigation` R2 completo: cada `evidence_references[]` contém
   exatamente `reference_id`, `source_type`, `locator`, `sha256`, `supports` e
   `verification_status`, sem propriedades legadas ou extras.
8. `Score` exige 30 inteiros (três julgadores por dez dimensões) e calcula sem
   arredondamento `100*soma/120`; trinta notas 3 resultam exatamente em 75.
9. `Analyze` compara apenas braços do mesmo caso e perfil. Vinte e sete pares
   com A=50 e B=75 produzem 27 wins e mediana 25 sob intenção de tratar.
10. `EvaluateGates` aceita os limites inclusivos congelados (25 válidas em B,
    taxa ponderada ao menos 0,65, mediana 10, estratos 0, agreement 0,67) somente
    quando todos os gates causais, de segurança, budget, refinamento e mutation
    score passam. Agreement opaco tem precedência inconclusiva.

## Oráculos de opacidade

- usage ou duração não exposta permanece `Observation{State: UNKNOWN, Reason}`;
- default de runtime não exposto permanece `UNKNOWN_RUNTIME_DEFAULT`;
- primeira saída do modelo desconhecida nunca é inferida como `false` para retry;
- run Bateia inconclusivo com schema válido permanece `VALID_OUTPUT`, enquanto a
  alegação dependente permanece `UNKNOWN`;
- falha de infra sem resultado válido permanece no denominador, com qualidade
  desconhecida e experimento confirmatório inconclusivo.

## Rubrica congelada

`00_nucleo/benchmark/0029/rubric.json` fixa dez dimensões, âncoras inteiras 0–4,
pesos iguais, 30 notas por saída válida, categorias críticas separadas,
Krippendorff alpha intervalar com distância quadrática e todos os thresholds do
veredito. A rubrica não pode ser alterada após abertura do split `TEST`.

## Limites do teste

Os casos são fixtures sintéticas públicas e não contêm casos, expectativas ou
artifacts do split `TEST` protegido. O harness de oráculo não prova execução de
sessões reais, integridade do store ou ganho de produto; ele cobre apenas o
fragmento determinístico e binding-free da API contratada. Casos negativos e
mutation score são responsabilidade segregada do adversário/verificador.

## Revisão focal O6

- causa pública: `ANONYMIZATION_POSITIVE_USED_PRE_R2_EVIDENCE_REFERENCE`;
- hipótese: transportar os mesmos identificadores opacos por campos válidos R2
  restaura apenas a validade estrutural do controle e preserva as expectativas
  preexistentes de alias consistente e texto não reparado;
- resultado candidato observado: nenhum;
- execução da suíte: não realizada nesta autoria.
