# 0123 — Laya real e plano de treinamento

Laya foi instalado e executado em ambiente isolado /tmp/bateia-laya-0123.
Fonte 3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63, pacote 0.4.0; checkpoint inglês
convaiinnovations/laya na revisão 7b928d828b7b0e022f929d9bd2e44165aa270148.
Cinco arquivos, 846.195.574 bytes, hashes registrados; safetensors conferido
contra SHA256 LFS publicado pelo Hub. Nenhum arquivo de peso no repositório.

Python 3.11.15, torch 2.14.1+cpu, transformers 5.19.0, CPU float32 com oito
threads, AMD Ryzen 7 5800H. Instalação por uv em venv temporária; não mudou
Python global. Inferência SDK offline após download, sem servidor HTTP.

## Diagnósticos conhecidos

| Caso | Rótulo autoral | Laya na ordem original |
| --- | --- | --- |
| tunny durante processamento | supports | supports |
| conc depois do primeiro Go | contradicts | undetermined |
| pond interface sobre zero | contradicts | contradicts |
| pond README sobre zero | supports | supports |

Três concordâncias em quatro referências autorais conhecidas, sem estimar
precisão independente. Seis permutações das três opções e repetição da primeira
ordem por caso: 28 execuções correlacionadas, não 28 exemplos. Todas as escolhas
ficaram estáveis entre ordens; respostas repetidas foram idênticas.

Nenhum texto-fonte foi cortado; opções permaneceram distintas. Auditoria do
tokenizador também confirmou instruções e descrições inteiras: 43 tokens de
instrução, opções 14/15/21, orçamento de instrução 139, head_max_len=192 e
max_len=512. Não se aplicou limiar de confiança. Inicialização ~5,7 s, mediana
exploratória SDK quente ~0,87 s; não é medição HTTP nem previsão de produção.

O carregador avisou de temperatura inválida em choice:11+, ajustada para 0,5.
Os pedidos usam três opções, bucket choice:3-5 com temperatura 1,76015; essa
temperatura também não é calibração no domínio Bateia. Log preservado.

A divergência conc requer revisar referência e contexto: a declaração fala
em família With*, a alegação nomeia WithMaxGoroutines. Isso é hipótese de
dificuldade de escopo, não causa isolada demonstrada. Os quatro casos não
incluem referência undetermined: qualidade dessa saída não foi medida.

## Como treinar

training-protocol.json define o plano, não um treino executado. Unidade é uma
alegação com condições e trecho documental, com três rótulos exclusivos.
Não treinar um juiz geral de código ou de segurança concorrente.

Meta inicial de planejamento: 600 exemplos em 12 famílias de projetos novos,
seis famílias para treino, duas desenvolvimento, duas calibração, duas teste.
Nenhum desses exemplos foi coletado/rotulado ainda. Família inclui forks,
versões e paráfrases; todas ficam no mesmo split, atribuído antes de aumentar
dados. Projetos conhecidos conc/tunny/pond ficam fora do teste final.

Dois revisores sem acesso às previsões classificam cada exemplo. Discordância
não resolvida exige adjudicação e não vira automaticamente undetermined.
Texto insuficiente é classe semântica; revisão incompleta é estado do dataset.
Rótulo de modelo não é gabarito. Proveniência/licença e escopo de autoridade
ficam em sidecar, ligados por example_id.

Treinador nativo Laya aceita JSONL state/questions/expected. Manter a pergunta
congelada, sem instrução gerada por classe. Começar comparando base e regra
simples, depois avaliar ajuste com encoder congelado. Não escolher parâmetros
por estes quatro diagnósticos. Usar desenvolvimento para seleção; desativar
split aleatório interno de calibração (calib_frac=0.0), fornecer eval_data de
desenvolvimento e ajustar temperaturas/política separadamente no grupo reservado.
Validar implementação desse caminho antes de qualquer execução de treino.

Medir especialmente falso suporte, matriz de confusão, abstenção/cobertura,
calibração e fatias de negação, condições, inicialização, família de métodos e
texto irrelevante/injetado. Calibrar não é treinar; testar não é escolher limiar.
Teste final permanece reservado. Português exigirá experimento com checkpoint
adequado; a rodada atual usou inglês técnico, não validou português.

## Reprodução

Ambiente/pesos em /tmp precisam ser recriados se removidos. Pacotes estão em
results/python-packages.txt; hashes em checkpoint.json e verification.json.
Dentro desse ambiente:

```sh
HF_HOME=/tmp/bateia-laya-0123/hf-cache HF_HUB_OFFLINE=1 TRANSFORMERS_OFFLINE=1 OMP_NUM_THREADS=8 /tmp/bateia-laya-0123/venv/bin/python lab/bancada/decisoes/0123/run_diagnostics.py
HF_HOME=/tmp/bateia-laya-0123/hf-cache HF_HUB_OFFLINE=1 TRANSFORMERS_OFFLINE=1 /tmp/bateia-laya-0123/venv/bin/python lab/bancada/decisoes/0123/audit_question.py
python3 lab/bancada/decisoes/0123/verify.py
```

Scripts substituem apenas os resultados desta bancada. Não reescrevem 0122.
Resultados brutos em predictions.jsonl; síntese diagnostics.json; verificação
diagnostic-verification.json. Modelo executado; fontes Go de conc/tunny/pond
não foram executadas nesta rodada. Nenhuma integração produtiva ou treino.
