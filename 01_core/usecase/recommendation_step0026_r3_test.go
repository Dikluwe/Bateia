/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_step0026_r3_test.md
 * @layer L1
 * @updated 2026-09-01
 */
package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

type step0026R3Reporter struct {
	events []domain.RecommendationProgressEvent
}

func (r *step0026R3Reporter) ReportRecommendationProgress(event domain.RecommendationProgressEvent) {
	r.events = append(r.events, event)
}

type step0026R3CancellingDiscovery struct {
	cancel context.CancelFunc
	calls  int
}

func (d *step0026R3CancellingDiscovery) Discover(ctx context.Context, _ domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	d.calls++
	d.cancel()
	<-ctx.Done()
	return domain.ProviderDiscovery{}, ctx.Err()
}

func TestStep0026R3RecommendAndRecommendObservedUseTheSameRealServicePipeline(t *testing.T) {
	semantic := &step0025SemanticSpy{}
	discovery := &step0025DiscoverySpy{}
	content := &step0025ContentSpy{}
	service := NewRecommendationService(semantic, discovery, content)
	syncResult, syncErr := service.Recommend(context.Background(), step0025Request())
	if syncErr != nil {
		t.Fatal(syncErr)
	}
	beforeSemantic, beforeDiscovery, beforeContent := semantic.calls, len(discovery.calls()), content.count()
	reporter := &step0026R3Reporter{}
	asyncResult, asyncErr := service.RecommendObserved(context.Background(), step0025Request(), reporter)
	if asyncErr != nil {
		t.Fatal(asyncErr)
	}
	if !reflect.DeepEqual(syncResult, asyncResult) {
		t.Fatalf("observed path forked result:\nsync=%#v\nasync=%#v", syncResult, asyncResult)
	}
	if semantic.calls != beforeSemantic*2 || len(discovery.calls()) != beforeDiscovery*2 || content.count() != beforeContent*2 {
		t.Fatalf("observed path did not delegate through the same ports: semantic=%d/%d discovery=%d/%d content=%d/%d", semantic.calls, beforeSemantic, len(discovery.calls()), beforeDiscovery, content.count(), beforeContent)
	}
	if len(reporter.events) == 0 {
		t.Fatal("observed path emitted no controlled progress")
	}
}

func TestStep0026R3CancellationObservedInsideProviderStartsNoLaterProviderOrContentCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	discovery := &step0026R3CancellingDiscovery{cancel: cancel}
	content := &step0025ContentSpy{}
	service := NewRecommendationService(&step0025SemanticSpy{}, discovery, content)
	_, err := service.RecommendObserved(ctx, step0025Request(), &step0026R3Reporter{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("provider-boundary cancellation = %v", err)
	}
	if discovery.calls != 1 || content.count() != 0 {
		t.Fatalf("new provider/content call started after cancellation observation: discovery=%d content=%d", discovery.calls, content.count())
	}
}
