# Exemplo sequencial: Bateia encontra a fonte, XERJ recupera o mecanismo

**O exemplo foi executado.** Bateia descobriu e leu `arclabs561/postings`; XERJ indexou a mesma fonte no mesmo commit e recuperou a organização do leitor, sua definição e o helper que lê uma faixa de postings com validação. O projeto-fonte não foi compilado nem executado. Não se mediu melhoria do XERJ.

Fonte: [arclabs561/postings](https://github.com/arclabs561/postings/tree/edb02683c6359efe2034a6cc71b64c9aaed90cb0), pin `edb02683c6359efe2034a6cc71b64c9aaed90cb0`.

## 1. Descoberta pelo Bateia

A primeira consulta foi `compressed postings blocks random access inverted index`. Ela encontrou candidatos pertinentes, incluindo um exemplo de índice invertido em Go, mas teve leituras recusadas/incompletas e uma fonte de acesso aleatório a gzip. A consulta foi refinada por mecanismos, sem informar URL ou nome de projeto:

```sh
/tmp/bateia-0155 idea insights discover \
  --question 'inverted index postings skip block maximum score top k search' \
  --algorithm-ideas --max-per-query 5 --max-references 4 \
  --timeout 120s --attempt-timeout 20s \
  --output /tmp/bateia-xerj-demo-0156/bateia/discovery-refined.json
```

O retorno selecionou `arclabs561/postings` e coletou seu README. [README:109](https://github.com/arclabs561/postings/blob/edb02683c6359efe2034a6cc71b64c9aaed90cb0/README.md#L109) descreve um leitor que mantém diretórios fixos em memória e lê payloads por faixa conforme os termos pedidos. Essa foi a razão da escolha manual para o exemplo. O recibo de descoberta está em [audit/002](bateia/audit/002.json), com os bytes em [snapshot](bateia/audit/002.output.snapshot).

O operador já conhecia o tema e a referência de um ensaio anterior. Esta rodada não é descoberta cega: demonstra que o projeto voltou a aparecer em uma consulta por mecanismos. Não atribuímos à ferramenta a escolha semântica feita pelo operador. O relatório completo foi PARTIAL: outras leituras tiveram INTEGRITY_ERROR ou CONTENT_LIMIT; isso não foi tratado como conteúdo efetivamente lido.

## 2. Handoff com identidade preservada

O [audit de aquisição](bateia/audit/003.json) registra URL e testemunho nativo anterior. O clone estava limpo e seu HEAD coincidia com o pin do README. Foi copiado para a entrada do XERJ sem modificar código; ambos os lados conservam o mesmo commit. Ver [handoff](handoff.json) e [resultado estruturado](result.json).

## 3. Indexação e consultas no XERJ

Foi usado XERJ rc.101 em nó temporário exclusivo, sem feedback, modelos neurais ou grafo. A indexação concluiu com 30 arquivos, 20 arquivos de código e 913 registros. Um fixture binário ficou fora do índice; exit 3 indicou `ok=true`, searchable, completed-with-junk. Isso não demonstra cobertura do fixture nem execução de testes. [Audit da indexação](xerj/audit/005.json).

```sh
xerj autoindex /tmp/bateia-xerj-demo-0156/xerj/acquired/postings \
  --url http://127.0.0.1:19701 --prefix demo-postings-0156 \
  --no-semantic --no-graph --workers 2 --max-index-size 512mb --yes
```

A primeira pergunta, sem arquivo ou símbolo fornecido, foi:

```text
How do file-backed postings readers keep directories resident and read posting payloads on demand?
```

O primeiro hit trouxe `RawSegmentFile`, sua descrição e os campos `term_dir`, `doc_meta`, `block_dir` e `data_offset`, em [src/raw.rs:1637](https://github.com/arclabs561/postings/blob/edb02683c6359efe2034a6cc71b64c9aaed90cb0/src/raw.rs#L1637). O mesmo retorno também trouxe a documentação de metadados residentes. [Audit/006](xerj/audit/006.json).

Depois, usamos os símbolos e vocabulário retornados pelo próprio XERJ. Foram seis ações de recuperação no total: cinco buscas e uma definição. A sequência inteira está em [queries.json](queries.json), com todos os hits nos audits.

| Etapa | Retorno observado |
|---|---|
| Pergunta ampla sobre diretórios e leitura sob demanda | Leitor RawSegmentFile e documentação correspondente |
| Pergunta ampla sobre validação | Um guard de capacidade pouco pertinente, tipo de erro e documentação útil; não foi resposta perfeita |
| `def from_file_range` | Declaração em src/raw.rs:1701 e lista de definições, que revelou os helpers de leitura |
| Busca pelo helper `read_posting_block_range_into` | Assinatura e wrapper; ainda não mostrou seu corpo |
| Busca por offset/leitura/checksum | Corpo de abertura com validação da faixa externa, header/layout/footer e diretórios |
| Busca combinando helper e `data_offset` | Corpo de `read_postings_range`, com faixa, offset físico, leitura e checksum condicional |

A última consulta foi:

```sh
xerj search --url http://127.0.0.1:19701 --prefix demo-postings-0156 \
  'read_posting_block_range_into data_offset' -k 3 --full 7000 --json
```

Ela recuperou [read_postings_range](https://github.com/arclabs561/postings/blob/edb02683c6359efe2034a6cc71b64c9aaed90cb0/src/raw.rs#L3017), mostrando a chamada `checked_range`, cálculo de offset relativo com `segment_file_offset`, leitura com `read_exact_at_positional` e verificação dos blocos quando existe mapa de checksums. [Audit/012](xerj/audit/012.json). O trecho é do helper de faixa de postings; não o apresentamos como corpo do helper de bloco cujo resultado anterior só tinha assinatura.

Os resultados contêm metadados de localização diferentes: em hits de arquivo inteiro, start_line pode ser 1, enquanto match_line e offsets/ordinal localizam a passagem. Os cinco exemplos citados foram conferidos contra o Git fixado, com linhas verificadas separadamente em [evidence.json](evidence.json). Não reescrevemos os locators originais como se fossem linhas exatas.

## O que o exemplo permite concluir

| Elemento | Conclusão sustentada |
|---|---|
| Descoberta | Bateia forneceu o candidato, a leitura inicial e o contexto da investigação |
| Recuperação | XERJ navegou da pergunta para tipos, definições, implementação e verificações da mesma fonte |
| Princípio | Separar metadados residentes e payloads endereçáveis permite ler a faixa pedida |
| Condições | Metadados/diretórios continuam consumindo RAM; offsets/layout e integridade precisam ser válidos; checksums são condicionais ao formato |
| Decisão | Investigar um leitor com acesso seletivo, comparando-o ao formato, caches e contratos atuais do alvo antes de propor migração |
| Limite | Biblioteca experimental/feature específica, sem integração executada; nada aqui mede latência, memória ou benefício no XERJ |

No ensaio anterior, o XERJ-alvo já tinha cache FTS e abertura de estatísticas separada. Este exemplo não altera essa condição: a hipótese de mudança continua ligada a aberturas frias/misses com custo relevante. Também não transfere automaticamente scoring, posições, total exato, compactação ou lifecycle de segmentos. A biblioteca declara que tokenização, armazenamento de documentos e scoring são responsabilidades do consumidor.

Este é um uso complementar escolhido para a demonstração, não prova de que uma ferramenta seja incapaz de fazer outra etapa ou de superioridade frente a um agente diferente.

## Texto utilizável no email

> Fiz um experimento complementar: usei o Bateia para encontrar uma referência sobre postings e depois o XERJ para investigar essa mesma fonte. O Bateia trouxe o projeto e a descrição da estratégia; no XERJ, parti de uma pergunta e cheguei aos tipos e aos trechos de leitura por faixa e validação. Isso deu um exemplo concreto de investigação de uma ideia seguida de recuperação do código que a sustenta, ainda sem alegar ganho ao aplicá-la no XERJ.

## Ocorrências e reprodução

O nó temporário recusou indexar antes de qualquer escrita porque o volume estava com cerca de 95% de ocupação, apesar de aproximadamente 15,8 GB livres. Uma alteração da configuração dinâmica foi aceita pelo endpoint, mas o preflight continuou reportando 95% nessa execução. Não investigamos a causa como bug. Reiniciamos **somente esse nó** com `limits.disk_flood_stage_percent = 97` e limitamos a estimativa de índice a 512 MB; o plano estimou 1,4 MB. A indexação então concluiu. Não apagamos dados nem modificamos outros nós. As duas recusas, resposta da configuração e TOML permanecem nos artefatos.

O [nó foi encerrado](xerj/node-lifecycle.json) e [audit/013](xerj/audit/013.json) confirmou que o endpoint ficou inacessível. Logs brutos de startup não foram copiados porque podem conter tokens de onboarding de uso único; comandos e lifecycle foram preservados. Nenhum programa do projeto postings foi executado; foram executados os binários Bateia/XERJ para descobrir, indexar e pesquisar. Nenhuma publicação ou mensagem externa foi enviada.

Para repetir a leitura, fixe o commit acima após adquirir o repositório; resultados de uma nova descoberta podem variar. O [runner](run.py) conserva ações e snapshots; os hashes dos binários estão em [protocol.json](protocol.json). O [controle do nó](node.py) registra PID/argv e só sinaliza o processo próprio. Clones, dados de índice e binários ficam em /tmp; não entram no repositório. O manifesto SHA-256 acompanha os artefatos finais.
