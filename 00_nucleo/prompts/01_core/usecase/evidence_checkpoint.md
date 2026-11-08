# evidence_checkpoint

**Camada:** L1. **Contrato:** passo 0130 e ADR-0013.

Expor EvidenceCheckpointCriteria e CanReuseEvidenceCheckpoint(previous, criteria, reader). EvidenceCheckpointReader.ArtifactSHA256 recebe somente caminho lógico. Política conforme passo 0130: estados concluídos vigentes, schema, hashes, presença opt-in, oito versões, primeira ocorrência não vazia dos dois artefatos obrigatórios, leitura ordenada e short-circuit. Sem os/filepath/tempo/rede. Não consultar artefatos arbitrários, não alterar entradas ou comparar policy/target para invalidar evidência.

## Diagnóstico — 0132

Implementar o contrato de `0132-explicar-retomada-de-evidencias.md`: `EvaluateEvidenceCheckpoint` retorna avaliação tipada com Reusable, Reason e Subject. O predicado legado delega à avaliação, sem segunda implementação. Motivos estáveis distinguem status, schema, entrada, catálogo, modelo, versão, leitor, artefato ilegível/divergente/não registrado e elegibilidade. MANIFEST_UNAVAILABLE é atribuído pela jornada ao falhar a leitura do manifesto. Subject contém somente campo/artefato lógico constante; não incluir erros ou valores recebidos. Preservar ordem, short-circuit e elegibilidade existentes. Ausência de dois artefatos escolhe discovery-snapshot.json primeiro, deterministicamente.
