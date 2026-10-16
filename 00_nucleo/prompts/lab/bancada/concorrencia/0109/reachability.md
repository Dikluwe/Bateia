# 0114 — Alcance estruturado limitado e ordem de defer

Lab puro v6. Acrescentar call_site_reachability e defer_order separados,
preservando relações/callbacks anteriores. Analisar cada função/closure como
escopo independente. Return e chamada direta ao builtin panic encerram sequência;
if com ambos ramos terminantes encerra continuação. Analisar sequências de
case/select/switch e corpos de loops; não provar viabilidade de guardas ou saída
de loops/select. POSSIBLY_REACHABLE é sobreaproximação, não prova de execução.

Goto, label e demais branch statements tornam todo escopo UNSUPPORTED_CONTROL_FLOW;
não emitir fatos de ordem nesse escopo. Closure não herda inacessibilidade de
sua criação. Chamada em return é avaliada antes do terminador. panic sombreado,
defer panic e go panic não são terminadores da sequência chamadora.
UNREACHABLE_AFTER_TERMINATOR deve vincular blocker físico/hash; os demais
sites permanecem POSSIBLY_REACHABLE. Modo DIRECT/DEFERRED/GO preservado:
classificação trata site/registro, não execução de defer/goroutine.

PREFIX_DEFER_BEFORE_CALL_SITE liga defer diretamente no corpo do escopo a site
DIRECT em statement posterior desse corpo, inclusive dentro de if/loop/case,
sem atravessar closure. Defer condicional não vira prefixo obrigatório; não
usar ordem textual entre funções/escopos. Somente site não bloqueado e escopo
sem jumps. Relação condicional: se o site for atingido normalmente pelo corpo,
o registro do defer prefixo foi atravessado. Não prova execução posterior do
callee diferido, ausência de panic/block antes do site ou balanceamento.

Autoridade LIMITED_STRUCTURED_CONTROL, até 10.000 fatos novos (sites+ordens),
erro antes de publicação. Ordenar determinístico. Preservar lacunas de
viabilidade/execução/instância e limitações. Não remover vínculo histórico
lexical quando seu site é inacessível; nova projeção explicita a distinção.
