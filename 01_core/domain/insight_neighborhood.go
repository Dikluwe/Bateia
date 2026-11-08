/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_neighborhood.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import "time"

type InsightNeighborhoodRequest struct {
	Input          InsightSiteReport
	Memory         *InsightExplorationMemory
	Focus          string
	AllTrails      bool
	FollowLinks    *int
	MaxPassages    int
	TotalTimeout   time.Duration
	AttemptTimeout time.Duration
}

type InsightNeighborhood struct {
	InputReportID    string                        `json:"input_report_id"`
	Focus            *ProductIdentity              `json:"focus,omitempty"`
	FrontierEntryIDs []string                      `json:"frontier_entry_ids"`
	SelectionPolicy  string                        `json:"selection_policy,omitempty"`
	Neighbors        []InsightNeighborhoodNeighbor `json:"neighbors,omitempty"`
}

// A neighbor is a pre-read selection snapshot, not a current evidence verdict.
type InsightNeighborhoodNeighbor struct {
	EntryID         string          `json:"entry_id"`
	Identity        ProductIdentity `json:"identity"`
	Distance        *int            `json:"distance,omitempty"`
	ReadAttempts    int             `json:"read_attempts"`
	SelectionStatus string          `json:"selection_status"`
	RelatedEntryID  string          `json:"related_entry_id,omitempty"`
}
