# ADR-0009 Admissão de conectores e projeção MCP

**Status**: ACEITO PARA O RECORTE 0104  
**Data**: 2026-10-06  
**Relacionado**: ADR-0008, passos 0027/0028/0034/0103/0104

## Contexto

A quantidade de conectores da primeira composição foi incorporada como validação em L4. Com novos catálogos, o registro padrão ficou válido em L1 mas rejeitado pelo construtor de síntese. O MCP anunciava um schema basal, diferente do registro efetivamente recebido, tanto no servidor de mercado quanto no de síntese.

## Decisão

1. L1 possui a admissão do registro configurado por `NewMarketDiscoveryRegistry`: ao menos um conector, capacidades válidas e IDs únicos, sem número fixo ou fornecedores obrigatórios. A admissão não consulta provedores nem comprova sua disponibilidade.
2. Preservar `NewDiscoveryConnectorRegistry` como representação genérica que admite um registro vazio. Ativação de descoberta e representação de registro são contratos distintos.
3. L4 usa o construtor configurado e conserva as assinaturas públicas, recuperação, inicialização e encerramento existentes. Nenhuma regra de seleção ou evidência é deslocada para wiring.
4. L2 projeta os IDs do registro nos schemas anunciados. A extensão de síntese compartilha a mesma projeção de mercado, sem restaurar silenciosamente o enum antigo.
5. Preservar o schema do conjunto basal Apple/GitHub/Sourcegraph e o comportamento sem registro das funções legadas. Para outros conjuntos, a projeção existente baseada no registro passa a ser efetivamente usada.
6. Manter limite de seleção de dez conectores por pedido, independentemente da quantidade registrada. Não confundir catálogo disponível com seleção de uma execução.

## Consequências

Composições válidas unitárias e ampliadas passam a iniciar sem depender da quantidade histórica de três fontes. O schema anunciado passa a corresponder às fontes configuradas. Listas de ferramentas, posições, artefatos, erros, orçamento e dispatcher permanecem os mesmos, salvo o enum/limite de seleção necessários para refletir o registro.

A apresentação continua dependendo do tipo de registro de L1 já existente. Separar esse contrato em uma porta menor é um trabalho possível, não parte deste corte. O runtime e as regras de investigação permanecem como no passo 0103.

## Verificação

Contrato [0104](../passos/0104-compor-descoberta-por-registro-validado.md): admissão pura e composição real de ambos os servidores, schemas, capabilities, rejeição pré-recursos e shutdown/reabertura. A aprovação local não demonstra funcionamento das APIs públicas remotas nem do ambiente de credenciais de produção.
