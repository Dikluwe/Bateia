# Prompt: Interface CLI e Controle de Flags (cli/app)

**Camada**: L₂ — Casca  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `02_shell/cli/app.go`  

---

## Contexto

O Bateia é uma ferramenta de linha de comando. Este componente trata a recepção dos argumentos da CLI, validação de parâmetros obrigatórios (`-query`), parsing de opções (`-lang`, `-max`, `-output`, `-must-have`, `-cache-dir`) e delegação para o orquestrador do funil.

---

## Restrições Estruturais

- **L₂ (Casca)**: Lida exclusivamente com interação de usuário/I/O de linha de comando.
- Pode importar `01_core/domain` e `01_core/usecase`.
- **Proibido Importar $L_3$**: A CLI não deve instanciar diretamente clientes de banco, rede ou filesystem de infraestrutura. Ela recebe a função runner ou o caso de uso pronto via $L_4$ (Fiação).

---

## Instrução

1. Definir o struct `CLIApp` que armazena referências para o executor do funil.
2. Criar a função `Run(args []string) error` usando o pacote `flag`:
   - Flag `-query` / `-q` (obrigatória): termo de busca arquitetural.
   - Flag `-lang` (default: `"go"`).
   - Flag `-max` (default: `20`).
   - Flag `-output` (default: `"relatorio_bateia.md"`).
   - Flag `-must-have`: lista de pastas separadas por vírgula.
   - Flag `-cache-dir`: diretório customizado do cache.
3. Se `-query` estiver vazia, exibir o menu de ajuda e encerrar amigavelmente.
4. Invocar o executor e passar o relatório resultante para a geração em disco.

---

## Critérios de Verificação

```
Dado a invocação da CLI sem o argumento -query
Quando Run for executada
Então retorna erro informando a ausência do argumento obrigatório.

Dado argumentos validos -query "hexagonal" -lang "go"
Quando Run for executada
Então monta a struct QueryFilter e invoca o executor.
```

---

## Resultado Esperado

- `02_shell/cli/app.go`: Manipulador da CLI do Bateia.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | app.go |

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.

## Passo 0062
Adicionar somente ajuda para effects verify e seu recorte offline.

## Passo 0070
Encaminhar `algorithm history analyze` para a casca offline definida em
`00_nucleo/prompts/02_shell/cli/algorithm_history.md`, sem lógica histórica no dispatcher.
