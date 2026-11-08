# 0145 — Ownership da investigação sob Tekt

Contrato: [passo0145](../../../../00_nucleo/passos/0145-reorganizar-fronteiras-de-investigacao-sob-tekt.md). Revisão arquitetural após a expansão/memória, sem novos recursos de produto ou consultas públicas.

O recorte separa fachada, admissão, fontes congeladas, fronteira, memória, ordenação, execução e resultado em L1. A admissão entrega origens verificadas à fachada, eliminando uma segunda interpretação ignorada. L2 possui o probe literal de --no-token; L4 compõe ambiente e cliente. Sem renomeação cosmética de camadas, mudança de runtime ou imports produtivos do Lab.

`initial-architecture.log` registra a inspeção estrutural anterior. `baseline.json` conserva hashes dos arquivos e das evidências históricas; `before/` contém fontes anteriores do recorte, como texto de bancada. `materialize.py` documenta a extração única e recusa sobrescrever módulos novos; não é ferramenta produtiva nem deve ser executado novamente sobre a árvore já reorganizada.

`audit_test.go` integra os gates ativos: verifica proprietários dos símbolos, ausência de interpretação do argumento em L4, chamadas fora do responsável e controles negativos. Também compara corpos AST normalizados dos consumidores realocados e de regras movidas, preservando seus nomes/asserts. Fontes, plano e resultado têm testes próprios; ordenação e probe de entrada possuem controles de contrato.

Comparações offline do binário0144 e do reorganizado usam relatórios reais congelados, leituras desativadas e feedback local. Resultados novos e comandos ficam nesta pasta, sem reescrever os relatórios anteriores. Replays anteriores exercitam somente seus contratos declarados; fixtures sintéticas não representam leituras públicas novas. Logs da pilha completa ficam em final/.

A auditoria cobre fronteiras verificáveis e o ownership do recorte. Não prova ausência universal de regras mal classificadas nem qualidade semântica dos insights. Schemas, hashes, ordenação e resultados devem permanecer iguais nos controles de preservação.

## Resultado verificado

final/verification.json registra lint, auditorias, testes completos, vet, build e diff aprovados. A auditoria global cobre398 arquivos produtivos; o novo controle de ownership passou com seus controles negativos. A prova AST confirmou os corpos dos consumidores realocados e das regras movidas. Regras canônicas, arquivos AGENTS e configuração estrutural permaneceram intactos.

comparison.json registra nove casos offline contra o binário0144 preservado: sete relatórios foram idênticos byte a byte (foco, herança, geral, alias, outra trilha, feedback e overrides); duas rejeições conservaram mensagens e a saída anterior. Nenhuma consulta ou leitura remota foi feita. Os comandos/saídas ficam em comparisons/. candidate-bateia e final/bateia possuem o mesmo hash.

Seis replays históricos passaram, e workflow-0142-replay.json reproduziu exatamente os dois relatórios e o feedback anteriores. Os120 arquivos de evidência0142–0144 no baseline foram preservados. A fixture0142 continua sintética; esses controles não promovem novidade ou benefício a fatos.

verification.json consolida os resultados e lineage.json liga fontes aos prompts. Revalidar sem rede: `python3 lab/bancada/arquitetura/0145/verify.py`. O gate ativo pode ser reproduzido com `make verify`; os gates selados continuam exigindo seus módulos/revisões próprios.
