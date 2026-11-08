# Escrever prompts que governem a materialização

Um prompt deve permitir que outra pessoa implemente o comportamento e construa um oráculo de teste sem deduzi-los do corpo da implementação. O caminho espelhado identifica o proprietário; a suficiência do contrato depende das regras escritas. A presença do arquivo e o lint não comprovam essa suficiência.

## Antes de implementar

Localize capacidade e camada no [mapa](../arquitetura/README.md). Leia [.agentrules](../../.agentrules), o ADR aplicável, o contrato compartilhado e os consumidores. Escreva o comportamento pretendido e seus limites antes de materializar a mudança. Se houver ambiguidade entre documentos, explicite e resolva a decisão em um passo; não escolha silenciosamente a versão conveniente.

## Conteúdo mínimo do prompt proprietário

| Parte | Pergunta que deve responder |
| --- | --- |
| Responsabilidade | Que regra ou estado este componente possui, em qual camada? |
| Interface | Quais entradas recebe, o que devolve e quem consome? Quais campos são opcionais? |
| Comportamento | Quais pré-condições, transições, ordem e invariantes são obrigatórios? |
| Limites e erros | Quais quotas, rejeições, efeitos permitidos e garantias de imutabilidade existem? |
| Dependências | Qual contrato compartilhado governa a regra e qual porta realiza cada efeito? |
| Aceitação | Para quais entradas e controles deve ocorrer qual resultado observável? |
| Origem | É intenção nova, consolidação de contratos existentes ou descrição derivada do código? |

Números e regras transversais têm um proprietário declarativo. Os prompts locais referenciam a seção exata desse contrato e acrescentam obrigações específicas do componente. Evite tanto copiar a mesma política em dezenas de arquivos quanto reduzir o prompt a “conforme o passo anterior”. O leitor deve conseguir chegar à regra por um link explícito.

## Testes como consumidores de um contrato

Defina pré-condição, ação, resultado e controle negativo. Inclua ausência versus zero, limites inclusivos/exclusivos, adulteração, falha parcial, cancelamento e preservação de estado quando forem relevantes. Atribua IDs estáveis aos critérios para relacionar prompt, teste e evidência. Nomes de testes existentes ajudam a localizar cobertura, mas não são o oráculo permanente.

Um teste que apenas repete os passos internos não demonstra a regra. Use portas controladas e conte efeitos quando a garantia for “rejeitar antes de ler”. Diferencie forma serializada, comportamento operacional e utilidade investigativa: uma fila estável ou um README lido não prova novidade ou benefício.

## Consolidação de acervo existente

Declare a origem posterior ao código quando houver reconstrução. Preserve o texto anterior e seus passos; organize primeiro o contrato vigente, depois a história. Não reescreva evidência selada ou datas de linhagem para simular causalidade anterior. Uma exigência de preservação de corpos de testes durante refatoração continua sendo histórica, sem impedir uma futura evolução contratada.

Registre separadamente obrigação, comportamento observado e lacuna. Uma divergência no código não deve virar requisito apenas por estar implementada. Sua correção demanda uma fatia própria. O [contrato de exploração](../contratos/exploracao-de-vizinhancas-v1.md) e sua [matriz](../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) exemplificam essa separação.

## Verificação e revisão

Confirme que os links e critérios existem, confronte contrato com consumidores e execute o gate canônico conforme o [guia de engenharia](engenharia.md). Ferramentas podem verificar estrutura, referências e preservação de arquivos; revisão semântica continua necessária. Registre resultados brutos no Lab e mantenha L0 declarativo. Só declare conformidade no recorte e nos critérios efetivamente verificados.
