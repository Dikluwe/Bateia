# Testes independentes da projeção de requisitos — 0046

Autor B, regime A/B sem atestação física. Entrada: contrato do passo 0046 e API pública de domínio; sem leitura da implementação candidata, testes anteriores, diffs ou gates.

Antes do código, fixa-se: linhas não vazias deduplicadas e ordenadas; evidência lexical exclusivamente de interseções entre termos canônicos e MatchedTerms; termos restantes como MissingTerms; origem direct_call e identidade incompleta/inválida não fundamentam fontes. Exigir status evidence_located apenas para fontes válidas; demais casos not_verified. A saída deve ser determinista sob permutação de hits, deduplicar fontes, preservar símbolo/linha/hash e não alterar o relatório de entrada. Cobrir consultas sem termos discriminantes, hits sem interseção e relatórios sem suporte/parciais sem fontes.

Os hashes dos arquivos de teste e prompts são congelados antes de qualquer execução. Hash identifica artefatos, não prova isolamento físico.

## Revisão B2 — orçamento por caminho

Após diagnóstico do baseline real e após implementação da correção pelo autor A, B recebe somente a política pública: RankGoImplementations visita um hit por caminho por rodada antes do limite, mantendo todos os símbolos exatos antes dos demais e sem mudar scores. Congelar esta revisão antes de executá-la, registrando explicitamente que seu congelamento sucede à implementação (sem atestação).

Caso sintético: várias funções com queue/store/read no nome em um arquivo e outro arquivo com somente queue; com limite 2 ambos os caminhos devem estar presentes, e a permutação dos arquivos não muda hits. Comparar scores e identidades ao limite amplo. Controle separado: duas declarações QueueTask no mesmo caminho prevalecem sobre alternativas lexicais, inclusive quando o limite é 2. Não consultar corpus0045 nem corpos da candidata.
