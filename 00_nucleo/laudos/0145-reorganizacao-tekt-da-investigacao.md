# 0145 — Reorganização Tekt da investigação

Materializado o [contrato0145](../passos/0145-reorganizar-fronteiras-de-investigacao-sob-tekt.md). A revisão estrutural do programa foi acompanhada de um corte de ownership na expansão recente: admissão, fontes congeladas, fronteira, estado, execução e projeção do resultado passaram a ter responsáveis explícitos. A fachada reutiliza as origens entregues pela admissão; não faz uma segunda interpretação ignorada das fontes. O probe de argumento saiu de L4 para L2; infraestrutura e composição continuam em seus papéis.

O [mapa atualizado](../arquitetura/README.md) liga capacidade, responsável e limite. A reorganização conserva APIs e políticas; não cria arquitetura ou runtime novo. A guarda de ownership fica no Lab e integra os gates, com controles negativos para impedir a volta de regras ao responsável errado.

Os [recibos e provas](../../lab/bancada/arquitetura/0145/README.md) registram a pilha completa aprovada, preservação dos corpos de consumidores e regras movidas, nove comparações offline contra o binário anterior e replays históricos. Sete relatórios permaneceram idênticos byte a byte; duas rejeições conservaram erro e saída anterior. A evidência histórica e os arquivos de orientação foram preservados.

Esta revisão verifica gravidade, efeitos e ownership explícito do recorte. Não equivale a prova universal da semântica do programa ou da utilidade dos insights. Estados operacionais e força de evidência permanecem distintos.
