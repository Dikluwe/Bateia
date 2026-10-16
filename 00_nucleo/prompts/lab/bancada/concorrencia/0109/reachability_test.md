# 0114 — Contrastes de alcance

Return antes de callback e no mesmo case bloqueia site. Return avaliando callback
não o bloqueia. Ambos if branches return bloqueiam continuação, um só não.
Builtin panic termina; nome panic sombreado, go panic e defer panic não.
Closure após return conserva escopo independente. Goto/label/break produzem
UNKNOWN por escopo. Defer prefixo antes de site gera ordem; condicional, após
site ou em outro escopo não. Parsing/tipos inválidos suprimem fatos novos.
Posições/hash, ordem, limite e fonte congelada devem ser verificados. Expectativas
autorais não constituem semântica independente.

Expressões em tipos de case (incluindo chamada builtin len constante) precisam
ter projeção lexical sem erro; não afirmar invocação runtime de builtin constante.
