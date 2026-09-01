# Contrato L0 — cliente GitHub autenticado

## Credencial e fronteira

O cliente recebe o token somente por injeção efêmera feita a partir da variável de ambiente na composição. Não lê `.env`, não persiste a credencial e não oferece representação textual que a exponha. Envia autenticação exclusivamente ao host HTTPS autorizado da API GitHub, sem redirecionamento que preserve o header para outro host.

Use `Authorization: Bearer <token>` apenas na request em memória. Nunca inclua esse header, cookies ou o token em URL, erro, log, receipt, hash, métrica ou snapshot. Erros de transporte e provedor são normalizados sem corpo/header bruto. O token não integra request hash.

## Request

Aceite somente busca de repositórios públicos, query não vazia e limite exatamente 10 para este experimento. Aplique timeout de 20 segundos através do contexto/cliente. Não consulte endpoints de conteúdo, código ou repositórios privados. O cliente não implementa retry; cada chamada corresponde a uma tentativa, e o caso de uso decide uma eventual repetição contabilizada.

Desabilite cookies e autenticação implícita. Redirecionamentos devem ser recusados ou reconstruídos sem credencial e somente para destino previamente autorizado; qualquer ambiguidade falha fechada. A URL sanitizada do receipt contém apenas endpoint base, sem query.

## Resposta e privacidade

Classifique:

- HTTP 200 com JSON válido e zero itens: `EMPTY`;
- HTTP 200 válido com itens comprovadamente públicos: `SUCCESS`;
- HTTP 403/429 relacionado a cota: `UNKNOWN_RATE_LIMIT`;
- deadline: `UNKNOWN_TIMEOUT`;
- transporte, status, JSON ou campos não classificáveis: `UNKNOWN_PROVIDER_RESPONSE`.

Decodifique a propriedade de privacidade/visibilidade fornecida pela API. Item `private=true` ou visibilidade não pública causa `REJECTED_PRIVATE`; nenhum item da resposta é devolvido para persistência. Campo de visibilidade ausente, contraditório ou não suportado causa `UNKNOWN_PROVIDER_RESPONSE`, também sem hits. Não deduza publicidade apenas pela URL.

Preserve apenas rank, URL de repositório, URL canônica, nome/título e descrição. Não colete código, README, owners além do que já compõe o nome público, ou campos desnecessários.

## Observabilidade redigida

Capture status HTTP e, quando disponíveis, apenas valores parseados dos contadores de rate limit permitidos pelo domínio. Não preserve nomes/valores brutos de headers. Receipt contém endpoint base sanitizado, duração e contadores redigidos. Request e response são selados apenas depois da normalização e da verificação de privacidade.

Uma varredura antes do retorno rejeita o token exato, `Bearer`, header de
autorização, `Set-Cookie`, `cookie=` ou derivado reconstruível. Palavras comuns
como `authorization` e `cookie` em metadados públicos não são segredo por si só.
O erro de rejeição é constante e não ecoa o achado.
