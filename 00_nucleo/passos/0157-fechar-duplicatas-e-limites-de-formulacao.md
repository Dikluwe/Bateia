# 0157 — Fechar duplicatas de flags e limites de formulação

## Intenção e contrato antes da implementação

Corrigir a divergência entre os prompts CLI e o parser: explore deve rejeitar repetições de --algorithm-ideas e --review-alternatives; discover deve rejeitar repetições de --algorithm-ideas. O responsável é L2 nos adaptadores existentes; nenhuma regra de investigação muda em L1. O contrato de unicidade inclui formas com um/dois hífens, valores implícitos e explícitos true/false, inclusive repetições com valores conflitantes. Uma ocorrência continua admitindo as formas booleanas válidas e ausência conserva false; familiar continua repetível.

Rejeitar duplicatas e valores booleanos inválidos durante o parse, antes de leitura de bundles, chamadas de catálogo/factory ou publicação. Conservar stdout vazio, output anterior e ausência de arquivo novo nas rejeições. Reutilizar a política de parse único existente, sem varrer argumentos textualmente ou confundir valores de pergunta com flags.

## Aceitação da formulação em L1

Completar os consumidores do prompt algorithm_insight: permutar referências de fonte/alvo, bundles e ideias deve preservar relatório/IDs sem mutar entrada; mudança válida nos bytes citados deve mudar a identidade causal. Testar limites inclusivos e primeiro valor inválido de campos textuais, steps/condições, fontes/ideias/referências, pedido de 4 MiB, grupos de exploração de 2 MiB e 512 itens, além de citação com até 20 linhas, faixa até 4096 bytes e literal até 2048 bytes. Campos obrigatórios vazios, controles e duplicatas de referências são rejeitados, sem relatório utilizável. Nenhum teste deve usar o resultado de um helper de produção como oráculo de identidade ou aceitação.

Os testes adicionais usam domínio de agendamento fictício, bytes em memória e a API pública FormulateAlgorithmInsights. O núcleo continua sem I/O, sem execução de fontes ou promoção semântica. As relações aceitas permanecem declarações e efeito continua UNVERIFIED.

## Preservação e verificação

Preservar os testes anteriores, contratos de formatos, passos, laudos, recibos e bancadas históricas. O novo prompt de testes registra intenção ex-ante desta ampliação; não atribui essa cobertura ao passo 0155. Executar controles de regressão antes da correção, testes focados e make verify; conferir diff e linhagem. Evidência bruta nova pode ficar em diretório temporário ou em bancada própria, nunca em L0. O gate raiz não reproduz gates selados de módulo próprio.

## Resultado e recorte de verificação

Os controles de duplicatas falharam antes da correção: explore alcançava a leitura do bundle e discover alcançava a factory, em vez de rejeitar o parse. A correção usa BoolFunc com o mesmo registro de unicidade das demais flags escalares. Os testes novos cobrem valores implícitos/true/false, valores conflitantes, um/dois hífens, saída inexistente e saída anterior preservada; controles positivos conservam modos e familiar repetível.

Os [consumidores adicionais da formulação](../../01_core/usecase/algorithm_insight_contract_test.go) passaram pela API pública, com hashes calculados independentemente. Cobrem permutações de referências de fonte/alvo e bundles, causalidade dos bytes nos dois papéis, duplicatas, campos obrigatórios/controles, fronteiras textuais em bytes ASCII/UTF-8, coleções, citações e tetos serializados exatos de 2 MiB e 4 MiB. Não houve alteração produtiva em L1 nem mudança de formatos ou algoritmos de identidade.

Testes focados passaram e make verify concluiu com código zero em lint, auditorias arquiteturais, go test ./..., go vet ./..., build de ./04_wiring e diff-check. O lint registrou No violations found e os 11 controles do parser passaram. O build e os logs foram direcionados a /tmp/bateia-0157-verification; os controles anteriores à correção e os testes focados ficaram em /tmp/bateia-0157-before.log e /tmp/bateia-0157-focused.log. Esses caminhos são temporários locais, não evidência histórica selada. Os corpos dos testes anteriores foram conferidos como preservados; laudos, recibos e bancadas existentes não foram alterados. Gates selados com módulo próprio não foram reproduzidos por esta execução.
