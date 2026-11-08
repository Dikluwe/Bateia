# ADR-0012 — Portas explícitas de execução e arquivos

**Status:** ACEITO PARA O RECORTE 0129
**Data:** 2026-10-28
**Governante:** Tekt, ADR-0002 e ADR-0008.

## Decisão e fronteiras

L1 declara contratos de runtime temporal e acesso a arquivos. L3 implementa relógio, deadlines, sinais, ambiente, entropia e operações de arquivos/raízes confinadas. L2 adapta pedidos e apresenta resultados usando essas portas; L4 configura e injeta o host. Não existe fallback produtivo que recrie I/O concreto em L1 ou L2.

O relógio da orquestração é obtido por porta explícita ou pelo contexto de execução configurado. Sem runtime, metadados temporais permanecem desconhecidos e durações zero; um timeout solicitado sem runtime configurado é erro explícito. Um runtime real é instalado nos caminhos produtivos; testes podem fornecer runtimes controlados. Parsers continuam com instantes de entrada, como no ADR-0011.

O host de arquivos expõe handles abstratos e raízes confinadas, sem tipos os.File/os.Root nem flags syscall no contrato. Operações de não seguir links e abertura não bloqueante são opções portáveis realizadas em L3. Tipos io/fs são dados de metadados. Escritas/leituras continuam respeitando os limites, retenção de descritores, checagens de identidade/hardlinks, atomicidade e ordens de cleanup existentes. Não trocar leituras confinadas por validação de caminho seguida de abertura insegura.

Composição do runtime conserva um único lifecycle/store/fila. Configuração de mercado e mecanismos passa a containers internos por capacidade, sem alterar fachadas públicas ou política de investigação. Testes de transporte/parse de conectores pertencem a L3; testes da composição permanecem em L4. As regras de lint devem reconhecer construções Go e fronteiras reais, sem suprimir a verificação de lógica de negócio em L4.

## Compatibilidade e verificação

Preservar formatos, IDs, hashes de conteúdo, erros de segurança, defaults produtivos e comportamento offline. Chamadas de biblioteca que usam filesystem precisam de host explicitamente configurado; APIs CLI públicas permanecem disponíveis. O executável compõe os serviços; testes de integração instalam o mesmo host em sua composição de teste. Não inicializar drivers por import lateral.

A conclusão exige inventário de efeitos sem chamadas concretas em L1/L2, testes funcionais de segurança/retomada/composição, vet, build e lint sem novas exceções genéricas. Alterações históricas e gates selados permanecem intactos.
