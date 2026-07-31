# ADR-0003: Atomização de Parsers — Um Arquivo Dedicado por Linguagem

**Status**: ACEITO  
**Data**: 2026-07-31  
**Documentos Relacionados**: ADR-0002 (estrutura Tekt), ADR-0001 (revisada)  
**Contexto Externo de Inspiração**: ADR-0109 do projeto typst-crystalline (atomização de módulos de layout)

---

## Contexto

Durante a expansão do suporte a linguagens no Bateia, os parsers foram inicialmente agrupados em arquivos temáticos por ecossistema:

```text
01_core/parser/
├── scripting.go          # Julia, R, Lua, OCaml, Elixir, Ruby
├── web_ecosystems.go     # JavaScript, TypeScript, PHP
├── modern_ecosystems.go  # Rust, Swift, Haskell, Dart, Scala
├── classic_ecosystems.go # Java, C#, C++
├── dynamic_ecosystems.go # Clojure, Zig, Pascal, Fortran, BASIC
├── emergent_ecosystems.go# F#, Perl, COBOL, Nim, Crystal, V, Erlang
├── extended_ecosystems.go# D, Gleam, ColdFusion, ActionScript
└── specialized_ecosystems.go # Ada, Mojo, Objective-C, Haxe, Raku, Racket, Chapel, Pony, Odin
```

Este agrupamento criou três problemas:

1. **Difusa rastreabilidade L0→L1**: O Prompt L0 `scripting.md` era a origem de seis parsers distintos. Ao alterar o parser de Ruby, era necessário ler a especificação de Julia, Lua e outros para garantir que não havia colisão de nomes ou lógica partilhada não intencional.

2. **Impacto de mudança desproporcional**: Um erro de compilação em `scripting.go` bloqueava a construção de todos os parsers do grupo. Um `git blame` numa linha de Julia apontava para um commit que também modificava OCaml.

3. **Violação do princípio de coesão L1**: Em L1, cada arquivo deve representar **uma unidade de domínio coesa**. Uma linguagem de programação é uma unidade de domínio natural. Um "ecossistema de scripting" é uma categoria de conveniência humana, não uma fronteira de domínio.

---

## Decisão

**Cada linguagem suportada pelo Bateia tem um arquivo dedicado** em `01_core/parser/` e um Prompt L0 dedicado em `00_nucleo/prompts/01_core/parser/`.

A forma canônica é:

```text
01_core/parser/<lang>.go        ← implementação L1
01_core/parser/<lang>_test.go   ← testes da unidade
00_nucleo/prompts/01_core/parser/<lang>.md  ← Prompt L0 perene
```

### O que esta decisão **é**

- Mover a lógica de cada parser para o arquivo da **unidade dona** (`<lang>.go`).
- Cada par `(<lang>.go, <lang>.md)` é legível e rastreável **sozinho**, sem precisar abrir outros arquivos do mesmo pacote.
- O dispatcher central (`manifest.go`) permanece como **ponto único de entrada**, com um `switch` que roteia para o parser correto — estático, exaustivo e sem despacho dinâmico.

### O que esta decisão **não é**

- **Não** é desacoplar `manifest.go` do restante do pacote. O `switch` e os imports ficam.
- **Não** é usar `interface{}`, vtable ou qualquer forma de despacho dinâmico para "zerar" o `switch`.
- **Não** é agrupar linguagens relacionadas num mesmo arquivo por conveniência (ex.: "linguagens funcionais"). A fronteira é a linguagem, não a família.
- **Não** é mover lógica de triagem ou scoring para os arquivos de parser — essa lógica permanece em `insights.go` e `relevance.go`.

---

## Forma Canônica do Switch Dispatcher

O `switch` em `manifest.go` é o único ponto onde a atomização se conecta ao sistema. A sua forma é **magra**: cada `case` delega imediatamente para a free function no arquivo dedicado.

```go
// manifest.go — correto
case "ruby":
    return ParseGemfile(ref, content)

// manifest.go — ERRADO (lógica no dispatcher)
case "ruby":
    if strings.Contains(content, "gem ") {
        // ... lógica aqui — isto viola a ADR-0003
    }
```

---

## Regra para Novas Linguagens

Ao adicionar suporte a uma nova linguagem `X`:

1. Criar `00_nucleo/prompts/01_core/parser/x.md` com o Prompt L0.
2. Criar `01_core/parser/x.go` com as funções `ParseX*(...)`.
3. Criar `01_core/parser/x_test.go` com a suíte de testes.
4. Adicionar a detecção em `DetectManifestType` e o roteamento em `ParseManifest`.
5. Se o manifesto for buscado remotamente, adicionar à lista `manifestFiles` em `03_infra/github/client.go`.

**Proibido**: criar um arquivo agrupador (ex.: `novos_ecosystems.go`) mesmo que temporariamente.

---

## Consequências

### Positivas

- **Rastreabilidade unívoca**: `00_nucleo/prompts/01_core/parser/ruby.md` → `01_core/parser/ruby.go`. Biunívoco.
- **Impacto cirúrgico**: modificar o parser de Ruby não afeta os testes de Julia.
- **Legibilidade isolada**: um revisor lê `ada.go` e entende o parser de Ada sem contexto adicional.
- **Consistência com ADR-0002**: espelha a estrutura de pastas canônica Tekt.

### Negativas / Compensações

- Número de arquivos em `01_core/parser/` cresce linearmente com as linguagens suportadas (atualmente 45+). Aceite: a navegabilidade por arquivo supera a densidade por arquivo.
- O `switch` em `manifest.go` cresce a cada nova linguagem. Aceite: é a única lista canônica do sistema e deve ser exaustiva.

---

## Métricas de Validação (estado em 2026-07-31)

| Métrica | Antes | Depois |
|---------|-------|--------|
| Arquivos de parser | 8 agrupadores | 45 arquivos dedicados |
| Prompts L0 de parser | 8 genéricos | 45 perenes |
| `go test ./...` | PASS | PASS |
| Violações de camada | 0 | 0 |

**Commit de referência**: `12ca73a`  
**Estado do working tree no momento da medição**: clean (`git status` = nothing to commit).

---

## Histórico de Revisões

| Data | Motivo |
|------|--------|
| 2026-07-31 | Criação — formaliza a atomização completa concluída em P16 |
