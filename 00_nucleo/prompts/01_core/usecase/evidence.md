# Deterministic evidence use case

Selecionar candidatos da peneira, criar planos limitados, coletar bundles por uma
interface de conteúdo e avaliar regras versionadas por correspondência lexical.
Deduplicar fontes, reduzir força para vendor/fixture/tradução, produzir explicações
por templates e nunca converter falha ou orçamento incompleto em ausência.

O plano recebe `IdeaSpec`; a coleta registra resultado por alvo e árvore completa;
símbolos exigem declaração relacionada; evidência auxiliar não herda força
estrutural. Calcular métricas somente de bundles e julgamentos congelados. Uma
falha de candidato não aborta os demais.
Verificar o contexto antes de cada chamada ao provedor e depois que ela retorna;
um cancelamento observado não pode iniciar uma nova coleta.

Rejeitar catálogos v2 com identidade, termos ou padrões genéricos/duplicados,
proveniência ausente ou hash divergente. Quando a regra exigir independência,
uma única fonte forte não produz `PROVEN`; conteúdo duplicado conta uma vez e
documentação isolada permanece abaixo de `PROVEN`.

## Passo 0040

Aplicar o contrato de autoridade, sinais lexicais sem comprovação comportamental, procedência local, seleção de produção, negação conservadora e lacunas acionáveis de `00_nucleo/passos/0040-corrigir-autoridade-da-evidencia.md`. Preservar determinismo e ausência de I/O no núcleo.

Limites de quantidade/profundidade devem registrar SKIPPED_BUDGET e Complete=false para arquivos elegíveis omitidos. Arquivo encontrado não equivale a arquivo examinado.

## Revisão 0043 — autoridade e ligação
Aplicar o contrato comum de evidence_authority antes de autorizar capacidade ou operação. Evidência lexical/documental não ganha PROVEN. Negação deve conservar polaridade. Na extração de mecanismos aceitar SHA-256 hexadecimal produzido pela avaliação e formato prefixado legado, verificando conteúdo e projeções de tipo file→symbol e test→test_implementation. Preservar IDs reais: catálogo sem regra compatível retorna lacuna explícita, nunca traduz nomes em prova.

# 0044 — Recuperação antes da seleção
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md.
L1 somente orquestra porta opcional, sem IO concreto. RetrieveImplementationEvidence usa budget finito padrão e reporta not_evaluated em provedor ausente. SelectImplementationPaths preserva manifesto/README e antepõe hits de caminhos relativos presentes na árvore, com limite de código. ValidateImplementationHits copia relatório, descarta fontes não extraídas ou hashes/intervalos divergentes e marca partial; não apaga identidade do índice.
Funnel/analyze registra relatório; CollectEvidence busca requisitos e prioriza árvore antes do orçamento, sem promoção do status de capacidades. Provedores sem porta preservam contrato antigo. A seleção por conteúdo é distinta da análise de comportamento. Relatório mostra símbolos com linhas/hash e snippet dos bytes revalidados, no máximo20 linhas por hit/8hits apresentados, JSON retém todas as evidências.

0046: conservar query agregada da busca de implementação no escopo CollectEvidence e usar ValidateImplementationRequirements após coleta. Projeção de requisitos acompanha bundle, sem alterar detectores/decisão.

## Passo 0072 — recuperação dirigida anterior à amostra

Compilar consultas somente para requisitos cobertos por uma regra `DECLARED`, usando
termos técnicos canônicos da própria regra. Enumerar a árvore antes da seleção,
reservar de modo justo até `RequirementFiles` por requisito dentro dos limites globais
e só então preencher a amostra geral. Casamento de caminho é exato por basename, stem
ou token identificador, nunca substring. Headers C++ sem extensão são elegíveis apenas
quando seu basename completo é termo da consulta. Deduplicar caminhos compartilhados
sem dupla cobrança e registrar candidatos, examinados e omitidos em ordem canônica.
Superfície plausível não lida produz `COLLECTION_INCOMPLETE`, confiança zero e nenhum
gap do candidato; o nome do arquivo não promove evidência.
Quando a cobertura dirigida de um requisito estiver `TARGETED_EXAMINED`, incompletude
da amostra geral ainda impede concluir `ABSENT`, mas não apaga evidência positiva ou
conflitante encontrada nos bytes dirigidos. Bundles legados sem cobertura por requisito
continuam sujeitos à completude global.

## Passo 0073 — ocorrência exata dentro do arquivo

Separar termos de recuperação (`TargetHints`) de padrões probatórios. Localizar padrões
fortes e auxiliares nos bytes coletados, ignorando comentários e literais em código,
respeitando fronteiras de identificadores e sem interpretar macros ou declarações
divididas. Fontes de regra carregam linha baseada em 1, snippet normalizado, padrão
canônico e classe da mesma ocorrência. Ordenar deterministicamente por autoridade da
classe, path, linha, padrão e snippet; path-only permanece diagnóstico e nunca fonte.
