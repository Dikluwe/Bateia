# Idea run adversarial tests

Congele o pipeline observável antes da implementação: replay deve produzir a
cadeia completa sem dependências externas; bundle ausente deve falhar; decisão
INCONCLUSIVE nunca pode promover o run a COMPLETE.

## 0158 — Double de conteúdo com identidade

No recorte de integridade0158, os doubles que retornam extração bem-sucedida devem preencher Ref com owner/repo/commit do candidato efetivamente solicitado. Conservar conteúdo, falhas, quotas, testemunhos, mutantes e assertivas funcionais anteriores. Ref ausente ou divergente passa a ser controle negativo explícito nos testes próprios de integridade, sem dispensar a validação produtiva.
