# Passo 0062 — Piloto executável de contrato de efeito

## Obrigação
Verificar uma função Go autônoma que recebe []int64 e devolve []int64. Contrato sort-preserve-v1: termina normalmente nas entradas exercitadas, devolve os mesmos elementos/multiplicidades em ordem não decrescente e mantém os valores da entrada após retornar. Não prova ausência de escrita transitória, não compartilhamento de memória, terminação geral ou efeitos externos. Não alterar busca, scores ou ActionAssessments.

## Interface e observáveis congelados
Novo comando offline `bateia effects verify --source FILE --symbol NAME --output FILE`, output opcional para stdout. `--timeout` duração positiva até30s (default5s), `--seed` uint64 default1. Ajuda deve declarar recorte/limites. Argumentos inválidos falham antes de executar. Saída JSON, inclusive unknown, exit0 para avaliação concluída; erros CLI/I/O exit não zero. Despachar em L4 antes de .env/clientes remotos, injetando adaptador L3 numa função CLI. Nenhum efeito eleva conclusão da busca automaticamente.

API L1 para B testar sem produção:
- domain.EffectObservation: InputBefore, InputAfter, Output []int64; Completed bool; Failure string (JSON snake_case).
- domain.EffectVerdict: Status, Reason string (JSON snake_case).
- usecase.EvaluateSortEffect(o domain.EffectObservation) domain.EffectVerdict: Completed=false => unknown; Failure não vazio com Completed=true => counterexample; caso contrário igualdade de valores antes/depois + ordenação/multiplicidades da saída => supported_in_scope ou counterexample. Nil/vazio equivalentes quanto aos valores. Nunca mutar observação. Contrato só int64, sem tolerância numérica.
Demais tipos/interface do executor podem ser definidos por A. JSON público do relatório deve conter schema_version=effect-run-v1, contract_id=sort-preserve-v1, status (supported_in_scope/counterexample/unknown), reason, source_sha256, source_path, symbol, seed, cases (lista de {observation,verdict}), limitations (lista), execution (metadados). Oráculo deve ser algoritmo/expressão independente da função candidata. Qualquer contraexemplo torna status global counterexample; ausência de contraexemplo com alguma execução incompleta => unknown; supported somente com todos os casos executados completos e válidos. Nenhum resultado vazio vira sucesso.

## Executor real e recorte
L3 recebe fonte local por bytes, congela cópia num diretório temporário com módulo independente. Fonte deve ser package candidate, símbolo exportado de assinatura exata func([]int64) []int64 (sem receiver ou type params), apenas declarações de funções (helpers permitidos), sem init, sem variáveis/tipos globais, sem diretivas //go:. Imports permitidos somente sort e slices; rejeitar demais, import aliases e dot/blank imports. Código fora do recorte=>unknown sem execução. Não seguir dependências do projeto original. Compilar wrapper próprio que chama a função e captura cópias independentes antes/depois e retorno; não confiar em observações fornecidas pelo candidato. Panic é contraexemplo nas entradas válidas, com recuperação instrumentada. Protocolos malformados/saída excedida/processo interrompido=>unknown. Timeout=>unknown, nunca sucesso.

Executar compilação e subprocesso com contexto/cancelamento e orçamento de saída (máximo1MiB por stdout/stderr), sem shell, em diretório temporário. Ambiente mínimo sem credenciais, GOPROXY=off, GOSUMDB=off, GOWORK=off, GOENV=off, GOTOOLCHAIN=local; caches/temporários graváveis. Build pode ter orçamento separado de até60s explicitado. Diretório removido ao terminar. Processo separado não é sandbox de segurança ou limite rígido de memória; declarar que não é para código hostil. Executar nesta tarefa somente fixtures locais do piloto.

Casos: conjunto fixo com vazio, unitário, repetidos, negativos, já ordenado, reverso e extremos int64; mais32 vetores pseudoaleatórios determinísticos de tamanho0..8 e valores-8..8, seed registrada. Todos os vetores e observações devem constar no recibo para reprodução sem depender do gerador. Wrapper e fonte efetivamente compilados: conteúdo e hashes no execution; hash binário, versão Go, GOOS/GOARCH, argv sem segredos, limites, motivo de parada e duração identificada como telemetria da ferramenta. Não alegar benchmark do candidato. Não executar go test do repositório analisado.

## Regime, papéis e gate anterior ao candidato
Ensaio do protocolo Tekt completo, sem atestação física de isolamento, ambiente compartilhado. Integrador escreve obrigação/contrato e executa verificação sem editar implementação ou oráculos. B cria primeiro gate independente com observações positivas, mutações semanticamente negativas e opacas, expectativa e oráculo próprio; cada negativo deve ser counterexample, positivo supported_in_scope e opaco unknown. B também fornece fixtures Go autônomas (referência e mutantes) para validar depois o executor real. Nenhum algoritmo do candidato A existe antes do gate. Gate discriminatório exige todos controles corretos; seu selo cobre apenas esse fragmento observável, não o executor futuro. Congelar contrato/baseline/oráculo/fixtures antes de A. B continua testes independentes sem ler A. A lê contrato e produção pertinente, não gate/testes/fixtures/relatórios. A atualiza prompts antes de código. Integrador roda testes/gates, não altera expectativas para fabricar sucesso.

Budget de calibração: até2 revisões focais por mesma causa; sem ganho observável, rever contrato em vez de expandir vocabulário. Não executar corpus de busca, pois sem alterações na recuperação. Verificar build, suíte Go, lint, teste negativo/positivo/opaco do executor, repetição de seed e rejeição de execução fora do recorte. Preservar todos os recibos e revisões.

## Escopo de arquivos
A: novos domain/effect_verification.go, usecase/effect_verification.go, infra/effects/runner.go, shell/cli/effects.go; integração wiring/main.go e cli/app.go (somente ajuda); prompts correspondentes. Caminhos completos sob01_core,03_infra,02_shell,04_wiring. B: testes companheiros dos novos módulos e prompts, fixtures fora do código produtivo. Não alterar legado para fabricar aprovação. Incluir documentação de uso e limitações no experimento. Nenhum commit solicitado.

### Assinaturas públicas para integração independente (antes de A)
- domain.EffectRequest {SourcePath string; Symbol string; Seed uint64; Timeout time.Duration}.
- domain.EffectCase {Observation EffectObservation; Verdict EffectVerdict}.
- domain.EffectReport {SchemaVersion, ContractID, Status, Reason, SourceSHA256, SourcePath, Symbol string; Seed uint64; Cases []EffectCase; Limitations []string; Execution map[string]any}. JSON snake_case.
- domain.EffectVerifier interface { Verify(context.Context, EffectRequest) (EffectReport, error) }.
- effects.NewRunner() *Runner; método Verify satisfaz interface.
- cli.RunEffectsCommand(args []string, out io.Writer, verifier domain.EffectVerifier) error; args começa com verify, -h/--help devem funcionar sem executar.
Cache de compilação pode ser reutilizado em diretório controlado sob os.TempDir, sem incorporar demais variáveis do ambiente. Execução da candidata usa diretório e ambiente separados.
