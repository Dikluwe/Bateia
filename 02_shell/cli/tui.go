/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/tui.md
 * @layer L2
 * @updated 2026-08-04
 */

package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tuiState int

const (
	stateInput tuiState = iota
	stateSearching
	stateResults
)

type searchFinishedMsg struct {
	report *domain.SearchReport
	err    error
}

type tuiModel struct {
	funnel    *usecase.SearchFunnel
	state     tuiState
	query     string
	lang      string
	report    *domain.SearchReport
	err       error
	width     int
	height    int
	styles    tuiStyles
}

type tuiStyles struct {
	title       lipgloss.Style
	label       lipgloss.Style
	box         lipgloss.Style
	accent      lipgloss.Style
	stats       lipgloss.Style
	repoItem    lipgloss.Style
	errStyle    lipgloss.Style
}

func newTUIStyles() tuiStyles {
	return tuiStyles{
		title:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).MarginBottom(1),
		label:    lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true),
		box:      lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7D56F4")).Padding(1, 2),
		accent:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")).Bold(true),
		stats:    lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Bold(true),
		repoItem: lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).MarginLeft(2),
		errStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true),
	}
}

func newTUIModel(funnel *usecase.SearchFunnel, initialQuery string, initialLang string) tuiModel {
	if initialLang == "" {
		initialLang = "go"
	}
	return tuiModel{
		funnel:    funnel,
		state:     stateInput,
		query:     initialQuery,
		lang:      initialLang,
		styles:    newTUIStyles(),
	}
}

func (m tuiModel) Init() tea.Cmd {
	if m.query != "" {
		return m.executeSearchCmd()
	}
	return nil
}

func (m tuiModel) executeSearchCmd() tea.Cmd {
	m.state = stateSearching
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		filter := domain.QueryFilter{
			Query:      m.query,
			Language:   m.lang,
			MaxResults: 10,
			Depth:      1,
		}

		rep, err := m.funnel.Execute(ctx, filter)
		return searchFinishedMsg{report: rep, err: err}
	}
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == stateInput && strings.TrimSpace(m.query) != "" {
				m.state = stateSearching
				return m, m.executeSearchCmd()
			} else if m.state == stateResults {
				m.state = stateInput
				m.report = nil
				m.err = nil
				return m, nil
			}
		case "backspace":
			if m.state == stateInput && len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
			}
		default:
			if m.state == stateInput && len(msg.String()) == 1 {
				m.query += msg.String()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case searchFinishedMsg:
		m.state = stateResults
		m.report = msg.report
		m.err = msg.err
	}

	return m, nil
}

func (m tuiModel) View() string {
	var b strings.Builder

	b.WriteString(m.styles.title.Render("🔍 Bateia — Motor de Busca Arquitetural Interativo (TUI)"))
	b.WriteString("\n\n")

	switch m.state {
	case stateInput:
		inputBox := fmt.Sprintf(
			"%s %s_\n\n%s %s\n\n%s",
			m.styles.label.Render("Digite a Query:"),
			m.query,
			m.styles.accent.Render("Linguagem alvo:"),
			m.lang,
			m.styles.stats.Render("Pressione ENTER para buscar | Ctrl+C ou Q para sair"),
		)
		b.WriteString(m.styles.box.Render(inputBox))

	case stateSearching:
		searchBox := fmt.Sprintf(
			"⏳ Executando Funil de Redução de 3 Estágios para: '%s' (%s)...\n\nPor favor, aguarde...",
			m.query,
			m.lang,
		)
		b.WriteString(m.styles.box.Render(searchBox))

	case stateResults:
		if m.err != nil {
			errBox := fmt.Sprintf("%s\n\n%s",
				m.styles.errStyle.Render(fmt.Sprintf("Erro ao executar funil: %v", m.err)),
				m.styles.stats.Render("Pressione ENTER para tentar novamente | Q para sair"),
			)
			b.WriteString(m.styles.box.Render(errBox))
		} else if m.report != nil {
			var resBuilder strings.Builder
			resBuilder.WriteString(m.styles.stats.Render(fmt.Sprintf(
				"📊 Funil Concluído: %d descobertos ➔ %d triados ➔ %d extraídos\n\n",
				m.report.TopCandidates, m.report.TriagedCount, m.report.ExtractedCount,
			)))

			resBuilder.WriteString(m.styles.label.Render("Top Repositórios Extraídos:\n"))
			for i, r := range m.report.Results {
				resBuilder.WriteString(fmt.Sprintf(
					"  %d. %s/%s (⭐ %d) — Score: %d/100 [%s]\n",
					i+1, r.Candidate.Owner, r.Candidate.Repo, r.Candidate.Stars, r.RelevanceScore, r.RelevanceReason,
				))
			}

			resBuilder.WriteString("\n" + m.styles.accent.Render("Pressione ENTER para nova busca | Q para sair"))
			b.WriteString(m.styles.box.Render(resBuilder.String()))
		}
	}

	return b.String()
}

// RunTUI inicia a interface TUI do Bateia.
func RunTUI(funnel *usecase.SearchFunnel, initialQuery string, initialLang string) error {
	p := tea.NewProgram(newTUIModel(funnel, initialQuery, initialLang))
	_, err := p.Run()
	return err
}
