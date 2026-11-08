# 0154 — Descoberta nativa corrigida, ensaio sem fontes de memória

**A descoberta nativa do Bateia passou a fornecer fontes úteis neste problema sem URLs ou projetos lembrados pelo agente.** As quatro aquisições têm testemunho anterior em retorno do produto. O material retornado pelo Bateia foi mais próximo de implementação de inventário/cliente offline que o obtido pelo agente XERJ nesta execução. O parecer anônimo encontrou famílias operacionais compartilhadas e uma diferença complementar de merge colaborativo; nenhuma solução completa exclusiva ou família inteiramente nova foi demonstrada.

O teste usa o mesmo cenário concreto de 0152/0153: oito oficinas, 45 mil peças individuais, 60 voluntários, Android simples e até 72 horas sem rede, com fotos, histórico, exportação e dois mantenedores. É um cenário hipotético; não foi implantado nem se executou projeto-fonte descoberto.

## Correção realizada

O diagnóstico encontrou busca excessivamente ampla em README, perda de contexto nas rotas, variedade antes de correspondência e resolução de revisão para todos os hits. A [documentação GitHub](https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories#search-by-repository-name-description-or-contents-of-the-readme-file) distingue a busca em metadados da inclusão explícita de README; isso também foi confrontado com respostas públicas da API preservadas em `diagnosis`.

| Ponto | Política contextual nova |
|---|---|
| Consultas | Cinco conjunções de termos da pergunta, sem acrescentar palavras genéricas como alternative/tradeoffs/failure e sem cadastrar projetos esperados |
| Escopo GitHub | Nome, descrição e tópicos; README é lido depois da seleção, sem participar da busca ampla de catálogo |
| Seleção | Quantidade de termos literais nos metadados antes de variedade; zero match não é escolhido; stars não determinam o ranking |
| Revisão | Busca conserva os metadados e adia resolução; as referências escolhidas recebem pin imutável pela porta existente antes de árvore/conteúdo |
| Compatibilidade | CLI fornece `context-v2`; `--query-policy lexical-v1` permite repetir a política anterior. Pedidos L1 sem política e campos opcionais omitidos preservam o contrato/hashes históricos |

Regras e estado continuam em L1, transporte em L3, flags/publicação em L2. A política chega à ponte GitHub e à seleção/memória multifonte. Não há modelo, tradução automática, interpretação semântica de relevância ou glossário de soluções embutido. Consultas podem precisar de reformulação por capacidades; parte do catálogo continua ruidosa.

O [patch dos componentes](product-change.patch), as cópias passivas anteriores em `before`, os [fingerprints](source-fingerprints.json), os prompts proprietários e o [passo 0154](../../../../00_nucleo/passos/0154-corrigir-descoberta-nativa-e-testar-sem-memoria.md) permitem revisar o corte. Os testes novos usam fontes fictícias e domínio de backup, fora do cenário de estoque: contexto preservado, sem vocabulário inventado, correspondência antes de popularidade/variedade, flags corretas ao catálogo, pin apenas dos escolhidos e identidade divergente sem leitura. O gate completo [passou](verification/verification.json): lint, arquitetura, testes, vet, build e diff. Seu binário coincide com o SHA-256 congelado para o ensaio.

## Como excluímos fontes de memória

Agentes novos receberam a mesma [pergunta](problem.txt), sem histórico, referências ou mecanismos sugeridos. Escreveram baseline antes de qualquer chamada. Mesmos limites: 20 minutos, 40 pesquisas, seis aquisições; setup separado. Ambos registraram opções consideradas/descartadas e lacunas. Não podiam buscar nomes de projetos lembrados, inventar URLs nem usar web ou busca genérica fora do produto. Reformular por problemas/capacidades era permitido.

O [runner](run.py) bloqueou aquisição com origem `model`, URL não observada e coleta de identidade não elegível. O registro de elegibilidade é derivado de catálogo/retornos e links dos conteúdos lidos, excluindo planos e perguntas ecoados. Saídas foram conservadas por chamada em snapshots antes de poderem ser sobrescritas. Cada aquisição grava seus testemunhos. A [conferência independente de proveniência](native-provenance-check.json) verificou URL no retorno e ordem descoberta → aquisição para os quatro casos. Os três corpora XERJ têm ordem descoberta → leitura [conferida](xerj-native-provenance-check.json); não foram inventadas URLs de GitHub a partir dos nomes.

A ajuda do produto e os guias de interface eram instrumentos permitidos. No XERJ, catálogo novo, feedback desativado e diretórios próprios; status global de journals foi excluído para evitar o vazamento de nomes observado em 0153. Nós próprios foram encerrados e o endpoint foi testado. Os pesquisadores não receberam o smoke nem resultados anteriores. Nenhum ajuste de produto foi feito a partir de seus resultados durante esta rodada.

## Diferença efetiva do retorno

| Braço | Material encontrado e lido nativamente | Pertinência observada |
|---|---|---|
| Bateia | Nove repositórios com conteúdo retornado, quatro adquiridos para aprofundar: inventário Android, outbox/sync, merge, reserva local, rede isolada e analogia de journal | Componentes e código/testes mais próximos da implementação pretendida, com condições que impedem adotá-los intactos |
| XERJ | Três corpora nativos úteis: caniuse, Cassandra e BookKeeper | Capacidade de armazenamento/inicialização web, limites de LWW/relógios e semântica de log/escritor único; não um cliente/inventário pronto |

Nove repositórios lidos não equivalem a nove fontes pertinentes: o viewer incidental não foi usado como desenho, e várias referências serviram para limitar/rejeitar adoção. A diferença de adequação é julgamento do avaliador sobre passagens e mecanismos, não um score semântico calculado pelo produto. A descoberta por URL de memória foi zero em ambos.

### Exemplos de decisões alteradas pelas fontes Bateia

- **synchronize_cache:** documentação de outbox e integração backend; a recomendação de cache de deduplicação de 24 horas não basta como contrato para reenvio após 72 horas ou replay de backup. O modo de atualização incondicional exige validação de domínio ao integrar saldo/reserva. Isso virou condição e caso de rejeição do piloto; não foi reproduzido como bug do projeto.
- **datum:** last-write-wins padrão versus estruturas colaborativas explícitas, com teste lido de identidades de réplica distintas. A opção final limita merge a observações/estruturas apropriadas e mantém disputa de disponibilidade separada. Não se atribui exclusividade de reserva à convergência.
- **benchstock:** arquitetura usa serialização de escritores com `BEGIN IMMEDIATE`; um teste de duas reservas espera uma reserva e um conflito. É evidência estática de banco único, não de réplicas desconectadas. O projeto se declara protótipo, trabalha por SKU e tem lacunas de histórico/exportação. A proposta exige adaptação e um protocolo de custódia entre oficinas ainda não encontrado pronto.
- **hot-stuff e inventory:** a primeira referência explicita requisito de Android 13 e risco de perda local; outra mostra rede isolada, o que não equivale a qualquer telefone operar desconectado dela. São restrições para implantação, não capacidades gerais inferidas pelo nome.

O XERJ também trouxe restrições úteis: Cassandra documenta a resolução por última escrita dependente de relógios; caniuse só sustenta investigar capacidades de cliente web; BookKeeper sustenta a distinção entre log imutável e escritor ativo único. Transferir essas regras para trabalho físico e Android permanece hipótese/síntese. Não se recomendou implantar Cassandra ou BookKeeper para a cooperativa.

## Alternativas e parecer anônimo

Um terceiro agente novo recebeu problema, rubrica prévia e pacotes X/Y anonimizados, sem mapping, workspace ou consulta externa. X=Bateia e Y=XERJ; o mapa só foi usado na síntese principal. O anonimato é editorial e não garante cegamento perfeito. O [parecer](assessment.md) e seu JSON preservam notas, pertinência por fonte e lacunas.

| Decisão | Bateia | XERJ | Parecer |
|---|---|---|---|
| Registrar offline e confirmar compromisso depois | A1 | A1 | Família compartilhada; Bateia concretiza outbox/validação central, XERJ deixa a unicidade dessa confirmação menos especificada |
| Custódia local e transferência com aceite | A3 | A2 | Família compartilhada; fonte Bateia apoia reserva no mesmo banco/rede, fonte XERJ oferece analogia de escritor único |
| Merge colaborativo seletivo com identidade de réplica | A2 complementar | Baseline contempla convergência; proposta final não tem apoio específico equivalente | Exclusiva **parcial** observada no resultado Bateia, relevante sob condições e confiança moderada; não solução ponta a ponta |

As cinco propostas atingiram o limiar da rubrica. Notas relevância/evidência/especificidade: Bateia A1=4/3/2, A2=3/3/2, A3=4/3/2; XERJ A1=4/3/2, A2=3/3/2. Notas não foram somadas; evidência 3 indica fonte específica de mecanismo/limite, não software integrado testado. Quantidade de títulos, chamadas ou bibliotecas não declara vencedor.

Nenhuma família nova foi demonstrada: diário, réplicas e custódia já estavam nos baselines. A contribuição da pesquisa foi tornar concretos e restringir os contratos de deduplicação/replay, merge/identidade, reserva local versus remota e adoção/exportação. A família da exclusiva parcial de merge já era imaginada pelos dois agentes; a pesquisa Bateia deu a ela suporte e limites mais específicos nesta execução. Não se atribui a descoberta da ideia geral à ferramenta.

**O objetivo de eliminar a dependência de fontes lembradas foi atingido neste caso.** Já a pergunta sobre exclusividade de soluções completas continua com resposta negativa: os dois mantiveram famílias operacionais centrais equivalentes. A diferença demonstrada está na proximidade do material de implementação e em restrições concretas recuperadas, não numa superioridade geral ou solução inédita.

## Auditoria e limites

| Registro | Bateia | XERJ |
|---|---:|---:|
| Chamadas de pesquisa | 7 | 27 |
| Setup | 0 | 9 |
| Aquisições com origem nativa comprovada | 4 | 0 (leitura dos corpora do Hub) |
| Aquisições de memória | 0 | 0 |
| Trechos distintos citados nas alternativas presentes no retorno registrado | 15/15 | 5/5 |

As contagens descrevem percursos diferentes e não estimam tokens, eficiência ou custo. HTTP de aprofundamento do Bateia teve 403 e exigiu transporte de URLs já descobertas para leitura pelo produto local. O XERJ teve ruído, limites de tamanho de corpus e ausência de cliente hospedado nesta build; ainda recuperou material nativo útil. Esses limites ficam na cobertura dos braços, sem resgate por memória.

A [conferência de evidência](evidence-check.json) verificou literalidade e os hashes dos baselines em todos os audits. Treze dos 15 trechos Bateia também foram conferidos em blobs Git dos clones fixados; os dois restantes são do mesmo README, conferido por HTTP no pin e com bytes idênticos aos retornados. Todos os localizadores Bateia cabem nas linhas declaradas. Os cinco trechos XERJ foram conferidos nos retornos; pins de corpus são os reportados pelo produto, sem afirmar verificação Git independente. Literalidade não comprova a interpretação ou comportamento executado.

Relatórios originais: [Bateia](bateia/report.md), [alternativas](bateia/alternatives.json), [fontes](bateia/sources.json), [cobertura](bateia/coverage.json); [XERJ](xerj/report.md), [alternativas](xerj/alternatives.json), [fontes](xerj/sources.json), [cobertura](xerj/coverage.json). Queries, snapshots e audits estão nas respectivas pastas. Clones, corpora baixados e binários não foram incorporados ao repositório. `setup.json` registra seus hashes; `metrics.json` registra contagens/pins e `comparison.json` a conclusão estruturada.

Houve um smoke inicial no cwd do workspace que falhou; a CLI carrega `.env` nesse cwd. O probe sem credenciais e a execução em diretório limpo passaram. Não inspecionamos nem copiamos os valores do `.env`; os ensaios usam diretórios limpos e ambiente sem tokens. Essa diferença de contexto foi eliminada antes dos pesquisadores. Os resultados preliminares permanecem em `diagnosis`, sem tratá-los como rodada de comparação.

É um problema, uma execução por braço e baselines curtos. A proibição de memória também muda as condições frente a 0153; não se estima causalidade de criatividade por comparação simples de números. Ausência de uma fonte não prova incapacidade de encontrá-la. Nenhum projeto-fonte ou piloto foi executado; não se mediram Android real, 45 mil peças/fotos, exportação/restauração completa, manutenção, custos ou redução de viagens. As fontes orientam um piloto e permitem rejeitar suposições, não aprovar adoção.
