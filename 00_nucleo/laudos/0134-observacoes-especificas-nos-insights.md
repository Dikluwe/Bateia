# 0134 — Observações específicas nos insights

**Data:** 2026-10-28. **Contrato:** [passo 0134](../passos/0134-conservar-observacoes-especificas-nos-insights.md). **Decisão:** [ADR-0015](../ADR/ADR-0015-observacoes-localizadas-na-investigacao-de-insights.md). **Evidências:** [bancada](../../lab/bancada/insights/0134/README.md).

## Entrega

`idea insights analyze` aceita pedidos com `schema_version: "2"` e produz `insight-source-analysis-v2`. Cada fonte pode fornecer observações de estratégia, pré-condição, modelo de falha, trade-off e limitação. O relatório conserva declaração, origem, identidade causal e referências a arquivo/hash/linhas, acrescentando trechos dos bytes fornecidos.

`source_analyses` separa definição genérica do catálogo, cobertura do extrator e observações específicas. Linguagens e ecossistemas não ocupam as pré-condições operacionais das hipóteses v2; essas condições e os trade-offs declarados vêm das observações da mesma fonte/classe, com referências aos seus IDs. Sem pré-condições declaradas, o campo fica vazio e a lacuna é explícita. Sem mecanismo classificado, observações são conservadas como inconclusivas, sem produzir hipóteses por si sós.

L1 possui validação, canonicalização e projeção. Reutiliza a classificação existente por delegação sobre cópia v1 sem observações; os mecanismos mantêm suas identidades. L2 seleciona o contrato pela versão de entrada e usa as proteções existentes de JSON estrito, segredo, alias e publicação. Não há novo runtime, I/O em L1 ou alteração da política de reúso. README, visão do produto, mapa de responsabilidades e índice foram atualizados.

## Caso conhecido: aprender com conc

O pedido novo conserva os 29 itens e as avaliações congeladas de `sourcegraph/conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b`, acrescentando seis observações manuais com onze referências. Agora o relatório explicita a necessidade de configurar `WithMaxGoroutines` antes do uso, o padrão sem teto, submissão bloqueante e ausência de contexto na assinatura de `Go`. Também conserva que o canal não tem buffer e que o teste inspecionado observa concorrência máxima/conclusão.

Esses detalhes ajudam a distinguir teto de workers, fila limitada e política de admissão. O resultado continua PARTIAL para o mecanismo e INCONCLUSIVE para a hipótese, com alvo UNKNOWN e NOT_AN_ADOPTION_DECISION. O comando não analisou automaticamente a fila do Bateia e não verificou benefício de uma alteração no alvo. A comparação manual anterior está no [laudo 0133](0133-foco-em-insights-e-primeira-investigacao.md).

## Compatibilidade e verificação

A saída v1 do caso 0106 foi reproduzida byte a byte, mantendo SHA-256 `1de0f841b35191410e1415a1e25a9673eb5e3225051d237c2a438252d4f01527`. Os cinco artefatos históricos comparados à baseline não mudaram. Dos 804 arquivos Go preexistentes, apenas os três componentes declarados da superfície de insights foram alterados; todos os testes anteriores permaneceram iguais. Cinco arquivos Go foram adicionados, incluindo três arquivos de testes. A [linhagem](../../lab/bancada/insights/0134/lineage.json) vincula os oito arquivos aos respectivos prompts.

Os testes da capacidade verificaram condições/trechos localizados, integridade e limites de entrada, determinismo e imutabilidade, causalidade dos IDs, isolamento entre fontes/classes, ausência de promoção de evidência e proteções da CLI. O replay do caso verificou reordenação/prefixo de hash, recusa de adulteração preservando a saída e conservação das observações sem avaliações, sem hipótese. Conferir os bytes não confirma a interpretação.

`make verify` passou integralmente: lint sem violações, onze controles de regressão do linter, auditorias arquiteturais sobre 356 arquivos produtivos, suíte raiz, vet, build e diff-check. O [recibo final](../../lab/bancada/insights/0134/final/verification.json) registra os comandos e códigos de saída. A execução final precisou de loopback para servidores HTTP de fixture; as tentativas restritas falharam somente nessa abertura de sockets e foram preservadas, sem relaxamento de testes. Não houve nova coleta da fonte investigada.

## Limites

O recorte conserva e localiza interpretações fornecidas; não faz descoberta automática de observações, comparação semântica geral, autenticação externa ou execução do código investigado. Declarações contraditórias permanecem declaradas e não viram contraevidência formal por palavras-chave. Catálogo e síntese durável compartilhados conservam seus contratos anteriores. O caso conhecido é diagnóstico/calibração, sem holdout ou demonstração de utilidade geral. Gates selados com módulos próprios não foram promovidos pela suíte raiz.
