Bundles repetidos da mesma identidade/revisão são erro, não coleções independentes
capazes de dividir um inventário de licença. Schema desconhecido nunca gera
interpretação verificada, mesmo que os bytes/digest sejam autoconsistentes.

# Orçamento da API técnica pura

Limitar antes de hash/traversal: cem bundles, cem itens por bundle, mil itens
totais, arquivo de 2 MiB, conteúdo total de 16 MiB e campos lidos de 8192 bytes.
Excesso é erro, não truncamento nem evidência parcial promovível. Identidade
inválida é UNVERIFIED. Chave de predicado estruturada, sem colisões por
delimitadores, deve usar a mesma equivalência casefold da licença do verificador;
SBOM preserva sensibilidade a capitalização.
