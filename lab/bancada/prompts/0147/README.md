# Bancada0147 — fechamento de geração e revisão de contratos

O [passo0147](../../../../00_nucleo/passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md) declara a intenção antes da correção. O container L1 possui o limite; feedback, admissão de discover e admissão de expand consomem a mesma guarda. Snapshot999999 continua válido para inspeção, mas não para avanço. Nenhuma regra é transferida à CLI e os schemas/algoritmos de ID permanecem.

## Artefatos e verificações

- `baseline.json` e `before/`: identidade dos arquivos no início desta fatia e cópias dos proprietários alterados. O checkout já continha mudanças anteriores; comparar contra Git HEAD confundiria escopos.
- `generation-before.log`: regressões antes da correção. Feedback devolvia sucesso com1000000; discover executava buscas mesmo com geração999999. Expand já rejeitava corretamente. A falha é o controle negativo do defeito, não resultado esperado do produto corrigido.
- `generation-after.log` e `core-regressions.log`: último avanço999998→999999, rejeição no limite e memória revalidável, sem mutação da entrada. O teste de feedback compara todos os campos restantes, usando entrada com Lead e LastRead.
- `selection-and-sources.log` e `sources-with-pending-destination.log`: desempates de geração/referência/ID em seis permutações, com distância ausente ou empatada; três tipos de origem, documentos de32 KiB e32 KiB+1, oito/nove documentos e128 KiB/128 KiB+1 agregados. Recibo/LastRead/duplicidade inconsistentes rejeitados antes de ler um destino elegível pendente. Os oráculos não alteram o comparador ou o validador de fontes.
- `cli-generation.log`: feedback no último avanço e no limite, publicação revalidável versus rejeição sem stdout/arquivo, input preservado; discover rejeita antes da factory.
- `validate.py` e `scope-and-docs.json`: links/âncoras, linhagem para prompts modificados, recorte permitido, inventário Go e preservação das regras e bancadas anteriores.
- `replay_generation.py` e `generation-replay/`: aplicam o binário corrigido à entrada sintética preservada de0146 e a um controle999998. Comandos são offline; o ensaio não modifica0146 ou o binário antigo.
- `final/`: logs, executável e recibo do gate canônico desta árvore.

Reproduzir os consumidores com `go test` dos pacotes usecase e cli, executar `python3 lab/bancada/prompts/0147/validate.py` e rodar `BATEIA_VERIFY_OUTPUT=lab/bancada/prompts/0147/final make verify`. O runner ajusta assinatura Git apenas no processo; fixtures podem exigir loopback. O build é gravado no Lab pelo runner, preservando a alternativa autorizada a /dev/shm. Replays usam diretório novo para conservar os próprios recibos.

## Limites da conclusão

A correção fecha MEM-04 e os novos consumidores ampliam MEM-02/ORD-02/SRC-01. Não demonstram toda combinação de fontes, limites de capacidade em bytes da memória, todos os erros de transporte ou utilidade/novidade de insights. A matriz mantém esses limites explícitos. A história0146 e seus scripts continuam sendo evidência daquela versão; não executar seu materializador para sobrescrever os contratos posteriores.
