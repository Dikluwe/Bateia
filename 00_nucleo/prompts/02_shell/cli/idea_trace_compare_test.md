# Testes CLI de comparação consultiva — passo 0100

Criar dois snapshots congelados a partir de relatórios com decisões declaradas
diferentes. `idea trace compare --input request.json [--output result.json]`
deve vincular ambos os hashes, distinguir a decisão histórica da candidata,
manter decision_impact NONE e preservar a decisão histórica dentro do plano.
O snapshot candidato não é automaticamente adotado nem confirmado.

Verificar stdout versus arquivo byte a byte, substituição explícita de resultado
regular distinto, entrada inalterada, ausência de temporários e zero revisão para
snapshots idênticos. Recusar schema futuro, hash ausente/adulterado em qualquer
snapshot, declaração adulterada, campos desconhecidos/duplicados/casefold e JSON
trailing sem criar/substituir arquivos ou emitir stdout parcial.

Reutilizar as garantias trace para recusar saída igual à entrada, hardlink,
symlink, FIFO e input FIFO; recusar argumentos ambíguos/desconhecidos. Nenhum teste
executa candidato, consulta rede ou altera fluxos recommend/default.
