# Prompt: Avaliador de Estrutura de Diretórios (parser/structure)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `01_core/parser/structure.go`, `01_core/parser/structure_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia, a árvore de diretórios raiz de um repositório precisa ser avaliada em relação aos filtros estruturais passados pelo usuário (`must_have`, como `cmd`, `internal`, `pkg`).

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função puramente determinística. Zero I/O.
- Recebe um slice de `domain.DirectoryNode` e a lista `mustHave []string`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

1. Definir a struct `StructuralScore` com campos `HasCmd`, `HasInternal`, `HasPkg`, `FoundPaths`, `Approved` (bool) e `Reason` (string).
2. Criar a função:
`func EvaluateStructure(nodes []domain.DirectoryNode, mustHave []string) StructuralScore`
3. Mapear os caminhos e subpastas de primeiro nível dos nós de diretório.
4. Se todos os itens de `mustHave` estiverem presentes na árvore, aprovar (`Approved: true`).
5. Se algum item de `mustHave` estiver ausente, reprovar (`Approved: false`) registrando o motivo em `Reason`.

---

## Critérios de Verificação

```
Dado uma lista de nós com caminhos ["cmd/app/main.go", "internal/domain/user.go"] e mustHave ["cmd", "internal"]
Quando EvaluateStructure for chamada
Então retorna Approved == true

Dado uma lista de nós sem a pasta "internal" e mustHave ["internal"]
Quando EvaluateStructure for chamada
Então retorna Approved == false e Reason informando a ausência da pasta "internal"
```

---

## Resultado Esperado

- `01_core/parser/structure.go`: Struct `StructuralScore` e função `EvaluateStructure`.
- `01_core/parser/structure_test.go`: Testes cobrindo aprovação e reprovação estrutural.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | structure.go, structure_test.go |
