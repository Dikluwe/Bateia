# Orientação de trabalho no Bateia

## Propósito e arquitetura

O Bateia investiga software para produzir oportunidades de reúso e insights de engenharia. Tekt governa todo o programa. Leia as [regras estruturais](.agentrules), os [ADRs de hierarquia](00_nucleo/ADR/ADR-0008-hierarquia-de-arquiteturas-sob-tekt.md) e o [mapa de responsabilidades](00_nucleo/arquitetura/README.md) antes de alterar uma capacidade. As regras históricas em `.agents/AGENTS.md` continuam aplicáveis; esta entrada acrescenta orientação de navegação e o propósito atualizado por instrução do usuário.

## Como localizar uma mudança

1. Identifique a capacidade, seu contrato, responsável e camada pelo mapa. Investigação, runtime e rastreabilidade são escopos de organização; não substituem L0–L4.
2. Leia o prompt proprietário e os testes consumidores. Declare comportamento esperado e limites antes da materialização relevante.
3. Aplique a menor fatia necessária. Regras e estado pertencem a L1, efeitos concretos a L3, entrada/apresentação a L2 e composição/lifecycle a L4.
4. Preserve formatos públicos e hashes históricos. Não mova arquivos apenas para fazer a árvore se parecer com o diagrama.
5. Mantenha a linhagem causal e execute lint, testes, vet e build conforme `.agentrules`.

## Distinções obrigatórias

- Estado de execução não é força de evidência: sucesso operacional não promove uma hipótese.
- Elegibilidade para reúso não determina todo o valor de um insight; as políticas atuais só mudam com contrato próprio.
- Presença em catálogo, observação de código, hipótese e efeito verificado possuem escopos diferentes.
- Arquitetura governante, organização de subsistema, mecanismo e vista são classificações diferentes. Classifique pelo contrato e invariantes, não pelo nome de um tipo.

## Verificação

Executar `crystalline-lint --fail-on warning .`, `go test ./...`, `go vet ./...` e build de `./04_wiring`. Testes de fixture podem precisar de loopback e assinatura Git desativada somente no processo; consultar o guia de engenharia. Gates selados possuem módulo próprio e não são promovidos por testes do módulo raiz.
