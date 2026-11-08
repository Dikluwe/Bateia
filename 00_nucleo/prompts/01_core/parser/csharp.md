# Prompt: Parser de Manifestos C# / .NET .csproj (parser/csharp)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/csharp.go`, `01_core/parser/csharp_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para o ecossistema C# / .NET, o conteúdo do arquivo de projeto `.csproj` precisa ser analisado para extrair o nome do projeto e os pacotes NuGet declarados em tags `<PackageReference Include="..." Version="..." />`.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar a função:
`func ParseCsProj(ref domain.CommitRef, content string) (*domain.Manifest, error)`

Regras do parsing:
1. Extrair todas as ocorrências de `<PackageReference Include="PackageName" Version="1.2.3" />` ou elementos `<PackageReference>` expandidos.
2. Identificar pacotes de teste (ex: `xunit`, `nunit`, `MSTest.TestFramework`, `Microsoft.NET.Test.Sdk`) e marcar como tipo `"dev"`, demais como `"direct"`.
3. Preencher `Language: "csharp"`.
4. Manter `ConsultadoEm` zerado; o instante de coleta é fornecido explicitamente no dispatcher (ADR-0011).

---

## Critérios de Verificação

```
Dado um .csproj com <PackageReference Include="Microsoft.EntityFrameworkCore" Version="8.0.0" /> e <PackageReference Include="xunit" Version="2.6.0" />
Quando ParseCsProj for executada
Então retorna um Manifest com Language "csharp", 1 dependência "direct" (EntityFrameworkCore) e 1 dependência "dev" (xunit).
```

---

## Resultado Esperado

- `01_core/parser/csharp.go`: Implementação do parser de .csproj.
- `01_core/parser/csharp_test.go`: Suite de testes cobrindo arquivos .csproj.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para C# / .NET | csharp.go, csharp_test.go |

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.
