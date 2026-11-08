# recommendation_dossier

Possuir a jornada do dossiê em L1: captura via porta, identidade causal, versões, retomada, pipeline/recomendação, alternativas, trace, estágio/manifesto e publicação ordenada. Preservar byte preimages, campos, erros, cleanup, sidecars e contagens. Porta de arquivos usa nomes lógicos e delega renderização JSON/Markdown à casca; nenhuma operação os/filepath/io concreta. A fachada recebe dependências compostas por L4; o caminho de API compatível usa a mesma jornada. Não duplicar regra no CLI.
Contrato 0131/ADR-0014.

A fachada também aplica o predicado semântico opt-in antes de qualquer publicação, para que o contrato não dependa exclusivamente da chamada pela CLI. Não alterar a admissão de modelos calibrados de outras jornadas.

A publicação da mensagem legível usa callback de apresentação antes do cleanup diferido, mantendo a ordem observável da CLI anterior inclusive quando Close falha. L1 não formata nem escreve stdout.

## Diagnóstico — 0132

Expor Checkpoint opcional no resultado em memória. Apenas req.Resume solicita avaliação: ler run.json uma vez; erro produz MANIFEST_UNAVAILABLE/run.json; sucesso usa EvaluateEvidenceCheckpoint uma vez. Reusable dessa avaliação escolhe a fonte. Preservar diagnóstico no retorno e no anúncio final, inclusive se uma etapa posterior falhar. Não persistir diagnóstico nem incluir em identidade/versões/artefatos; sem flag manter nil. Aplicar contrato do passo 0132 e preservar o callback de apresentação somente após publicação final.
