# Jornada fria — somente XERJ

Pergunta: How can XERJ corpus evidence provenance verification distinguish retrieval relevance from tested behavior, expose applicability conditions and support reproducible counterexamples without confusing signatures or feedback with correctness?

O percurso produziu três propostas de composição, apoiadas em contratos e recibos do próprio alvo e em um corpus externo descoberto pelo Hub. São hipóteses fundamentadas, não correções de defeitos estabelecidos. Não executei código dos projetos nem os experimentos propostos.

## Origem e isolamento

Comecei só com a pergunta, abordagens familiares e URL alvo https://github.com/xerj-org/xerj. A URL alvo é origem **target**, não descoberta de XERJ. Nenhuma referência externa foi fornecida pronta; não usei referências da memória (**model: nenhuma**). Não consultei Bateia, histórico, pesquisa web ou arquivos do clone diretamente. Conteúdos de pesquisa foram obtidos somente por comandos XERJ no runner auditado. Shell/Python prepararam entradas, consolidaram resultados e encerraram meu próprio nó.

O corpus **admission-compare** é origem **tool**: apareceu em sexto lugar na busca nativa da pergunta, audit/006.json. Sua identidade recuperada é corpus + pin 7a439382d306ac8599b9cf613e967bf54c86d0f2, data 2026-10-24, Apache-2.0. O resultado não trouxe URL upstream; não inventei owner ou endereço. Seus caminhos r2r-jev/... identificam o conteúdo do pack. Os demais cinco candidatos receberam só metadados, sem leitura de conteúdo: kkp-connected-spectral-counterexamples, agentic-security-ietf-corpus, xerj-search, xerj-storage e xerj-vector. Seus nomes não são evidência de aplicabilidade.

## Caminho completo observado

1. Li ajuda geral, hub/find e ask (001–003). A primeira tentativa hub find falhou (004, exit 1): nenhum catálogo local e o comando não o busca automaticamente.
2. Executei hub sync (005): catálogo gerado em 2026-11-05T04:20:12Z, 514 corpora, nenhum pack prévio. A busca repetida (006) retornou seis candidatos; o primeiro era sobre contraexemplos espectrais e não foi tratado como relevante só por sua posição. Selecionei admission-compare por sua descrição de experimento de admissão versus threshold, não por ranking superior.
3. Adquiri o alvo exato via wrapper (007, 69,774 s). Aquisição não conta como descoberta do produto. O runner não expôs commit do clone; não recorri a Git/leitura direta para obtê-lo. Isso limita a reprodutibilidade dos localizadores do alvo, embora os resultados auditados preservem os trechos observados.
4. Consultei autoindex/search help (008–009) e iniciei nó próprio com dados novos, loopback 19591 e modo lexical (010). Sem índices anteriores. A primeira pergunta externa tentou hosted com --no-fetch (011, exit 1): este build não possui cliente hosted; o fallback não podia baixar pack com --no-fetch.
5. Repeti ask com hosted off, pack permitido e URL do nó (012). Download de 740.654 bytes restaurou 4 índices e 142 definições. Obtive cinco passagens com pin/licença, incluindo história versus autoridade, determinismo e limites. Outra pergunta ao mesmo pack (015) recuperou condições de admissão, contraexemplos e trecho de src/admission.rs.
6. Uma busca no prefixo alvo ainda vazio (013) retornou zero hits. É efeito da sequência, não ausência de conteúdo do projeto.
7. Indexei o alvo com --no-semantic --no-graph --workers 2 --yes (014). O produto mediu um piso de extração cliente de 6,6–6,8 s e explicitou que não era previsão do tempo total. O comando completo levou 112,235 s. Terminou exit 3, **ok: indexed and searchable**, geração 1, 33 datasets e 258.201 registros vivos. Não houve falha geral de indexação.
8. Buscas 016–025 ocorreram parcialmente enquanto o índice ainda se construía. 016 recuperou PinnedCitation e ruído de páginas de respostas; 017 por nome de módulo mostrou cabeçalhos; 018 misturou feedback útil e material não pertinente; 019 JSON retornou só declaração da struct; 020 encontrou recibo CJK útil e ruído de issues; 021 por EvidenceManifest não encontrou símbolo e avisou que o autoindex ainda não terminara; 022 encontrou recibo de startup e documentação da diferença entre ocorrência de frase e verificação; 023 foi ruído para a questão de assinatura; 024 recuperou from_bytes com hashes/offsets; 025 por correctness foi dominada por issues Kibana. Não converti misses, ruído ou orçamento esgotado em ausência de capacidades.
9. Encerrei exclusivamente meu nó por argumentos exatos de diretório e PID 949494, SIGTERM. O runner de longa duração retornou exit 0; receipt em node-shutdown.json e audit/010.json. Não publiquei nem enviei feedback.

