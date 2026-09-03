# Idea run orchestrator

Orquestre os casos de uso existentes, sem duplicar regras, nos modos live, local
e replay. Valide configuração antes de chamar dependências. Replay jamais chama
provedores. Falhas parciais são preservadas e INCONCLUSIVE nunca é promovido a
COMPLETE. Ordene candidatos locais e não execute código de terceiros.
Entre candidatos e etapas externas, respeite cancelamento já observado e não
inicie novas chamadas a provedores.
