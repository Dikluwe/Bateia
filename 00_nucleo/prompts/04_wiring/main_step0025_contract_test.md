# Prompt espelhado — oráculo A/B L4 do Passo 0025

Exija um helper de composição L4 `newMCPServer` que receba o funil legado, expansor,
provedores de descoberta e conteúdo e streams; construa `RecommendationService`, injete-o
em `mcp.NewMCPServer` e retorne o servidor pronto. Com fakes inteiramente offline, a
composição real deve anunciar exatamente três ferramentas e rotear uma chamada válida de
`bateia_recommend` até o provedor injetado, sem rede, ambiente ou estado global.

Não leia, adapte nem corrija implementação candidata.
