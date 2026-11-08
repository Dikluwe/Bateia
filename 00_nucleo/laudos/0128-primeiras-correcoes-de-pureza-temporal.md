# Laudo 0128 — Primeiras correções de pureza temporal

**Data:** 2026-10-28
**Contrato:** [passo 0128](../passos/0128-explicitar-instante-de-coleta-nos-manifestos.md)
**Decisão:** [ADR-0011](../ADR/ADR-0011-instante-de-coleta-explicito-nos-manifestos.md)

## Correções

Foram removidas 48 consultas ao relógio em 43 arquivos de parsers de manifestos. ParseManifest e parsers específicos conservam assinaturas e conteúdo interpretado; ConsultadoEm fica zero quando não existe metadado de coleta. ParseManifestAt recebe instante explícito e o converte a UTC, preservando zero desconhecido. Os coletores GitHub/local continuam preenchendo o instante real em L3. A reanálise do funil usa a data recebida no pacote extraído, em vez de fabricar uma data nova.

A revisão encontrou também ordem variável de dependências nos objetos JSON de package.json, composer.json, box.json, haxelib.json e dub.json. Esses parsers agora ordenam chaves dentro dos grupos já existentes. Nomes, versões, tipos e exclusões são preservados; somente a ordem anteriormente indeterminada passa a ser definida. Nenhum relatório, recibo, gate selado ou hash histórico foi reescrito.

crystalline.toml foi adaptado ao formato de tabela esperado pelo lint instalado. Nenhuma regra foi suprimida e o limiar warning foi mantido. O comando padrão consegue carregar a configuração; ainda reprova por 46 avisos V12 anteriores, sem novos diagnósticos.

## Verificação

- Corpus autoral de 49 tipos: digest do parsing, excluindo somente ConsultadoEm, congelado antes da correção e preservado depois. O corpus cobre os tipos, sem pretensão de cobrir todas as sintaxes dos ecossistemas.
- Metadado conhecido/desconhecido, conversão UTC, duas chamadas idênticas e erros: aprovados.
- Ordenação de múltiplas dependências e grupos: falhou nos cinco formatos antes da correção e passou depois; valores e exclusões preservados.
- Helper de chaves: entrada não mutada, lista independente, Unicode e mapas vazios verificados.
- Coletores GitHub/local preservam data de coleta não zero e UTC; reanálise preserva tanto data fixa quanto zero.
- `go test ./...`: aprovado. Assinatura Git desativada somente no processo das fixtures, com cache e loopback autorizados. Testes finais de manifestos também aprovados após adicionar a verificação própria do helper exigida pelo lint.
- `go vet ./...`: aprovado.
- `go build -o /tmp/bateia-arquitetura-0128 ./04_wiring`: aprovado, usando destino temporário autorizado.
- `git diff --check`: aprovado.
- `crystalline-lint --fail-on warning .`: executa, com zero erros e 46 avisos anteriores; gate ainda não aprovado.

Registros: [bancada 0128](../../lab/bancada/arquitetura/0128/README.md) e [verificação](../../lab/bancada/arquitetura/0128/verification.json).

## Compatibilidade e pendências

Chamadores diretos que esperavam uma hora gerada pelo parser precisam usar ParseManifestAt com o metadado de coleta. Chamadores sem essa informação, incluindo caminhos legados de parsing, recebem zero desconhecido. O esquema JSON permanece igual. A nova ordem de dependências é estável; não se afirma equivalência byte a byte com a ordem aleatória anterior para todos os inputs.

Relógio ainda existe em funnel.go, idea_run.go e lab_extractor.go. Deadlines em L1, persistência e captura local em L2, mistura de responsabilidades do runtime, localização de testes de adaptadores em L4 e divergência de orientação documental continuam fora deste recorte. Os 46 avisos exigem análise individual; não são apresentados como 46 violações comprovadas. Gates selados possuem módulo próprio e não foram reproduzidos. As alterações locais anteriores de 0126/0127 e artefatos de outras investigações foram preservados.
