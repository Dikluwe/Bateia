# Testes independentes — domínio da verificação de efeito

Raiz causal: `00_nucleo/passos/0062-verificar-efeito-executado.md`, contrato congelado sort-preserve-v1. Autor B, posterior ao gate discriminatório e sem leitura da implementação A. Materializar testes externos `domain_test` em `01_core/domain/effect_verification_test.go`.

Verificar campos e nomes JSON públicos de EffectObservation, EffectVerdict, EffectCase e EffectReport, inclusive preservação exata de extremos int64 via round-trip tipado. Verificar compatibilidade da interface EffectVerifier com context.Context, EffectRequest e retorno EffectReport/error. Os testes não executam processos, I/O de arquivos ou código candidato. Expectativas vêm do contrato e não devem ser adaptadas à implementação.
