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
