# 0106 Corrigir esperas e investigar GitHub

**Data**: 2026-10-06  
**Estado**: VERIFICADO NO RECORTE  
**Base**: passos 0104/0105 e contratos de container 0079/gate ativo R5

## Contrato

Preservar a exigência de SUCCEEDED e os vereditos dos 60 mutantes do gate ativo.
A espera de lifecycle é watchdog diagnóstico, sem SLA de dois segundos. M28
precisa observar conclusão antes da abertura de recovery ou rendezvous no
caminho que espera essa abertura; timeout não determina o veredito.
As correções pertencem aos testes L3 e à bancada ativa, com prompts prévios.
Não alterar runtime produtivo nem recibos históricos.

Coletar um repositório público real pelo adaptador GitHub existente e registrar
SHA imutável, plano, bytes, hashes, limites e falhas no Lab. Escolher
sourcegraph/conc para investigar limitação de concorrência e backpressure no
Bateia. Congelar o bundle e avaliações explícitas, selecionando testemunhas de
implementação e teste vinculadas a linhas reais. A coleta online não substitui
a avaliação semântica; a inspeção do código de um teste não significa executá-lo.

Executar insights com essa fonte e repetir offline para conferir identidade e
bytes da saída. Testar adulteração de um item e exigir recusa, preservando a
saída anterior. Distinguir limites opcionais, bloqueio da submissão e
rejeição/cancelamento; não preencher lacunas do catálogo com promessa do autor.
Nenhum efeito sobre o alvo está demonstrado; continuar com Target UNKNOWN.

## Verificação

Repetir os testes corrigidos, rodar a suíte inteira com paralelismo padrão,
lint cristalino, vet e build. Registrar resultado e limites em laudo próprio;
evidência bruta e instruções de reprodução pertencem a lab/bancada/insights/0106.

## Execução

O [laudo 0106](../laudos/0106-correcoes-e-ensaio-github.md) registra as correções,
a coleta real, os controles e a suíte concorrente aprovada. Aplicabilidade e
efeito no alvo permanecem pendentes.
