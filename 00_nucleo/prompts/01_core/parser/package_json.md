# Prompt: Parser de Manifestos package.json (parser/package_json)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/package_json.go`, `01_core/parser/package_json_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para o ecossistema TypeScript/Node.js, o conteúdo de texto do arquivo `package.json` lido de um repositório candidato precisa ser analisado para extrair o nome do pacote, versão e dependências (de produção, desenvolvimento e runtime peers).

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `encoding/json` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar a função:
`func ParsePackageJson(ref domain.CommitRef, content string) (*domain.Manifest, error)`

Regras do parsing:
1. Parsear o JSON usando `encoding/json` em uma struct interna `PackageJSON`.
2. Mapear o campo `name` para `ModuleName` e `version` para `Version`.
3. Preencher `Language: "typescript"`.
4. Mapear `dependencies` como dependências do tipo `"direct"`.
5. Mapear `devDependencies` como dependências do tipo `"dev"`.
6. Mapear `peerDependencies` como dependências do tipo `"peer"`.
7. Tratar erros de unmarshal com mensagem clara.
8. Manter `ConsultadoEm` zerado; o instante de coleta é fornecido explicitamente no dispatcher (ADR-0011).

---

## Critérios de Verificação

```
Dado um JSON com "name": "express-app", "dependencies": {"express": "^4.18.2"}, "devDependencies": {"jest": "^29.0.0"}
Quando ParsePackageJson for executada
Então retorna um Manifest com ModuleName "express-app", Language "typescript", 1 dependência "direct" (express) e 1 dependência "dev" (jest).
```

---

## Resultado Esperado

- `01_core/parser/package_json.go`: Implementação da função `ParsePackageJson`.
- `01_core/parser/package_json_test.go`: Suite de testes cobrindo parsing de pacotes válidos, devDependencies, peerDependencies e erros de JSON.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para TypeScript | package_json.go, package_json_test.go |

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.

## Ordem de dependências — 0128

Iterar os objetos JSON por chaves lexicograficamente ordenadas, usando sortedManifestKeys. Preservar grupos, tipos, versões, exclusões e política de erros. A ordem anterior dos mapas Go era indeterminada; a nova ordem não modifica artefatos históricos.
