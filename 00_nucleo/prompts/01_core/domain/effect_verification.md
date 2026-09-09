# Prompt: 01_core/domain/effect_verification

Tipos públicos EffectObservation, EffectVerdict, EffectRequest, EffectCase, EffectReport e EffectVerifier conforme assinaturas do passo 0062. Campos JSON snake_case; nenhuma operação de I/O.

Contrato autoritativo: `00_nucleo/passos/0062-verificar-efeito-executado.md`. Implementação A; testes e veredito independentes pelo protocolo segregado.

Passo 0063: EffectRequest acrescenta ExpectedSourceSHA256 (expected_source_sha256). EffectReport acrescenta SourceBinding (source_binding), com Status, ExpectedSHA256, ActualSHA256, Symbol e Declaration *EffectDeclaration, sempre serializada (null quando indisponível). EffectDeclaration contém StartByte/EndByte, StartLine/EndLine, Content e SHA256, JSON snake_case. Coordenadas físicas: bytes zero-based [início,fim) e linhas one-based inclusivas. Contrato adicional autoritativo: `00_nucleo/passos/0063-vincular-fonte-ao-efeito.md`.

Passo 0064: acrescentar EffectAnalyzerIdentity, EffectTraceInputs e EffectTrace, todos campos JSON snake_case sem omitempty, e EffectReport.Trace. Identidade contém Status, Reason, BinarySHA256, ModulePath, ModuleVersion, GoVersion, VCSRevision string e VCSModified *bool (ausência null). Inputs contém ContractID, AnalyzerBinarySHA256, SourceSHA256, Symbol string, Seed uint64 e ExecutionTimeoutNS, BuildTimeoutNS, StdoutLimitBytes, StderrLimitBytes int64. Trace contém SchemaVersion, Status, Reason, StartedAt, FinishedAt, RequestSHA256 string, Analyzer e Inputs. Contrato autoritativo: `00_nucleo/passos/0064-identificar-analisador-e-parametros.md`.
