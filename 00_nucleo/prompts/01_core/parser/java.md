# Prompt: Parser de Manifestos Java/Kotlin (parser/java)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/java.go`, `01_core/parser/java_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para o ecossistema Java/Kotlin, o conteúdo de texto de manifestos como `pom.xml` (Maven) ou `build.gradle` / `build.gradle.kts` (Gradle) precisa ser analisado para extrair o artefato/projeto e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParsePomXml(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `<artifactId>` principal ou `<name>`.
   - Extrair dependências sob `<dependency>` (`<groupId>`, `<artifactId>`, `<version>`, `<scope>`).
   - Marcar dependências com `<scope>test</scope>` como tipo `"dev"`, demais como `"direct"`.
   - Marcar `Language: "java"`.

2. `func ParseBuildGradle(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas sob blocos `dependencies { ... }` (ex: `implementation 'org.springframework.boot:spring-boot-starter-web'`, `testImplementation 'org.junit.jupiter:junit-jupiter'`).
   - Marcar `testImplementation` / `testCompile` como tipo `"dev"`, `implementation` / `api` como `"direct"`.
   - Marcar `Language: "java"`.

---

## Critérios de Verificação

```
Dado um pom.xml com artifactId "spring-boot-demo" e dependências "spring-boot-starter-web" e "junit-jupiter" (test)
Quando ParsePomXml for executada
Então retorna um Manifest com ModuleName "spring-boot-demo", Language "java", 1 dependência "direct" e 1 dependência "dev".
```

---

## Resultado Esperado

- `01_core/parser/java.go`: Implementação dos parsers Maven e Gradle.
- `01_core/parser/java_test.go`: Suite de testes cobrindo `pom.xml` e `build.gradle`.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Java/Kotlin | java.go, java_test.go |

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.
