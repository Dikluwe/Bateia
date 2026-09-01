# Descoberta regional em idiomas naturais

## Intenção congelada

O Bateia deve ampliar a descoberta de alternativas além do vocabulário anglófono. Uma
mesma ideia pode ser consultada em idiomas naturais e mercados diferentes para revelar
aplicativos locais cujas soluções, fluxos ou práticas possam melhorar produtos de outras
regiões.

Idioma natural e região são dimensões de descoberta. Não são linguagem de programação,
não provam adequação funcional e não autorizam tradução silenciosa de evidência.

## Contrato observável

1. A execução aceita zero ou mais locales BCP 47 por `--search-locale`, repetível. A
   ausência preserva exatamente o comportamento anterior. O campo `--language` continua
   sendo exclusivamente o filtro opcional de linguagem de programação.
2. A primeira materialização reconhece exatamente estes locales canônicos:
   `pt-BR`, `en-US`, `es-ES`, `fr-FR`, `de-DE`, `ja-JP`, `ko-KR`, `ru-RU`, `ar-SA` e
   `hi-IN`. A entrada `all` representa esse conjunto. Locales desconhecidos são rejeitados;
   repetição e ordem da entrada não alteram o plano canônico.
3. Consultas localizadas só podem nascer de um catálogo versionado e auditável. Cada
   consulta preserva locale, regra, consulta-fonte e estratégia de derivação. Tradução
   livre, implícita ou sem origem é inválida.
4. A primeira prova temática é `organização de receitas culinárias`. Para cada locale
   solicitado, o catálogo produz uma consulta nativa equivalente a "gerenciador de
   receitas", sem transliterar scripts japonês, coreano, russo, árabe ou hindi para ASCII.
5. O plano registra locales solicitados, realizados e lacunas. Se um conceito não possui
   tradução controlada para um locale solicitado, o estado permanece `Unknown` com motivo;
   a execução não finge cobertura.
6. A requisição ao provedor, o snapshot, a tentativa efetiva e a origem de cada candidato
   preservam o locale. Quando o mesmo repositório aparece em vários mercados, ele é
   deduplicado por identidade e mantém todas as origens regionais em ordem determinística.
7. O relatório de alternativas expõe os locales de descoberta em JSON e Markdown. Locale
   é contexto de procedência, nunca bônus de ranking, evidência de capacidade, licença ou
   manutenção.
8. Reordenação de locales, propostas e respostas do provedor produz os mesmos artefatos de
   domínio. IDs localizados incluem o locale para impedir colisões entre traduções iguais.
9. No máximo dez locales podem ser solicitados. Consultas localizadas contam normalmente
   para orçamento e limites; não existe busca regional oculta.
10. Replay e retomada vinculam os locales ao hash da execução. Um checkpoint criado para
    outro conjunto de locales não pode ser reutilizado.

## Política de opacidade

- Locale inválido, tradução sem regra, origem sem consulta-fonte ou catálogo inconsistente:
  erro de contrato.
- Conceito válido sem tradução catalogada: `Unknown` explícito no plano, não erro e não
  sucesso.
- Texto localizado encontrado: apenas descoberta. A avaliação funcional continua exigindo
  evidências rastreáveis do repositório e pode permanecer `INCONCLUSIVE`.

## Limite da alegação

Este contrato demonstra descoberta regional controlada para o tema congelado e a
infraestrutura extensível de locales. Ele não afirma tradução geral, cobertura de lojas
proprietárias, qualidade equivalente entre idiomas nem compreensão automática de toda
documentação multilíngue.
