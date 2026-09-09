# Passo 0045 — Recuperar implementações por perguntas funcionais

2026-09-07. Contrato congelado antes da implementação. Resultado em relatório separado.

## Objetivo e recorte

Melhorar a seleção local Go para perguntas funcionais sem promover coincidência textual a capacidade comprovada. Baseline: worktree 0044, não apenas HEAD. Reutilizar AST; não instalar serviços ou ampliar para outras linguagens. Regime Tekt A/B, autores distintos de implementação e testes, filesystem compartilhado sem atestação física.

## Contrato observável

Preservar API RankGoImplementations e limites de varredura/extração do 0044, identidade por bytes, exclusões e não avaliado para linguagens sem suporte. Identificador exato continua prioritário. Para perguntas funcionais, dar preferência a funções com corpo e identificadores operacionais pertinentes sobre tipos, comentários e literais isolados. Correspondência lexical continua heurística; não interpretar negação como prova de comportamento. Glossário PT/EN genérico e explícito, sem nomes de projetos ou funções esperadas. Não elevar documento repetitivo a mecanismo por contagem de ocorrências.

Distribuir resultados entre requisitos separados por linha, mantendo determinismo sob reordenação e duplicação de entradas. Expandir somente um salto de chamadas diretas por identificador, para funções de pacote no mesmo diretório/pacote, com destino único e sem sombreamento local. Não resolver seletores, interfaces, function values ou chamadas ambíguas. Limitar expansão a no máximo um terço do orçamento de hits; manter espaço para sementes lexicais. Não apresentar relação de execução dinâmica.

Adicionar a ImplementationHit os campos opcionais Retrieval (lexical/direct_call) e RelatedTo []ImplementationReference. ImplementationReference contém Path, Symbol, Line, SHA256, com JSON snake_case. A referência indica a função chamadora e linha da chamada. Hits expandidos identificam definição e referência; referências devem corresponder a arquivos extraídos com hash e intervalo válidos na releitura. Referência obsoleta remove a relação; um hit exclusivamente direct_call sem referência válida é descartado. A saída explica pertinência lexical e chamadas estáticas separadamente, sem afirmar que o mecanismo satisfaz a pergunta.

Versionar método/analisador e checkpoint local. Retomar não reutiliza versão antiga; replay preserva os dados congelados. Núcleo L1 sem IO, prompts prévios e linhagem.

## Avaliação congelada antes da candidata

Autor B prepara fixtures Go de pelo menos três pequenos projetos sintéticos independentes (por exemplo fila, armazenamento, autenticação), perguntas funcionais e arquivos esperados, incluindo negativos: documentação/strings repetitivas, tipos sem corpo, função irrelevante, chamada ambígua ou sombreada, ausência e parser inválido. Não ler implementação candidata. Congelar corpus e hashes, executar no baseline 0044 antes da alteração e depois na candidata. Medir recall de arquivos esperados e precisão dos arquivos retornados nos primeiros três caminhos distintos; apresentar contagens e escopo, não benchmark global. Exigir ganho agregado de recall ou precisão sem perda de recall e preservação de negativos e símbolo exato. Uma relação estática limitada não transforma ausência/negação em aprovação.

Testes independentes verificam prioridade funcional, chamada a helper sem termos, profundidade máxima, ambiguidade/sombreamento, limites, determinismo e invalidação de hashes de referências. O coordenador não lê corpos dos testes para implementar. Mudança de fixture exige motivo e hashes anteriores; não ajustar expectativa ao resultado candidato.

Repetir consultas reais do 0044, inclusive natural e símbolo exato. Reportar honestamente quais implementações foram recuperadas; melhorias em fixtures não demonstram ganho em projetos reais. Rodar suíte completa, lint, build do pacote ./04_wiring e replay local. Não alterar detectores para forçar capacidades aprovadas.

## Fora do escopo

Equivalência semântica, execução/benchmark do candidato, grafo interprocedural completo, C/C++, detecção universal de capacidades, gRPC/Protobuf e Bleve. Não reescrever gates/certificados existentes. Entregar fontes, passo e relatório; commit não faz parte desta solicitação.
