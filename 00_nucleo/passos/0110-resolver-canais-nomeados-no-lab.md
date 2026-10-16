# 0110 — Resolver canais nomeados no Lab

Capacidade: investigação de propriedades. Responsável: observador exploratório
no Lab, sob Tekt A0; não há promoção para L1–L4.

Contrato ex-ante nos prompts observer, type_resolution e testes espelhados.
Versão v2 checa ASTs em memória sem imports e só atribui tipos em grupos sem
erros. Observações sintáticas conservam autoridade própria. Tipos não comprovam
limites globais. Corpus/recibos v1 preservados; baseline v2 explícita.

Código em lab/bancada/concorrencia/0109; corpus e resultados novos em
[0110](../../lab/bancada/concorrencia/0110/README.md).

Atribuição por tipos validada em exemplos autorais, incluindo aliases,
multiarquivo, sombras, ciclos, capacidade dinâmica e canais direcionais.
Na fonte conc congelada, dependências ausentes impedem atribuição por tipos;
o próximo contrato deve tratar fornecimento rastreável de dependências/build.
Qualidade semântica independente e análise de fluxo continuam pendentes.
