# Prompt: Dispatcher de Manifestos Multi-Linguagem (parser/manifest)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/manifest.go`, `01_core/parser/manifest_test.go`  

---

## Contexto

Para permitir que o Bateia suporte múltiplos ecossistemas sem acoplar a infraestrutura a parsers específicos, é necessário um Dispatcher Genérico que identifique o tipo de manifesto pelo nome do arquivo (ex: `go.mod`, `package.json`) e direcione a execução para o parser correspondente.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP.
- Depende apenas de `01_core/domain`, `path/filepath`, `strings` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func DetectManifestType(filename string) string`
   - `"go.mod"` → `"gomod"`
   - `"package.json"` → `"packagejson"`
   - Outros / Desconhecido → `""`

2. `func ParseManifest(filename string, ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Identifica o tipo do manifesto via `DetectManifestType(filename)`.
   - Roteia para `ParseGoMod` ou `ParsePackageJson`.
   - Retorna erro caso o arquivo não seja suportado.

---

## Critérios de Verificação

```
Dado o arquivo "package.json" e um conteúdo JSON válido
Quando ParseManifest for executada
Então o Dispatcher chama ParsePackageJson e retorna um Manifest com Language "typescript".

Dado o arquivo "go.mod" e um conteúdo go.mod válido
Quando ParseManifest for executada
Então o Dispatcher chama ParseGoMod e retorna um Manifest com GoVersion preenchido.
```

---

## Resultado Esperado

- `01_core/parser/manifest.go`: Implementação do Dispatcher de manifestos.
- `01_core/parser/manifest_test.go`: Suite de testes validando roteamento e detecção.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Dispatcher Pattern | manifest.go, manifest_test.go |
