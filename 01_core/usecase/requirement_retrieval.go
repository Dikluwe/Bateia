/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/requirement_retrieval.md
 * @layer L1
 * @updated 2026-09-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// ValidateImplementationRequirements discards stale projections before rebuilding
// them exclusively from the hits whose extracted bytes have been confirmed.
func ValidateImplementationRequirements(query string, report domain.ImplementationSearchReport, files []domain.ExtractedFile) domain.ImplementationSearchReport {
	out := ValidateImplementationHits(report, files)
	out.Behavior = parser.InspectImplementationBehavior(out, files)
	out.QueryInterpretation = parser.InterpretImplementationQuery(query)
	out.Requirements = parser.DescribeImplementationRequirements(query, out)
	out.ActionAssessments = parser.AssessImplementationActions(query, out, files)
	return out
}
