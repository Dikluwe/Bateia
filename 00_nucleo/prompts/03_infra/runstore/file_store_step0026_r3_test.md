# Prompt espelhado — store R3 do Passo 0026

Adicione testemunhos públicos L3 sem congelar filenames ou helpers privados. Em um store
real com diretório temporário, leia o registro imediatamente após cada confirmação:
create, publicação do request, oito progressos, cada publicação restante, terminal e
reopen. O terminal `SUCCEEDED` só pode existir com os seis artefatos publicados.

Para cursores, obtenha um cursor real, decodifique o payload público e reproduza a
assinatura HMAC a partir do segredo público injetado no construtor. Confirme primeiro que
o cursor re-assinado sem alteração é válido; então re-assine payloads com versão, hash e
offset realmente alterados e exija `CURSOR_INVALID`.

Encha um store real com `CreateOrResolve` até `QUOTA_EXCEEDED`, feche, reabra e tente novo
create. Exija a mesma rejeição e confirme que todos os runs admitidos continuam presentes;
nenhuma eviction é permitida. Use somente fakes offline e dirs temporários.
