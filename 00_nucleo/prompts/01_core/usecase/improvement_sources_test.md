# Caracterização da resolução de fontes da síntese

**Camada:** L1
**Proprietário:** `01_core/usecase/improvement_sources_test.go`
**Contrato:** passo 0127.

Caracterizar o código anterior à extração, com fontes em memória e digest fixado da projeção resolvida e dos cinco conjuntos da análise, incluindo seus IDs e estado de negócio. Cobrir sucesso, degradação de fonte terminal, Hash/ArtifactHash, erros propagados, run incompatível/não terminal, índice ausente, hash divergente, página com identidade/hash incorretos, JSON inválido, mecanismos duplicados, cursor inválido, cópia de bytes e teto de 1000 páginas.

Testes não realizam I/O, não introduzem relógio nem consultam provedores. Reutilizar os mesmos cenários após extração. Asserções comportamentais devem provar que os resultados e falhas são preservados; a restrição da porta é expressa pela interface de produção.
