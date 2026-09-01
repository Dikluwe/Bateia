# Software alternative domain

Modele um alvo conhecido, filtros e um relatório determinístico de alternativas.
Separe alternativas elegíveis de excluídos; preserve identidade imutável,
cobertura, capacidades, lacunas, licença, manutenção, operação, bloqueios,
razões e evidências. Estados: RECOMMENDED, PARTIAL, REFERENCE, BLOCKED,
INCONCLUSIVE e NOT_ALTERNATIVE. Popularidade nunca substitui evidência.

API: SoftwareAlternativeTarget{Name, *RepositoryIdentity};
SoftwareAlternativeFilters{AllowedLicenses, RequireLocal, AllowArchived,
RecommendedMinCoverage, PartialMinCoverage}; SoftwareAlternative{Repository,
URL, Description, Stars, Fit, Coverage, ProblemAlignment, Capabilities, Gaps,
License, Maintenance, Operation, Reasons, Blockers, Evidence}; e
SoftwareAlternativeReport{SchemaVersion, Target, Filters, TopChoice,
Alternatives, Excluded}. Use tags JSON snake_case.
