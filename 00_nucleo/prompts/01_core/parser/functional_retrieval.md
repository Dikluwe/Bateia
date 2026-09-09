# Recuperação funcional AST — passo 0045
Origem: contrato congelado `00_nucleo/passos/0045-recuperar-perguntas-funcionais.md`.
Este prompt precede `01_core/parser/functional_retrieval.go`; testes são de outro autor.
L1 puro sem I/O. Extrair termos somente de identificadores AST do corpo de funções,
sem percorrer corpos de closures. Pontuar conjuntos: função com corpo e nome/corpo
pertinente precede declaração não funcional, sem superar identificador exato.
Construir índice por diretório/pacote/nome de declarações pacote. Destino elegível é
função sem receiver com corpo e única declaração desse nome; rejeitar colisões,
parsing inválido e caminhos com conteúdo conflitante para relações. Resolver somente
CallExpr.Fun Ident sem objeto local (ou objeto ast.Fun de declaração pacote). Ignorar
seletores, function values e FuncLit. Chamada deve estar no corpo do caller lexical;
referência contém caminho, símbolo caller, linha da chamada e hash dos bytes caller.
Expandir sementes selecionadas por um salto apenas. Reservar pelo menos dois terços
para sementes se disponíveis; expansão tem teto floor(limit/3). Ordenar destinos e
referências deterministamente e deduplicar. Destino já lexical preserva Retrieval
lexical; novo destino recebe direct_call e motivo explícito de relação estática.

## Revisão 0046 após baseline real

Perda observada: limite de hits consumido com múltiplos símbolos de poucos arquivos antes da seleção por caminhos. Antes de expansão e truncamento, distribuir hits lexicais em rodadas por arquivo preservando ordem relativa das primeiras ocorrências e prioridade de TODOS os exatos. Uma rodada admite no máximo um hit de cada caminho; rodadas seguintes preenchem orçamento. Não mudar pesos/glossário conforme corpus. Sementes multirrequisito preservam precedência de suas primeiras fontes; expansão direta mantém reserva e caller retido. Método passa go_ast_functional_v3. Avaliação compara corpus congelado e regressões0045, sem alegar ganho garantido.

0046 — Veredito experimental: a tentativa v3/diversificação acima foi descartada após comparação congelada sem ganho de recall@12 e com mais retornos não rotulados. Contrato ativo de ranking permanece v2/0045; histórico experimental preservado, não implementação ativa.
