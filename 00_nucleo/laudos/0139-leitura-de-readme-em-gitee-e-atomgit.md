# Laudo 0139 — Leitura de README em Gitee e AtomGit

**Data:** 2026-10-28. **Resultado:** materialização e verificação aprovadas para o [contrato0139](../passos/0139-ler-readmes-de-gitee-e-atomgit.md), sob [ADR-0020](../ADR/ADR-0020-leitura-de-readme-com-contexto-por-sitio.md).

## Comportamento entregue

A descoberta multifonte recebe leitores próprios para Gitee/AtomGit. A resolução exige metadata pública e identidade correspondente, fixa commit da branch padrão e requisita árvore/conteúdo nesse SHA. README regular raiz, validação de blob Git e texto, quatro GETs máximos por referência, limites de32KiB/documento e128KiB/conjunto, sem redirects/download_url/retries ou execução de candidatos.

L1 conserva política, seleção de README, admissão, quotas e contexto; L3 possui HTTP e protocolo; L4 injeta leitores/token correspondentes. L2 e contratos globais de evidência não foram alterados. Contextos novos carregam identidade do sítio e cartões com IDs próprios; tasks.reading_card_refs permite à IA usar seus trechos automaticamente. Bundles/exploration existentes mantêm escopo GitHub. Max-passages limita cartões agregados, com precedência GitHub preservada e lacuna explícita para contexto sem cartão.

## Verificação e preservação

O gate canônico make verify passou: lint com11 controles negativos, auditoria de378 arquivos produtivos, testes do módulo raiz, vet, build e diff-check. A primeira execução encontrou restrição de socket local no sandbox; o gate foi repetido com loopback autorizado e passou integralmente. As evidências da execução limitada foram conservadas separadamente.

Fixtures usam conectores e leitores reais com transporte HTTP controlado. Duas forges recebem credenciais distintas, fazem cinco buscas e quatro GETs de leitura por host, fixam revisão e produzem contextos/cartões distintos para mesmo owner/repo/commit/conteúdo. Controles negativos verificam identidade/auth, metadata/branch/SHA inválidos, truncamento/paginação, symlink, redirects, corpos grandes, base64/tamanho/blob inválidos, binário, cancelamento e budgets; limite exato32KiB com SHA256 também passou.

Seis verificadores históricos (0106,0134–0138) passaram com o binário final, preservando replays e contextos congelados. Os477 arquivos Go de testes anteriores e195 arquivos anteriores de bancada permanecem idênticos. Foram acrescentados três arquivos produtivos e quatro de testes; cinco arquivos produtivos existentes foram alterados. Prompts e hashes de linhagem estão no recibo novo.

## Limites

Não houve ensaio autenticado online dessas forges neste recorte; as fixtures não certificam disponibilidade remota. Árvores com100 entradas/página incompleta são recusadas conservadoramente. A leitura admite identidades ASCII owner/repo e os quatro nomes de README já aceitos pelo coletor. Demais fontes permanecem catálogo. Metadados continuam não verificados; leitura de documentação não comprova código, autoria, relevância, novidade ou efeito.

A [bancada0139](../../lab/bancada/insights/0139/README.md) contém uso, protocolos consultados, logs e [recibo consolidado](../../lab/bancada/insights/0139/verification.json). Nenhuma evidência histórica foi reescrita e nenhum candidato foi executado.
