# Runner 0053 — contrato de autoria A

Autor A implementa usando apenas o schema do passo 0053 e tipos JSON de domínio.
Não lê corpus, freeze, fontes de projetos nem testes/benchmarks antigos. Este documento
precede o código. Não modifica ranking, catálogo, seleção ou julgamento de produção.

CLI: `python3 benchmark/0053/runner.py --corpus CORPUS --freeze FREEZE --binary BIN --probe PROBE --output NOVO_DIR [--discovery]`.
FREEZE deve ter `corpus_sha256`. Executáveis são caminhos; não se usa shell.
CLI local tem timeout 120s; probe local 60s; descoberta 45s. Saídas preservadas por
caso, stdout/stderr e recibos de comando inclusive em erro. Snapshot é inventariado
com hashes antes e depois de cada caso, todos arquivos regulares; symlink é recusado.
Hashes de corpus, executáveis e argumentos são registrados. Runner não executa fontes.

API pura pública para verificação independente:
`evaluate_case(case, project, raw, report, snapshot, *, errors=None, inventory_valid=True, discovery=None)`
e `summarize(cases)`. snapshot mapeia caminhos relativos para bytes exatos; discovery
omitido significa not_evaluated. Erro, report malformado, hash de evidência adulterado,
inventário diferente ou falha de execução torna local unknown, fora do denominador.
Resultados parciais/limitados são observáveis finitos, com status/reason preservados.

Cada alvo positivo é caminho+símbolo exato. Métricas por caminho e símbolo são separadas:
arquivo examinado (raw.examined), símbolo recuperado (raw.hits), top3/12 caminhos
distintos na ordem raw.hits, arquivo extraído (results.files), símbolo confirmado
(results.implementation_search.hits com hash validado), corpo disponível (linha final
válida e bytes extraídos iguais ao snapshot com hash). 'Confirmado' é localização de
símbolo, não satisfação funcional. Canais lexical/documentation/direct_call são
separados; canal ausente é lexical, canais desconhecidos preservados como other.

Coverage, requisitos, reasons e estados originais são preservados; satisfação não é
inferida de score/evidence_located. verified_behavior sempre not_evaluated. Negativas
registram indícios e alegações de satisfação explícitas separadamente, sem classificar
indício como violação. Descoberta registra candidatos/tentativas e acerto de identidade
owner/repo nos primeiros 10, sem atestar commit ou sucesso remoto integrado. A execução
de jornada remota integrada é responsabilidade do integrador, documentada separadamente.

cases.json contém observações individuais e receipts.json inputs, comandos, hashes e
inventários. summary.json contém numeradores/denominadores de positivos avaliáveis,
recall de símbolos e caminhos, unknown separado, negativas e descoberta separadas.
Reexecuções offline devem comparar summary.json (sem timestamps/telemetria).
Após implementação, runner-freeze.json congela SHA256 de runner.py e deste contrato.
