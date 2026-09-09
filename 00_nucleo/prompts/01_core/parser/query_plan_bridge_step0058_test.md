# Oráculos independentes B — passo0058

Autoridade: contrato congelado em `00_nucleo/passos/0058-diagnosticar-perda-do-alvo.md`.
Regime A/B executado sem atestação física de isolamento. B consulta contrato, AGENTS,
tipos de domínio e testes existentes; não consulta candidato, corpus, benchmark ou patch de A.

Antes de executar, congelar este prompt e `01_core/parser/query_plan_bridge_step0058_test.go`
com SHA256. Não modificar expectativas legadas. O integrador executa os gates.

Verificar, pela API pública, que linhas naturais com plano hypothesis integram todos os
termos action/object, sem duplicatas, mantendo Original, Symbols e Constraints. Exigir
exatamente dois recibos0058, Source igual à linha inteira, Kind translation, Status hypothesis,
RuleID composto da versão e regra action/object do plano, Terms exatamente do grupo.
Não gerar novos recibos para negação, ambiguidade, sintaxe, mixed ou symbol; não cruzar
linhas. Preservar recibos0052 existentes byte a byte por RuleID. Verificar determinismo
por repetição/permutação/deduplicação de linhas e independência entre resultados mutáveis.

Fixtures autorais, sem projeto externo: uma unidade Go em caminho neutro contém operações
WatchDirectory, BlendColour e MarshalDirectory, além de rotina alheia. Perguntas portuguesas
monitorar arquivos, interpolar cores e serializar arquivos devem recuperar a operação
correspondente por termos do plano existente, incluindo objeto traduzido. Controle dirigido
com os termos explícitos do plano confirma recuperabilidade. Não ajustar score nem limiar,
nem inferir satisfação funcional a partir de match lexical. Fontes e entrada ficam intactas.
