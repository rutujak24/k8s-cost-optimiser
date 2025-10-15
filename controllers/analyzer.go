package controllers

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/yaml"
	"k8s.io/apimachinery/pkg/api/resource"

	finopsv1 "github.com/yourname/k8s-cost-optimizer/api/v1"
)

// Analyzer is a lightweight analyzer for demo/college project use.
// It compares a pod's memory request to its observed average usage and
// emits a CostRecommendation when usage is much lower than request.
type Analyzer struct {
	// In a full operator this would hold a k8s client and logger.
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzePodMemory compares memory request (e.g. "2Gi") with averageUsageBytes
// and returns a CostRecommendation pointer if a recommendation should be made.
// The function uses a simple heuristic: if average usage < thresholdPercent of request
// (default 30%), recommend lowering to safeTargetPercent (default 50% of observed peak/request balance).
func (a *Analyzer) AnalyzePodMemory(namespace, podName, memRequest string, averageUsageBytes float64) (*finopsv1.CostRecommendation, error) {
	// parse memory request
	q, err := resource.ParseQuantity(memRequest)
	if err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}
	reqBytes := float64(q.Value())

	if reqBytes <= 0 {
		return nil, fmt.Errorf("invalid request size: %v", reqBytes)
	}

	usagePct := (averageUsageBytes / reqBytes) * 100.0

	// threshold to consider overprovisioned
	const thresholdPct = 30.0

	if usagePct >= thresholdPct {
		// no recommendation
		return nil, nil
	}

	// recommended target: round up to nearest Mi and set to 50% above observed average
	safeTargetBytes := averageUsageBytes * 1.5
	// convert to Mi
	targetMi := int((safeTargetBytes / (1024.0 * 1024.0)) + 0.5)
	if targetMi < 4 {
		targetMi = 4
	}

	rec := &finopsv1.CostRecommendation{
		Spec: finopsv1.CostRecommendationSpec{
			ResourceType:     "Pod",
			ResourceName:     podName,
			Namespace:        namespace,
			Issue:            fmt.Sprintf("Memory request %s, average usage %0.1fMi", memRequest, averageUsageBytes/(1024.0*1024.0)),
			Recommendation:   fmt.Sprintf("Reduce to %dMi", targetMi),
			PotentialSavings: "$X/month", // placeholder — integrate cloud pricing later
			Confidence:       "medium",
			AutoApplicable:   false,
		},
	}

	return rec, nil
}

// EmitYAML serializes a CostRecommendation to YAML (for demo purposes).
func (a *Analyzer) EmitYAML(rec *finopsv1.CostRecommendation) ([]byte, error) {
	if rec == nil {
		return nil, fmt.Errorf("nil recommendation")
	}
	// add apiVersion/kind for nicer YAML output
	rec.APIVersion = "finops.yourname.dev/v1"
	rec.Kind = "CostRecommendation"
	b, err := yaml.Marshal(rec)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Run performs a single analysis pass (placeholder for periodic runs).
func (a *Analyzer) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Analyzer: placeholder run")
	}
	return nil
}
