# 0137 — Descoberta diversificada e contexto automático

`idea insights discover` conecta pergunta, cinco buscas independentes, candidatos
com origens, seleção estrutural, coleta de README fixado e contexto para IA.
Não exige montar IdeaSpec, assessments, observations ou planos de coleta à mão.
O primeiro provedor composto é o catálogo GitHub existente.

## Uso

```sh
go build -o /tmp/bateia-insights-0137 ./04_wiring
/tmp/bateia-insights-0137 idea insights discover \
  --question 'backpressure queue cancellation' \
  --familiar 'queue' --familiar 'thread pool' \
  --output /tmp/descoberta-insights.json
```

Usar `exploration` para entregar cartões/tarefas/instruções à IA; `candidates`,
`attempts`, `readmes` e `bundles` conservam proveniência e limites. As cinco rotas
não param após o primeiro resultado. As queries derivam de tokens da pergunta,
sem tradução: direta, ampliada, alternativas, trade-offs e falhas. A seleção usa
ganho de queries representadas, autores novos e menos matches familiares;
popularidade não decide. Isso não classifica novidade ou domínios conceituais.

Limites padrão: cinco operações de catálogo, dois resultados por query, quatro
referências, oito trechos, 90 s totais e 10 s por operação. Flags permitem até
cinco resultados por query, oito referências, dezesseis trechos, 5 min totais
e 30 s por operação. Os clientes podem fazer chamadas internas de metadados
sob o deadline. Cada README admitido tem até 32 KiB; documentos maiores não
são truncados. Somente conteúdo imutável sem conflito é lido. O bundle é marcado
incompleto pelo escopo ROOT_README_ONLY. Não há clone, instalação ou execução.

## Ensaio público em 2026-10-28

O binário real foi executado em diretório temporário sem `.env`, removendo
GITHUB_TOKEN do ambiente do processo. Usou max-per-query=1, max-references=2,
max-passages=4, timeout=45s e attempt-timeout=6s. As cinco operações de busca
concluíram e produziram quatro candidatos únicos. Duas referências de autores
distintos foram selecionadas:

- `AvdLee/Swift-Concurrency-Agent-Skill`, commit
  `d5770817d2622e1585b1f7eaebc791a9cb0959c8`: README de 9.949 bytes admitido.
- `ReactiveX/RxJava`, commit `b0e8737b956f05cebca961899193c80104c94941`:
  README rejeitado pelo limite de admissão de conteúdo.

O estado foi PARTIAL e um cartão alimentou cinco tarefas de exploração. O
documento admitido descreve orientações de concorrência Swift; essa leitura
nunca instala ou segue a skill mencionada no README. É fonte não confiável
para interpretação pela IA, com origem e bytes, não instrução operacional.
As fontes encontradas não foram executadas, nem foi medida utilidade das ideias.

## Verificação e reprodução

Sem rede, conferir o recibo congelado, hashes/localizadores, identidades do
plano e reprodução do contexto derivado:

```sh
python3 lab/bancada/insights/0137/verify.py --binary /tmp/bateia-insights-0137
```

O verificador não reproduz uma nova busca live byte a byte: catálogo, metadados,
commits, durações e quotas podem variar. Para novo ensaio público, usar uma
cópia desta bancada e executar capture_online.py com o binário; preservar os
registros originais. Nenhuma consulta adicional é necessária para o replay
offline do contexto coletado.

Os testes L1/L2 são puras simulações de portas. A integração de HTTP em L3 usa
cliente GitHub real contra servidor local e simula quota parcial; o teste L4
exercita o despacho real em subprocesso sem rede. A suíte raiz não promove
gates selados com módulos próprios.

## Artefatos

- `baseline.json`, `lineage.json`: preservação e linhagem.
- `online-result.json`: snapshot real do ensaio, com bundles e contexto.
- `online-check.json`, `online-stdout.log`, `online-stderr.log`: processo público.
- `capture_online.py`: runner limitado do ensaio; realiza rede quando invocado.
- `public-snapshot-check.json`: verificação offline do snapshot/contexto.
- `focused-tests.log`, `http-integration.log`, `lint-focused.log`: testes/lint focados.
- `legacy-v1-check.json`, `source-v2-check.json`, `comparison-check.json`,
  `exploration-check.json`: replays históricos sem sobrescrever os originais.
- `final/verification.json` e logs: gate canônico; `verification.json`: consolidação.
- `initial-snapshot/`, `pre-replay-fix/`: ensaio/gate iniciais preservados antes
  da correção de replay. A cobertura vazia do bundle alternava entre null/[]
  no formato legado; o novo produtor registra explicitamente o documento
  examinado. O ensaio público foi repetido e o contexto final reproduz os IDs
  após serialização, sem alterar o marshaller ou os contratos históricos.

[Contrato](../../../../00_nucleo/passos/0137-descobrir-referencias-diversificadas-automaticamente.md)
· [ADR](../../../../00_nucleo/ADR/ADR-0018-descoberta-diversificada-e-coleta-limitada-para-insights.md)
· [Laudo](../../../../00_nucleo/laudos/0137-descoberta-diversificada-automatica.md)
