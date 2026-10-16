# Laudo 0120 — Avaliação externa de configuração

Observador v8 congelado por hashes de fontes/binário; referências/protocolo
registrados antes das respostas. Dois projetos externos (tunny e pond), seis
perguntas, três perfis por projeto. Seis entradas geram dezoito linhas de
comparação, sem contabilizar variantes como amostras independentes. conc é
calibração fora dos denominadores. Candidatos externos não foram executados.

Os perfis completos passaram na checagem de tipos. A chamada New → SetSize de
tunny foi localizada (uma referência estrutural). As três perguntas sobre
contrato não receberam veredito: capacidade não implementada. As duas perguntas
de atomicidade não foram promovidas a prova. Nenhuma guarda do padrão restrito
foi emitida nos projetos externos, sem inferir ausência de proteção. Na
calibração, conc manteve uma guarda e atomicity=UNRESOLVED.

Quatro perfis incompletos produziram diagnósticos de tipos e suprimiram relações
de chamada tipadas. Replays foram byte-idênticos. Localizadores/hash e cinco
controles de contabilidade/promoção do avaliador passaram. O status emitido é
OBSERVED_SYNTAX, que não representa conclusão contratual. O leitor do avaliador
corrigiu nomes de campos from/to e essa interpretação do status, sem mudar
referências, observador ou pontuação esperada.

A referência publicada distingue políticas: tunny e pond apresentam ajustes
durante uso, enquanto conc impede reconfiguração após Go no recorte anterior.
Em pond, README admite zero, interface diz maior que zero e implementação trata
zero. Esse conflito de declarações foi identificado na revisão manual; não é
detecção automática nem defeito de execução demonstrado.

A avaliação confirma uma relação estrutural e delimita cobertura contratual
zero neste recorte. Não mede precisão/recall semânticos: interpretações autorais,
seleção intencional, apenas dois projetos e nenhuma revisão independente. Docs
do README ficam na referência, fora da entrada suportada pelo observador.
Abstenção não foi automaticamente pontuada como sucesso.

Produto e observador preservados. Lint, testes raiz sem cache, vet, build e
diff check passaram; fixtures usaram assinatura Git desativada apenas no
processo. Módulos selados não reproduzidos, sem import produtivo do Lab.

Próximo contrato recomendado: observação de declarações de API, ligada a fonte,
escopo e autoridade, antes de classificar políticas. Estes projetos agora são
diagnóstico conhecido; uma futura medição de generalização exige novos casos
reservados e revisão semântica independente.

[Referência](../../lab/bancada/concorrencia/0120/testdata/reference.json),
[avaliação](../../lab/bancada/concorrencia/0120/results/evaluation.json),
[verificação](../../lab/bancada/concorrencia/0120/results/verification.json)
e [reprodução/limites](../../lab/bancada/concorrencia/0120/README.md).
