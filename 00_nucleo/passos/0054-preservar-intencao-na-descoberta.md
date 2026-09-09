# Passo 0054 — Preservar intenção na descoberta

## Problema observado
O benchmark 0053 registrou fallbacks para “Como criar” e “Como obter”, retirando a linguagem e retornando candidatos sem o alvo esperado. Há a mesma redução em SearchCandidatesWithTrace e Discover, incluindo substituição por um discriminador isolado.

## Decisão e contrato congelado
A busca pública deve manter a pergunta recebida. Este passo desativa o fallback destrutivo; não cria um tradutor semântico nem promete aumento de recall.

1. Ambos os caminhos Sourcegraph enviam uma única consulta: texto original integral, precedido de lang: quando a linguagem foi fornecida e do escopo repo: quando solicitado na porta Discover. Nenhum discriminador substitui o texto. Não gerar lang: vazio. Manter a conversão de nomes de linguagens existente.
2. Consulta vazia ou só espaços não causa chamada ao provedor e retorna erro explícito. Contexto já cancelado não causa chamada e retorna erro. Erro, timeout e cancelamento nunca viram sucesso vazio, inclusive na porta Discover. Uma resposta válida vazia é avaliada como vazia, sem novas tentativas.
3. Manter rastros com a consulta efetivamente enviada, resultados e erros redigidos sem token; na busca direta conservar OriginalQuery. Candidatos seguem sendo indícios, sem promoção de relevância ou alteração da intenção nas etapas seguintes.
4. Manter compatibilidade dos campos públicos. Versionar a evidência para invalidar caches da estratégia anterior. Não alterar ranking local, corpus, runner ou expectativas do benchmark 0053.
5. Testes devem distinguir resultados vazios de falhas, conservar termos tardios/negação/símbolos, linguagem e escopo, resistir a discriminador genérico e verificar ambas as portas por servidor HTTP controlado.

## Execução Tekt
Regime A/B, sem atestação de isolamento físico: autor A implementa produção sem ler testes; autor B escreve testes a partir deste contrato sem ler o patch. Coordenação escreve contrato, integra e executa. Registrar hashes antes da execução. Testes legados que exigem remoção de linguagem são explicitamente substituídos, pois esse comportamento é o defeito corrigido.

## Verificação e limites
Executar controles independentes, suíte Go, linteiro e build. Repetir benchmark 0053 com corpus/runner congelados em diretório 0054. Rede indisponível/timeout é não avaliado. Comparar preservação de consulta e recall separadamente; zero resultados não prova qualidade semântica. Telemetria é tempo do Bateia/provedor, não benchmark dos projetos.

## Próximo limite
Tradução de intenção para consultas de código com termos discriminantes e validação da intenção original permanece pendente. Só reintroduzir relaxamento com contrato próprio e ganho medido em perguntas congeladas.
