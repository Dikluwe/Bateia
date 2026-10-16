# Contrato da CLI de insights

**Camada:** L2  
**Arquivo proprietário:** `02_shell/cli/idea_insights_test.go`  
**Contrato:** passo 0105

Verificar despacho público, stdout JSON ou arquivo atômico, rejeição sem saída/mutação para inputs inválidos, aliases de output e writers curtos. Fixtures de fontes congeladas sem chamadas de provedores ou runtime. Nenhum teste deve alterar saídas legadas ou goldens históricos.
