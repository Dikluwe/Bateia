# Passo 0057 — Observações de comportamento por função

0056 mostrou que ação/objeto no mesmo arquivo não demonstra comportamento. Acrescentar observações AST sobre funções recuperadas, sem alterar descoberta, ranking ou prometer verificação semântica geral.

## Contrato congelado
API pura parser `InspectImplementationBehavior(report domain.ImplementationSearchReport, files []domain.ExtractedFile) []domain.ImplementationBehavior`. Novo domínio implementation_behavior.go: ImplementationBehavior com Source ImplementationReference, EndLine int, Status string, Reason string, Facts []BehaviorFact, Limitations []string. BehaviorFact com Kind string, Line int, EndLine int, Target string, Snippet string. JSON snake_case sem omitir campos necessários.

Somente hits Kind=definition de arquivo .go com SHA256 igual aos bytes realmente extraídos e referência exata à declaração FuncDecl (nome, linha inicial e final) são avaliáveis. Nome do método segue representação já utilizada na recuperação; correspondência não pode confundir métodos de receptores diferentes. Exigir arquivo único no input (duplicação ambígua => não avaliado), tamanho Content igual Size quando Size positivo, parse sem erros. Ausência/hash adulterado/language não suportada/range ou símbolo divergente => registro status not_evaluated com motivo, zero fatos. Não usar comentários, documentação, metadados ou snippets do provedor como corpo analisado.

Para função completa válida: status observed; Source path/symbol/line/SHA256 e EndLine correspondem aos bytes. Emitir fatos sintáticos em ordem do fonte: assignment para cada lado esquerdo de AssignStmt (inclui := e =), increment_decrement para IncDecStmt, return para ReturnStmt. Target é expressão LHS impressa por go/format (return usa Target vazio), Snippet é trecho exato da instrução nos bytes fonte; Line/EndLine reais. Não inspecionar corpos de FuncLit aninhadas (não são execução garantida da função exterior). Calls não são interpretadas. Fatos são sintaxe observada, não prova de execução, efeito externo, imutabilidade ou equivalência.

Limitations em TODO registro observado devem explicitar: fluxo/execução não verificados; aliases, chamadas e efeitos indiretos não verificados; ausência de escrita direta não prova preservação das entradas; relação com intenção original não verificada. Status nunca satisfied/exact/verified. Nenhuma regra de aprovação/rejeição de candidato/score depende desses fatos neste passo.

Integrar em ValidateImplementationRequirements após validação de hits e bytes, substituindo projeções anteriores: ImplementationSearchReport.Behavior []ImplementationBehavior json:behavior,omitempty. Chamada independente do parser também valida toda identidade. Não mutar entrada. Relatório Markdown exibe no máximo10 funções e8 fatos por função com contagens de omissão; JSON mantém tudo. Mostrar caráter observado e limitações junto aos fatos, incluindo zero fatos/not_evaluated. Não modificar arquivos extraídos.

EvidenceVersion evidence-v18-function-behavior e checkpoint implementation_search0057 para modos local/live. Descoberta/checkpoint discovery_query permanecem como0056. Prompts antes código, cabeçalhos Tekt.

## Verificação A/B
A produção sem ler testes/corpus; B testes a partir contrato sem patch/produção. Sem atestação física ou selo de refinamento. Controles: escrita em receptor/index/parâmetro, retorno novo versus retorno de entrada apenas como sintaxe, chamada opaca, closures, comentário enganoso, hash errado, função homônima com outro receptor/range, código inválido, C++ não avaliado, ausência de arquivo, recomputação sem fatos antigos, limites Markdown/JSON. Nenhum código-alvo executado.

Reexecutar localmente corpus0053 congelado; não repetir rede pois descoberta não mudou. Medir funções com observações e preservar métricas de recuperação. Corpus conhecido, não holdout. Testes completos/build/linteiro e hashes. Resultado não é ganho comprovado de recall nem comportamento verificado. Próximo passo dependerá dos fatos disponíveis, sem inventar análise onde não há implementação recuperada.
