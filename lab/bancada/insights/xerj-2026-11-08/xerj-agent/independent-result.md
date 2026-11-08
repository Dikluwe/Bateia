# Ideias independentes obtidas usando o XERJ

## Método e alcance

Usei somente sete buscas pelo wrapper autorizado `search.py`, sem ler fontes diretamente do disco, relatórios do Bateia, propostas de outro agente ou web. O wrapper informou consultas lexicais BM25 sobre `semantic_text`; não executei busca semântica nem testes do produto. As observações abaixo são sobre trechos recuperados de um índice limitado a 19 arquivos. Uma lacuna nesses resultados não demonstra ausência em todo o XERJ. Os prefixos de fonte distinguem revisões: `xerj-org--xerj--1557de2f` (main) e `xerj-org--xerj--708e9745` (corpus-hub). Citações referem-se às linhas devolvidas pelo mecanismo, e não a uma inspeção independente.

## 1. Recibo por alegação com verificações separadas

**Observado.** O formato de pack verifica hash do arquivo, assinatura e versão antes de extrair (`xerj-org--xerj--1557de2f/engine/crates/xerj-hub/src/pack.rs`, linhas 20–25). O contrato já explica que checksum prova integridade e assinatura prova autoria (`xerj-org--xerj--708e9745/tools/xerj-code/hub/PROGRAM-100.md`, linhas 159–160, passagem iniciada em 150). G7 é um spot-check manual de recuperação: cinco consultas, top-5, relevante/parcial/miss; não é execução do comportamento descrito (`PROGRAM-100.md`, linhas 29–35). Feedback contém sucesso/lentidão de pull e avaliação do agente, com envio automático em MCP (`xerj-org--xerj--1557de2f/engine/crates/xerj-hub/src/feedback.rs`, linhas 5–18).

**Proposta concreta.** Acrescentar ao resultado um recibo versionado vinculado a uma alegação explícita, fonte/pin/span e, quando houver, execução. Expor campos distintos: `artifact_integrity`, `publisher_authenticity`, `retrieval_judgment`, `claim_execution` e `user_feedback`. Cada campo identifica seu verificador e estado (`not_checked`, `passed`, `failed`, `not_applicable`), sem um booleano geral `correct`. `claim_execution: passed` exige identificador do caso, entradas, ambiente, esperado e observado; sem execução, permanece `not_checked`, mesmo com pack assinado, G7 alto e nota 5. Uma execução aprovada sustenta apenas a alegação e as condições registradas.

**Já existe.** Além das verificações de pack e G7, Query-1 oferece programa, fatos e spans para lookup/comparação, e distingue resultado final de programa não checado (`xerj-org--xerj--1557de2f/docs/XERJ_QUERY_1.md`, linhas 3–8, 34–46 e passagem 61–89). Portanto, a novidade proposta é uma apresentação e contrato comum entre superfícies com escopo explícito por alegação; não inventar assinatura, proveniência ou executor novos.

**Limite.** O recibo organiza o que foi verificado; não transforma o oráculo ou fatos de entrada em verdade. Não encontrei nos trechos um contrato transversal como esse, mas não provei que inexiste.

**Experimento de rejeição proposto, não executado.** Preparar um pack válido e assinado contendo função que falha num limite de entrada; obter hit relevante e registrar feedback positivo. Rejeitar a implementação se qualquer desses três sinais tornar `claim_execution` aprovado antes do caso, ou se a falha do caso apagar o fato de que a assinatura é válida. Comparar com o mesmo caso corrigido: só o campo e a alegação cobertos pela nova execução podem mudar.

## 2. Condições de aplicabilidade como parte verificável da resposta

**Observado.** O catálogo registra domínio, uma descrição `ask_for`, pin, dependências e uma nota G7 (`xerj-org--xerj--1557de2f/engine/crates/xerj-common/src/hub.rs`, passagem 140–171, especialmente linhas 147 e 167). O roteador consulta lockfiles e fornece ranking por evidência (`xerj-org--xerj--1557de2f/engine/crates/xerj-hub/src/router.rs`, linhas 162–180, passagem iniciada em 137). Há um contraexemplo documentado mais forte que mera inferência: Query-1 conta membros ao perguntar quantos dispositivos existem e devolve a contagem como verificada, embora o namespace não tenha fatos de dispositivos (`xerj-org--xerj--1557de2f/docs/XERJ_QUERY_1.md`, linhas 215–221, passagem iniciada em 208). O documento já prevê correção para rc.102.

**Proposta concreta.** Expor uma avaliação de aplicabilidade por resposta: sujeito/tipo de entidade, operação solicitada, escopo temporal ou revisão, versão do software e condições de ambiente relevantes. Para cada condição, retornar origem e `satisfied`, `violated` ou `unknown`. Um resultado de execução pode ser matematicamente válido e ainda ter `applicability: violated`. Não servir como resposta final uma contagem quando o tipo da entidade ou a operação não corresponde à pergunta; não presumir compatibilidade de versão apenas porque o nome da dependência coincide. Começar com as condições estruturadas que o consumidor já fornece; permitir desconhecido em vez de inferir suporte universal.

