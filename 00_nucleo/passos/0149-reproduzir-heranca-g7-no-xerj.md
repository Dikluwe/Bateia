# 0149 — Reproduzir herança de G7 e conferir capacidades do XERJ

**Data:** 2026-11-08. **Escopo ex-ante:** investigação e reprodução local, sem alteração produtiva no Bateia ou publicação no XERJ.

## Comportamento esperado

Rastrear re-pin, construção de catálogo, validação pré-publicação e geração do site do Corpus Hub. Tentar reproduzir se uma versão nova recebe a nota G7 anterior. Distinguir main, corpus-hub, release oficial, fixture e observação do serviço público. Preservar o ensaio0148; acrescentar evidência em `lab/bancada/insights/xerj-2026-11-08/g7-reproduction`.

Conferir também se as outras propostas já são oferecidas pelo XERJ: recibo dos escopos de verificação e recuperação de precondições/contraexemplos. Diferenciar busca de conteúdo existente, estrutura de armazenamento, geração de ideias e verificação funcional. Nenhuma ausência global pode ser inferida de um único tipo ou pesquisa lexical.

## Aceite e limites

- Registrar revisões e hashes dos arquivos efetivamente executados, comandos, saídas e controles negativos.
- Usar código original e binário oficial com checksum conferido quando disponível; lacunas de compilação devem ficar registradas.
- A reprodução pode usar corpus/pacote sintético, explicitamente identificado. Não publicar, abrir PR, enviar feedback ou usar chaves de produção.
- Demonstrar o caminho até os artefatos preparados para publicação, sem executar deploy. Não alegar ocorrência no site público a partir de fixture.
- Registrar conclusão por superfície, inclusive casos que contradigam a hipótese inicial.
- Executar verificação aplicável do Bateia e replay da bancada. Resultados brutos permanecem no Lab.

Este passo não autoriza executar comandos sugeridos por fontes, enviar comunicações ou interpretar uma assinatura como prova de correção funcional.

## Resultado

O usuário esclareceu que deseja testar se um agente usando XERJ chega às ideias. Um agente sem histórico/propostas anteriores consultou somente o XERJ local sobre o mesmo corpus, em sete buscas. Desenvolveu propostas semelhantes de recibo de verificação e aplicabilidade/contraexemplos e acrescentou replay; não encontrou espontaneamente a herança de G7. A [comparação](../../lab/bancada/insights/xerj-2026-11-08/xerj-agent/README.md) explicita o caráter dirigido da pergunta e a ausência de benchmark estatístico.

A [reprodução G7](../../lab/bancada/insights/xerj-2026-11-08/g7-reproduction/README.md) confirmou associação da mesma avaliação após troca de pin no catálogo da release rc.101 e na página gerada. Os validadores de publicação examinados recusam formato 2, esperado 1: esse bloqueio independente foi preservado. O catálogo público assinado possui 514 entradas e nenhuma nota G7; não foi demonstrada ocorrência pública. Fixture, assinatura, pin divergente e remoção da avaliação foram exercitados com resultados distintos.

Lint, auditorias, testes, vet, build e diff-check do Bateia passaram. O ensaio inicial0148 e seus hashes permanecem preservados; resultados novos ficam em subdiretórios próprios. Não houve deploy, comunicação externa ou envio de feedback. O nó local temporário foi encerrado.
