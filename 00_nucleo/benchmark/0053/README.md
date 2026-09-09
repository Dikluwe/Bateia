# Avaliação 0053

Executor por estágios, sem alterar algoritmo. Corpus e executor têm congelamentos separados; perguntas novas não foram entregues ao autor do executor durante sua implementação. Não há atestação física de isolamento.

## Reexecutar

A partir da raiz do Bateia, compile o programa e o auxiliar em diretório temporário. `probe.go.txt` é um auxiliar de experimento; não é código de produção L4.

```sh
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0053-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0053 ./04_wiring
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0053-probe /tmp/bateia0053-probe.go
python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0053 --probe /tmp/bateia0053-probe --output /tmp/bateia0053-new-run
```

O diretório de saída deve ser novo. Os caminhos dos snapshots constam no corpus; os arquivos precisam ter os mesmos hashes. Mudanças ou arquivos adicionais invalidam a elegibilidade da medição. Não troque o conteúdo do corpus para acomodar outro snapshot; crie outra avaliação.

Sem `--discovery`, não há descoberta pública: a etapa fica `not_evaluated`. Para executar também as três buscas congeladas, acrescente `--discovery`. Essa etapa usa rede e pode exigir permissão do ambiente. Sourcegraph pode relaxar a consulta; tentativas e relaxamentos são preservados. Erro do provedor não equivale a zero candidatos. Segredos de autenticação vêm do ambiente e não devem ser copiados para recibos.

A análise local é feita pela CLI real; o auxiliar mede a recuperação anterior à seleção. Resultados locais sobre snapshots fixos e descoberta pública são observações separadas, não uma jornada remota integral reproduzida no mesmo commit. Campos de comportamento ficam não avaliados: nenhum código dos projetos é executado. Contagem de indícios em negativas não é automaticamente contagem de capacidades comprovadas incorretamente.

## Reavaliar as respostas já coletadas

Sem repetir rede/CLI, mantendo os snapshots disponíveis e os mesmos hashes:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/replay_evaluation.py --input 00_nucleo/benchmark/0053/public-collected --output /tmp/bateia0053-recomputed
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/test_runner.py
```

A revisão r2 do avaliador corrige a classificação de erro em tentativa de descoberta. `public-collected/summary.json` preserva a contagem anterior; `public-final/summary.json` é a contagem corrigida, derivada dos mesmos bytes. Não compare tempos como medida de desempenho dos projetos. `replication.json` registra igualdade das métricas locais entre as duas coletas.
