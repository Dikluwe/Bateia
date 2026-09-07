# Passo 0040 — Corrigir autoridade da evidência na jornada de reúso

Estado: CONCLUÍDO — correções implementadas e regressões verificadas.

## Contrato

Corrigir os achados do passo 0039 em RecommendationService/RunIdeaPipeline, preservando a finalidade de decidir reúso entre projetos e linguagens.

1. Correspondência lexical, símbolo localizado e teste localizado não demonstram comportamento. O coletor heurístico pode produzir PARTIAL, nunca PROVEN por acumulação desses sinais. Registrar limite explicitamente e versionar a mudança. Não inventar um verificador comportamental nem executar código de candidatos.
2. Licença exige caminho de autoridade no escopo raiz do projeto; conteúdo de testes, exemplos e dependências não autoriza licença global.
3. Candidato local fornecido pelo usuário não recebe descrição construída da ideia.
4. Negação explícita de não objetivo não vira conflito. Menções ambíguas ficam não verificadas; manter trechos e linhas reais.
5. Priorizar código de produção pertinente sobre temporários, gerados e artefatos de validação; permitir alvo explícito desses diretórios. Não prometer análise semântica onde só existe seleção lexical.
6. INCONCLUSIVE deve separar capacidades não comprovadas de condições de adoção desconhecidas e indicar evidências faltantes por candidato.

## Validação

Testes públicos de regressão com contraexemplos do 0039 e controles positivos de fonte, localização e decisão. Atualizar expectativas públicas que antes exigiam PROVEN por nomes, registrando a substituição desse contrato. Não modificar oráculos protegidos. Reexecutar candidatos sintéticos e autoanálise; preservar dossiês históricos. Executar crystalline-lint, go test ./... e build.

## Limites

Este incremento reduz confiança indevida; não implementa prova comportamental geral, descoberta pública avaliada ou cobertura de todas as capacidades. Regra ausente permanece explícita. Não é uma avaliação independente.

## Resultado

Implementadas as seis frentes do contrato, incluindo remoção das consultas do cálculo de cobertura preliminar e registro SKIPPED_BUDGET para arquivos não examinados. Cinco recomendações locais, replay e controles serviço/MCP verificados. Suíte completa, linteiro e compilação do pacote aprovados.

A autoanálise continua com cobertura insuficiente, agora explicitada; não foi implementada comprovação comportamental geral. Veja o [relatório de mudanças e limites](../relatorios/correcoes-jornada-principal-0040.md) e os [artefatos](../benchmark/0040/README.md).