Foram **24 chamadas de produto**, incluindo cinco ajudas e o nó de longa duração, mais uma aquisição. O arquivo journey.json preserva todos os argumentos, códigos e durações. As durações se sobrepõem: não devem ser somadas como tempo de jornada. Não medi tokens. O timestamp em journey.json mediu aproximadamente 387 s desde a primeira chamada até a consolidação JSON, excluindo esta edição final do relatório.

## Cobertura e trabalho manual

A indexação registrou 13 arquivos e 2 registros junk/skipped. engine/crates/xerj-api/src/es_compat.rs (2.252.994 bytes) e engine/crates/xerj-engine/src/index.rs (3.120.037 bytes) excederam o limite fixo de extração de código de 2 MiB. Arquivos binários e dois archives de fixtures também não foram abertos. Regras ignoraram 206 arquivos e podaram 7 diretórios, incluindo 26 *.md por .gitignore e arquivos ocultos. Uma linha informou conteúdo cortado em um arquivo; a saída auditada foi limitada pelo runner, e não tenho seu caminho completo. Datasets sem campo textual foram explicitamente sinalizados como consultáveis por estruturas, mas não por passage search. Portanto o corpus pesquisado não equivale a leitura integral do repositório.

Eu escolhi perguntas e termos, interpretei passagens, distingui contexto de ruído, comparei recibos e redigi as propostas manualmente. Não há um resultado único do produto que já seja a síntese abaixo. O Hub descobriu uma referência utilizável, e o produto recuperou trechos com localizadores; seleção, hipóteses, condições e desenho dos experimentos foram trabalho do agente. Não usei corpus/search de pares do target porque o orçamento foi gasto no caminho Hub + autoindex + search.

## Evidências recuperadas

- **T1 — engine/crates/xerj-common/src/hub_evidence.rs:60 e :98**, audits 016/019/024: PinnedCitation carrega corpus, repository, commit, arquivo, hash do arquivo, offsets de bytes originais, linhas e hash do span. from_bytes verifica intervalo não vazio dentro do arquivo e calcula hashes/linhas. Isso demonstra contrato de citação, não que todos os consumidores façam uma verificação completa ou que uma asserção semântica esteja correta.
- **T2 — tools/cjk-judge/VERIFICATION.md:167**, audit 020: recibo diferencia controle de comportamento de performance com builds comparáveis; fixa snapshots e comando. Scores de uma fixture lexical não são probabilidades calibradas. Afirma checks e gate anteriores; eu recuperei esse relato, não reexecutei os testes.
- **T3 — tools/startup-diagnostics/VERIFICATION.md:70**, audit 022: recibos locais não verificam release assets ou CI hosted; exclusões e condições não testadas ficam explícitas.
- **T4 — content/answers/verify-ai-report-claims-against-sources.md:23**, audit 022: match_phrase testa ocorrência de palavras no campo; miss se refere à ausência ali. A linguagem de “verificar” precisa conservar esse escopo.
- **T5 — engine/crates/xerj-autoindex/src/feedback.rs:36**, audit 018: proveniência de relato de agente e números somente se medidos. É uma base de atribuição, não validação automática do relato.
- **A — admission-compare@7a439382...**, audits 012/015: r2r-jev/docs/architecture.md:1/63/135/203 separa julgamento, admissão, governança e enforcement; Hold não altera relações, e determinismo depende de envelope/política/tick. README.md:120 traz casos de confiança alta rejeitada por fonte fraca ou expiração e corroboradores; src/admission.rs:148 mostra ramos correspondentes. docs/evidence-admission-semantics-v0.1.md:305 lista limites, incluindo calibração e conluio. Os limiares do demo não são transferidos para XERJ.

## Três propostas

### I1 — Separar ocorrência recuperada, integridade e comportamento testado na mesma alegação

**Fontes da observação:** T1, T2, T4, T5

**O que já existe observado:** PinnedCitation, hashes/offsets e recibos de verificação com limites explícitos já foram observados; não estabeleci que a composição proposta esteja ausente.

**Interpretação:** Uma citação íntegra torna a observação auditável; ranking lexical e frase encontrada não demonstram comportamento; feedback é testemunho atribuído.

**Hipótese:** Expor facetas independentes retrieval, integrity, execution e behavior, ligadas ao claim, reduz promoções indevidas sem descartar fontes úteis.

**Proposta:** Claim carrega consulta/campo/rank; citação pinada; resultado criptográfico; teste e asserção; status observed/not-run/pass/fail/inconclusive. Assinatura e utilidade ficam em campos separados e nunca promovem behavior.

**Condições:** O executor só registra pass quando a asserção específica foi observada sob condições registradas; ausência de execução permanece not-run.

