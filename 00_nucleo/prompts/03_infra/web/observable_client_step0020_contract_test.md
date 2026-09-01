# Prompt — Contrato adversarial do cliente observável no Passo 0020

Materialize ataques L3 independentes derivados do Passo 0020 para verificar o contrato
observável do cliente web autenticado, incluindo sigilo da credencial, rejeição de conteúdo
privado, preservação dos estados externos obrigatórios, selagem das respostas e tratamento
inteiro dos contadores de rate limit.

Use transporte HTTP falso exclusivamente em memória, sem servidor, socket ou dependência de
rede. Preserve a independência adversarial: este prompt autoriza somente o teste contratual
espelhado em `03_infra/web/observable_client_step0020_contract_test.go` e não autoriza criar,
alterar ou corrigir implementação.
