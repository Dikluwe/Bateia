# Laudo 0123 — Teste real do Laya e plano de treino

Laya no commit 3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63 e pesos ingleses na revisão
7b928d828b7b0e022f929d9bd2e44165aa270148 foram instalados em venv/cache temporários
isolados. Cinco arquivos de checkpoint (~846 MB) com hashes; safetensors contra
LFS publicado. Pacotes/ambiente registrados. Inferência offline por SDK em CPU,
float32, oito threads, sem alteração do ambiente global ou produto.

Quatro pedidos preparados em 0122 geraram três concordâncias com rótulos autorais.
tunny foi supports, conc undetermined, interface pond contradicts e README pond
supports. conc divergiu do rótulo contradicts. Referências conhecidas/autorais,
sem acurácia semântica independente ou avaliação representativa.

Vinte e oito execuções correspondem a quatro casos com seis ordens de opções
e uma repetição. Escolhas estáveis entre ordens e respostas repetidas idênticas.
Texto, instruções e opções couberam sem corte/collapse. Não foi definido limiar
de confiança. Mediana SDK quente ~0,87 s e carga ~5,7 s neste ambiente: medidas
exploratórias, não serviço HTTP. Um aviso de temperatura clamped em choice:11+
foi preservado; a rodada usa choice:3-5 e não fez calibração no domínio.

A discrepância conc envolve declaração de família With* e alegação de método
específico. Precisamos revisar gabarito, associação/condições e modelo; não se
demonstrou causalmente que ajustar pesos resolve. Nenhum gold undetermined no
piloto, portanto nem o comportamento dessa classe foi avaliado adequadamente.

Plano de treino: relação alegação/declaração com condições, três categorias,
novos exemplos rastreáveis e revisão independente das previsões. Meta inicial
de planejamento de 600 exemplos/12 famílias; dataset ainda não produzido.
Split por família antes de versões/paráfrases, com grupos separados de treino,
desenvolvimento, calibração e teste. Fontes conhecidas excluídas do teste final.

Comparar modelo base/regra e primeiro ajuste com encoder congelado; decidir
descongelamento pelo desenvolvimento. Treinador nativo tem split aleatório de
calibração: plano exige calib_frac=0.0 no train já particionado e ajuste posterior
no grupo reservado. Usar dev para seleção, nunca teste final. Requisitos desse
caminho precisam de verificação antes de executar treino. Nenhum treino ou
limiar foi produzido nesta fatia. Confiança não muda autoridade das fontes.

Lint, testes raiz sem cache, vet, build e diff check passaram. Assinatura Git
desativada apenas no processo de fixtures; módulos selados não reproduzidos.
Produto/observadores preservados, sem imports do Lab. Execução real foi do Laya,
não dos candidatos Go. Banco de dados de produção, serviço HTTP e integração
de classificação não foram materializados.

[Bancada/reprodução](../../lab/bancada/decisoes/0123/README.md),
[resultados](../../lab/bancada/decisoes/0123/results/diagnostics.json),
[plano estruturado](../../lab/bancada/decisoes/0123/training-protocol.json)
e [verificação](../../lab/bancada/decisoes/0123/results/verification.json).
