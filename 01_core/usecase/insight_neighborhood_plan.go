/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_plan.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
)

type insightNeighborhoodAdmission struct {
	Request domain.InsightNeighborhoodRequest
	Origins []insightLeadOrigin
}

// ValidateInsightNeighborhood admits declared history, without authenticating it.
func ValidateInsightNeighborhood(request domain.InsightNeighborhoodRequest) (domain.InsightNeighborhoodRequest, error) {
	admission, err := admitInsightNeighborhood(request)
	return admission.Request, err
}

func admitInsightNeighborhood(request domain.InsightNeighborhoodRequest) (insightNeighborhoodAdmission, error) {
	raw, err := json.Marshal(request.Input)
	if err != nil || len(raw) > 4<<20 {
		return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
	}
	var input domain.InsightSiteReport
	if json.Unmarshal(raw, &input) != nil || input.SchemaVersion != "insight-site-discovery-v1" || input.ExplorationMemory == nil {
		return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
	}
	id := input.ID
	input.ID = ""
	if id != "site_discovery_"+insightHash(input) {
		return insightNeighborhoodAdmission{Request: request}, domain.ErrIntegrity
	}
	input.ID = id
	original, err := ValidateInsightExplorationMemory(input.ExplorationMemory)
	if err != nil {
		return insightNeighborhoodAdmission{Request: request}, err
	}
	input.ExplorationMemory = &original
	if err := validateInsightNeighborhoodReceipts(input); err != nil {
		return insightNeighborhoodAdmission{Request: request}, err
	}
	origins, err := insightNeighborhoodOrigins(input)
	if err != nil {
		return insightNeighborhoodAdmission{Request: request}, err
	}
	memory := original
	if request.Memory != nil {
		memory, err = ValidateInsightExplorationMemory(request.Memory)
		if err != nil {
			return insightNeighborhoodAdmission{Request: request}, err
		}
		if memory.Generation < original.Generation {
			return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
		}
		for _, entry := range original.Entries {
			if insightMemoryFind(&memory, entry.Identity) == nil {
				return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
			}
		}
	}
	if request.AllTrails && request.Focus != "" {
		return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
	}
	if request.Focus == "" && !request.AllTrails && input.Neighborhood != nil && input.Neighborhood.Focus != nil {
		request.Focus = input.Neighborhood.Focus.CanonicalURL
	}
	if request.Focus != "" {
		found := false
		identity, _ := insightLeadTarget(request.Focus)
		for _, entry := range memory.Entries {
			if entry.ID == request.Focus || entry.Identity == identity {
				if entry.State == "DISCARDED" {
					return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
				}
				target := insightNeighborhoodAliasTarget(&memory, entry.Identity)
				if canonical := insightMemoryFind(&memory, target); canonical != nil && canonical.State == "DISCARDED" {
					return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
				}
				request.Focus, found = target.CanonicalURL, true
				break
			}
		}
		if !found {
			return insightNeighborhoodAdmission{Request: request}, domain.ErrInvalidRunInput
		}
	}
	if err := validateInsightMemoryAdvance(memory); err != nil {
		return insightNeighborhoodAdmission{Request: request}, err
	}
	request.Input, request.Memory = input, &memory
	options := neighborhoodSiteRequest(request)
	options, err = ValidateInsightSiteOptions(options)
	if err != nil {
		return insightNeighborhoodAdmission{Request: request}, err
	}
	request.Memory, request.FollowLinks = options.Memory, options.FollowLinks
	request.MaxPassages = options.Discovery.MaxPassages
	request.TotalTimeout, request.AttemptTimeout = options.Discovery.TotalTimeout, options.Discovery.AttemptTimeout
	return insightNeighborhoodAdmission{Request: request, Origins: origins}, nil
}

func neighborhoodSiteRequest(request domain.InsightNeighborhoodRequest) domain.InsightSiteRequest {
	options := request.Input.Plan.Request
	options.Memory = request.Memory
	if request.FollowLinks != nil {
		options.FollowLinks = request.FollowLinks
	}
	if request.MaxPassages != 0 {
		options.Discovery.MaxPassages = request.MaxPassages
	}
	if request.TotalTimeout != 0 {
		options.Discovery.TotalTimeout = request.TotalTimeout
	}
	if request.AttemptTimeout != 0 {
		options.Discovery.AttemptTimeout = request.AttemptTimeout
	}
	return options
}
