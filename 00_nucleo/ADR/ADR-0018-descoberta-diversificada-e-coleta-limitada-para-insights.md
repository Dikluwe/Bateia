# ADR-0018 — Descoberta diversificada e coleta limitada para insights

**Status:** ACEITO PARA O RECORTE 0137. **Data:** 2026-10-28. **Relacionados:** ADR-0008, ADR-0017; passos 0136/0137.

## Contexto e decisão

O modo explore prepara contexto automaticamente, mas dependia de bundles escolhidos/coletados pelo chamador. O usuário autorizou iniciar descoberta diversificada para ampliar o repertório da IA.

Adicionar `idea insights discover`: pergunta/familiaridade originam cinco consultas lexicais independentes. Todas são tentadas dentro dos deadlines, mesmo após sucesso inicial. O catálogo GitHub existente é o primeiro provedor composto; outros conectores não são implicitamente incluídos. Candidatos mantêm query IDs e conflitos de revisão. A seleção prioriza cobertura das consultas, autores diferentes e menos matches familiares, sem usar estrelas como critério.

Coletar somente um README de raiz de cada referência selecionada, em SHA imutável sem conflito, com validação de identidade/path/bytes e limite de admissão de 32 KiB. Entregar bundles e contexto de explore no mesmo relatório, dispensando curadoria manual de avaliações e montagem de planos de coleta. Coleta limitada, revisões opacas, falhas e quotas permanecem explícitas.

## Fronteiras e consequências

Plano, execução, seleção e política de coleta têm owners L1 separados. Catálogo e leitura são portas; o cliente GitHub L3 existente realiza HTTP e seus controles. L2 trata flags/sinais/JSON protegido; L4 injeta as duas portas no ramo online específico. Os outros comandos de insights permanecem offline, sem clientes/cache/runtime novo.

Cobertura de consultas/autores não comprova diversidade conceitual ou utilidade. READMEs não comprovam comportamento em execução. As heurísticas não traduzem a pergunta nem inventam sinônimos/domínios. O relatório separa estado operacional de novidade/efeito não avaliados; resultados parciais são publicados antes de erro de interrupção ou falha total.

O corte automatiza busca e coleta inicial, com no máximo cinco operações de catálogo e oito coletas. O transporte pode fazer requisições internas de metadados sob o mesmo deadline. Não clona, instala ou executa candidatos. O limite pequeno pode rejeitar READMEs grandes; não há truncamento que altere sua identidade. Ampliar provedores, estratégias semânticas ou seleção de arquivos exige contratos próprios.
