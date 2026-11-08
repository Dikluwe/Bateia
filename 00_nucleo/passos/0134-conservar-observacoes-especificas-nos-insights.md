# 0134 — Conservar observações específicas nos insights

**Data:** 2026-10-28. **Contrato ex-ante.** Origem: foco e diagnóstico 0133; autorização do usuário para começar a materialização.

## Escopo e ownership

Acrescentar investigação offline de fontes com observações declaradas e localizadas, separando condições operacionais da cobertura do extrator. L1 possui validação, ligação aos bytes, canonicalização e projeção; L2 somente escolhe o contrato de transporte conforme a versão e apresenta JSON pelas proteções existentes. Sem busca nova, execução de fonte, runtime, IA ou classificação semântica automática. O alvo continua UNKNOWN; hipóteses nunca SUPPORTED. Não mudar catálogo, extração compartilhada, política de reúso ou síntese durável.

## Compatibilidade e superfície

`idea insights analyze --input ... [--output ...]` continua aceitando versão `1` com as mesmas saídas/IDs/preimages e erros existentes. Uma entrada `schema_version: "2"` chama `InvestigateSourceInsights` e retorna `insight-source-analysis-v2`, com identidade própria. `InvestigateInsights` continua dono de v1 e rejeita v2; v1 rejeita observações não vazias. O campo opcional `observations` em cada fonte é omitido na canonicalização quando vazio, preservando hashes v1. Nenhum artefato histórico é reescrito.

## Observações e referências

Cada observação contém `mechanism_class`, `kind` (`STRATEGY`, `PRECONDITION`, `FAILURE_MODE`, `TRADEOFF`, `LIMITATION`), `statement` e `references`. A classe deve ser conhecida e solicitada. A fonte é a do bundle que contém a observação. Cada referência contém `path`, `content_hash`, `start_line` e `end_line`, com caminho/hash ligados a um item do mesmo bundle e intervalo real. A saída acrescenta o trecho exato desses bytes e um ID determinístico. Hash sem ou com prefixo `sha256:` é canonicalizado para o prefixo. Nenhum texto de erro, URI física ou fonte arbitrária é usado como referência.

Limites: até 32 observações por fonte, 128 no pedido, 1–8 referências por observação, 256 no total, até 20 linhas e 4096 bytes por trecho textual não vazio, declaração de até 1000 bytes; manter o teto de 2 MiB e demais limites do contrato de investigação. Rejeitar tipos desconhecidos, texto inválido, referências vazias/duplicadas, caminhos/hashes/linhas inválidos e observações duplicadas após canonicalização. Não truncar declarações silenciosamente. Ordenar fontes, observações e referências de forma determinística, sem mutar as entradas; mudança de declaração/tipo/referência muda as identidades causais de v2. Não descartar contradições entre declarações diferentes nem resolvê-las por regra lexical.

## Resultado v2

O relatório contém `source_analyses`, hipóteses, lacunas e proveniência. Cada análise identifica fonte/classe, definição do catálogo explicitamente `CATALOG_TEMPLATE`, `extraction_scope` com idiomas/ecossistemas, referências e estado da classificação, observações declaradas com trechos e origem `DECLARED_NOT_REASSESSED`. Não apresentar limites de linguagem como pré-condições operacionais.

Pré-condições e trade-offs nas hipóteses vêm das observações daquela fonte/classe, com IDs em `evidence_refs`; sem declaração, pré-condições ficam vazias e a lacuna é explícita. Estados da classificação continuam os da extração anterior; observações não os promovem. Sem classificação, conservar observações numa análise INCONCLUSIVE com lacuna de mecanismo não observado e não gerar hipótese. Fontes/classes sem observações não recebem estratégia operacional inventada; a definição do catálogo continua rotulada como tal. Observações não são reavaliadas semanticamente, autenticadas externamente ou tomadas como contraevidência formal apenas por seu texto/tipo. Licença/adopção/efeito conservam suas restrições anteriores.

As hipóteses derivadas preservam plano conservador e limite de quantidade. A identidade v2 inclui contrato, pedido normalizado e catálogo. Não duplicar a classificação: usar a jornada v1 sobre cópia sem observações para obter mecanismos/gaps, sem materializar artefatos legados. Preservar a ligação do mecanismo original e seus testemunhos; a análise v2 é uma projeção distinta, não uma reescrita dos mecanismos históricos.

## Aceite

Testes puros de tipos e análise: vínculo/localizadores, rejeições e limites, trechos exatos, ordenação/imutabilidade/causalidade, fontes/classes múltiplas sem mistura, observações não promovendo avaliações, falta de observações ou assessments, campos semânticos distintos da cobertura. Testes CLI: v2 pela superfície pública, JSON estrito, saída/entrada protegidas, short writes, segredo e nenhuma publicação em falha. Manter oráculos existentes e replay 0106 idêntico. O caso conhecido `conc` deve expor configuração do limite, submissão bloqueante e ausência de contexto como declarações localizadas, sem execução da fonte ou efeito medido no Bateia. Executar `make verify`; registros novos em `lab/bancada/insights/0134`.
