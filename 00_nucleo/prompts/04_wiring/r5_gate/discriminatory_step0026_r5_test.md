# Prompt espelhado — mutation gate substituto R5 do Passo 0026

Materialize, em pacote isolado pertencente ao oráculo, as portas públicas congeladas do
Passo 0026 e uma implementação substituta compilável de referência. Materialize cada
mutação negativa como uma alteração executável dessa mesma implementação ou como um
fault point real selecionado na construção do substituto. O pacote não pode importar,
ler ou executar a implementação candidata.

Para cada M01–M60, aplique o estímulo executável correspondente pela mesma interface
black-box à referência e ao mutante. O veredito `Violated` deve depender exclusivamente
de diferença entre observáveis normalizados; igualdade é sobrevivência. O ID serve
somente para construir o mutante e escolher o estímulo, nunca para produzir o resultado.
Vincule `mutation`, `stimulus` e `witness` do caso ao digest da evidência, sem aceitar
vocabulário, frases ou bindings declarativos como prova.

Use concorrência real para Start/CAS/cancelamento/recovery, diretórios temporários e
OpenFile exclusivo/Write/Sync/Rename/sync de diretório para persistência, fault points
reais para os prefixos de crash, corrupção e links físicos reais, HMAC real para cursor,
e parser AST sobre fixtures Go compiláveis para as violações de camada. Para legado,
execute por lista explícita somente o servidor predecessor certificado e invoque os
goldens observáveis de search, benchmark e recommend. Os controles P01–P05 também devem
executar portas black-box; O01–O05 permanecem `Unknown` com construção ambiental.

Execute os 70 casos em três ordens determinísticas, três vezes por ordem, exigindo
veredito e digest idênticos nas nove observações. O gate autoritativo deve alcançar
60/60 mutações mortas sem compilar candidato; se um observável não diferir, reporte o
caso como sobrevivente.


## Adendo 0098 — Execução do legado após extração do transporte

O recorte antigo de quatro arquivos MCP deixou de compilar depois da extração
anterior dos helpers de dispatch/segurança. Executar o pacote completo
./02_shell/mcp, selecionando explicitamente os mesmos goldens R3 e certificados,
mais o teste real TestStep0025LegacySuccessfulSearchUsesEvidenceContract
(sucessor já existente do nome MatchesFrozenBaselineBytes). Não alterar hashes,
goldens, expectativas, mutantes ou classificação dos casos. Essa correção é de
invocação: recompõe helpers ausentes e impede seletor legado vazio; não credita
qualidade do produto ou nova prova aos testes de infraestrutura.
