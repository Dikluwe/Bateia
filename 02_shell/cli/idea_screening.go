/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_screening.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func readJSONFile(path string, destination any) error {
	data, err := readIdeaJSONInput(path)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", path, err)
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("JSON inválido em %s: %w", path, err)
	}
	return nil
}

func runIdeaScreen(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea screen", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotPath := fs.String("snapshot", "", "Snapshot de descoberta")
	modelPath := fs.String("model", "", "Modelo de peneira; omita para baseline não calibrado")
	judgmentsPath := fs.String("judgments", "", "Julgamentos opcionais para replay diagnóstico")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea screen inválidos: %w", err)
	}
	if *snapshotPath == "" {
		return fmt.Errorf("idea screen exige --snapshot")
	}
	var snapshot domain.DiscoverySnapshot
	if err := readJSONFile(*snapshotPath, &snapshot); err != nil {
		return err
	}
	model := usecase.DefaultScreeningModel()
	if *modelPath != "" {
		data, err := readIdeaJSONInput(*modelPath)
		if err != nil {
			return err
		}
		var calibration domain.CalibrationReport
		if err := json.Unmarshal(data, &calibration); err == nil && calibration.Model.SchemaVersion != "" {
			model = calibration.Model
		} else {
			model = domain.ScreeningModel{}
			if err := json.Unmarshal(data, &model); err != nil || model.SchemaVersion == "" {
				return fmt.Errorf("arquivo não contém modelo de peneira")
			}
		}
	}
	report, err := usecase.ScreenDiscovery(snapshot, model)
	if err != nil {
		return fmt.Errorf("peneira inválida: %w", err)
	}
	if *judgmentsPath != "" {
		var judgments domain.DiscoveryJudgments
		if err := readJSONFile(*judgmentsPath, &judgments); err != nil {
			return err
		}
		report = usecase.DiagnoseScreening(report, judgments)
	}
	return writeIdeaJSON(report, *outputPath, stdout)
}

func runIdeaScreening(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea screening <calibrate|audit>")
	}
	switch args[0] {
	case "calibrate":
		return runIdeaScreeningCalibrate(args[1:], stdout)
	case "audit":
		return runIdeaScreeningAudit(args[1:], stdout)
	default:
		return fmt.Errorf("subcomando screening desconhecido %q", args[0])
	}
}

func runIdeaScreeningCalibrate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea screening calibrate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotsDir := fs.String("snapshots", "", "Diretório de snapshots")
	judgmentsDir := fs.String("judgments", "", "Diretório de julgamentos")
	strategy := fs.String("strategy", "lexical", "Estratégia dos snapshots")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de calibração inválidos: %w", err)
	}
	if *snapshotsDir == "" || *judgmentsDir == "" {
		return fmt.Errorf("calibração exige --snapshots e --judgments")
	}
	pattern := filepath.Join(*snapshotsDir, "*."+strings.ToLower(*strategy)+".snapshot.json")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	sort.Strings(paths)
	groups := make([]domain.ScreeningDatasetGroup, 0, len(paths))
	for _, snapshotPath := range paths {
		base := filepath.Base(snapshotPath)
		suffix := "." + strings.ToLower(*strategy) + ".snapshot.json"
		ideaID := strings.TrimSuffix(base, suffix)
		var snapshot domain.DiscoverySnapshot
		if err := readJSONFile(snapshotPath, &snapshot); err != nil {
			return err
		}
		var judgments domain.DiscoveryJudgments
		if err := readJSONFile(filepath.Join(*judgmentsDir, ideaID+".judgments.json"), &judgments); err != nil {
			return err
		}
		groups = append(groups, domain.ScreeningDatasetGroup{IdeaID: ideaID, Snapshot: snapshot, Judgments: judgments})
	}
	report, err := usecase.CalibrateScreening(groups, usecase.DefaultScreeningModel())
	if err != nil {
		return fmt.Errorf("calibração falhou: %w", err)
	}
	return writeIdeaJSON(report, *outputPath, stdout)
}

func runIdeaScreeningAudit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea screening audit", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	revisionsPath := fs.String("revisions", "", "Arquivo de revisões com evidências")
	screeningPath := fs.String("screening", "", "Relatório de peneira para selecionar auditoria")
	judgmentsPath := fs.String("judgments", "", "Julgamentos atuais")
	screeningDir := fs.String("screening-dir", "", "Diretório com relatórios *.screening.json")
	judgmentsDir := fs.String("judgments-dir", "", "Diretório com julgamentos correspondentes")
	lowCount := fs.Int("low-count", 10, "Controles entre os menores scores")
	boundaryCount := fs.Int("boundary-count", 10, "Controles irrelevantes próximos ao avanço")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de auditoria inválidos: %w", err)
	}
	if *revisionsPath != "" {
		var revisions []domain.JudgmentRevision
		if err := readJSONFile(*revisionsPath, &revisions); err != nil {
			return err
		}
		report, err := usecase.ValidateAuditRevisions(revisions)
		if err != nil {
			return fmt.Errorf("auditoria inválida: %w", err)
		}
		return writeIdeaJSON(report, *outputPath, stdout)
	}
	var screening domain.ScreeningReport
	var judgments domain.DiscoveryJudgments
	if *screeningDir != "" && *judgmentsDir != "" {
		paths, err := filepath.Glob(filepath.Join(*screeningDir, "*.screening.json"))
		if err != nil {
			return err
		}
		sort.Strings(paths)
		judgments.SchemaVersion = "1"
		for _, path := range paths {
			var partial domain.ScreeningReport
			if err := readJSONFile(path, &partial); err != nil {
				return err
			}
			screening.SchemaVersion = "1"
			screening.Decisions = append(screening.Decisions, partial.Decisions...)
			ideaID := strings.TrimSuffix(filepath.Base(path), ".screening.json")
			var labels domain.DiscoveryJudgments
			if err := readJSONFile(filepath.Join(*judgmentsDir, ideaID+".judgments.json"), &labels); err != nil {
				return err
			}
			judgments.Judgments = append(judgments.Judgments, labels.Judgments...)
		}
	} else {
		if *screeningPath == "" || *judgmentsPath == "" {
			return fmt.Errorf("auditoria exige --revisions, par de arquivos ou par de diretórios")
		}
		if err := readJSONFile(*screeningPath, &screening); err != nil {
			return err
		}
		if err := readJSONFile(*judgmentsPath, &judgments); err != nil {
			return err
		}
	}
	selection, err := usecase.SelectAuditCandidates(screening, judgments, *lowCount, *boundaryCount)
	if err != nil {
		return err
	}
	return writeIdeaJSON(selection, *outputPath, stdout)
}
