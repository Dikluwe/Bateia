# Laudo 0017 — Consultas por frames funcionais

**Veredito**: `INCONCLUSIVE`  
**Estado do passo**: `EM VALIDAÇÃO`

## Resultado

Foi materializado, sem LLM no pipeline, o caminho determinístico entre a
ontologia do Passo 0014 e consultas funcionais auditáveis:

- seis frames funcionais;
- doze consultas geradas por templates congelados;
- cobertura de slots obrigatórios de 100%;
- doze consultas aceitas e nenhuma ocorrência do registry;
- contrato, oráculo, léxico e templates selados antes do experimento;
- mutation score 1,0 para as obrigações offline.

O experimento externo executou doze consultas no braço nominal e doze no braço
funcional, no mesmo provedor e com o mesmo limite declarado. O braço funcional
encontrou GitMVP e também pistas novas como IdeaProof e Reposeek. O nominal
encontrou alternativas genéricas, como `udit-001/repo-discovery`, mas nenhuma
das quatro sementes do gate.

## Veredito causal

Há sinal direcional, mas não aprovação: F obteve 1/4 de recall de produto contra
0/4 de N. O ganho observado foi 25 pontos percentuais, abaixo dos 30 exigidos, e
nenhuma semente repo-backed foi resolvida. O resultado agregado do provedor não
preserva rank e snippet por consulta em formato reexecutável; por isso precisão
em 10 e conversão em prova permanecem `null`.

A frase rígida `tool that evaluate` também mostrou uma falha própria: o template
concatena slots, mas não flexiona verbos. Isso trouxe literatura acadêmica e
homônimos. Corrigi-la agora, depois de ver os resultados, contaminaria o
experimento congelado. Ela deve entrar como nova hipótese versionada.

## Limites

O Bateia não usa LLM, embeddings ou reranking neural nessa implementação. O
índice web é externo e opaco; ele pode empregar métodos semânticos próprios.
Logo, somente a formulação e a auditoria das queries pertencem ao Bateia.

O protocolo Tekt foi executado com contrato, oráculo, ataque e veredito
separados em artefatos, porém sem atestação de isolamento: a mesma autoridade
conversacional realizou os papéis.

## Próxima hipótese segura

Versionar um Passo 0018 para corrigir morfologia e diversificar templates por
intenção, acrescentar uma porta de busca com resposta bruta reexecutável e só
então repetir o A/B. O novo léxico deve ser construído a partir das classes de
erro deste laudo, nunca dos nomes das sementes que faltaram.
