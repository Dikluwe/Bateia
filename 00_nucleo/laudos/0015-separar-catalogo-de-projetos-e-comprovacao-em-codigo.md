# Laudo 0015 — Catálogo de projetos separado da comprovação em código

**Estado**: `EM VALIDAÇÃO`  
**Data**: 2026-08-26  
**Regime**: protocolo completo Tekt, executado sem atestação de isolamento

## Materializado

- papéis nominais `PROJECT_CATALOG` e `CODE_PROOF`;
- recusa de provedor com papel incompatível antes da rede;
- porta `ProjectCatalog` independente de HTTP, GitHub CLI e Sourcegraph;
- GitHub como catálogo de repositórios por nome, descrição e README;
- resolução de todos os candidatos para commits imutáveis;
- preservação de campos casados, discriminadores e origens por query;
- deduplicação canônica por `owner/repo`;
- peneira determinística com `ADVANCE`, `REVIEW` e `REJECT`;
- estrelas isoladas e suporte multi-query sem discriminator não avançam;
- plano de prova gerado somente para candidatos `ADVANCE`;
- `repository_scope` propagado até o Sourcegraph;
- todos os fallbacks de prova preservam `repo:^owner/repo$`;
- resultados retornados fora do scope são descartados pelo núcleo;
- CLI `idea catalog discover` e `idea catalog screen --proof-plan`;
- manifesto, selo, receipts e certificado reproduzíveis.

O status do arquivo congelado do Passo 0015 permanece literalmente `PRONTO PARA
IMPLEMENTAÇÃO` porque alterá-lo depois do selo invalidaria a cadeia. O índice e
este laudo registram o estado operacional atual.

## Ataques

O contrato rejeitou 3/3 mutações previstas, mutation score 1,0:

- GitHub promovido a `CODE_PROOF`;
- Sourcegraph promovido a `PROJECT_CATALOG`;
- identidade opaca promovida implicitamente a sucesso.

O primeiro run público expôs outra mutação: o adaptador considerava qualquer
token da query como match de campo. Isso promoveu 24/34 candidatos. O snapshot
foi preservado com sufixo `invalid-token-match`, o matcher passou a exigir o
discriminator completo e uma regressão rejeita token genérico isolado.

## Braço C público

Política congelada: 12 queries ontológicas, `max-per-query=3`, concorrência 2,
timeout total de três minutos e attempt timeout de vinte segundos.

| Observável | Resultado |
|---|---:|
| queries concluídas | 12/12 |
| candidatos únicos | 34 |
| commits resolvidos | 34/34 |
| sementes repo-backed | 0/4 |
| recall competitivo | 0% |
| `ADVANCE` | 0 |
| `REVIEW` | 34 |
| queries de prova | 0 |

GitHub resolveu identidade muito melhor que Sourcegraph, mas a busca pública por
essas queries conceituais não recuperou Replen, OSS Compass, ARIS ou Starboard.
Como nenhum candidato tinha correspondência discriminatória defensável, a
peneira não chamou Sourcegraph. Isso comprova a separação de papéis e evita
gasto de prova em ruído.

## Veredito

O fragmento arquitetural está implementado e preservou `Unknown/REVIEW`. O gate
de produto falhou: recall 0% é inferior ao piso de 40% que o passo definiu para
continuar com GitHub como catálogo único.

O próximo passo deve comparar outro catálogo público de produtos ou uma busca
híbrida de metadados/README, usando exatamente a mesma ontologia e registry. Não
é permitido afrouxar discriminadores, promover reviews ou adaptar queries às
sementes retrospectivamente.

A execução é descrita como **executada sem atestação de isolamento**: houve
segregação causal por manifesto, selo e hashes, mas contrato, implementação e
veredito compartilharam a mesma autoridade conversacional.
