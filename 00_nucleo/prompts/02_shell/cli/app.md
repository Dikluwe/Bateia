## Passo 0094 — ajuda do scanner

Incluir security scan --binary CAMINHO_ABSOLUTO [--output arquivo.json] na ajuda.
Informar uso de govulncheck instalado e que o binário candidato não é executado.

# Prompt: Interface CLI e Controle de Flags (cli/app)

0135: anunciar idea insights compare na ajuda, como comparação de observações de fontes e baseline do alvo.

0136: anunciar idea insights explore como preparação de exemplos variados e tarefas para IA a partir de pergunta/bundles, sem curadoria de avaliações.

0137: anunciar idea insights discover como descoberta online diversificada com coleta de README fixado e contexto para IA; output obrigatório.

0138: indicar --sources all para descoberta entre os conectores configurados, com metadados de catálogo separados dos READMEs.

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

## Passo 0090

Inclua ajuda e exemplo para `idea landscape analyze`, mantendo o comando opcional e offline.

## Passo 0093

Adicionar ajuda e exemplo para `idea evidence recheck`, releitura opcional
no provedor GitHub com entrada explícita e bundles exportáveis.


## Passo 0105

Acrescentar na ajuda o comando idea insights analyze, sem alterar comandos ou defaults existentes.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

## Evolução0142

Aplicar passo0142: memória validada antes de efeitos, followups de um salto com limites e origem própria, --memory/--follow-links na casca e feedback offline. FollowReader apenas GitHub, demais forges reaproveitam RepositoryReader. Memória em snapshot único de saída, sem releitura automática de examinados/descartados, sem promover histórico a evidência atual. Modo sem --sources preservado.


## Passo0143

Exibir expand na ajuda como continuação por relatório e foco opcional, sem nova busca. Preservar comandos anteriores.
