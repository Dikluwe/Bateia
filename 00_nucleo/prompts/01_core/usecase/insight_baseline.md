# Normalização de baseline de conteúdo

**Camada:** L1. **Contrato:** passo 0135.

Implementar validação/canonicalização de repository, arquivos, hashes e digest de manifesto canônico do snapshot conforme contrato. Aplicar quotas de arquivos/conteúdo e UTF-8 válido. Reutilizar validação/projeção de observações de 0134 sobre adaptação interna com digest do snapshot; expor apenas identidade de snapshot, namespace target_obs_ e escopo SUPPLIED_FILES_ONLY. Ordenação determinística e cópias profundas, sem mutar entradas nem autenticar origem. Baseline nil conserva desconhecido. Sem I/O.
