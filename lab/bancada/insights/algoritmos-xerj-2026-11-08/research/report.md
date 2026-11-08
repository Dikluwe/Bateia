# Investigação independente: XERJ

Decisão proposta: considerar **um protótipo condicional de leitura de postings por faixa para abertura fria seletiva**. Não recomendar migração sem perfil e paridade. O princípio genérico de poda WAND já existe no alvo e foi descartado como novidade.

## Método e limites

Baseline escrito antes de qualquer chamada do produto (hash em coverage.json e audits). Alvo fornecido: xerj-org/xerj, pin 1557de2f5be981d9db441969201c92a8af6d2de3. Toda descoberta/coleta/transporte passou pelo runner Bateia; consulta externa por capacidade após ler o alvo. Fontes não foram executadas nem abertas diretamente. Duas aquisições externas. Formulate valida estrutura/bytes, não interpretação nem efeito.

Os providers locais registram explicitamente que o pin é metadado Git consultivo: hashes dos itens/snapshots testemunham os bytes lidos, mas não autenticam esses bytes pelo commit. Não foram lidos dados de treinamento/modelos proprietários nem testado serviço do XERJ. Testes citados foram lidos, nunca executados.

A comparação docs/XERJ_VS_LUCENE.md é histórica, pinada em outro commit; serviu somente para orientar caminhos. As decisões abaixo citam arquivos no snapshot atual retornado pelo produto. A cadeia interna FTS foi verificada; configuração do engine/HTTP, cache e frequência real de abertura permanecem desconhecidos.

## I1 — Experimentar payloads de postings endereçáveis por bloco para abertura fria seletiva

Relação: **DIFFERENT_APPROACH**. Separar diretórios pequenos residentes do payload e endereçar blocos independentes permite carregar somente listas e blocos necessários, em vez de descomprimir todos os postings do campo na abertura. É uma hipótese de transferência de formato e leitor, não um ganho medido.

Observação e decisão atual: No recorte atual, FtsIndexReader::open lê o .post inteiro e descomprime o envelope; a busca usa slices sem descompressão por consulta. Existe open_stats_only que já evita postings no pré-passe de estatísticas.

Alternativa: Decidir primeiro por um protótipo de diretório + framing independente e leitor por faixas para abertura fria seletiva; manter o leitor atual como baseline/fallback. Não aprovar migração nem mudança de total sem experimento.

Passos e invariantes:

1. Persistir diretório com intervalo de doc IDs, offset e comprimento de cada bloco; manter a base necessária para decodificar gaps independentemente.
2. Usar offsets relativos a uma unidade acessível: no XERJ, um novo framing teria de substituir ou particionar o envelope Zstd inteiro; gravar apenas o SkipEntry atual não basta.
3. Abrir diretórios e metadados; ler e validar por faixa os blocos dos termos efetivamente pedidos.
4. Preservar codecs legados, posições, frequências, normas e invariantes de checksum; manter o percurso completo onde a API exige total exato.
5. Medir separadamente abertura fria, cache aquecido e percursos completos antes de decidir a migração.

Evidência literal da fonte:

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `docs/raw-segments.md:3-5`, `sha256:c4db83f4f97cf8c5360b619ea1d72380f52ce986821ded5cd860d0ca2f1ff872`

```
`raw-segment` exposes numeric-term byte and file segment formats for large
lexical or learned-sparse indexes. The readers keep fixed directories resident
and range-read posting payloads for query terms.
```

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `src/raw.rs:609-615`, `sha256:698db589830dfc475eb6b2f82b72c00a016e69f7ab74207b2369f75af31f09dd`

```
pub struct RawPostingBlockMeta {
    base_doc_id: DocId,
    last_doc_id: DocId,
    postings_offset: u64,
    postings_len: u32,
    max_weight: u32,
}
```

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `src/raw.rs:3308-3318`, `sha256:698db589830dfc475eb6b2f82b72c00a016e69f7ab74207b2369f75af31f09dd`

