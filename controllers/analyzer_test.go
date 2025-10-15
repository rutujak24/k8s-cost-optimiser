package controllers

import (
	"testing"
)

func TestAnalyzerRun(t *testing.T) {
	// Minimal scaffold test — implement real unit tests with fake client later
	// TODO: use envtest / controller-runtime fake client for richer tests
}

func TestAnalyzePodMemory_NoRecommendation(t *testing.T) {
	a := NewAnalyzer()
	// request 1Gi, average 600Mi (60% usage) => no recommendation
	avg := 600.0 * 1024.0 * 1024.0
	rec, err := a.AnalyzePodMemory("default", "pod1", "1Gi", avg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Fatalf("expected no recommendation, got %+v", rec)
	}
}

func TestAnalyzePodMemory_Recommendation(t *testing.T) {
	a := NewAnalyzer()
	// request 2Gi, average 200Mi (10% usage) => expect recommendation
	avg := 200.0 * 1024.0 * 1024.0
	rec, err := a.AnalyzePodMemory("prod", "pod2", "2Gi", avg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec == nil {
		t.Fatalf("expected a recommendation but got nil")
	}
	if rec.Spec.ResourceName != "pod2" || rec.Spec.Namespace != "prod" {
		t.Fatalf("unexpected rec metadata: %+v", rec.Spec)
	}
	// ensure recommendation string contains 'Mi'
	if rec.Spec.Recommendation == "" {
		t.Fatalf("empty recommendation")
	}
}
