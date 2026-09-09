# Passo 0052 — Expansão de vocabulário com proveniência

2026-09-08. A pergunta portuguesa de descompressão segue sem fontes, apesar da leitura do arquivo. Implementar expansão lexical delimitada na interpretação da consulta, sem alterar a consulta original nem o corpus de código.

## Contrato congelado

ImplementationQuery ganha Expansions []ImplementationQueryExpansion (json expansions,omitempty). Novo tipo: Source string, Terms []string, RuleID string, Kind string, Status string. Status sempre hypothesis; Kind translation ou mechanism. Source é o token natural original normalizado; Terms lista ordenada sem duplicatas; RuleID identifica versão e família. Não afirmar equivalência comportamental.

Aplicar somente a tokens naturais, nunca a identificadores explícitos, incluindo tokens entre backticks, CamelCase/snake ou consulta de identificador único. Não aplicar regra a substring, nem encadear resultados como gatilhos. Se a linha possuir qualquer marcador de restrição da interpretação 0048, não aplicar novas expansões nesta linha: escopo da negação desconhecido. Preservar termos e símbolos anteriores. Terms acrescenta os termos derivados; Expansions registra todas as derivações aplicadas em ordem determinística. Repetição do mesmo gatilho não duplica recibo. Tokens naturais com inicial maiúscula devem funcionar igualmente; símbolos explicitamente reconhecidos permanecem protegidos.

Catálogo inicial finito 0052-v1 (não vocabulário geral):

| Gatilhos exatos naturais | translation | mechanism (alternativas, não equivalências) |
| --- | --- | --- |
| descomprimir, descompressão, descompressao | decompress, decompression | decode, decoded, decompresses |
| comprimir, compressão, compressao | compress, compression | encode, encoded, compresses |
| serializar, serialização, serializacao | serialize, serialization | marshal |
| desserializar, desserialização, desserializacao | deserialize, deserialization | unmarshal |
| ordenar, ordenação, ordenacao | sort, sorting | — |
| comparar, comparação, comparacao | compare, comparison | — |
| bloco, blocos | block, blocks | — |
| corrompido, corrompidos, corrompida, corrompidas | corrupt, corrupted | — |
| dados | data | — |
| destino | destination | dst |
| reutilizar | reuse | — |
| rejeitar | reject | — |

Cada família/tipo tem ID estável prefixado 0052-v1. Ausência no catálogo não vira erro nem equivalência inferida. Este catálogo é uma hipótese de engenharia escolhida antes de testar a candidata; não é extraído dos projetos nem aprendido automaticamente.

Ranking e documentos usam Terms como antes. Não mudar IDF, limites, peso por termo, glossário existente de código, classificação de capacidades ou requisitos. Revalidação reconstrói Interpretation da consulta original como antes. Reason/projeção pode continuar lexical; Markdown de interpretação deve exibir as expansões (até 6 por linha, até8 termos cada, com contagens restantes) e avisar que correspondência derivada não prova intenção/capacidade. JSON preserva recibos completos. Atualizar versão e checkpoints local/live.

## Validação e critério de ganho

B congela testes independentes a partir deste contrato, sem ler candidata: gatilhos, capitalização natural, não-substring/não-cascata, símbolos explícitos, negação, determinismo, fonte e tipo/estado, nova recuperação por API funcional em corpus de nomes neutros e controles negativos. Conferir Markdown e reconstrução após revalidação. Não tratar hits documentais como prova.

Integrador executa antes/depois na CLI para pergunta portuguesa de descompressão, mantém a pergunta original e verifica se Decode em decode.go é localizado e revalidado nos primeiros12 caminhos. Controlar ordenação ou serialização em outro domínio e consultas negativas, incluindo alternativa técnica irrelevante. Publicar hits/termos e limites: crescimento de hits pode incluir falsos positivos; ausência de benchmark amplo impede alegar precisão geral. Preservar controles locais/remotos, build, testes e lint. Sem alteração dos testes congelados de passos anteriores.

## Regime e limites

Tekt A/B sem atestação física: integrador contrato/presentação/versões, A interpretação/domínio, B testes independentes; filesystem compartilhado. Nenhum selo de refinamento. Expressões com negação continuam sem novas expansões; palavras desconhecidas não ganham tradução automática. Codificar não equivale a comprimir, marshal não equivale a qualquer serialização: são alternativas para leitura, rotuladas hypothesis. Somar aliases pode alterar ranking; calibração por conceitos e classificação conjunta ação/objeto/condição são trabalho posterior. Sem LLM, rede adicional ou alteração da consulta ao provedor neste passo.
