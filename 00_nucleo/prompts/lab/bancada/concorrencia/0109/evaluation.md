# Prompt 0109 — evaluation

**Regime**: Lab, sem promoção para L1–L4  
**Data**: 2026-10-06

Executar corpus JSON exploratório com referência de autoria declarada, casos/famílias e expectativas de observações ou lacunas. Conferir presença e perfis completos fornecidos: observação extra é erro, não ignorar falso positivo. Reportar checks/casos/famílias, cobertura sintática e erros, sem chamar a taxa de acerto de qualidade semântica. Corpus autoral não é independente. Todos os resultados vêm da execução do observer; nunca interpretar ID como veredito.

0110: aceitar channel_type_identity e as lacunas de resolução de tipos;
conservar corpus e resultados históricos v1. Corpus adicional autoral verifica
nomes, aliases, encadeamento, múltiplos arquivos, sombras, ciclos, importação
ausente, pacote inválido, grupos distintos e tipos slice/map.

0112: aceitar lacunas de relações na validação de corpus; perfis históricos
continuam apenas property/value, com testes próprios para relations.

0113: admitir CALLBACK_VALUE_FLOW_UNRESOLVED em corpora; sintaxe/relations
históricas conservam contratos, callbacks possuem testes separados.
