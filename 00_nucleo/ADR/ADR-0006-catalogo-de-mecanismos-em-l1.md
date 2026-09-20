# ADR-0006: O catálogo ativo de mecanismos pertence a L1

**Status**: ACEITO
**Data**: 2026-09-18
**Relacionado**: passo 0078, passo 0028, ADR-0004

## Contexto

O catálogo `mechanism-rules.json` descreve classes, evidências aceitas e limites
da extração determinística. Seu decoder, seu hash canônico e todos os consumidores
semânticos já vivem em L1, mas o blob ativo e uma segunda cópia literal do hash
permanecem em L4.

Essa divisão faz a camada de composição possuir política de domínio e permite
divergência entre bytes, identidade e validação.

## Decisão

1. Mover o JSON byte a byte para `01_core/usecase/mechanism-rules.json`.
2. Materializá-lo em L1 por `go:embed`, sem I/O em runtime.
3. Expor somente `ActiveMechanismRuleCatalog() ([]byte, string)`.
4. Retornar um clone novo dos bytes em cada chamada e o hash canônico constante.
5. Fazer L4 obter uma vez o par bytes/hash e apenas injetá-lo na fachada MCP.

## Consequências

L1 passa a possuir conteúdo, identidade e validação do catálogo. L4 perde embed,
arquivo JSON, literal e cálculo de hash. O formato v3, regras, outputs, ordem de
bootstrap e superfície MCP não mudam.
