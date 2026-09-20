# Piloto 0062 — Ordenação com preservação dos valores de entrada

Contrato: passos/0062-verificar-efeito-executado.md. O comando público recebe um arquivo Go autônomo package candidate e uma função exportada func([]int64) []int64. Apenas imports sort/slices e declarações de funções no recorte. Não carrega nem executa projetos encontrados pela busca.

O resultado supported_in_scope refere-se somente às entradas executadas. Não comprova ausência de escrita transitória, não compartilhamento de memória, terminação geral ou efeitos externos. Processo separado e ambiente mínimo não constituem sandbox de segurança para código hostil.

## Cadeia registrada

- baseline.json identifica estado anterior; manifest.json registra papéis e limites.
- gate/ contém cópias documentais dos artefatos B; arquivos Go/Python terminam em .txt para não virarem código produtivo L0.
- gate-export.json mapeia originais temporários para cópias e hashes.
- gate-receipt.json e gate-seal.json registram o gate anterior a A: 6 positivos, 10 negativos e 4 opacos; selo limitado às observações, sem atestação física.
- Gate não comprova que wrapper/executor capturam a execução. Essa parte é verificada depois por testes independentes e execução das oito fixtures Go.

Para restaurar artefatos do gate, copiar cada stored de gate-export.json para original (remover somente sufixo .txt segundo o mapeamento). Conferir gate-freeze.sha256 antes de usar. Nenhuma expectativa pode ser alterada para fabricar aprovação.

Uso pretendido após build a partir da raiz:

```sh
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0062 ./04_wiring
/tmp/bateia-0062 effects verify --source /tmp/bateia0062-gate/fixtures/reference/candidate.go --symbol Transform --seed 1 --timeout 2s --output /tmp/effect-reference.json
python3 /tmp/bateia0062-gate/oracle.py --report /tmp/effect-reference.json --expected supported_in_scope
```

O recibo conserva entradas e estados observados, identidade do código, wrapper e compilador, limites e telemetria. Contraexemplo prevalece se uma execução completa viola a propriedade; timeout/incompletude não viram sucesso. Exit code zero da ferramenta não significa que a propriedade foi satisfeita: consultar status.

Sem alterações na recuperação ou avaliações semânticas de busca. Não confundir a rejeição dos mutantes controlados com taxa de detecção geral. Tempos são telemetria do executor, não benchmark do candidato.

Execução concluída: 8/8 fixtures com status esperado, referência repetida com 39 observações iguais, suíte completa/build/lint aprovados. Ver verification.json, reproduction.json e executions/. O script documental run-fixtures.py.txt reproduz a execução após restaurar o gate; passar binário e um diretório novo. A separação de autores é procedimental e não atesta isolamento físico.
