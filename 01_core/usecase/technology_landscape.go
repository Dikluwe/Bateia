/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technology_landscape.md
 * @layer L1
 * @updated 2026-09-24
 */
package usecase

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"bateia/01_core/domain"
)

// AnalyzeTechnologyLandscape classifies declared evidence without fetching or
// asserting that an opaque source reference proves its claim.
func AnalyzeTechnologyLandscape(input domain.TechnologyLandscapeRequest) (domain.TechnologyLandscapeReport, error) {
	report := domain.TechnologyLandscapeReport{
		SchemaVersion: domain.TechnologyLandscapeSchemaVersion,
		AsOf:          input.Current.AsOf,
		Need:          append([]string(nil), input.Current.Need...),
		Assessments:   []domain.LandscapeAssessment{},
		Changes:       []domain.LandscapeChange{},
		Counts: map[domain.LandscapeRole]int{
			domain.LandscapeSubstitute: 0, domain.LandscapeComplement: 0,
			domain.LandscapeAnalogue: 0, domain.LandscapeUnresolved: 0,
		},
		Caveats: []string{"referências de fonte são declaradas, não verificadas materialmente", "movimentos de sinais são descritivos, não notas de maturidade"},
	}
	if input.SchemaVersion != domain.TechnologyLandscapeSchemaVersion {
		return report, fmt.Errorf("landscape schema_version deve ser %s", domain.TechnologyLandscapeSchemaVersion)
	}
	if err := validateLandscapeSnapshot(input.Current); err != nil {
		return report, fmt.Errorf("current: %w", err)
	}
	if len(input.Current.Candidates) == 0 {
		report.Caveats = append(report.Caveats, "nenhum candidato fornecido; isso não comprova ausência de alternativas")
	}
	if input.Previous != nil {
		if err := validateLandscapeSnapshot(*input.Previous); err != nil {
			return report, fmt.Errorf("previous: %w", err)
		}
		if input.Previous.AsOf >= input.Current.AsOf {
			return report, fmt.Errorf("previous.as_of deve anteceder current.as_of")
		}
		if !sameLandscapeNeed(input.Current.Need, input.Previous.Need) {
			return report, fmt.Errorf("previous.need difere de current.need; mudanças não são comparáveis")
		}
	} else {
		report.Caveats = append(report.Caveats, "sem snapshot anterior; movimentos do ecossistema indisponíveis")
	}
	current := assessLandscapeSnapshot(input.Current)
	for _, assessment := range current {
		report.Assessments = append(report.Assessments, assessment)
		report.Counts[assessment.Role]++
	}
	if input.Previous != nil {
		report.Changes = compareLandscapeSnapshots(assessLandscapeSnapshot(*input.Previous), current)
	}
	return report, nil
}

func validateLandscapeSnapshot(snapshot domain.LandscapeSnapshot) error {
	if _, err := landscapeDate(snapshot.AsOf); err != nil {
		return fmt.Errorf("as_of: %w", err)
	}
	if len(snapshot.Need) == 0 || len(snapshot.Need) > 64 || len(snapshot.Candidates) > 500 {
		return fmt.Errorf("need deve ter 1..64 itens e candidates no máximo 500")
	}
	need := make(map[string]bool, len(snapshot.Need))
	for _, id := range snapshot.Need {
		if !landscapeID(id) || need[id] {
			return fmt.Errorf("need contém ID vazio ou duplicado: %q", id)
		}
		need[id] = true
	}
	seen := map[string]bool{}
	for _, c := range snapshot.Candidates {
		if !landscapeID(c.ID) || seen[c.ID] {
			return fmt.Errorf("candidate ID vazio ou duplicado: %q", c.ID)
		}
		seen[c.ID] = true
		if strings.TrimSpace(c.Name) == "" || (c.Kind != domain.LandscapeRepository && c.Kind != domain.LandscapeProduct && c.Kind != domain.LandscapeTechnology) {
			return fmt.Errorf("candidate %s: name/kind inválido", c.ID)
		}
		if len(c.Capabilities) > len(snapshot.Need) || len(c.Signals) > 16 {
			return fmt.Errorf("candidate %s: orçamento de claims/signals excedido", c.ID)
		}
		claims := map[string]bool{}
		for _, claim := range c.Capabilities {
			if !need[claim.ID] || claims[claim.ID] {
				return fmt.Errorf("candidate %s: capability fora de need ou duplicada: %q", c.ID, claim.ID)
			}
			claims[claim.ID] = true
			if err := validateLandscapeClaim(claim.Status, claim.Sources); err != nil {
				return fmt.Errorf("candidate %s capability %s: %w", c.ID, claim.ID, err)
			}
		}
		for _, a := range []struct {
			name string
			v    domain.LandscapeAssertion
		}{{"standalone", c.Standalone}, {"integration", c.Integration}, {"analogy", c.Analogy}} {
			if err := validateLandscapeClaim(a.v.Status, a.v.Sources); err != nil {
				return fmt.Errorf("candidate %s %s: %w", c.ID, a.name, err)
			}
		}
		signals := map[domain.LandscapeSignalKind]bool{}
		for _, signal := range c.Signals {
			if signals[signal.Kind] {
				return fmt.Errorf("candidate %s: signal duplicado: %s", c.ID, signal.Kind)
			}
			signals[signal.Kind] = true
			if err := validateLandscapeSignal(signal, snapshot.AsOf); err != nil {
				return fmt.Errorf("candidate %s signal %s: %w", c.ID, signal.Kind, err)
			}
		}
	}
	return nil
}

