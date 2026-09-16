# Contrato canônico 0068 — compilação técnica finita

Autoridade normativa:
`../../passos/0068-compilar-capacidade-em-consultas-tecnicas.md`.

## Fragmento observável

Para a frase portuguesa canônica, o expansor preserva a expansão literal e acrescenta,
na ordem, a tradução inglesa exata, `orient2d incircle` e `Shewchuk exact predicates`.
As três propostas usam `technical-lexicon-v1`, origem/rule ID rastreáveis e pares de
discriminadores definidos no passo. Nenhuma entrada parcial ou fora do catálogo ativa a
regra. As regras anteriores mantêm saída byte a byte sob o mesmo limite explícito.

O limite produtivo injetado é três. `DiscoveryRequest` recebe clone da derivação. Requests
técnicos válidos executam uma única query com linguagem, escopo e `count` limitado; não
usam o planejador 0055 nem fallback. Somente evidência retida do mesmo arquivo contendo
os dois discriminadores admite candidato. Todas as propostas independentes são
executadas pelo plano e a união conserva as origens de queries distintas.

## Baseline e métricas

Os hashes dos dossiês congelam o diagnóstico real: v1 trouxe Whisk na redução genérica;
v2 não trouxe candidato para o requisito robusto; o controle dirigido `Shewchuk` trouxe
quatro candidatos brutos, três extraídos, incluindo Carve com evidência textual de exact
predicates. Isso não prevê o resultado atual da rede.

No replay contrafactual C01, baseline: recall@5 0, contamination@5 1,
evidence-addressable 0, uma tentativa. Candidato: recall@5 1, contamination@5 0.5,
evidence-addressable 1, quatro tentativas, RobustKernel com duas origens técnicas. C02
controla split-file/metadata; C03 controla frase parcial e entrada fora do domínio; C04
controla acentos, caixa, ordem, deduplicação e compatibilidade legada.

`OBSERVED` somente se todas as melhorias e guardas do passo forem atendidas em pelo menos
duas famílias. Undefined/Unknown/erro/budget/replay ausente é `INCONCLUSIVE`; outro vetor
válido é `NOT_OBSERVED`. “Evidence-addressable” não é “capacidade comprovada”.

## Limites e não regressão

Máximo: três propostas/origem, 64 expansões, dois discriminadores técnicos de até 64
bytes, query técnica de até 256 bytes, quatro tentativas por origem no experimento; todos
os limites 0056 continuam. Nenhum LLM/I/O em L1, mudança de ranking/veredito, edição de
dossiê ou alteração de 0065–0067. Oráculos eliminam M01–M16 com score 1.0 e confirmam
20/20 permutações, custo e hashes.
