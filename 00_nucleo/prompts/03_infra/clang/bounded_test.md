# Limite real da saída Clang

Passo 0125. Testar io.Copy com leitor sem WriterTo, exercitando ReadFrom do
destino usado por os/exec. Entrada acima do limite deve falhar e sinalizar
exceeded sem armazenar mais bytes que o limite; saída no limite permanece íntegra.
