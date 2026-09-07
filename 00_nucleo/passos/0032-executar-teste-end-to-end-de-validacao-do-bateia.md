# Passo 0032 — Executar teste end-to-end de validação do Bateia

**Estado**: `IMPLEMENTADO`  
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)  
**Corrige**: ausência de oráculo integrado ponta a ponta (E2E) que exercite simultaneamente a CLI local e o ciclo completo de ferramentas MCP (síncrona e assíncrona durável com defaults V1) sob determinismo  
**Baseline obrigatório**: Passo 0031 (autopreenchimento de defaults V1 e gates R1–R6 consolidados)  
**Camadas previstas**: L0 e L4; L1, L2 e L3 para consumo e verificação estrita  
**Natureza**: verificação e teste integral determinístico da cadeia completa

## Problema

O Bateia possui verificações pontuais ricas por camada (parsers puros em L1, servidores e CLI em L2, adaptadores de rede e cache em L3, e gates de integridade R1–R6 em L4). No entanto, não havia um oráculo formal de validação ponta a ponta que exercitasse, em uma única batida reprodutiva e sem dependência de rede externa flutuante:

1. A **CLI local offline** processando uma worktree real através do funil de 3 estágios (descoberta, triagem sintática AST Go e extração) e gerando o relatório Markdown íntegro;
2. O **Servidor MCP via stdio** processando a sequência canônica de protocolo:
   - `initialize` e negociação de capacidades;
   - `tools/list` com declaração estrita das ferramentas duráveis (`bateia_recommend`, `bateia_start_recommendation`, `bateia_get_run`, `bateia_get_artifact`, `bateia_cancel_run`);
   - `bateia_recommend` com payload mínimo de negócio (`idea` com `problem` e `capabilities`), comprovando a aplicação transparente dos defaults V1 consolidados no Passo 0031;
   - `bateia_start_recommendation` com chave de idempotência e payload mínimo, persistindo a execução durável em disco com `run_id` auditável;
   - `bateia_get_run` recuperando o estado estruturado da execução persistida;
   - `bateia_get_artifact` inspecionando o artefato de recomendação emitido com hash sha256;
   - `bateia_cancel_run` garantindo a preservação do estado terminal da execução.
3. O cumprimento das invariantes de isolamento: zero I/O em L1, pureza da composição em L4 e conformidade com o linteiro `crystalline-lint`.

## Matriz de Cobertura E2E

| Subsistema | Entrada | Comportamento Esperado | Oráculo |
| :--- | :--- | :--- | :--- |
| **CLI Local** | `-path . -query "funnel"` | Funil de 3 estágios executado, triagem AST Go concluída, relatório Markdown gerado | Arquivo de saída não-vazio, contendo seções de candidatos e evidências |
| **MCP Handshake** | RPC `initialize` | Resposta JSON-RPC 2.0 com `capabilities` e `serverInfo` | Sucesso sem erro de protocolo |
| **MCP Descoberta** | RPC `tools/list` | Listagem das ferramentas duráveis suportadas com esquemas fechados | Existência garantida das ferramentas essenciais |
| **MCP Síncrono V1** | RPC `bateia_recommend` com payload mínimo | Autopreenchimento de `policy`, `target` e `search` aplicados com sucesso | Objeto de recomendação retornado contendo defaults V1 aplicados |
| **MCP Assíncrono Durável** | RPC `bateia_start_recommendation` | Criação de job durável, geração de `run_id` e persistência em disco | `run_id` válido em estado inicial `QUEUED` |
| **MCP Consulta Durável** | RPC `bateia_get_run` por `run_id` | Recuperação do snapshot da execução persistida até resolução | Estado terminal `SUCCEEDED` |
| **MCP Artefatos Duráveis** | RPC `bateia_get_artifact` | Consulta segura paginada do artefato publicado | Retorno com `artifact_hash` prefixado por `sha256:` |

## Fora de Escopo

- Depender de credenciais reais ou chamadas ativas de rede para GitHub/Sourcegraph durante o teste E2E;
- Alterar as estruturas internas de domínio de L1;
- Afrouxar qualquer gate anterior (R1 a R6).

## Critério de Conclusão

O passo é considerado concluído quando:
1. O teste E2E `step0032_e2e_test.go` for implementado e passar 100% verde sem I/O não controlado em L1;
2. Todas as ferramentas MCP e a CLI local forem exercitadas com sucesso;
3. O linteiro `crystalline-lint .` passar sem erros;
4. O laudo em `00_nucleo/laudos/0032-executar-teste-end-to-end-de-validacao-do-bateia.md` for emitido e congelado.
