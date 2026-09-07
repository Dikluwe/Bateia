# Oráculo independente A/B — Passo 0035

Derivado exclusivamente de `.agents/AGENTS.md`, do Passo 0035 e do contrato público congelado comunicado pelo coordenador. Não consultar implementação candidata nem adaptar expectativas ao seu resultado.

Criar testes externos `clang_test` com Clang real e diretórios temporários. Uma unidade C com include e macro definidos via `compile_commands.json` deve produzir definição e chamada resolvida com arquivo/linha, versão e argumentos efetivos. Verificar hashes SHA-256 reais de fonte, header e configuração e sua alteração após edição independente de cada entrada. Configuração/include ausentes, flags de plugin/resposta, entradas apenas `command` e contexto cancelado não podem resultar em completude; falhas devem explicar o impedimento por erro, motivo ou diagnóstico. O compilador indicado pelo repositório e comandos de shell jamais são executados; usar arquivos sentinela para detectar execução.

Não inferir cobertura global, execução correta ou adequação à consulta da AST. Este oráculo avalia apenas os estados e evidências públicos da análise local. Não atesta isolamento de capacidades.
