# Comando offline de insights

**Camada:** L2  
**Arquivo proprietário:** `02_shell/cli/idea_insights.go`  
**Contrato:** passo 0105

Expor idea insights analyze --input pedido.json [--output resultado.json]. JSON estrito e limitado, rejeição de desconhecidos/duplicatas/valores extras. Proteger input contra sobrescrita/aliases e usar writer atômico/bloqueio de segredos/short writes já existentes na casca. Apenas chamar o caso de uso puro; sem clientes L3, runstore ou política de reúso.
