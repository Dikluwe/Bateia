# Prompt espelhado — gate causal R4 do Passo 0026

Construa um gate pré-candidato isolado que leia integralmente os 70 casos adversariais.
É proibido despachar por `M01`–`M60` ou concluir pelo binding. Selecione programas pela
categoria semântica declarada e passe a eles o caso completo. Cada mutação deve produzir
duas observações executadas: referência e mutante. `Violated` só existe quando o campo
causal indicado difere e termos materiais de `mutation`, `stimulus` e `witness` foram
consumidos; observações iguais ou textos insuficientes são `SURVIVED`.

Para mutações de persistência use diretórios temporários e operações reais: arquivos
exclusivos, escrita parcial/completa, `Sync`, `Rename`, sync de diretório, confirmação,
reopen lógico, corrupção, symlink, hardlink e modos. Não use um modelo de disco em
memória. Para cursor use HMAC real; para concorrência use barreiras/atomics; para AST use
o parser Go; para legado verifique hashes do certificado e fonte 0025 e mantenha a
invocação observável predecessor-only como gate complementar. Casos opacos conservam
`Unknown` apenas com construction/witness explícitos.

Execute os casos em ordem lexical, reversa e intercalada, três repetições por ordem, e
exija o mesmo digest normalizado nas nove rodadas. O pacote deve compilar sem importar
qualquer implementação candidata.
