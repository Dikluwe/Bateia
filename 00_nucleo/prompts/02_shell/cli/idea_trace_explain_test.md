# CLI explain — passo 0099

Testar execução offline stdout/arquivo, seleção de nó e vínculo ao graph_hash
original. Localização de razões deve apontar ao JSON Pointer no documento report.
Recusar build input diferente, hash ausente/adulterado, seleção desconhecida ou
duplicada, schema desconhecido, chaves desconhecidas/duplicadas/casefold e trailing
data. Toda falha deve preservar resultado anterior, não criar saída nem emitir
stdout. Recusar input FIFO e aliases de output por path/inode/link simbólico.
