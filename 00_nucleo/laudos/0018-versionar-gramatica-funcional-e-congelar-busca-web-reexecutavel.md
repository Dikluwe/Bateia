# Laudo 0018 — Gramática funcional e busca reexecutável

**Status**: `EM VALIDAÇÃO`
**Veredito**: `INCONCLUSIVE_EXTERNAL_RUN_PENDING`

## Correção realizada

A execução anterior foi invalidada porque usava consultas placeholder,
`web-mock`, repositórios fictícios e denominador competitivo fixo. Os snapshots,
julgamentos e percentuais sintéticos foram removidos; a promoção anterior não é
evidência do produto.

A implementação refeita gera 12 consultas F2 a partir dos seis frames reais,
seleciona regras por `concept_id`, registra morfologia e hashes, falha fechado
para verbo ou slot desconhecido, injeta o provedor por L4 e rejeita adulteração
de snapshots. Anúncios ficam fora do orgânico e recall usa somente
`expected_id` explícito; sem denominador, permanece `null`.

## Evidência offline

- N contém 12 rótulos nominais verdadeiros;
- F1 reproduz o artefato congelado do Passo 0017;
- F2 contém 12 queries da gramática v2.1.0;
- 12/12 queries F2 passaram a auditoria gramatical;
- mutation score do fragmento auditado: 1,0;
- suíte completa: verde.

## Pendência externa

A execução N/F1/F2 contra a API pública foi bloqueada antes da rede porque o
processo poderia carregar `GITHUB_TOKEN` do `.env` e enviar 36 consultas
derivadas do projeto. Não houve tentativa de contorno. Portanto não existem
novos números de precisão ou recall, e a gramática não foi promovida.

Para concluir, é necessária autorização explícita para enviar as 36 consultas.
A opção mais segura é executá-las sem token; eventual rate-limit será preservado
como `UNKNOWN`.

## Limitação Tekt

Contrato, ataques e recibos foram materializados, mas a mesma autoridade
conversacional executou os papéis. O resultado é executado sem atestação de
isolamento e cobre somente os observáveis registrados.
