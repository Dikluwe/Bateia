# Prompt: Parser de Manifestos go.mod (parser/gomod)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `01_core/parser/gomod.go`, `01_core/parser/gomod_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia, o conteúdo de texto do arquivo `go.mod` lido de um repositório candidato precisa ser analisado para extrair o nome do módulo, versão do Go e dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain` e `strings`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar a função:
`func ParseGoMod(ref domain.CommitRef, content string) (*domain.Manifest, error)`

Regras do parsing:
1. Identificar o nome do módulo a partir de linhas com o prefixo `module `.
2. Identificar a versão do Go a partir de linhas com o prefixo `go `.
3. Extrair dependências declaradas tanto inline (`require github.com/foo/bar v1.0.0`) quanto dentro de blocos `require (...)`.
4. Marcar dependências que possuam o comentário `// indirect` com `Indirect: true`.
5. Ignorar linhas em branco e comentários de cabeçalho.
6. Preencher `ConsultadoEm` com a data/hora UTC atual.

---

## Critérios de Verificação

```
Dado o texto de um go.mod com módulo "github.com/example/bateia", go "1.22" e dependência "github.com/pkg/errors v0.9.1 // indirect"
Quando ParseGoMod for executada
Então retorna um Manifest com ModuleName "github.com/example/bateia", GoVersion "1.22" e 1 dependência indireta marcada.
```

---

## Resultado Esperado

- `01_core/parser/gomod.go`: Implementação da função `ParseGoMod`.
- `01_core/parser/gomod_test.go`: Suite de testes cobrindo parsing de blocos, dependências diretas e indiretas.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | gomod.go, gomod_test.go |
