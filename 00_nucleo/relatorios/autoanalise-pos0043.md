# Autoanálise após 0043 — Execução real do Bateia

2026-09-07; commit `9056cd259e40d6360a06e6e21a157aa2e5b9eac7`. Foram executados `idea recommend` sobre a ideia do próprio produto, `analyze` com uma pergunta funcional e `analyze` com o nome exato de uma função. [Comandos, identidade e inventário](../../lab/bancada/history/benchmarks/autoanalise-pos0043/observations.json).

**Conclusão:** a saída conserva melhor a incerteza, mas ainda não fornece evidência suficiente para uma boa análise funcional do próprio projeto. Este diagnóstico é local e dirigido; não mede recall público nem qualidade global.

## O que funcionou

- A recomendação permaneceu INCONCLUSIVE, sem transformar ausência de comprovação em BUILD ou saúde perfeita.
- O relatório distingue conteúdo examinado, trecho exibido, hash, commit-base e análise não realizada. Tempos aparecem como telemetria da ferramenta.
- A decisão contém 24 ações de investigação e distingue falta de detector de coleta incompleta.

Fontes: [decisão canônica](../../lab/bancada/history/benchmarks/autoanalise-pos0043/reuse-decision.json) e [relatório funcional](../../lab/bancada/history/benchmarks/autoanalise-pos0043/analysis.md).

## Onde a qualidade ainda falha

1. **Recuperação por símbolo insuficiente.** `SanitizeAssessmentEvidence` existe em `01_core/usecase/evidence_authority.go`, linha 125. Porém a consulta exata selecionou somente README e go.mod e não encontrou definição. O controle demonstra falha de recuperação neste caso, não inexistência de AST no produto. [Saída do controle](../../lab/bancada/history/benchmarks/autoanalise-pos0043/symbol-control.md).
2. **A pergunta funcional não chegou às implementações centrais.** Foram selecionados 14 arquivos, incluindo `evidence_authority.go`, mas ficaram fora `idea_discovery.go`, `reuse_decision.go` e `investigation_plan.go`. A saída não apresentou definição ancorada em chamadas estáticas. A relevância 26/100 foi justificada por termos do README, incluindo “como” e “sem”; é pontuação pouco discriminatória para esta pergunta.
3. **A jornada principal não reconhece a maior parte da descrição natural do produto.** As dez capacidades ficaram INCONCLUSIVE: nove sem regra versionada aplicável e uma por coleta incompleta. Isso não demonstra falta dessas funções no projeto. Demonstra alcance insuficiente do avaliador para essa formulação.
4. **Orçamento consumido com seleção pouco dirigida.** Foram registrados 43 alvos coletados, 837 ignorados por orçamento e quatro não encontrados. O inventário inclui numerosos tipos de domínio, documentos e fixtures; implementações centrais ficaram fora da coleta. A lista é evidência de cobertura, não revisão dos conteúdos de todos esses arquivos. Simplesmente aumentar o limite pode aumentar ruído.
5. **Plano ainda repetitivo.** Das 24 ações, 14 pedem definir verificação/detector. É uma explicação honesta, mas insuficiente como plano de engenharia: falta apontar qual implementação examinar e qual teste específico responderia à pergunta.
6. **Uma inferência ainda aparece como observação.** A matriz anuncia gRPC; o manifesto contém `google.golang.org/protobuf` como dependência indireta. `polyglot_matrix.go`, linhas 23 e 40–50, mapeia substring de dependência para gRPC. Isso não demonstra uso de gRPC pelo Bateia.
7. **Ruído operacional.** O modo analyze tentou análise Clang num projeto Go e declarou não avaliado. Também houve falha de inicialização do índice Bleve com fallback para cache em disco neste ambiente; não foi diagnosticada a causa do índice.

## Prioridade indicada pela execução

Primeiro, recuperar implementação por símbolos e relações antes de gastar orçamento com documentos/tipos; depois, ligar requisitos naturais a verificações delimitadas sem fabricar equivalência semântica. O plano deve usar as lacunas e fontes realmente encontradas para propor verificações específicas. Classificação de dependências deve permanecer indício, e backends não aplicáveis não devem ocupar o relatório principal.

Nenhuma correção de produção foi feita neste diagnóstico. Os resultados não são benchmark dos candidatos nem certificação independente. Não se consultaram provedores remotos.
