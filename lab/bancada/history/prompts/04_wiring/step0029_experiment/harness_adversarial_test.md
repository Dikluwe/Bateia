# Prompt espelhado — harness adversarial do Passo 0029

## Linhagem e papel

- intenção: `00_nucleo/passos/0029-medir-ganho-da-ia-usando-bateia.md`;
- contrato público: `00_nucleo/benchmark/0029/contract.md`;
- papel autoral: `ADVERSARIO`;
- regime: protocolo Tekt completo, pré-candidato;
- alvo espelhado: `04_wiring/step0029_experiment/harness_adversarial_test.go`;
- alegação máxima: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

O adversário recebe apenas intenção, contrato, snapshot público de perfis e
registry público do corpus. Não recebe implementação candidata, oráculos,
positivos, opacos, fixtures protegidas, receipts, resultados ou vereditos.

## Objetivo observável

Construir uma suíte black-box em package externo que invoque somente a API Go
pública de `step0029_experiment`. Nenhum teste realiza rede, filesystem,
subprocesso, relógio implícito, aleatoriedade implícita ou qualquer outro I/O.

Materializar mutações negativas unitárias para:

1. modelo trocado;
2. esforço trocado;
3. perfil indisponível substituído;
4. Bateia presente no braço A;
5. artifact Bateia entregue ao controle;
6. tutorial entregue apenas a subconjunto tratado;
7. prompt divergente fora da instrução de tratamento;
8. budget divergente dentro do par;
9. contexto ou histórico reutilizado;
10. retry de falha comportamental;
11. célula despachada removida do denominador ITT;
12. identidade experimental visível ao julgador;
13. aliases inconsistentes para a mesma referência;
14. correção ou normalização de conteúdo no anonymizer;
15. rubrica alterada após abertura de `TEST`;
16. violação crítica escondida pelo composto;
17. pareamento entre perfis diferentes;
18. descarte de estrato que regrediu;
19. popularidade usada como qualidade;
20. não significância convertida em equivalência;
21. `Unknown` convertido em zero/fato;
22. uso omitido em vez de observação `UNKNOWN` justificada;
23. sessão comportamental 109 sem autorização;
24. terceira revisão integral da mesma falha sem ganho discriminatório;
25. artifact `product-only` promovido a arquitetura provada.

Cada mutação deve produzir ao menos um código de violação estável e específico,
registrado em `adversarial-cases.json`. Falha por panic, texto de log ou erro não
estruturado não conta como rejeição.

## Observabilidade mínima exigida pela suíte

Além dos campos explicitamente tipados no contrato, os tipos auxiliares públicos
devem tornar observáveis, como dados fornecidos, os fatos necessários para as
validações acima. Em particular, a suíte procura campos públicos equivalentes a
`UnavailableProfiles`, `MountedBateiaArtifacts`, `TutorialIncluded` e
`HistoryInherited`; ausência dessa fronteira observável é falha do contrato
executável, não autorização para inferir ou ignorar o ataque.

Inputs auxiliares de retry, anonimização, scoring, análise, receipt e gates são
preenchidos por JSON usando os nomes de campo públicos. A implementação pode
conter mais campos, mas não pode depender de estado externo para reconhecer os
negativos.

## Mutation score

Calcular ao fim da suíte:

```text
mutation_score = negativos_rejeitados_com_codigo_esperado / negativos_totais
```

O único valor aceito é `1.0` (`25/25`). Positivos e opacos pertencem à autoria
dos oráculos e não são criados nem contados por esta suíte adversarial.

## Restrições

- não implementar nem corrigir o candidato;
- não redefinir o contrato, a rubrica ou o veredito;
- não ler ou materializar conteúdo do split `TEST` protegido;
- não escrever positivos, opacos, receipts, selo, certificado ou laudo;
- não chamar agentes, Bateia ou qualquer serviço externo;
- preservar independência da implementação e emitir somente ataques e autoria.
