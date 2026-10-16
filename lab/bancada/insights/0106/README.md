# 0106 — Ensaio real com GitHub

Em 2026-10-06, o coletor do Bateia consultou publicamente `sourcegraph/conc`,
resolveu HEAD para `5f936abd7ae87036af1f75c95fb9d0daaf00116b` e coletou
29 itens/95.340 bytes, sem falhas. `complete=true` descreve o escopo e os
classificadores do coletor, não auditoria de todos os arquivos do repositório.
Não foram usados tokens, clone, instalação ou execução de código remoto.

A fonte permite investigar reserva de capacidade antes de iniciar workers.
A avaliação é manual e PARTIAL: WithMaxGoroutines configura teto de workers,
mas o padrão é ilimitado; o canal de tarefas não tem buffer, a submissão
bloqueia e a assinatura Go não recebe contexto de cancelamento. Não é uma
fila durável com rejeição explícita. O teste da fonte observa concorrência
máxima; seus bytes foram inspecionados, sem execução do teste.

O catálogo agrega esses testemunhos em BOUNDED_QUEUE_BACKPRESSURE PARTIAL.
As descrições de estratégia no relatório vêm do catálogo e não descrevem
integralmente a implementação inspecionada. `source-review.json` conserva
as diferenças concretas. A hipótese permanece INCONCLUSIVE, alvo UNKNOWN,
REFERENCE_ONLY e NOT_AN_ADOPTION_DECISION. O efeito sobre o Bateia ainda não
foi medido; avaliações fornecidas não são refeitas pelo comando.

## Artefatos

- `github-plan.json`: plano inicial com HEAD; é registro da coleta, não fonte imutável.
- `github-plan-pinned.json`: plano resolvido para reprodução online.
- `github-bundles.json`: bytes, commits, hashes e recibos do coletor.
- `github-insight-input.json`: bundle congelado e avaliações manuais localizadas.
- `source-review.json`: observações, condicionantes e limites da classificação.
- `github-insight-result.json`: saída original, reproduzível byte a byte.
- `online-check.json`: releitura online dos dois arquivos testemunha e controles.
- `go-test.log`, `go-race.log`, `verification.json`: verificações desta árvore.

## Reprodução

Na raiz do Bateia, compilar e analisar sem rede:

```sh
go build -o /tmp/bateia-insights-0106 ./04_wiring
python3 lab/bancada/insights/0106/verify.py --binary /tmp/bateia-insights-0106
```

Adicionar `--online` relê somente os dois arquivos testemunha do SHA fixado
no GitHub, compara bytes e executa os mesmos controles locais. O script
verifica todos os hashes, localizadores, replay, recusa de adulteração com
preservação da saída e controle sem avaliações (INCONCLUSIVE). Não executa
o código do repositório consultado. Requer rede pública e Python 3.9+.

Para repetir a coleta completa, usar um diretório temporário sem `.env`,
sem tokens, e uma saída nova, preservando os registros originais:

```sh
mkdir -p /tmp/bateia-online-0106
cd /tmp/bateia-online-0106
env -u GITHUB_TOKEN -u SOURCEGRAPH_TOKEN /tmp/bateia-insights-0106 idea evidence collect \
  --plan /repos/Antigravity/Bateia/lab/bancada/insights/0106/github-plan-pinned.json \
  --output /tmp/bateia-online-0106/recollected.json --timeout 120s
```

Metadados consultados podem variar; os dois arquivos testemunha no SHA
fixado devem conservar hashes. Inspecionar complete/failures/targets, pois
o coletor também pode retornar JSON com falhas sem exit code diferente de zero.
Limites da API pública e indisponibilidade de rede não validam uma hipótese.

Os testes do Bateia foram executados localmente. Este ensaio usa o GitHub
como fonte pública real; não é uma execução de GitHub Actions.
