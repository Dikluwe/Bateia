# Corpus 0097 — captura, identidade e ciclo de vida

Fixtures temporárias sintéticas, sem rede/execução de candidato. Hash com sink
deve ser byte-idêntico ao oráculo v2 e ao hash sem sink, inclusive DFS, binários,
.git arquivo, diretórios vazios, exclusões e symlinks estáticos ignorados.

Mutar/renomear/remover origem após captura e extrair exatamente bytes originais
pela cópia, com URL/CommitID públicos estáveis, paths privados ausentes e método
explicitamente congelado. Arquivos0400, sem execução. Close remove somente o
container criado. Falhas de FIFO, limite64MiB, cancelamento, temp sobreposto,
sink ocupado/fechado e output não-vazio não retornam hash parcial nem deixam
captura/dossiê. Não sobrescrever arquivos no sink.

Git hints congelados não mudam após alteração do config; config64KiB+1 é desconhecido.
Sem target não ler Git config. Nova execução com hints diferentes mantém digest
dos arquivos, mas muda local_metadata e RunID. Em resume, reaproveitar somente
evidências de arquivos válidas/mesma versão; recalcular ranking consultivo.

Superseder explícita e somente a expectativa do spy de coleta local0096: CLI
agora substitui fake local pela captura real. Controle positivo deve inspecionar
README efetivamente coletado/hash e checkpoint_reused0/1; versão de coleta
ausente força recoleta. Preservar testes/counters live e certificados históricos.
