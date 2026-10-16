# Contrastes 0116

Nome alterado mantém padrão; guard ausente/condicional/adiado/posterior,
panic sombreado, nil sombreado, método com defer recover, base diferente ou
reassign/address/capture/jump retiram projeção. Comparação invertida aceita;
!= vs == não confundir. Writes de sentinela por outra instância não autorizam
inicialização daquela base. Erro de tipos suprime. Localização/ordem/limites;
fonte real guard e sentinel write sem certificar atomicidade.
