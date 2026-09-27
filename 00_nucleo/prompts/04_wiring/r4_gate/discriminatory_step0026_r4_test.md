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

## Supersessão explícita de segurança do MCP — 2026-09-26

O guard comum corrige bypass do orçamento e framing ilimitado nas extensões MCP,
exigindo alterar server.go. Nos dois witnesses de compatibilidade deste gate,
substitua somente o hash esperado dessa fonte:
e0c8605039740c895e733c8edc7cb7b6e6137a23b91ce93a1dbfa34fe18a0d42
→ ec8fc8c4461d73984c29a73b8b5d65f5ffcb6c3a35b7ad907529b5fc613c4ff9.

O certificado histórico 0025 permanece inalterado, assim como os 70 ataques,
categorias, expected outcomes, programa de mutação de retrocompatibilidade,
comparação causal e nove rodadas. A revisão não é recertificação histórica da
nova implementação. Contratos/goldens públicos são mantidos pelos testes MCP;
o hash identifica a revisão sucessora de segurança para permitir que este gate
continue exercitando o mesmo ataque, em vez de falhar antes de executá-lo.
