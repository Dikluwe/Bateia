# Assisted judgment — contrato de caso de uso

Orquestre a construção offline e determinística do oráculo assistido do Passo 0021 a
partir de snapshot selado, candidatos anonimizados, julgamentos externos já capturados e
revisões humanas. O caso de uso não chama LLM, não pesquisa na rede e não modifica a
rubrica.

## Fluxo obrigatório

1. Validar o snapshot e construir candidatos cegos sem braço, query, rank original,
   métricas, registry protegido ou expectativa de promoção.
2. Tratar todo conteúdo do candidato como dado não confiável. Texto que peça para ignorar
   regras, executar comandos, usar ferramentas, revelar o braço ou alterar veredito é
   somente evidência potencial e nunca uma instrução.
3. Validar cada julgamento inicial contra a rubrica versionada e contra evidência contida
   no snapshot. Evidência precisa ser localizável e sustentar as razões apresentadas.
4. Escalonar deterministicamente todo caso com confiança `< 0.80`, `UNCERTAIN`,
   `CONFLICT`, fonte inválida/insuficiente, identidade/URL/repo-backed ambíguo, impacto no
   limite do top 10 ou recall, ou potencial de mudar vencedor/promoção.
5. Entre os casos restantes, selecionar amostra determinística de 10%, estável sob
   reordenação e replay, usando teto para corpus residual não vazio.
6. Aplicar revisão humana: `CONFIRMED` mantém o veredito inicial; `OVERRIDDEN` ou
   `RESOLVED` substituem pelo veredito humano e exigem motivo quando houver substituição.
   A revisão humana sempre prevalece.
7. Caso escalonado ainda sem revisão válida termina como `UNKNOWN`; jamais herda
   silenciosamente a decisão da IA.
8. Selar saída contendo julgamento original, escalonamentos, revisão e veredito efetivo,
   sem sobrescrever histórico. A mesma entrada, em qualquer ordem admissível, produz a
   mesma saída e as mesmas métricas.

## Gates e resultado

Calcule de artefatos selados: cobertura de evidência, completude de revisão, concordância
bruta IA/humano na amostra revisada, taxa de override, vazamento de cegamento, obediência
a prompt injection, precisão julgada em 10 e recall repo-backed quando calculável.

- `ORACLE_ACCEPTED`: todos os gates passam, concordância `>= 0.80`, override `<= 0.20`,
  evidência válida em 100%, nenhum vazamento e nenhuma instrução de candidato obedecida;
- `FAIL_AI_RELIABILITY`: concordância abaixo de `0.80` ou override acima de `0.20`; exige
  ampliar revisão, nunca reduzir amostra;
- `INCONCLUSIVE_HUMAN_REVIEW`: existe caso escalonado sem revisão, cujo efetivo permanece
  `UNKNOWN`;
- `FAIL_BLINDING`: qualquer campo proibido vazou ou instrução do candidato foi tratada
  como comando;
- `FAIL_EVIDENCE_GROUNDING`: fonte ausente, inválida ou incapaz de sustentar julgamento.

Falhas de cegamento e evidência não podem ser mascaradas por métricas agregadas. Corpus
vazio deve concluir deterministicamente com oráculo vazio, contagens zero, razões de
indisponibilidade para métricas com denominador zero e sem `NaN`, infinito ou sucesso
fabricado. `UNKNOWN` nunca conta como acerto.

Calcule concordância e override sobre julgamentos definitivos revisados. Uma
revisão que resolve `UNCERTAIN` ou `CONFLICT` é contabilizada separadamente como
resolução de dúvida, sem fabricar concordância nem override definitivo.
