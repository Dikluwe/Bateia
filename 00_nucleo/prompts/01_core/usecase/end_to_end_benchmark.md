# End-to-end benchmark engine

Valide correspondência antes de medir o motor. Trocas comprovadas são VIOLATED;
identidade, commit, regra ou proveniência opaca são UNKNOWN; UNKNOWN nunca conta
como sucesso. Verifique hashes selados, rejeite mutações, reexecute somente casos
PRESERVED offline e compute métricas nullable sem ajustar oráculos.
