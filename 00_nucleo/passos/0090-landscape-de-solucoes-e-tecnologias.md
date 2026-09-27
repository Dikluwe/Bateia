# Passo 0090 — Landscape de soluções e tecnologias

## Objetivo

Separar similaridade de papel no problema: substituto completo, complemento
integrável, analogia arquitetural ou irresolvido. O modo simples de `search`,
o radar conceitual e a decisão de reúso existente permanecem inalterados.

## Comando

`bateia idea landscape analyze --input landscape.json --output landscape-report.json`

`lab/bancada/corpora/technology-landscape-illustrative.json` demonstra o
contrato, mas contém claims inventados para teste, **não achados reais**.

## Critérios determinísticos

- `SUBSTITUTE`: todas as capacidades da necessidade são `PROVEN` e a operação
  independente é `PROVEN`.
- `COMPLEMENT`: pelo menos uma capacidade `PROVEN`, uma `ABSENT`, nenhuma
  desconhecida e integração `PROVEN`.
- `ANALOGUE`: todas as capacidades `ABSENT` e analogia `PROVEN`.
- `UNRESOLVED`: qualquer outro caso. Não encontrado não é ausência comprovada.

Claims com status conclusivo exigem `sources`, referências opacas a avaliações,
testes ou documentos. O avaliador valida forma e consistência, **não verifica
se a fonte realmente prova o claim**. O operador deve vincular esses IDs a
artefatos de evidência do Bateia antes de usar a saída para decisão.

Sinais suportados: `RELEASES` (`count`), `MEDIAN_RESPONSE_HOURS` (`hours`),
`BUS_FACTOR` (`people`) e `ARCHIVED` (`boolean` 0/1). Os três primeiros exigem
`window_days` positivo; `ARCHIVED`, janela zero. Mudanças só são comparadas
quando tipo, unidade e janela coincidem. A direção é descritiva; maior número
de releases, por exemplo, não implica automaticamente melhor software.

Snapshots anterior e atual exigem a mesma lista de necessidades (ordem pode
variar) e datas `YYYY-MM-DD` crescentes. O relatório mostra candidatos
adicionados/removidos, mudança de papel e movimento de sinais. Não há
pontuação global de maturidade, inferência jurídica de patentes ou anéis
automáticos de adoção.

## Próxima validação

Congelar um corpus de achados reais com rótulos independentes de papel e
comparar o classificador com triagem por similaridade. Medir falsos substitutos,
complementos perdidos, custo de adjudicação e impacto na decisão final.