func landscapeDate(value string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", value)
	if err != nil || d.Format("2006-01-02") != value {
		return time.Time{}, fmt.Errorf("data deve usar YYYY-MM-DD: %q", value)
	}
	return d, nil
}

func landscapeID(id string) bool {
	return id != "" && id == strings.TrimSpace(id) && len(id) <= 160
}

func validateLandscapeClaim(status domain.EvidenceStatus, sources []string) error {
	switch status {
	case "", domain.EvidenceInconclusive:
		if len(sources) != 0 {
			return fmt.Errorf("claim inconclusivo não deve declarar fontes como prova")
		}
		return nil
	case domain.EvidenceProven, domain.EvidencePartial, domain.EvidenceAbsent, domain.EvidenceConflict:
		if len(sources) == 0 || len(sources) > 16 {
			return fmt.Errorf("claim exige 1..16 source refs")
		}
		for _, source := range sources {
			if !landscapeID(source) {
				return fmt.Errorf("source ref vazia ou inválida")
			}
		}
		return nil
	default:
		return fmt.Errorf("status inválido: %q", status)
	}
}

func validateLandscapeSignal(s domain.LandscapeSignal, asOf string) error {
	if _, err := landscapeDate(s.ObservedAt); err != nil || s.ObservedAt > asOf {
		return fmt.Errorf("observed_at inválido ou posterior ao snapshot")
	}
	if math.IsNaN(s.Value) || math.IsInf(s.Value, 0) || s.Value < 0 || len(s.Sources) == 0 || len(s.Sources) > 16 {
		return fmt.Errorf("value/sources inválidos")
	}
	for _, source := range s.Sources {
		if !landscapeID(source) {
			return fmt.Errorf("source ref vazia ou inválida")
		}
	}
	switch s.Kind {
	case domain.LandscapeReleases:
		if s.Unit != "count" || s.WindowDays <= 0 || s.Value != math.Trunc(s.Value) {
			return fmt.Errorf("RELEASES exige count inteiro e janela positiva")
		}
	case domain.LandscapeResponseHours:
		if s.Unit != "hours" || s.WindowDays <= 0 {
			return fmt.Errorf("MEDIAN_RESPONSE_HOURS exige hours e janela positiva")
		}
	case domain.LandscapeBusFactor:
		if s.Unit != "people" || s.WindowDays <= 0 || s.Value != math.Trunc(s.Value) {
			return fmt.Errorf("BUS_FACTOR exige people inteiro e janela positiva")
		}
	case domain.LandscapeArchived:
		if s.Unit != "boolean" || s.WindowDays != 0 || (s.Value != 0 && s.Value != 1) {
			return fmt.Errorf("ARCHIVED exige boolean 0/1 sem janela")
		}
	default:
		return fmt.Errorf("kind inválido")
	}
	return nil
}

func sameLandscapeNeed(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[string]bool, len(a))
	for _, id := range a {
		set[id] = true
	}
	for _, id := range b {
		if !set[id] {
			return false
		}
	}
	return true
}

