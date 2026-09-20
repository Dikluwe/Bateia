# Validação do passo 0037

## Resultado
O relatório abre com entrada candidata, assinatura, contexto único de compilação, trechos e chamadas citados. No Coin, a abertura contém seleção `SORTED_LAYERS_BLEND`, verificação de profundidade/suporte, fallback, laço de passagens e seleção FP/NV. As seis definições centrais continuam presentes.

A projeção contém 23 definições e 39 chamadas, exclusivamente da unidade `SoGLRenderAction.cpp`. As demais unidades analisadas permanecem no JSON integral. As âncoras genéricas de glue e 3DS não compõem a resposta inicial. Cada trecho foi conferido contra os bytes extraídos; cada chamada projetada existe no contexto escolhido.

Markdown original: 46.364 bytes (38.140 no passo anterior); abertura com evidências: 11.557 bytes. Aumento deliberado para mostrar condições, laço e fallback. Não é benchmark do Coin nem prova de qualidade universal. Cobertura Clang permanece 12/898 unidades.

## Verificação
- Suíte Go completa aprovada, incluindo novos testes de contexto e apresentação.
- Teste com Clang real: implementação inteira excluída por `#if 0` não sustenta resposta.
- Caso sintético: comentário pertinente distante das chamadas não promove função irrelevante; unidades incompatíveis não completam uma cadeia artificial.
- Similaridade agregada sem medidas retorna -1 e aparece “Não avaliado”; commit ausente é explicitamente indisponível.
- Lint sem violações; build aprovado; `git diff --check` sem achados.

## Arquivos
`coin-analysis.md`, `coin-analysis.json.gz`, `validation.json`, logs de testes/lint/CLI e `check_projection.py`. Executar o script com os caminhos Markdown e JSON/gzip para repetir as verificações de apresentação.

Comando de geração: `./bateia analyze -path /tmp/bateia-0035-coin -query 'order independent transparency opengl' -lang cpp -depth 1 -format markdown -output /tmp/bateia-0037-coin.md`.

O corpus externo e sua configuração são os mesmos do recibo `../0035/coin-inputs.json`; Coin não é dependência do Bateia. Snapshot sem `.git`: o commit-base original está nesse recibo externo. Somente o link lateral foi atualizado no Markdown arquivado.

## Limites
Esta é uma síntese estrutural determinística, não uma explicação universal da semântica do algoritmo. Proximidade lexical é uma heurística e ainda pode aceitar comentários enganosos próximos de chamadas. Condições são trechos textuais; não há análise de dominância, prova de caminho executável ou classificação completa de ramos desativados dentro de funções ativas. A escolha de uma unidade evita combinar contextos incompatíveis, mas pode perder mecanismos distribuídos entre unidades. Não houve execução gráfica ou avaliação independente.
