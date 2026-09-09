# Testes independentes do executor de efeitos

Papel C, contrato congelado: `00_nucleo/passos/0062-verificar-efeito-executado.md`.
Não consultar implementação candidata, gate, fixtures de B ou relatórios. Criar
fontes `package candidate` próprias em `t.TempDir`. Testar API pública em
`package effects_test`, sem executar testes neste papel.

Cobrir ordenação válida com preservação por oráculo próprio de multiplicidades,
contraexemplos reais (mutação da entrada e panic), timeout, cancelamento e fontes
fora do recorte. Validar recibos públicos, reprodução por seed, conteúdo/hash da
fonte e wrapper, hashes de artefatos, ausência de segredo sentinela e remoção dos
temporários mencionados. Não supor nomes internos em execution, pois o contrato
congela somente seu conteúdo sem fixar as chaves. Fontes positivas são compiladas
duas vezes exclusivamente para verificar reprodução. Nenhuma execução de código
hostil. O processo não constitui sandbox. Prompt precede o teste.
