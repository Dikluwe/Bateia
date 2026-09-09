# 0044 — Porta e evidência de recuperação de implementação
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md (contrato congelado).
Tipos de domínio para orçamento, relatório, inputs hashados e hits de AST; porta opcional de recuperação sem IO concreto em L1. Busca pode ser partial/not_evaluated; hits não são comprovação de capacidade. Preservar hashes da varredura e conferir identidade antes de exibir definição. Campos opcionais em resultados de extração/bundles mantêm compatibilidade com consumidores existentes.

Passo0045: adicionar ImplementationReference com Path, Symbol, Line, SHA256 e JSON
snake_case; Line é linha da chamada no arquivo chamador. ImplementationHit acrescenta
Retrieval (lexical/direct_call) e RelatedTo []ImplementationReference, opcionais.
Relação descreve chamada AST limitada, nunca execução dinâmica ou capacidade provada.

0046: adicionar RequirementRetrieval com Requirement,Status,Reason strings; MatchedTerms,MissingTerms []string e Sources []ImplementationReference, JSON snake_case. ImplementationSearchReport.Requirements []RequirementRetrieval omitempty. Fonte significa declaração associada por termos, não capacidade comprovada.

0047: ImplementationTermWeight expõe Term, Field string, DocumentFrequency,
CorpusDocuments int e Weight float64, JSON snake_case/omitempty. ImplementationHit
adiciona TermWeights []ImplementationTermWeight opcional. Weight é contribuição
finita efetivamente usada na soma lexical; não probabilidade. Chamadas sem termos
não inventam contribuição lexical.

0048: ImplementationQuery contém Original,Mode strings; Terms,Symbols []string;
Constraints []ImplementationQueryConstraint. Constraint contém Marker,Text,Status
strings. Todos com JSON snake_case; constraints opcionais. Relatório acrescenta
QueryInterpretation []ImplementationQuery omitempty; RequirementRetrieval acrescenta
Constraints []ImplementationQueryConstraint omitempty. Restrição permanece
not_verified, distinta de indícios lexicais evidence_located.

0051: ImplementationHit acrescenta DocumentationLine e DocumentationEndLine inteiros, JSON documentation_line/documentation_end_line omitempty, intervalo do grupo Doc AST. RequirementRetrieval.DocumentationSources []ImplementationReference JSON documentation_sources,omitempty separa documentação de Sources lexicais. Status documentation_located para apenas documentação; nunca comprova comportamento.

0052: ImplementationQuery acrescenta Expansions []ImplementationQueryExpansion,
JSON expansions,omitempty. Cada recibo contém Source string, Terms []string,
RuleID string, Kind string e Status string, JSON snake_case. Source conserva o
gatilho natural normalizado; Kind translation/mechanism e Status hypothesis
descrevem derivação lexical, nunca equivalência de comportamento.


## Passo 0057
ImplementationSearchReport adiciona Behavior []ImplementationBehavior com json behavior,omitempty; projeção observacional independente do ranking.


## Passo 0061
Adicionar ActionAssessments []RequirementActionAssessment com JSON action_assessments,omitempty sem alterar outros campos.
