# Resolução limitada 0110

Auxiliar puro do observador. Receber ASTs e file set já presentes na memória.
Agrupar por diretório/nome do pacote; ordem determinística. Não carregar imports.
Qualquer erro em go/types ou arquivo malformado do grupo suprime resolução.
Usar tamanhos fixos gc/amd64, sem alegar representar build declarado.
Somente builtin make aplicado a expressão de tipo canal e chamada com resultado
canal permite atribuição; exigir ausência de erros no grupo. Slice/map não
recebem observação nem lacuna específica de canal quando checados corretamente.

0111: importer recursivo opcional de pacotes mapeados + exports fornecidos.
Não usar lookup padrão de go/importer. Cache local por execução; ciclo falha;
retornar pacote ao importador apenas após sucesso e MarkComplete do checker.
Erro de dependência invalida consumidores. Versão sem ambiente mantém contrato.

Registrar até 100 diagnósticos go/types por grupo, texto normalizado/truncado,
com posição física vinculada ao SHA do arquivo fornecido. Reportar diagnostics
sem interpretar mensagens como defeitos demonstrados do projeto; dados/build
incompletos também geram erros. Ordem determinística.

0112: preservar Defs e Selections em types.Info para indexar funções e campos
por objeto resolvido. Não inferir chamadas por texto de nomes.
