# 0119 — Executar referência original do candidato

Intenção ex-ante: avaliar as alegações comportamentais do passo 0118 com testes
originais do candidato, escritos fora do Bateia e congelados antes da revisão.
Responsável: bancada de concorrência no Lab, subordinada a Tekt A0. Produto e
observador não mudam; resultados e scripts ficam no Lab, síntese em L0.

Reconstituir módulo temporário exclusivamente dos bytes coletados em 0106,
validando hashes e commit. Não substituir asserts nem produzir testes que
espelhem nossos julgamentos. Executar TestPool e ExamplePool com detector de
race e repetição, registrando eventos, ambiente e dependências efetivas.

Comparar configuração antes de Go, rejeição depois de Go, rejeição após Go/Wait
e reutilização com as expectativas originais. Demais alegações não recebem
promoção: os testes selecionados não decidem suporte de configuração concorrente,
atomicidade universal, defeito por transferência de modelo ou Wait antes de Go.

Independência aqui é a origem dos testes, não revisão humana independente de
toda a interpretação. Não calcular precisão do classificador usando esses casos
selecionados. Módulo temporário não constitui checkout completo certificado;
dependências resolvidas e versão Go precisam constar do recibo.

Resultado: testes originais passaram; quatro comportamentos concordam com a
revisão anterior e quatro alegações não foram avaliadas pelos testes selecionados.
[Laudo](../laudos/0119-referencia-comportamental-do-candidato.md) e
[bancada](../../lab/bancada/concorrencia/0119/README.md).
