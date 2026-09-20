# Execução do passo 0055

## Entrega
Implementado planejador puro e finito de ação + objeto, integrado às duas portas Sourcegraph. A consulta original é tentada primeiro. Uma hipótese derivada só é executada após resposta válida vazia, preservando linguagem, escopo e limites. Consultas negadas, ambíguas, multilinha ou com sintaxe explícita permanecem literais. Catálogo cobre seis famílias de ação e quatro de objeto; não é tradução semântica geral.

Candidatos derivados exigem indícios dos dois conceitos em snippets retidos do mesmo arquivo. Nomes/descrições de repositórios não bastam. Rastros registram regras, termos, consulta original/efetiva, contagens brutas e rejeitadas, arquivos/trechos admitidos e restrições não verificadas. Markdown mostra o recibo; JSON conserva estrutura. A avaliação posterior no funil continua recebendo filter.Query original, não a hipótese derivada.

Plano `0055-v2`, analisador `evidence-v16-planned-discovery-syntax` e identidade de retomada live atualizada impedem reuso da estratégia antiga. Local/replay mantêm suas regras de identidade anteriores. Snippets não comprovam implementação, comportamento ou solução exata.

## Resultado e limites

| Medida | 0054 | 0055 final |
|---|---:|---:|
| Projetos esperados admitidos na descoberta | 0/3 | 0/3 |
| Funções esperadas recuperadas localmente | 1/6 | 1/6 |
| Arquivos esperados extraídos | 5/6 | 5/6 |

Não houve ganho demonstrado de recall. Na pergunta de bits, a consulta derivada retornou um candidato bruto e rejeitou um. As outras duas permaneceram literais: negação em fsnotify e asteriscos na notação L*a*b* acionaram os bloqueios conservadores. Esses bloqueios são limites conhecidos do planejador, não compreensão dessas restrições.

[Métricas e contagens](../../lab/bancada/history/benchmarks/0055/comparison.json), [tentativas finais](../../lab/bancada/history/benchmarks/0055/collected-r2/cases.json), [recibos](../../lab/bancada/history/benchmarks/0055/collected-r2/receipts.json).

**Fato observado:** diagnóstico separado com repositório conhecido encontrou bits-and-blooms/bitset. O diagnóstico global da consulta derivada retornou um FileMatch em periph/devices, caminho ssd1306/image1bit/image1bit.go, com 24 previews; não o alvo esperado. Esses diagnósticos são consultas adicionais e não entram no recall nem substituem respostas coletadas. A coleta final registra rejeição por evidência lexical insuficiente nos snippets retidos; não demonstra ausência da capacidade no repositório. [Disponibilidade](../../lab/bancada/history/benchmarks/0055/availability-baseline.json), [formato do retorno](../../lab/bancada/history/benchmarks/0055/provider-shape-diagnostic.json).

**Inferência:** as traduções lexicais não resolvem por si só seleção/cobertura dos resultados do provedor. A triagem conservadora evita aprovar metadados como evidência, mas também pode descartar arquivos quando previews são insuficientes.

**Hipótese para próxima investigação:** medir separadamente cobertura do índice, seleção dos candidatos brutos e suficiência dos snippets; para candidatos pertinentes sem evidência suficiente, avaliar aquisição limitada de conteúdo antes da admissão. Não adicionar sinônimos apenas para favorecer estas perguntas.

**Não verificado:** generalização, atendimento às restrições originais, comportamento executável e jornada remota integrada. Este corpus é conhecido pelo autor do contrato, portanto calibração; não holdout. As perguntas, alvos, snapshots e todos os arquivos de 0053 mantêm hashes. Telemetria mede Bateia/provedor, não desempenho dos projetos.

## Correção do próprio experimento
A primeira serialização `content:"(a|b)"` não encontrou o alvo no diagnóstico controlado. A forma `patterntype:regexp (a|b) AND (x|y)` encontrou. O contrato foi revisado explicitamente antes da segunda implementação. Termos e expectativas não mudaram; a primeira coleta foi preservada e houve uma rodada da versão corrigida. A [documentação oficial](https://sourcegraph.com/docs/code-search/queries) fundamenta padrões e operadores; o comportamento da emissão concreta foi validado por diagnóstico. [Adendo](../passos/0055-adendo-sintaxe.md), [recibos da sintaxe](../../lab/bancada/history/benchmarks/0055/query-syntax-c.json).

## Verificação e Tekt
A/B executado sem atestação física. A implementou produção sem ler testes/corpus; B escreveu testes sem ler produção/patch/corpus; integração escreveu o contrato e a apresentação/checkpoint e executou os gates. O contrato final foi revisado após observar a falha de sintaxe: não alegamos congelamento final anterior a todas as observações nem selo de refinamento.

Build inicial falhou por acesso duplicado a Results, corrigido. Linteiro inicial rejeitou import slices no teste, trocado por laço equivalente. A primeira suíte encontrou fixture null indevida e bloqueio no cleanup HTTP: preservada e interrompida; fixtures passaram a arrays válidos, null continua erro explícito e cleanup passou a limitado. As expectativas de admissão/erro não foram afrouxadas.

Suíte Go completa final, compilação, linteiro sem violações e nove controles do avaliador passaram. Reavaliação offline reproduziu as métricas. [Verificação/hashes](../../lab/bancada/history/benchmarks/0055/verification.json), [log Go](../../lab/bancada/history/benchmarks/0055/full-tests-final.log), [reprodução](../../lab/bancada/history/benchmarks/0055/README.md).