func assessLandscapeSnapshot(snapshot domain.LandscapeSnapshot) []domain.LandscapeAssessment {
	out := make([]domain.LandscapeAssessment, 0, len(snapshot.Candidates))
	for _, c := range snapshot.Candidates {
		claims := map[string]domain.EvidenceStatus{}
		for _, claim := range c.Capabilities {
			claims[claim.ID] = claim.Status
		}
		a := domain.LandscapeAssessment{
			CandidateID:     c.ID,
			Name:            c.Name,
			Kind:            c.Kind,
			Role:            domain.LandscapeUnresolved,
			Capabilities:    append([]domain.LandscapeCapability{}, c.Capabilities...),
			Standalone:      landscapeReportAssertion(c.Standalone),
			Integration:     landscapeReportAssertion(c.Integration),
			Analogy:         landscapeReportAssertion(c.Analogy),
			Signals:         append([]domain.LandscapeSignal{}, c.Signals...),
			Proven:          []string{},
			Absent:          []string{},
			Unknown:         []string{},
			MissingEvidence: []string{},
		}
		for _, id := range snapshot.Need {
			switch claims[id] {
			case domain.EvidenceProven:
				a.Proven = append(a.Proven, id)
			case domain.EvidenceAbsent:
				a.Absent = append(a.Absent, id)
			default:
				a.Unknown = append(a.Unknown, id)
				a.MissingEvidence = append(a.MissingEvidence, "capability:"+id)
			}
		}
		sort.Strings(a.Proven)
		sort.Strings(a.Absent)
		sort.Strings(a.Unknown)
		if len(a.Proven) == len(snapshot.Need) && c.Standalone.Status == domain.EvidenceProven {
			a.Role = domain.LandscapeSubstitute
		} else if len(a.Proven) > 0 && len(a.Absent) > 0 && len(a.Unknown) == 0 && c.Integration.Status == domain.EvidenceProven {
			a.Role = domain.LandscapeComplement
		} else if len(a.Absent) == len(snapshot.Need) && c.Analogy.Status == domain.EvidenceProven {
			a.Role = domain.LandscapeAnalogue
		}
		if a.Role == domain.LandscapeUnresolved {
			if len(a.Proven) == len(snapshot.Need) && c.Standalone.Status != domain.EvidenceProven {
				a.MissingEvidence = append(a.MissingEvidence, "standalone_operation")
			}
			if len(a.Proven) > 0 && len(a.Absent) > 0 && c.Integration.Status != domain.EvidenceProven {
				a.MissingEvidence = append(a.MissingEvidence, "integration")
			}
			if len(a.Absent) == len(snapshot.Need) && c.Analogy.Status != domain.EvidenceProven {
				a.MissingEvidence = append(a.MissingEvidence, "analogy")
			}
		}
		sort.Strings(a.MissingEvidence)
		sort.Slice(a.Signals, func(i, j int) bool { return a.Signals[i].Kind < a.Signals[j].Kind })
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CandidateID < out[j].CandidateID })
	return out
}

func landscapeReportAssertion(a domain.LandscapeAssertion) domain.LandscapeAssertion {
	if a.Status == "" {
		a.Status = domain.EvidenceInconclusive
	}
	a.Sources = append([]string(nil), a.Sources...)
	return a
}

func compareLandscapeSnapshots(previous, current []domain.LandscapeAssessment) []domain.LandscapeChange {
	old := map[string]domain.LandscapeAssessment{}
	now := map[string]domain.LandscapeAssessment{}
	ids := map[string]bool{}
	for _, a := range previous {
		old[a.CandidateID] = a
		ids[a.CandidateID] = true
	}
	for _, a := range current {
		now[a.CandidateID] = a
		ids[a.CandidateID] = true
	}
	ordered := make([]string, 0, len(ids))
	for id := range ids {
		ordered = append(ordered, id)
	}
	sort.Strings(ordered)
	changes := []domain.LandscapeChange{}
	for _, id := range ordered {
		before, had := old[id]
		after, has := now[id]
		switch {
		case !had:
			changes = append(changes, domain.LandscapeChange{CandidateID: id, Kind: "ADDED", AfterRole: after.Role})
		case !has:
			changes = append(changes, domain.LandscapeChange{CandidateID: id, Kind: "REMOVED", BeforeRole: before.Role})
		default:
			movements := landscapeSignalMovements(before.Signals, after.Signals)
			if before.Role != after.Role || len(movements) > 0 {
				changes = append(changes, domain.LandscapeChange{CandidateID: id, Kind: "UPDATED", BeforeRole: before.Role, AfterRole: after.Role, Signals: movements})
			}
		}
	}
	return changes
}

func landscapeSignalMovements(previous, current []domain.LandscapeSignal) []domain.LandscapeSignalMovement {
	old := map[domain.LandscapeSignalKind]domain.LandscapeSignal{}
	now := map[domain.LandscapeSignalKind]domain.LandscapeSignal{}
	kinds := map[domain.LandscapeSignalKind]bool{}
	for _, s := range previous {
		old[s.Kind] = s
		kinds[s.Kind] = true
	}
	for _, s := range current {
		now[s.Kind] = s
		kinds[s.Kind] = true
	}
	ordered := make([]domain.LandscapeSignalKind, 0, len(kinds))
	for kind := range kinds {
		ordered = append(ordered, kind)
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i] < ordered[j] })
	out := []domain.LandscapeSignalMovement{}
	for _, kind := range ordered {
		before, had := old[kind]
		after, has := now[kind]
		m := domain.LandscapeSignalMovement{Kind: kind}
		switch {
		case !had:
			m.Movement = "NEW"
			v := after.Value
			m.After = &v
		case !has:
			m.Movement = "REMOVED"
			v := before.Value
			m.Before = &v
		default:
			b, a := before.Value, after.Value
			m.Before, m.After = &b, &a
			if before.Unit != after.Unit || before.WindowDays != after.WindowDays {
				m.Movement = "INCOMPARABLE"
			} else if after.Value > before.Value {
				m.Movement = "INCREASED"
			} else if after.Value < before.Value {
				m.Movement = "DECREASED"
			} else {
				continue
			}
		}
		out = append(out, m)
	}
	return out
}
