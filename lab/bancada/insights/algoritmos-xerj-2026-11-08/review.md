# Parecer independente — algoritmos XERJ

**Utilizável com condições; nenhuma melhoria demonstrada.** O ensaio produz uma hipótese concreta de formato/acesso e uma confirmação útil que elimina uma recomendação duplicada. São duas decisões de natureza diferente, não duas melhorias novas.

| Item | Evidência (0–3) | Transferência (0–3) | Decisão (0–2) | Uso no email |
|---|---:|---:|---:|---|
| I1 — payload por blocos/faixas | 3 | 2 | 2 | Hipótese condicional de pesquisa |
| I2 — poda por limite superior | 3 | 3 | 2 | Confirmação ALREADY_PRESENT |

As notas são separadas, sem soma. Os 15 vínculos constam como válidos em Git e retornos no evidence-check. Isso sustenta integridade dos recortes; formulate não reavalia semântica nem demonstra ganho. Consultei apenas os três documentos autorizados e os clones permitidos xerj e postings-rust. postings-go não existe: o diretório contém xerj, postings-rust e postings-java; não li postings-java. Não consultei histórico, outros ensaios, fontes web ou código produtivo Bateia; não executei fontes nem implementei algoritmos.

## I1 — candidato condicional, não migração aprovada

**Fonte → princípio/passos:** raw.rs da fonte materializa diretórios residentes e metadados de blocos com base, último doc, offset e comprimento; o leitor lê faixas e verifica checksums quando presentes. A abertura também carrega metadados por documento: “diretórios pequenos” deve ser hipótese relativa ao payload, não promessa de custo constante. Separar endereçamento do payload e ler apenas os termos necessários é uma estratégia operacional, não mero slogan.

**Condições → alvo:** o escritor XERJ em index.rs:1887–1932 já produz offsets por termo, blocos de codec e envelope ZPS2 integral. O leitor em :2690–2738 lê/descomprime o envelope e mantém bytes para slices; formatos legados incluem raw/mmap. Portanto, gravar SkipEntry sozinho não resolve acesso aleatório comprimido. open_stats_only já evita postings para estatísticas. Essas capacidades existentes foram consideradas e não equivalem à alternativa de blocos comprimidos independentes.

A revisão encontrou uma condição adicional relevante no clone: xerj-engine/src/index.rs:11704–11726 e :33420–33459 têm **fts_reader_cache/fts_reader_for**, com chave segmento+conjunto de campos, reutilização de Arc no hit e abertura com orçamento no miss. O loop FTS chama essa função em :27061–27065. Não rastreei toda rota HTTP nem medi taxa de hit/miss, orçamento efetivo ou frequência de reabertura. Os comentários contêm números históricos; não são resultados atuais desta revisão.

O pacote é prudente ao declarar cache/lifecycle desconhecidos e não atribuir descompressão por consulta ao searcher. Porém o email deve tornar o cache existente explícito: a oportunidade seria abertura fria, miss, recusa de retenção ou pressão de memória, a confirmar por perfil. Não se deve comparar contra reabrir artificialmente em toda consulta e apresentar isso como baseline de produção.

**Decisão → teste:** vale decidir primeiro pelo perfil; só então justificar protótipo de framing/diretório/leitor, com baseline/fallback atual. Custos de versão, writer/merge, codecs, posições, frequências, legado, integridade, I/O e compressão estão reconhecidos. A API atual expõe iteradores de slices, portanto integração não é uma troca local demonstrada. Contagem exata exige visitar correspondências pertinentes à consulta, não ler indiscriminadamente todos os termos do campo.

O experimento proposto é rejeitável: paridade semântica, legado/corrupção e comparação de bytes, latência, memória, espaço e ingestão em braços frio/aquecido. Acrescentar distinção entre cache de leitor e cache de páginas, frequência real de misses e memória transitória. Os limiares de 20% e 5% são critérios sugeridos a acordar antes da execução; não são ganho estimado ou exigência da fonte. Nenhuma execução ocorreu. Transferência fica em 2 porque compatibilidade completa e relevância operacional ainda não foram estabelecidas.

## I2 — confirmação utilizável; WAND genérico corretamente descartado

**Fonte → princípio/passos:** os limites superiores, heap cheio e comparação estrita para poda definem um algoritmo concreto. A fonte impede poda por bloco para peso não positivo/não finito e mantém igualdade competitiva; não se pode transferir max_weight diretamente para BM25.

**Condições → alvo:** conferi search → search_bounded → wand_should_bool em search.rs:594–630 e :684–704. O caminho exige bool apenas should com termos simples, min_should_match efetivo 1, sem explain e cap finito. Em :788–793, igualdade continua sendo pontuada; :824–836 conta cada documento e avança seus cursores. A decisão de não recomendar “adicionar WAND” está sustentada. A confirmação é do princípio de poda de score no recorte, não da presença de todo algoritmo block-max da fonte.

**Decisão → teste:** preservar a capacidade existente evita implementação duplicada. Li o teste :3375–3435: ele exige uma precondição não vazia de poda e compara total e bits de score para vários caps; raw.rs:7122 em diante compara a variante de blocos com brute force. São evidências estáticas específicas, não resultados executados. Conservadorismo de limites, domínio de boosts e aritmética continuam condições de correção; esta revisão não certifica todo domínio do alvo. Block-max + BM25 + total exato permanece NOT_ESTABLISHED, corretamente sem recomendação de implementação.

## Formulação utilizável ao Ivan

“Encontramos uma hipótese de leitura seletiva de postings com diretórios e blocos independentes, sustentada por uma implementação externa e por uma diferença no formato comprimido observado do XERJ. Antes de prototipar, é preciso medir o custo de abertura/miss frente ao cache FTS já existente e preservar resultados, scores, posições, formatos antigos e total exato. Ainda não há ganho medido. Também verificamos que poda de score por limite superior já existe para um shape específico de consulta; por isso descartamos adicionar WAND genericamente como novidade. A variante block-max não foi estabelecida como transferência viável.”

A baseline H1/H3 foi refinada por um mecanismo concreto; H2 foi confirmada. Não há comparação independente com outros ensaios ou ferramentas, prova de ausência global, adoção, exclusividade ou superioridade geral. O recorte sustenta um candidato de experimento e um limite útil; **zero melhorias demonstradas**.
