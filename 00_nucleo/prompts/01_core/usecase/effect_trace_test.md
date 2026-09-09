# Testes independentes do fingerprint 0064

Papel B, regime A/B em filesystem compartilhado, sem atestação física. Entradas: contrato0064 congelado, contratos0062/0063 e API domain baseline. Não consultar produção candidata, testes legados, fixtures ou relatórios.

Criar testes externos do usecase.EffectTraceFingerprint. Construir oráculo independente por serialização explícita das dez linhas e SHA-256, além de vetor conhecido. Cobrir normalização ASCII, sensibilidade individual de todos campos, seed máxima, limites inclusivos, identificação Go exportada Unicode, valores inválidos e entrada imutável. L1 sem I/O. Não executar antes do congelamento dos artefatos.
