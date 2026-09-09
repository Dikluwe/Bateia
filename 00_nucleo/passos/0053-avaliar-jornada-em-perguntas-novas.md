# Passo 0053 — Avaliar perdas por estágio em perguntas novas

2026-09-08. Entregar corpus congelado e executor reproduzível, sem alterar ranking, catálogo, seleção ou julgamento do Bateia para favorecer resultados.

## Regime e autoria

Tekt A/B sem atestação física: B escolhe novos projetos/perguntas/esperados sem ler ranking/catálogos/testes/benchmarks anteriores; A implementa executor a partir do schema sem ler o corpus até implementação congelada; integrador escreve contrato, auxiliares de execução e verifica métricas em controles sintéticos. Filesystem compartilhado, sem selo de refinamento. Benchmark tem três projetos Go de domínios distintos, diferentes de Bateia/Coin/Snappy/bbolt/oauth2/uuid/Bleve, dois positivos funcionais em português e uma negativa por projeto. Corpus de conveniência do cache de módulos, não amostra representativa.

B lê fontes somente dos projetos candidatos para formular casos, registra versão/licença/proveniência, hashes e intervalos que sustentam os alvos. Congelar corpus antes de executar Bateia. Não escolher consultas por palavras retornadas pelo Bateia. Não exigir nome de símbolo na consulta funcional. Perguntas negativas pedem funcionalidade de outro domínio; resultados brutos serão chamados indícios na negativa, sem afirmar automaticamente violação de uma capacidade.

## Schema corpus.json (0053-v1)

projects: lista {id, local_path absoluto, repo owner/repo, version, license, provenance, files:[{path relativo, sha256}]}. Listar todos arquivos regulares do snapshot (ou declarar recorte explícito em provenance); fonte não executada. cases: lista {id, project_id, query, polarity:positive|negative, expected:[{path,symbol}], rationale, evidence:[{path,start_line,end_line,sha256}], discovery:boolean}. Exatamente uma positiva por projeto marcada discovery. Negativa expected vazia. Positiva deve conter pelo menos um alvo; arquivo/símbolo exato observado, não apenas caminho genérico. freeze.json registra SHA256 de corpus.json antes de qualquer execução, inputs autor e data.

## Executor e observáveis

Implementar runner.py em benchmark/0053 com --corpus --freeze --binary --probe --output e --discovery opcional. Diretório de saída novo obrigatório. Probe é executável auxiliar: `probe local PATH QUERY` retorna ImplementationSearchReport JSON anterior à seleção; `probe discovery QUERY` retorna {candidates,attempts,error?}, máximo10 candidatos, contexto30s. CLI local: `binary analyze -path PATH -query QUERY -format json -output arquivo`. Falhas/timeout/JSON inválido são unknown, nunca zero sucesso. CLI limitado a120s por caso e probe local60s/discovery45s. Nenhum shell=True. Runner não altera corpus ou código.

Antes e depois de cada caso validar inventário/hash do snapshot contra corpus. Diferença torna caso unknown e não elegível para métricas. Hash de executáveis/args/corpus e saídas em receipt. Detecção de arquivo encontrado usa raw.examined; função recuperada usa raw.hits correspondendo path+symbol; arquivo extraído usa results.files de CLI; função confirmada usa results.implementation_search.hits e SHA do arquivo; corpo inteiro disponível exige EndLine válido no conteúdo extraído, hash e bytes iguais ao snapshot. Canais lexical/documentation/direct_call registrados separadamente. Top3/12 contam caminhos distintos no raw, não slots hipotéticos de extração. Além disso registrar recall de símbolos exatos, sem trocar por arquivo.

Conclusão sustentada: registrar estados/reasons de cobertura e requisitos, separando indício de prova. `verified_behavior` será not_evaluated porque não há execução de projeto ou oracle comportamental. O executor não converterá evidence_located/score/arquivo presente em satisfação. Nenhuma extração ou fonte equivale a execução.

Descoberta pública opcional só nas três perguntas discovery; registrar candidatos e tentativas/fallbacks. Descoberta não executada fica not_evaluated. Erro fica unknown. Expected repo encontrado é métrica top10 de identidade de projeto; commit do provedor não é automaticamente igual ao snapshot local. Não multiplicar descoberta pública por análise local e chamar isso de sucesso ponta a ponta reproduzido. O relatório deve expor separadamente esses observáveis. Para ao menos um caso executar jornada remota integrada se descoberta trouxer candidato pertinente e rede permitir; do contrário declarar etapa integrada não avaliada.

## Métricas, testes e entrega

Contagens com numerador/denominador de positivos avaliáveis; casos unknown separados. Negativas com indícios contadas separadamente de conclusões indevidas. Controle sintético do avaliador: erro não conta zero, path sem símbolo não conta símbolo, hash adulterado invalida, hit documental não conta comportamento, descoberta omitida não conta falha.

Congelar runner antes de ler corpus; rodar todos9 casos, executar três descobertas com permissão de rede se disponível, publicar recibos/JSON e relatório com perda dominante e próximo reparo indicado pelos dados. Nenhum alvo/pergunta muda após execução. Rerun idêntico offline deve reproduzir métricas locais sem comparar timestamps/telemetria. Documentar instruções de reexecução. Build/lint/testes pertinentes; algoritmo de produção deve conservar hashes da baseline.
