# Adaptador de configuração durável local

Implemente em L3 um adaptador que resolva `BATEIA_RUN_ROOT`, prepare e audite uma
raiz local privada e carregue ou crie uma chave de cursor persistente de 256 bits.

Rejeite symlinks, hardlinks, arquivos especiais e segredos com tamanho inválido
antes de mutar inodes externos. Crie diretórios e chave com modos restritos,
preserve a chave entre cargas e exponha somente os valores necessários à composição:
raiz absoluta, segredo, limites do store, workers e capacidade da fila.

Preserve os defaults e textos de erro congelados no passo 0077. Não implemente
regras de domínio, lifecycle da fila, parsing de CLI ou despacho de comandos.
