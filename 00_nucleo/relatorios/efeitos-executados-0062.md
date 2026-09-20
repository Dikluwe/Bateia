# Passo 0062 — Evidência de efeito executado

Implementado `bateia effects verify`: primeiro piloto offline de execução instrumentada. O contrato `sort-preserve-v1` exige saída ordenada, mesmas multiplicidades e valores da entrada preservados após retornar. O comando aceita um arquivo Go autônomo e uma função exportada `func([]int64) []int64`; não executa automaticamente repositórios recuperados pela busca.

## Resultado observado

Oito fixtures congeladas antes da implementação foram executadas pela CLI pública e seus recibos foram conferidos por um oráculo Python independente. Todas tiveram o status esperado:

| Fixture | Status observado | Evidência |
| --- | --- | --- |
| referência | supported_in_scope | 39 chamadas completas satisfazem o contrato |
| ordem reversa | counterexample | saída [2,2,2,1] para entrada [2,2,1,2] |
| multiconjunto errado | counterexample | entrada [7], saída [42] |
| mutação da entrada | counterexample | entrada passou de [2,2,1,2] para [1,2,2,2] |
| retorno sem ordenar | counterexample | saída [2,2,1,2] |
| panic | counterexample | panic capturado na chamada com entrada vazia |
| laço infinito | unknown | timeout; nenhuma chamada completa |
| import proibido | unknown | recorte AST rejeitado antes da compilação/execução |

As seis fixtures executadas até o fim registraram 39 observações cada. Os dois resultados unknown conservam os 39 casos planejados como incompletos: um caso listado não significa execução concluída. Fonte, wrapper, módulo, binário, toolchain, plataforma, argumentos e orçamentos estão identificados nos recibos quando essa etapa foi alcançada. Tempos são telemetria da ferramenta, não benchmarks dos candidatos.

A repetição da referência com seed 1 preservou os 39 casos/observações e os hashes de fonte, wrapper, módulo e binário, além da versão Go. Comparação em `../../lab/bancada/history/benchmarks/0062/reproduction.json`. Recibos e conferência por fixture em `../../lab/bancada/history/benchmarks/0062/executions/verification.json` e arquivos vizinhos.

## Verificação e autoria

O gate anterior à implementação distinguiu 6 controles positivos, 10 violações e 4 observações opacas: 20/20. A rejeição de 10/10 negativos vale apenas para esse conjunto controlado. O executor real distinguiu 5/5 implementações defeituosas; não é uma estimativa de detecção geral.

Contrato, oráculo/fixtures, implementação e testes de execução tiveram autores separados. Testes do núcleo incluem controles congelados e 121 saídas pequenas enumeradas; testes de CLI/infra incluem argumentos inválidos, reprodução, mutação, panic, timeout, cancelamento e rejeição fora do recorte. Congelamentos foram preservados antes da execução, sem ajustar expectativas para aprovar o candidato. Os contextos foram segregados, mas o filesystem é compartilhado: ensaio Tekt sem atestação física de isolamento.

Auditoria da implementação confirmou ambiente explícito de compilação e execução, sem herdar `os.Environ`, execução sem shell e despacho da CLI antes de carregar `.env`/clientes. Isso é evidência de leitura do código; não foi usada uma candidata que inspeciona credenciais, pois esse import está fora do recorte. Subprocessos têm orçamento de tempo e saída e diretórios temporários removidos. Não há sandbox de segurança nem limite rígido de memória para código hostil.

## Como interpretar

- **Fato observado:** entradas, valores antes/depois, saída e panic de chamadas completas, vinculados aos bytes compilados.
- **Inferência limitada:** supported_in_scope significa apenas que os casos executados satisfizeram o contrato.
- **Hipótese:** a função pode cumprir a mesma propriedade para outras entradas. Este piloto não a prova.
- **Não verificado:** execução incompleta, fora do recorte, efeitos externos, escritas transitórias, compartilhamento de memória e terminação geral.

O status global também considera `execution.stop_reason`. Uma falha de processo/protocolo posterior a chamadas válidas torna a execução unknown; não reescreve chamadas válidas como violações. Contraexemplo já observado prevalece. O oráculo do experimento confere os casos/expectativas das oito fixtures; não cobre exaustivamente todos os modos de falha do protocolo.

## Uso

```sh
go build -o /tmp/bateia-0062 ./04_wiring
/tmp/bateia-0062 effects verify --source /caminho/candidate.go --symbol Transform --seed 1 --timeout 2s --output /tmp/efeito.json
```

Consultar `status`, `cases`, `execution` e `limitations`; exit code zero também ocorre em counterexample e unknown. Fonte deve declarar package candidate, apenas funções e imports sort/slices sem aliases; demais dependências não são seguidas.

A busca e os scores permanecem fora deste experimento. O ganho demonstrado é conseguir refutar efeitos concretos e preservar evidência reproduzível num recorte pequeno, não ter resolvido entendimento semântico geral. Antes de conectar o recurso à busca, o próximo contrato deve definir como ligar a implementação extraída à função realmente compilada, sem perder dependências nem ampliar a conclusão além da cobertura.

Validação final: build, testes novos sem cache, suíte Go completa (com cache nos pacotes inalterados), crystalline-lint e git diff --check passaram. Contrato, implementação, oráculo e testes mantiveram seus hashes congelados. Comandos e limites em `../../lab/bancada/history/benchmarks/0062/verification.json`; logs preservados no mesmo diretório.
