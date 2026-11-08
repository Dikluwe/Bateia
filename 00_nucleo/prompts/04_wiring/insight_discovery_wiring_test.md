# Composição do comando de descoberta de insights

**Camada:** L4, testes de wiring. **Contrato:** passo 0137.

Verificar despacho real de main em subprocesso de teste isolado, sem replicar regra de plano/seleção e sem rede real. A chamada idea insights discover sem output deve alcançar a validação da CLI injetada (e não o fallback offline sem portas), sem conteúdo em stdout. Isolar ambiente/diretório do processo. Analyze/compare/explore preservam caminho offline; sem auditoria AST geral em L4.
