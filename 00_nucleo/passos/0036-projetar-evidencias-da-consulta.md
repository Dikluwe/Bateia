# Passo 0036 — Projetar evidências da consulta

## Problema
A coleta Clang do passo 0035 encontra implementações, mas a projeção Markdown perde chamadas e corpos relevantes. Vetores ausentes são confundidos com divergência. Recibos repetidos dominam a saída.

## Implementação prevista
1. Realimentar a extração com unidades completas selecionadas pelo Clang; conferir hashes antes de usar seus corpos.
2. Selecionar definições por termos no corpo/nome e expandir chamadas diretas, com limites determinísticos. Registrar intervalos integrais examinados, trechos exibidos, motivo e limites. Não declarar solução comprovada.
3. Não preencher seleção lexical específica com arquivos sem sinal. Separar score de descoberta da evidência estática. Similaridade com dimensões ausentes será indeterminada, sem domínio fabricado.
4. Resumir recibos no Markdown e preservar dados integrais em JSON lateral na CLI. Deduplicar evidências apresentadas sem apagar contextos distintos do JSON.
5. Corrigir reconhecimento de definições enum com enumeradores.

## Aceitação
- Coin: renderizador extraído, cadeia de composição exibida, cobertura parcial explícita; nenhuma alegação de execução gráfica.
- Caso sintético sem nomes Coin: corpo pertinente e chamadas recuperados; hash divergente não produz evidência válida; arquivos irrelevantes não completam a cota.
- Vetores desconhecidos não viram Divergent ou score geral numérico.
- Markdown limitado e recibo JSON preservado; lint, testes Go e build aprovados.

## Limites
Seleção continua heurística lexical e estática; chamadas não demonstram comportamento em execução. Não há certificação independente nem comprovação universal de qualidade semântica.

## Conclusão da execução
Implementado e validado. Evidências e comandos em [benchmark/0036](../benchmark/0036/README.md). As seis funções centrais do Coin e suas chamadas chegam ao Markdown; 12/898 unidades continuam explicitamente parciais. Documento passou de 476.063 para 38.140 bytes. Lint sem violações, suíte Go completa e build aprovados. A melhoria demonstrada é de seleção, rastreabilidade e projeção; adequação funcional continua não verificada.
