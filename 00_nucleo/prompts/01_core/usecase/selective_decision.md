# Prompt — Avaliar abstenção por risco-cobertura

Avalie somente casos adjudicados. Cobertura = decisões conclusivas / casos
adjudicados; risco seletivo = conclusões erradas / conclusões emitidas. Risco sem
conclusões é indefinido. Julgamentos aceitáveis contam como acertos. Diferencie
abstenções corretas das evitáveis sem tratá-las como erro seletivo. Gere pontos da
curva somente quando toda conclusão tiver score explícito, em limiares únicos,
sem dividir empates e sem extrapolar para além da cobertura operacional.
Indique intervalo superior unilateral aproximado de Wilson e amostras pequenas.
Rejeite IDs duplicados, decisões inválidas, score fora de [0,1], split inválido e
schema não suportado. Não use I/O nem ajuste limiares neste avaliador.
