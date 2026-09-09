# Passo 0059 — Preservar o canal documental sob limite de resultados

## Contrato congelado
A seleção atual pode ocupar todo o limite com código e chamadas antes de considerar documentação. Corrigir a competição por vagas; não promover comentário a evidência de comportamento.

Antes de expandir chamadas, calcular reserva documental entre candidatos documentais únicos que não possuem a mesma identidade (path, line, symbol) em nenhum hit lexical admitido. Ordenação documental existente permanece score desc/path/line/symbol. Reserva nominal: max(1, limit/4), limitada à quantidade desses candidatos e às vagas após todas as declarações exatas (score>=100000). Havendo hits de código, manter pelo menos uma vaga para código. Sem candidatos documentais novos, reserva zero. Sem hits de código, permitir preencher todo limite com documentação como antes.

Aplicar relações e corte ao orçamento de código limit-reserva; somente depois preencher vagas documentais pelo procedimento existente. Exatos não são desalojados pela reserva. Preservar ordem relativa, scores, campos, identidade e provenance de cada canal. Limite total nunca aumenta. Omissões de código/documentação continuam registradas e status partial. Havendo reserva efetiva, acrescentar Reason informando reserva documental limitada, sem comprovação de comportamento. Não criar chamadas a partir de documentação e não manter chamadas órfãs de sementes removidas pelo novo orçamento.

Não modificar consulta, catálogo, pesos, AST, validação de comportamento ou métricas do corpus. EvidenceVersion=evidence-v21-documentation-budget; checkpoint implementation_search=0059 local/live, discovery_query inalterado. Documentação continua retrieval=documentation, restrições não verificadas e comportamento not_evaluated.

## Escopo e regime
Tekt A/B sem atestação física; A lê contrato e produção, B contrato/tipos/testes legados sem produção candidata. A modifica documentation_retrieval.go, implementation_search.go, evidence.go, idea_recommend.go e respectivos prompts antes do código. B cria testes independentes, congela antes de execução. Integrador executa gates e corpus0053 congelado. Nenhuma alteração de testes legados sem registrar conflito. Não há selo de refinamento.

Validar saturação, dedup anterior ao corte, reserva no limite1/2/8, exatos, ausência de documentação, documentação isolada, determinismo, limite e provenance. Validar que comentários continuam sem fatos de comportamento. Comparar por caso e canal com0058, inclusive perdas em hits de código. Casos conhecidos são calibração, não holdout. A fração de um quarto é uma política de orçamento proposta, não um ótimo demonstrado.

Esclarecimento anterior ao candidato: o conjunto documental passado ao preenchimento também exclui identidades presentes no conjunto lexical completo anterior ao corte. Uma declaração lexical desalojada não pode reaparecer como Doc e consumir a reserva dos candidatos exclusivos.
