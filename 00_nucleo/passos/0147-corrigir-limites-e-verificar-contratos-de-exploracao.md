# 0147 — Corrigir limites e verificar contratos de exploração

**Data:** 2026-10-29. **Contrato ex-ante de correção e ampliação de verificação**, após a consolidação0146.

## Intenção e comportamento esperado

Corrigir MEM-04: memória selada de geração999999 continua válida para leitura/inspeção, mas não admite operação que avance geração. Feedback, discover e expand devem rejeitá-la com ErrInvalidRunInput antes de alterar entrada ou executar buscas/leituras. Geração999998 pode avançar uma vez para999999 e o resultado deve validar. Não saturar ou reiniciar geração silenciosamente, não remover histórico, não fabricar leitura e não alterar schemas/algoritmos de ID.

O container de memória em L1 possui o limite e uma guarda de avanço reutilizada por feedback e admissão. ValidateInsightSiteOptions recusa avanço impossível antes de compor/executar plano; a admissão de expand usa a mesma política. L2 conserva transporte e publicação, L3 efeitos e L4 composição. Nenhuma duplicação de política na CLI.

## Revisão de prompts e critérios

Atualizar o contrato compartilhado e os prompts proprietários de memória/admissão antes do código. Consolidar os prompts do plano de sítios, retirando a necessidade de interpretar precedências cronológicas para conhecer as regras atuais. Preservar os textos anteriores no Lab e distinguir a intenção nova da consolidação posterior ao código.

MEM-02: feedback nos dois estados permitidos preserva todos os campos não alterados, inclusive Lead, LastRead e tentativas. ORD-02: tabela de desempates por geração, presença/posição de referência e ID, sob permutações, sem efeitos/mutação de memória; não mudar o comparador. SRC-01: testar as três origens, recibos/LastRead incoerentes, duplicidade e limites inclusivos/excedidos de32 KiB por documento, oito documentos e128 KiB agregados, sem modificar a política ou unificar erros históricos. Não promover cobertura desses controles a cobertura de todos os limites do sistema.

## Verificação e preservação

Primeiro executar as novas regressões contra o código anterior e registrar a falha esperada de geração; depois corrigir e exigir passagem. Testar a entrada pública de descoberta com portas controladas e rejeição antes de busca/leitura. Conferir feedback pela CLI sem publicação de saída inválida. Reproduzir o contraexemplo sintético com o binário novo e manter intactos os resultados0146 e as bancadas seladas0142–0145.

Evidência em lab/bancada/prompts/0147. Executar lint, auditoria de arquitetura, testes completos, vet e build pelo gate canônico. L0 contém somente contratos/sínteses; scripts, logs e snapshots permanecem no Lab. Não alterar corpos de testes históricos selados da refatoração0145; acrescentar novos consumidores. Atualizar a matriz somente conforme o resultado efetivamente demonstrado.
