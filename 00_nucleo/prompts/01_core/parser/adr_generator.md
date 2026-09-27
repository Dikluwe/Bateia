# Specification: Architecture Decision Record Generator (L1 Core)

## Purpose
Synthesize domain SearchReport into a formal Architecture Decision Record (ADR) document (Context, Decision, Status, Consequences).

## Rules
- L1 package (`01_core/parser`).
- ZERO I/O (pure transformation in memory).
- Function: `GenerateADRDocument(report *domain.SearchReport) string`.
- Returns a Markdown ADR proposal organized as context, proposal, status and possible consequences; no acceptance workflow is inferred.
- Every report generated from search has status `PROPOSTO`; search relevance is not architectural or human approval.
- Candidate, relevance score/reason, URL and commit reference remain traceable. Describe the first result as a candidate for evaluation, not an adoption decision.
- Dependencies are observed declarations, not approved recommendations. Keep the existing direct-only, five-item limit.
- Benefits are explicitly hypotheses to validate. Never claim production testing, reduced technical debt or human acceptance without evidence in the report.
- With no results, report only that the report contains no extracted candidates; do not infer that all alternatives failed strict criteria or imply benefits of a nonexistent selection.
- Preserve the pure transformation and simple search mode; do not add approval states to the domain schema.
- A missing or floating (`HEAD`) commit reference must be identified as not pinned, never rendered as a verified commit. Preserve a repository reference when available.
