# Prompt: Ponto de Entrada e Fiação (wiring/main)

**Camada**: L₄ — Fiação  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `04_wiring/main.go`  

---

## Contexto

Conforme ADR-0002, a camada $L_4$ (Fiação) é o ponto de montagem único da aplicação. Ela substitui a subpasta redundante `cmd/bateia/` por uma fiação direta em `04_wiring/main.go`.

---

## Restrições Estruturais

- **L₄ (Fiação)**: Absorve a complexidade de injeção de dependências.
- **Pode Importar Todas as Camadas**: Importa $L_1$ (`domain`, `usecase`), $L_2$ (`cli`, `report`) e $L_3$ (`cache`, `sourcegraph`, `github`).
- **Zero Lógica de Negócio**: Não pode conter estruturas condicionais de negócio ou algoritmos de decisão. Apenas instanciação e repasse de controle.

---

## Instrução

1. Em `04_wiring/main.go`, definir a função `main()`.
2. Ler tokens de ambiente (`GITHUB_TOKEN`, `SOURCEGRAPH_TOKEN`).
3. Instanciar os adaptadores concretos de $L_3$:
   - `DiskCacheStore`
   - Sourcegraph `Client`
   - GitHub `Client`
4. Instanciar o caso de uso `SearchFunnel` de $L_1$, injetando os adaptadores de $L_3$ como implementações das interfaces de $L_1$.
5. Instanciar e executar a `CLIApp` de $L_2$, passando o orquestrador e o gerador de relatórios Markdown.

---

## Critérios de Verificação

```
Dado o comando "go run ./04_wiring/main.go --help"
Quando executado no terminal
Então exibe as opções de ajuda da CLI do Bateia sem erros de compilação ou pânico.
```

---

## Resultado Esperado

- `04_wiring/main.go`: Ponto de entrada executável do Bateia.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 (ADR-0002) | main.go |
