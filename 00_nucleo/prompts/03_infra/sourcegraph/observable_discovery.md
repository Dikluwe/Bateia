# Prompt Spec L0: Descoberta Observável no Sourcegraph

## Objetivo

Expor as queries completas e fallbacks realmente tentados pelo Sourcegraph, sem
vazar credenciais e preservando o contrato de busca existente.

## Restrições

- L3, adaptador da porta observável de L1.
- Registrar query, quantidade e erro resumido por tentativa.
- Distinguir resposta vazia de erro.
- Parar no primeiro resultado não vazio.
- Nunca serializar token, cabeçalhos ou corpo sensível.