**Já existe.** Há distinção entre fatos tipados, UNKNOWN/CONFLICT, fechamento do conjunto para contagens, identificadores de sujeito e consultas por revisão (`XERJ_QUERY_1.md`, linhas 34–46, 65–89, passagem 139–169 e linhas 205–212). A correção do substantivo em contagens já está anunciada; essa parte não é reivindicada como ideia inédita. A extensão proposta é um contrato público de condições e sua origem, também nos resultados de corpus e seleção por dependências.

**Limite.** Classificar texto livre em condições também pode errar. O MVP deve limitar-se a predicados explícitos e apresentar os demais como desconhecidos; não bloquear descoberta de código só porque falta uma condição, mas bloquear sua promoção a resposta aplicável.

**Experimento de rejeição proposto, não executado.** Usar o mesmo namespace com duas perguntas que diferem apenas em pessoas/dispositivos e outra dupla contagem/soma. Em seguida, usar o mesmo símbolo de biblioteca sob duas versões, apenas uma coberta pelo exemplo. Rejeitar se o caso incompatível mantém uma resposta final aplicável, ou se o caso compatível deixa de funcionar por um filtro lexical sem vínculo com o tipo/operação/versão. Medir também recusas indevidas, não só falsos positivos.

## 3. Exportação e replay de um contraexemplo mínimo

**Observado.** Lookup/comparação já trazem programa, fatos e spans, mas contagens identificam o programa por hash e não mostram programa nem registros contados (`xerj-org--xerj--1557de2f/docs/XERJ_QUERY_1.md`, linhas 3–6 e 205–212). Query-1 permite voltar a uma revisão e tem exigência de quote igual ao span (`XERJ_QUERY_1.md`, passagem iniciada em 139, até linha 169). Packs preservam o snapshot opaco byte a byte (`xerj-org--xerj--1557de2f/engine/crates/xerj-hub/src/pack.rs`, linhas 3–17). O processo do Hub já publica misses, kills e A/B pré-registrados (`xerj-org--xerj--708e9745/tools/xerj-code/hub/PROGRAM-100.md`, linhas 163–168). Isso oferece partes de reprodução, mas os trechos recuperados não oferecem um artefato único para contestar uma resposta.

**Proposta concreta.** Criar exportação de caso que reúna: pergunta original; requisição de busca e parâmetros; versão do engine/modelo; hash e versão do pack ou revisão do namespace; candidatos efetivamente entregues com localizadores; programa executado completo quando autorizado; condições; esperado, observado e logs. Um comando de replay deve separar duas fases: repetir a execução sobre entradas congeladas e repetir a recuperação original. Assim, mudar o ranking não é confundido com mudar o comportamento do executor. Para contagens, disponibilizar programa e conjunto/identidades necessárias a conferir a cardinalidade, ou declarar explicitamente quais dados privados impedem replay; o hash sozinho identifica, mas não permite reexecutar.

**Já existe.** A proposta compõe snapshots, pins, hashes, revisões e rastros já disponíveis. Não é propor commit pinning ou signed packs novamente. Também preserva a separação entre respostas de busca e o mecanismo Query-1, que é opt-in e não extrai fatos de documentos arbitrários (`XERJ_QUERY_1.md`, linhas 7–12 e 215–217).

**Limite.** Um subconjunto mínimo pode mudar fechamento de conjunto, conflito ou ranking; preservar o contexto necessário é parte do formato. Exportação precisa respeitar licença e privacidade. Se dados não podem sair, o replay deve aceitar execução local e marcar a limitação de auditoria externa, em vez de chamar o caso de reproduzível publicamente.

**Experimento de rejeição proposto, não executado.** Exportar o contraexemplo pessoas/dispositivos documentado; reexecutar offline na mesma revisão e exigir o mesmo resultado observado. Alterar apenas a revisão: o replay congelado deve manter o resultado antigo, enquanto uma execução nova pode mudar. Alterar um byte do artefato: exigir erro de integridade. Por fim, mudar apenas o ranking e verificar que o replay do executor permanece igual e que a divergência aparece apenas na fase de recuperação. Rejeitar se o pacote não contém dados suficientes para essas distinções ou se o minimizador elimina a falha.

## Consultas realizadas

1. `provenance relevance applicability tested behavior counterexample signatures feedback correctness`
2. `pack manifest sha256 signature verification provenance evidence`
3. `G7 queries graded suite relevance applicability wrong negative counterexample`
4. `router Evidence abstain ambiguous dependency licence ask_for confidence`
5. `lockfile frozen restore fingerprint replay snapshot query`
6. `feedback never ranking correctness usefulness found basis automatic`
7. `XERJ Query-1 limits provenance source quote verified entailment program hash`
