# Passo 0037 — Resposta fundamentada por mecanismo

## Contrato
Abrir a apresentação com uma síntese estática citada, antes da telemetria. Priorizar uma implementação ancorada na consulta e seus caminhos diretos. Cada observação deve apontar para símbolo, arquivo, intervalo e contexto de compilação; não inventar intenção nem resultados de execução.

## Alterações
- Selecionar a projeção em uma única unidade de tradução completa, com hashes confirmados; preservar demais unidades no JSON. Não combinar chamadas de configurações diferentes.
- Ranking favorece termos próximos e definições menores; comentários são indícios, nunca comprovação. Proximidade com chamada resolvida é necessária para âncoras apenas textuais.
- Exibir múltiplos trechos exatos: âncora, condições/laços textuais e chamadas; rotular contexto textual sem alegar análise de fluxo de controle. Limites explícitos por função.
- Resposta inicial mostra entrada candidata, assinatura, chamadas e trechos, com referência ao contexto. Ausência de prova funcional permanece explícita.
- Corrigir corpo “examinado integralmente”, commit-base vazio e similaridade agregada desconhecida.

## Validação
Coin: preservar mecanismo por camadas e expor seleção, suporte, laço e composição. Caso sintético: comentários distantes de chamadas não dominam; unidades incompatíveis não formam cadeia; trechos correspondem aos bytes. Similaridade desconhecida nunca é exibida como 0%. Executar suíte Go, lint e build.

## Limites
Síntese determinística de fatos estáticos; não interpretação universal de algoritmos. Trechos condicionais não provam caminho executável. Código excluído pelo pré-processador não produz chamadas AST; conteúdo textual isolado não é prova. Não há veredito independente.

## Execução concluída
Implementado e validado. [Resultados e limites](../../lab/bancada/history/benchmarks/0037/README.md): condições, fallback, laço e composição presentes na abertura; 23 definições e 39 chamadas em contexto único. Testes Go, Clang real, lint e build aprovados.
