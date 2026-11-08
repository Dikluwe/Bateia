# CLI de releitura técnica

Expor comando com entrada JSON, data de corte, timeout limitado, relatório
e exportação opcional de bundles. Receber provedor configurado pelo wiring,
nunca do JSON. Impedir entrada e saídas sobrepostas. O executável usa o cliente
GitHub existente; despacho offline não deve fingir uma releitura remota.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