```
            self.read_posting_block_range_into(
                block.postings_offset,
                block.postings_len as u64,
                &mut block_bytes,
            )?;
            for_each_posting_in_block(
                entry.term_id,
                &block_bytes,
                block.base_doc_id,
                block.last_doc_id,
                |doc_id, doc_weight| {
```

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `src/raw.rs:3042-3053`, `sha256:698db589830dfc475eb6b2f82b72c00a016e69f7ab74207b2369f75af31f09dd`

```
        checked_range(offset, len, self.file_len, "postings")?;
        let len = usize::try_from(len).map_err(|_| Error::SegmentTooLarge)?;
        bytes.clear();
        bytes.resize(len, 0);
        read_exact_at_positional_into(
            &mut self.file,
            segment_file_offset(self.data_offset, offset)?,
            bytes,
        )?;
        if let Some(checksums) = self.decoded_block_checksums()? {
            verify_block_slice(checksums, offset, bytes)?;
        }
```

Evidência literal do alvo:

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/index.rs:2700-2719`, `sha256:f91700ce299b3172fb3cefed572510ffb173273263f6c83825206359ddeaa92a`

```
            // The query path references `post_data` by slice in all
            // cases, so there's no per-query decompress cost.
            //
            // STATS-ONLY: skip the read + decompress entirely — the caller
            // only wants `field_stats`/`term_doc_freq`, both of which live
            // in `.meta`/`.fst`.
            let raw_post = if stats_only {
                Vec::new()
            } else {
                fs::read(&post_path).with_context(|| format!("reading postings {:?}", post_path))?
            };
            let zstd_envelope = raw_post.len() >= 8
                && (&raw_post[..4] == POST_MAGIC_ZSTD_V2 || &raw_post[..4] == POST_MAGIC_ZSTD);
            let post_block_v2 = raw_post.len() >= 4 && &raw_post[..4] == POST_MAGIC_ZSTD_V2;
            let post_data = if zstd_envelope {
                let mut len_buf = [0u8; 4];
                len_buf.copy_from_slice(&raw_post[4..8]);
                let uncompressed_len = u32::from_le_bytes(len_buf) as usize;
                let decompressed = zstd::bulk::decompress(&raw_post[8..], uncompressed_len)
                    .map_err(|e| anyhow::anyhow!("ZPS postings decompress failed: {e}"))?;
```

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/index.rs:2597-2609`, `sha256:f91700ce299b3172fb3cefed572510ffb173273263f6c83825206359ddeaa92a`

```
    /// Open ONLY what BM25 collection statistics need: the FST term
    /// dictionary (mmap'd) and the small `.meta` side-car.
    ///
    /// Skips the two expensive parts of a full [`Self::open`]:
    ///
    ///  * the postings envelope — a whole-file `fs::read` plus a zstd
    ///    decompress of every posting list in the field, O(index bytes);
    ///  * the norms table — a whole-file `fs::read` plus decode, O(docs).
    ///
    /// That makes the per-segment `field_stats` + `term_doc_freq` pre-pass
    /// the index-wide scorer needs (#188) cost, per field, one mmap plus one
    /// `.meta` read-and-decode — for the current ZFM4 format that decode is
    /// a zstd decompress of the `num_terms × 24`-byte records section, i.e.
```

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/postings.rs:74-83`, `sha256:35f56fe56978ca6437bf78571d88f647dd540aa7684f297d51a7a997d29cc5e3`

```
//! ## Skip-list acceleration: NOT implemented
//!
//! There is **no on-disk skip table**.  [`PostingsWriter::encode_term`] computes
//! a [`SkipEntry`] vector in memory (one entry every [`SKIP_INTERVAL`] blocks)
//! but it is **never serialised** — the caller in `xerj-fts::index` discards it,
//! and the whole `.post` blob is wrapped in a Zstd envelope, so intra-blob byte
//! offsets would not be usable as seek targets anyway.  Consequently
//! [`PostingsReader::advance_to`] performs a **linear scan**, not a skip-list
//! seek.  Traversal is correct, just unaccelerated; the [`SkipEntry`] machinery
//! is scaffolding for a future block-skip implementation.
```

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/search.rs:746-756`, `sha256:cacffa87fd1a7a4d67643ce88dd3d355164cf6532e32014d0943a20a3c2abade`

```
            // The factory reads the field's real positions layout and the
            // envelope's block codec — a docs-only reader must synthesise
            // tf=1 rather than misparse a positions-bearing stream, exactly
            // like `scan_term`; the decoded positions are simply never read.
            let Some(mut reader) = self.reader.postings_reader(&tq.field, &tp) else {
                continue;
            };
            let Some(first) = reader.next() else { continue };
            let bm25 = self.make_scorer(&tq.field);
            let score_df = self.scoring_df(&tq.field, &tq.term, tp.doc_frequency as u64);
            let ubound = bm25.idf(score_df) * (bm25.k1 + 1.0) * tq.boost;
```

Condições de aplicação:

- Leitores abertos com poucos termos consultados em relação ao vocabulário e payload do campo.
- Perfil confirmar custo material de fs::read/descompressão e retenção do .post; frequência de reabertura e caches da rota de produção ainda não foi estabelecida.
- Novo formato permitir blocos independentes e leitor compatível com arquivos históricos.

Condições de falha:

- Corpus pequeno, leitores mantidos aquecidos ou consultas que percorrem quase todo o payload podem perder para a leitura sequencial atual.
- Compressão por bloco pode piorar espaço, throughput de ingestão ou aumentar operações de I/O.
- Índice de offsets sem mudança do envelope Zstd inteiro não torna o payload comprimido endereçável.
- Pular correspondências para ranking pode quebrar total exato; esta proposta não autoriza essa mudança.

Benefício hipotético: Hipótese: diminuir bytes lidos e descomprimidos e memória retida quando uma abertura consulta poucos termos. Não há previsão quantitativa sustentada.

Custo: Versão nova de formato, escritor/leitor/merge compatíveis, checksums e testes de corrupção/recuperação, diretórios adicionais, possível perda de compressão e mais I/O. Integração em caches e lifecycle não estudada.

Experimento/observável: Experimento futuro, sem executar fonte nesta investigação: corpus real de código e logs, consultas com poucos/muitos termos, novo leitor e leitor atual com cache frio/aquecido; medir bytes realmente lidos/descomprimidos, RSS, p50/p95, ingestão e tamanho. Comparar IDs, scores, posições e total exato; testar arquivos legados e corrupção. Pré-registrar aceitação: pelo menos 20% menos bytes descomprimidos no braço frio seletivo, sem divergência e sem regressão >5% de p95 aquecido ou ingestão; revisar limites conforme orçamento do projeto antes de executar.

Rejeição: Rejeitar migração se houver qualquer divergência semântica/corrupção não detectada, se perfil mostrar leitores aquecidos dominantes sem custo de abertura relevante, ou se os ganhos frios não compensarem regressões/custos pré-registrados.

Baseline: H1/H3 refinadas por evidência para protótipo condicional

Alcance: Cadeia interna da crate FTS lida. Rotas HTTP/engine, configuração de cache, frequência real de abertura e workload não rastreados; aplicabilidade em produção desconhecida.

Testemunhos: audit/004.json e 004.output.snapshot; audit/009.json e 009.output.snapshot; audit/010.json e 010.output.snapshot. Identidades e SHA-256 dos snapshots em sources.json.

## I2 — Confirmar poda de score por limite superior já presente; descartar adoção genérica de WAND como novidade

Relação: **ALREADY_PRESENT**. Quando o heap top-k está cheio, uma contribuição cujo limite superior é estritamente menor que o pior score retido pode evitar avaliação; igualdade deve respeitar desempate. A unidade podada e o contrato de contagem continuam específicos de cada sistema.

Observação e decisão atual: XERJ já encaminha search_bounded para wand_should_bool, constrói ubound BM25, evita score abaixo do pior retido e avança next contando todo doc.

Alternativa: Descartar a recomendação genérica de adicionar WAND; preservar a implementação existente. Manter a variante block-max como relação específica não estabelecida, sem recomendação de implementação neste relatório.

Passos e invariantes:

1. Preencher o heap top-k antes de usar seu limiar competitivo.
2. Construir limites superiores conservadores das contribuições de score.
3. Podar apenas abaixo do limiar estrito; tratar igualdade conforme ordenação pública.
4. No XERJ, continuar visitando e contando cada correspondência para total exato; não transferir automaticamente a poda de blocos da fonte.

Evidência literal da fonte:

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `src/raw.rs:3290-3305`, `sha256:698db589830dfc475eb6b2f82b72c00a016e69f7ab74207b2369f75af31f09dd`

```
        self.validate_block_directory(block_directory)?;
        let mut ranked = Vec::with_capacity(k);
        let mut sorted = false;
        let can_prune_blocks = query_weight > 0.0 && query_weight.is_finite();
        let mut block_bytes = Vec::new();
        for block_index in 0..block_directory.block_count {
            let block = self.posting_block_at(entry, block_directory, block_index)?;
            if can_prune_blocks && ranked.len() == k {
                if !sorted {
                    ranked.sort_by(crate::cmp_doc_scores);
                    sorted = true;
                }
                let threshold = ranked.last().expect("top-k buffer is full").1;
                if query_weight * (block.max_weight as f32) < threshold {
                    continue;
                }
```

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `src/raw.rs:2502-2506`, `sha256:698db589830dfc475eb6b2f82b72c00a016e69f7ab74207b2369f75af31f09dd`

```
                    // Equality cannot be discarded: the smaller document id
                    // wins score ties in the crate's public ordering.
                    if upper_bound >= threshold {
                        pivot = Some(candidate);
                        break;
```

- `arclabs561/postings@edb02683c6359efe2034a6cc71b64c9aaed90cb0`, `src/raw.rs:7122-7128`, `sha256:698db589830dfc475eb6b2f82b72c00a016e69f7ab74207b2369f75af31f09dd`

```
    fn block_max_pruning_skips_a_block_and_matches_brute_force() {
        // Oracle test for the block-skip at top_k_single_raw_term: block 0
        // (128 entries, weights ~100) fills the top-k, then block 1's max
        // (weight 3) falls below the threshold and the block is skipped.
        // The result must still equal an unpruned brute force over the raw
        // postings; the skip precondition is asserted explicitly so the test
        // stays pinned to the pruning regime if constants drift.
```

Evidência literal do alvo:

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/search.rs:622-630`, `sha256:cacffa87fd1a7a4d67643ce88dd3d355164cf6532e32014d0943a20a3c2abade`

```
        // WAND fast path for the shape a `match`/`multi_match` query becomes
        // after planning (`bool{should: [Term..]}`, nothing else).  Only when
        // nothing downstream needs the generic walk: `explain` needs every
        // retained hit's breakdown, and `cap == usize::MAX` callers want the
        // complete match set.
        if !explain && cap != usize::MAX {
            if let Some(res) = self.wand_should_bool(query, cap)? {
                return Ok(res);
            }
```

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/search.rs:652-664`, `sha256:cacffa87fd1a7a4d67643ce88dd3d355164cf6532e32014d0943a20a3c2abade`

```
    /// Design: WAND-style upper-bound pruning (Broder et al., *Faster Top-k
    /// Document Retrieval Using Block-Max Indexes*; the same dynamic-pruning
    /// idea as tantivy's `block_wand`, MIT — quickwit-oss/tantivy
    /// `src/query/boolean_query/block_wand_union.rs`, whose block-max half
    /// needs impact-ordered postings and skip tables this postings format
    /// does not carry).  The walk still visits EVERY matching doc — postings
    /// have no skip lists, and the caller needs an EXACT total — but a doc
    /// is scored only when the sum of its matched terms' score upper bounds
    /// (`idf × (k1+1) × boost`, the supremum of `score_term` over any tf/dl)
    /// can still enter the heap.  For stopword-only docs that bound sum sits
    //    below the 10th-best score almost immediately, so the entire scoring
    ///   cost for them disappears; what remains per skipped doc is one
    ///   cursor compare and one `next()`.
```

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/search.rs:788-793`, `sha256:cacffa87fd1a7a4d67643ce88dd3d355164cf6532e32014d0943a20a3c2abade`

```
            let score_it = match top.worst() {
                None => cap > 0,
                Some(worst) => {
                    let bound: f32 = scorers[..k].iter().map(|s| s.ubound).sum();
                    bound >= worst.score
                }
```

- `xerj-org/xerj@1557de2f5be981d9db441969201c92a8af6d2de3`, `engine/crates/xerj-fts/src/search.rs:824-836`, `sha256:cacffa87fd1a7a4d67643ce88dd3d355164cf6532e32014d0943a20a3c2abade`

```
            total += 1;

            // Advance every scorer sitting on `d` — exactly `scorers[..k]`
            // (the sort invariant makes the at-d cursors a prefix).
            // Exhausted cursors are marked `u32::MAX`, which sorts last, and
            // pop off after the re-sort.  `sort_unstable` never allocates the
            // scratch buffer a stable sort spills to on >20-element slices —
            // this re-sort runs once per DOC, so that alloc was per-doc too.
            for s in scorers.iter_mut().take(k) {
                match s.reader.next() {
                    Some(next) => {
                        s.doc = next.doc_id;
                        s.tf = next.term_freq;
```

Condições de aplicação:

- O princípio amplo de evitar cálculo de score não competitivo já opera no XERJ no bool somente should com termos simples, min_should_match efetivo 1, sem explain e com cap finito.
- Transferência da variante em blocos só faria sentido com limites adequados ao BM25 e preservação explícita do total.

Condições de falha:

- Limites subestimados removem resultados válidos.
- Podar em igualdade pode mudar o vencedor do desempate.
- A fonte pontua produto interno esparso; max_weight não é automaticamente limite de BM25.
- Poda de blocos que deixa de visitar docs não mantém automaticamente a contagem exata exigida no alvo.

Benefício hipotético: Evitar trabalho duplicado e evitar confundir confirmação de princípio com novidade. Nenhum ganho novo alegado.

Custo: Nenhuma migração proposta; revisão de contratos seria necessária se alguém explorasse a variante block-max.

Experimento/observável: Leitura confirma a cadeia interna search -> search_bounded -> wand_should_bool -> postings_reader/next. Testes de paridade da fonte e alvo foram apenas lidos. Uma futura mudança teria de preservar IDs/scores/ordem/total, incluindo empates e fallback de shapes ricos.

Rejeição: Rejeitar a alegação de novidade do princípio amplo porque ele já está presente. Não inferir que a variante por blocos também esteja presente.

Baseline: Confirma H2; princípio concreto só identificado após ler fontes e chamadores.

Alcance: Cadeia interna da crate FTS lida. Rotas HTTP/engine, configuração de cache, frequência real de abertura e workload não rastreados; aplicabilidade em produção desconhecida.

Testemunhos: audit/004.json e 004.output.snapshot; audit/009.json e 009.output.snapshot; audit/010.json e 010.output.snapshot. Identidades e SHA-256 dos snapshots em sources.json.

## Cobertura e encerramento

A fonte arclabs561/postings tem pin edb02683c6359efe2034a6cc71b64c9aaed90cb0. A cadeia de consumo é top_k_weighted_u32 -> top_k_single_raw_term -> read_posting_block_range_into; o caminho multitermo usa diretórios e pode evitar decodificar blocos anteriores em advance_wand_cursor_to (2577–2605). Isso sustenta o mecanismo na fonte, sem provar adoção adequada no alvo. A matriz de benchmark inclui bounds úteis e deliberadamente inúteis, reforçando que layout e distribuição importam; seus números não foram executados nem transferidos.

Não perseguido: a fonte Java, cujo README isolado não sustenta algoritmo, nem a correção de mixed read-under-write baseada em documento histórico. Não foi feita busca exaustiva de capacidades do alvo. Ausência no recorte não demonstra ausência global. A evidência atual suporta diferença de leitor/formato e confirmação de poda de score; não suporta superioridade, economia quantificada ou benefício em produção.

Formulate concluído com sucesso (audit/011.json e respectivo snapshot); o resultado é registro estrutural, com efeito UNVERIFIED. Total: 9 chamadas de pesquisa/produto (incluindo formulate) e 2 aquisições.
