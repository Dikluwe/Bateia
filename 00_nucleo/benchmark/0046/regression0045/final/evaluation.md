STATUS: verificação FINAL após coordenador informar reversão do ranking experimental.

# Regressão independente do contrato 0045 contra candidata 0046

Aprovada no recorte congelado: recall4/4=100%, precisão4/4=100%, exatamente os mesmos resultados por caso da candidata aceita0045. Ausência, parser inválido e símbolo exato preservados. Todos os seis testes de parser e sete subcasos de referências passam (exit0). Os hashes do corpus, runner e dois arquivos de teste coincidem com manifest0045; nenhum artefato0045 foi modificado. Autor B não leu produção candidata nesta verificação.

Comandos executados a partir da raiz do repositório:

```sh
GOCACHE=/tmp/bateia-0046-cache go test ./01_core/parser ./01_core/usecase -run '^TestFunctional0045' -count=1 -v > 00_nucleo/benchmark/0046/regression0045/final/tests.log
python3 00_nucleo/benchmark/0045/compare.py capture 00_nucleo/benchmark/0046/regression0045/final/tests.log 00_nucleo/benchmark/0046/regression0045/final/candidate.json
python3 00_nucleo/benchmark/0045/compare.py compare 00_nucleo/benchmark/0045/baseline.json 00_nucleo/benchmark/0046/regression0045/final/candidate.json > 00_nucleo/benchmark/0046/regression0045/final/comparison-baseline0044.json
```

Adicionalmente, Python verificou SHA256 dos arquivos listados em manifest0045 e igualdade integral das métricas por caso entre candidate0045.json e candidate desta regressão; resultado em regression.json. A comparação com o baseline0044 mantém o ganho anterior (recall75%→100%, precisão30%→100%); não é ganho novo0046.

Escopo: seis casos sintéticos, três projetos, nove hits e primeiros três caminhos distintos; não mede resultados reais0046 nem implica equivalência semântica. A/B sem atestação física.
