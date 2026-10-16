# Pesquisa 0122 — Uso de modelos de decisão no Bateia

Data: 2026-10-07. Recomendação: começar com Laya local por HTTP no Lab, com uma
pergunta sobre permissão declarada de reconfiguração. O Bateia já prepara fonte,
localizador, dono e evidências. O modelo acrescentaria uma interpretação de
relação entre alegação e declaração; não valida comportamento concorrente.

## Fontes e alcance

O coletor do Bateia consultou commits fixados de Laya, Arbiter e Decision Index:
53 itens com hashes verificados. Laya/Arbiter ficaram parciais por orçamento;
Decision Index respondeu HTTP 403. Seu README foi recuperado no mesmo SHA por
HTTP público, com proveniência separada. Não se corrigiu complete para true.
Nenhum código remoto, modelo ou benchmark foi executado. Interpretação autoral.

O projeto pesquisado é [NandhaKishorM/laya](https://github.com/NandhaKishorM/laya),
modelo de decisões; não o aplicativo de notificações de laya.aay.sh.

## Escolhas de integração

| Componente | Papel possível | Avaliação necessária |
| --- | --- | --- |
| Laya | Classificador local ajustável, saídas choice/score/noul | Especialização na tarefa, calibração e cobertura de tokens |
| Servidor nativo laya-serve | HTTP para o Bateia em Go, POST /v1/systemone | Identidade de pesos, resposta e limites no commit escolhido |
| Arbiter | Serviço alternativo para modelos tipados | Só acrescentar se batching/operação justificar outro componente |
| Clef / Clef-flash | Comparadores de decisão, inclusive multimodais | Medir custo/qualidade; não presumir equivalência de confidence |
| Jev | Comparador hospedado com versão fixada | Serviço/chave e avaliação no mesmo corpus |
| Decision Index | Referência de protocolo e avaliação geral | Não substitui corpus de contratos de software |

Laya documenta execução em CPU, API HTTP, configuração explícita de orçamento
de tokens e abstenção opcional. O limiar usa answer_confidence; confidence tem
outra semântica. Os modelos exigem calibração no domínio. A versão declarada
no pyproject coletado é 0.4.0; isso não confirma disponibilidade de wheel ou
compatibilidade da instalação nesta máquina.
[Fontes no commit](https://github.com/NandhaKishorM/laya/tree/3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63).

Arbiter é camada de serviço, não modelo novo.
[README fixado](https://github.com/0xBakeer/arbiter/blob/bf71358774c85aeff2051ec7b1dcc4577a87706d/README.md).
Clef é anunciado como modelo multimodal de 27B e Clef-flash de 9B, com interface
de perguntas tipadas. Essa comparação é documental, sem ranking local.
[Clef](https://developers.cloudflare.com/workers-ai/models/clef/),
[Clef-flash](https://developers.cloudflare.com/workers-ai/models/clef-flash/).
Jev publica modelos/limites próprios; escolher versão explícita em vez de alias
mutável é parte da reprodução. [Modelos Jev](https://docs.typesafe.ai/models).

## Piloto preparado

Há quatro requisições diagnósticas prontas, sem inferência: permissão de ajuste
em tunny; proibição após Go em conc; declaração da interface pond sobre zero;
declaração do README pond sobre zero. As respostas esperadas são autorais e os
casos são conhecidos; não são holdout, treino ou medida independente de qualidade.

A pergunta avalia cada declaração isoladamente em relação a uma alegação, com
três categorias: supports, contradicts, undetermined. Uma interpretação sobre
texto não se torna garantia de execução. Probabilidades e texto-fonte permanecem
separados; o modelo não precisa gerar citações, porque o Bateia já as possui.

Fluxo proposto: vista documental 0121 → preparação curta da pergunta → serviço
local Laya → resposta registrada → política de abstenção/revisão. Em eventual
promoção, contrato/estado em L1, cliente HTTP em L3, apresentação em L2 e composição
em L4. O processo do modelo é recurso concreto, não regra de domínio.

## Antes de medir ganho

1. Fixar software, revisão/tokenizador dos pesos, calibração, dispositivo/dtype
   e schema da pergunta. Alias english sozinho não identifica os pesos usados.
2. Instalar em ambiente isolado, começar com um checkpoint e servidor loopback.
   Laya/torch/transformers não estão instalados no Python consultado; nada foi instalado.
3. Executar os quatro diagnósticos para validar transporte, orçamento e registro.
   Falha/truncamento não deve virar undetermined decidido pelo modelo: é erro operacional.
4. Criar corpus novo rotulado separadamente, com divisão por projeto para treino,
   desenvolvimento, calibração e teste. Incluir condições, negações, conflito,
   evidência ausente e instruções dentro do texto-fonte.
5. Comparar regra simples, Laya base/ajustado e outro candidato com entrada e
   política iguais. Medir matriz de confusão, falsos supports, calibração,
   cobertura/revisão, latência completa e consumo. Escolher limiar na calibração,
   congelá-lo antes do teste; não importar o número 0,85 de um exemplo da documentação.

O Decision Index ajuda a conferir a superfície de perguntas tipadas, mas possui
suíte pública com tarefas e termos próprios. Sua nota não mede permissão de
reconfiguração. [Kit no commit](https://github.com/apolinario/decision-index/tree/9eb2dbe2a358004c8782c66e40a83ac07b953fec).

## Verificação

Hashes de 53 itens e fallback conferidos; quatro payloads preparados. Lint,
testes raiz, vet, build e diff check passaram. Produto/observadores não alterados;
gates selados não reproduzidos. Testes de fixtures desativaram assinatura Git
somente no processo. Não houve teste de inferência, treino ou calibração.

[Bancada](../../lab/bancada/decisoes/0122/README.md),
[payloads](../../lab/bancada/decisoes/0122/pilot-requests.json) e
[recibos](../../lab/bancada/decisoes/0122/results/collection-review.json).
