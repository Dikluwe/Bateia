# 0120 — Avaliar configuração em fontes externas

Intenção ex-ante: congelar observador v8 e avaliar configuração/inicialização
em projetos que não participaram de seu desenvolvimento. Capacidade de
investigação no Lab sob Tekt A0; L0 contém contrato/síntese, bancada contém
fontes e medições. Não corrigir o observador após conhecer o resultado.

O observador emite fatos estruturais, não vereditos de contrato de uso. Medir
separadamente presença de evidência, cobertura de perguntas contratuais,
abstenção e extrapolação. Ausência de configuration_guards não é ausência de
proteção. Relação de guarda não comprova atomicidade. Não usar o status global
INCONCLUSIVE como classificação de cada operação.

Selecionar dois projetos externos de pools com configuração/inicialização
explícitas: Jeffail/tunny e alitto/pond. Congelar commits, arquivos e referências
publicadas antes de executar o observador. Seleção intencional, sem alegação de
representatividade. Usar recortes completos e com dependências ausentes para
distinguir evidência disponível de incompletude; variantes não são amostras
independentes. conc fica apenas como calibração, fora dos denominadores externos.

A referência combina contratos do upstream e interpretação autoral anterior
às respostas. Não há revisão semântica independente por outra pessoa/agente.
Conservar essa limitação explicitamente. Não calcular precisão geral de um
classificador contratual inexistente, nem considerar toda abstenção um acerto.

Verificação: identidades/hashes de fontes e observador, replay determinístico,
localizadores de evidência e separação de métricas; lint, testes raiz, vet e
build. Produto/observador preservados, sem import produtivo do Lab.

Resultado: uma referência de chamada localizada; três perguntas contratuais
sem veredito; nenhuma promoção de atomicidade nos dois projetos. Status real
OBSERVED_SYNTAX, sem semântica contratual; INCONCLUSIVE era apenas o estado
inicial do relatório. Referência anterior às respostas, porém autoral, não
independente. [Laudo](../laudos/0120-avaliacao-externa-de-configuracao.md) e
[bancada](../../lab/bancada/concorrencia/0120/README.md).