**Contraexemplo:** Pack assinado com texto relevante e feedback 5/5 descreve algoritmo que falha na fixture; integrity pode passar enquanto behavior falha.

**Pequeno experimento de rejeição:** Construir quatro recibos sintéticos cruzando citação íntegra/adulterada e asserção passa/falha; adicionar relevância alta e feedback positivo em todos. Rejeitar o desenho se qualquer assinatura, rank ou feedback alterar o resultado de comportamento, ou se adulteração não invalidar integridade.


### I2 — Expor aplicabilidade como predicados e desconhecidos antes de reutilizar evidência

**Fontes da observação:** T2, T3, A

**O que já existe observado:** Há condições/exclusões em documentação e Hold no corpus externo; não alego ausência de política equivalente no produto.

**Interpretação:** Os recibos do próprio alvo já delimitam build, fixture, ambiente e o que não foi testado; admissão externa exemplifica estado não autoritativo quando falta suporte.

**Hipótese:** Uma ficha de condições ligada ao claim evita extrapolar um recibo de fixture para produção ou release.

**Proposta:** Extrair manualmente predicados explícitos: build/revisão, configuração lexical/neural, corpus/fixture, plataforma, operação e validade. Comparar contexto solicitado com supported/contradicted/unknown; preservar a fonte de cada condição e não converter score de retrieval em confiança.

**Condições:** Predicados devem vir de recibos ou contrato, sem inferir limiares universais do demo de admissão. Falta de metadados dá unknown; requerimentos críticos contraditos impedem aplicação automática.

**Contraexemplo:** Teste lexical de três documentos com restart não sustenta probabilidade calibrada, performance matched-build, release distribuída ou outro sistema operacional.

**Pequeno experimento de rejeição:** Criar seis pedidos variando só uma condição conhecida: mesma fixture, outro build, outro SO, modo neural, exigência de calibração e release publicado. Revisão humana cega marca o que o recibo suporta. Rejeitar se o sistema disser supported em condição explicitamente não testada, esconder unknown, ou deixar de preservar suporte no caso idêntico.


### I3 — Anexar cápsulas pinadas de contraexemplo e replay à alegação

**Fontes da observação:** T1, T2, A

**O que já existe observado:** Alvo publica snapshots/comandos; corpus externo tem casos executáveis e baseline; não verifiquei que já exista formato unificado.

**Interpretação:** Hashes localizam bytes; replay precisa ainda de entradas, política, ambiente, asserção e saídas. Determinismo do demo é limitado ao envelope/política/tick, não prova geral.

**Hipótese:** Um pequeno pacote negativo reproduzível detecta sobre-generalização de relevância e evidência com menos ambiguidade.

**Proposta:** Cápsula imutável com claim, citações, fixture/hash, versão de política, ambiente/build, comandos, pré-condições, oracle esperado e resultado observado, incluindo failed/not-reproduced/inconclusive. Guardar variantes negativas e ligar falhas à alegação original sem apagar o histórico.

**Condições:** Replay apenas em ambiente controlado autorizado e sem segredos; processo reprodutível não é garantia de oracle correto. Plataformas ou dependências não disponíveis produzem inconclusive.

**Contraexemplo:** Mesmos bytes assinados e mesmo score de relevância podem levar a resultado diferente sob configuração neural ou regra de expiração diferente; não é corrupção da citação.

**Pequeno experimento de rejeição:** Preparar duas execuções idênticas e uma variante que altera um único predicado (evidência expirada ou modo de busca). Verificar que idênticas mantêm hashes/resultados e que a variante preserva citação mas muda applicability/resultado esperado. Rejeitar se replay não identifica dependência, se falha some no histórico ou se ferramenta chama não reproduzido de confirmado.

## Desconhecidos e limites finais

Não estabeleci a ausência de formato unificado de evidência, predicados de aplicabilidade ou replay no XERJ. O nome EvidenceManifest foi um palpite de consulta, não uma entidade exigida pelo contrato. Não examinei integralmente hub_evidence.rs nem seus testes, não verifiquei public keys ou assinatura por um teste de adulteração e não executei code dos projetos. A assinatura verificada no caminho do pack, os hashes das citações, a utilidade atribuída e a correção de uma asserção pertencem a observações diferentes; a separação proposta é inferência de engenharia sobre as fontes recuperadas.

A referência externa demonstra um mecanismo análogo e limites explícitos, não prova eficácia em XERJ. Os três experimentos são propostas não executadas. O orçamento chegou a 24 chamadas; isso encerrou busca adicional, não resolveu desconhecidos. Resultados relevantes coexistiram com ruído, arquivos ignorados, extração limitada e consultas durante indexação. Qualquer comparação com outro produto deve incluir esses custos de descoberta, preparação, seleção manual e cobertura — não apenas o tempo das últimas buscas.
