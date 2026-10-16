# 0104 Composição por registro validado

**Data**: 2026-10-06  
**Estado**: IMPLEMENTADO E VERIFICADO LOCALMENTE  
**Contrato**: [passo 0104](../passos/0104-compor-descoberta-por-registro-validado.md)  
**Decisão**: [ADR-0009](../ADR/ADR-0009-admissao-de-conectores-e-projecao-mcp.md)

## Resultado

A composição padrão de síntese já não é rejeitada por conter mais de três conectores. `NewMarketDiscoveryRegistry` em L1 possui a admissão de configuração não vazia, capacidades válidas e IDs únicos. Os dois construtores L4 usam esse contrato sem contagem fixa e preservam recursos/lifecycle compartilhados. O registro genérico vazio continua permitido para os consumidores que precisam dessa representação.

L2 agora anuncia o enum das fontes configuradas tanto no servidor de descoberta quanto na extensão de síntese. A função de projeção por registro já existia, mas não estava conectada à lista pública de ferramentas; a síntese também reconstruía a lista basal sem usar o registro. O limite de seleção por pedido permanece até dez, mesmo com onze fontes disponíveis.

## Evidência de regressão e preservação

Antes da correção, os testes novos reproduziram rejeição dos conjuntos unitário e ampliados e schema de três fontes nos registros com cinco, seis, dez e onze conectores. Depois da correção, os dois servidores passaram pelos conjuntos de 1, 3, 5, 6, 10 e 11 fontes, listagem de capacidades, enum, encerramento e reabertura da mesma raiz.

Dois testes históricos de captura/schema basal configuravam apenas um stub Apple, apesar de exigirem o schema do conjunto Apple/GitHub/Sourcegraph. O preparo desses dois casos foi corrigido para configurar o conjunto descrito. Os cenários unitários de outros testes mantiveram sua configuração. Hashes, nomes, posições, comparações profundas, capturas e evidências históricas permaneceram intactos. A suíte MCP passou com os digests congelados existentes.

## Verificação

| Verificação | Resultado | Recorte |
| --- | --- | --- |
| Testes puros do registro | Passaram | Configurações válidas, inválidas, canonicalização e admissão sem Discover |
| Testes reais de composição | Passaram | Ambos os servidores, seis perfis, enum/capacidades, shutdown/reabertura e rejeição antes de criar recursos |
| Startup do binário | Passou | MCP padrão com cinco fontes e padrão/Bitbucket com seis, onze ferramentas e enum correspondente |
| `go test ./...` final | Passou | Produto e gates ativos alcançados pelo módulo raiz; alguns pacotes em cache |
| `go vet ./...` | Passou | Inclui o ajuste final das fixtures |
| Build de `./04_wiring` | Passou | `/tmp/bateia-arquitetura-0104` |
| `crystalline-lint --fail-on warning .` | Passou | Sem violações |
| Links, fronteiras e whitespace | Passaram no recorte | Nenhum executor em L0, import de Lab no produto ou mudança em evidências congeladas |

A primeira execução da suíte detectou as duas fixtures basais inconsistentes. Após o ajuste de preparo, a suíte completa foi repetida e aprovada. Testes com servidores locais usaram loopback permitido e `commit.gpgsign=false` somente no ambiente do processo; a configuração persistente não foi alterada.

A [bancada](../../lab/bancada/arquitetura/0104/README.md) conserva os resultados do startup do binário. As duas execuções usaram diretório/home/cache/raiz temporários e ambiente sem credenciais. Enviaram apenas tools/list e list_connectors e encerraram com código zero. Nenhuma pesquisa remota foi solicitada.

## Limites

A validação do registro não comprova disponibilidade, cobertura ou qualidade das fontes. O startup foi verificado em configuração isolada, sem avaliar credenciais ou o ambiente real do usuário. Não foram alteradas políticas de evidência, reúso ou insights. O módulo separado de gates selados não foi reproduzido.

A hierarquia A0–A2 permanece um contrato em experimentação, não uma avaliação concluída de produtividade da IA. O próximo corte pode isolar resolução de fontes/publicação do runtime ou tornar a investigação de insights independente de decisões de reúso, cada um com contrato específico.
