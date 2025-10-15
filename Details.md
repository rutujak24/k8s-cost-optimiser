High-level summary
This repo is a small college-project scaffold for a Kubernetes Cost Optimizer. It is not yet a full operator, but it contains a working demo analyzer that:

Takes a pod memory request (e.g. "2Gi") and an observed average memory usage (bytes).
Uses a simple heuristic to decide whether the pod is over‑provisioned.
If over‑provisioned, builds a CostRecommendation object (in-memory) and can emit it as YAML.
Provides a CLI demo mode (--sample) and a script / Makefile target to produce a sample YAML file (sample-recommendation.yaml).
Includes unit tests for the analyzer logic.
The real operator features (controller runtime reconciliation, automatic creation of CRs in-cluster, Prometheus-driven sampling, cloud pricing) are scaffolds / TODOs.

Where to look (important files)
costrecommendation_types.go

Defines the CRD Go type CostRecommendation, with Spec and Status structs. Fields in Spec:
resourceType, resourceName, namespace, issue, recommendation, potentialSavings, confidence, autoApplicable.
Note: this file defines types only; code generation (deepcopy, scheme registration) is intentionally not run in this scaffold.
analyzer.go

Implements the demo analyzer logic:
func (a *Analyzer) AnalyzePodMemory(namespace, podName, memRequest string, averageUsageBytes float64) (*finopsv1.CostRecommendation, error)
Parses memRequest with resource.ParseQuantity.
Computes usagePct = average / request * 100.
If usagePct >= 30% → returns nil (no recommendation).
If usagePct < 30% → builds a CostRecommendation with:
Recommendation target = int(round(1.5 * average / Mi)) Mi, min 4Mi.
Confidence "medium", PotentialSavings "$X/month" (placeholder).
func (a *Analyzer) EmitYAML(rec *finopsv1.CostRecommendation) ([]byte, error)
Adds apiVersion: finops.yourname.dev/v1 and kind: CostRecommendation and marshals the object to YAML using sigs.k8s.io/yaml.
Run(ctx context.Context) — a placeholder periodic run (no real work).
prometheus.go

Thin wrapper around the Prometheus HTTP client.
NewPrometheusClient(address string) returns a PrometheusClient.
QueryInstant(ctx, query) calls Prometheus API and returns result.String(). It currently uses time.Now() as the timestamp.
Not currently integrated with the analyzer.
controllers/recommender.go

Placeholder for converting analysis into CRD creation/upserts. Currently only a scaffold with NewRecommender() and Recommend(...) (no in-cluster creation).
main.go

Simple CLI program (not controller-runtime manager):
Flags:
--prometheus-addr (default http://localhost:9090)
--sample (bool) — when set, runs a sample analyzer invocation and prints YAML.
Sample behavior: with --sample the program calls AnalyzePodMemory("production","my-app-abc123","2Gi", 200Mi) and prints YAML to stdout.
Otherwise it creates a Prometheus client (if reachable) and runs a placeholder ticker.
Makefile

make build — builds binary to bin/k8s-cost-optimizer.
make run — runs the program (ticker placeholder).
make sample — runs the sample analyzer and writes output to sample-recommendation.yaml.
hack/run-sample.sh

Helper script that runs the demo and writes YAML to sample-recommendation.yaml.
sample-recommendation.yaml

The YAML output produced by the sample run (example produced by your run). Example content:
controllers/analyzer_test.go

Unit tests:
TestAnalyzePodMemory_NoRecommendation — ensures no recommendation when usage is high (60% case).
TestAnalyzePodMemory_Recommendation — ensures a recommendation for a low-usage pod (10% case).
Tests run with go test ./... and currently pass.
Concrete behavior details (the contract)
Inputs:
namespace, podName (strings)
memRequest string, valid Kubernetes quantity (e.g. "2Gi", "512Mi")
averageUsageBytes float64 (bytes)
Output:
(*finopsv1.CostRecommendation, error)
nil, nil means no recommendation
rec, nil means a recommendation was produced
_, err on parsing or other error
Heuristic constants:
Over-provisioned threshold: 30% (if observed average < 30% of request → recommend)
Target recommended = 1.5× observed average, rounded to Mi, minimum 4Mi
Potential savings: placeholder string (not calculated)
How to run / test this repo locally
(You already ran these successfully; repeating for convenience.)

Run unit tests:

Create the sample YAML (CLI demo):

Run the program interactively (ticker mode):

Run the demo directly:

Current limitations and notable omissions
Not a full Kubernetes operator:
No controller-runtime Manager or Reconciler is wired.
The CostRecommendation objects are not created in the cluster — the analyzer only emits YAML as demo output.
Type generation missing:
The API types are defined, but no controller-gen was used to generate deep-copy helpers / CRD manifests. For a production operator you should run the code generation steps.
Prometheus integration:
controllers/prometheus.go exists but is not yet used to fetch real metrics for the analyzer.
Pricing:
potentialSavings is a placeholder; there is no cloud pricing adapter yet (AWS/GCP/Azure).
Heuristic is simple:
Uses a single average value and a fixed threshold. Production FinOps needs more robust statistics (percentiles, percentile over time windows, request vs observed 95th percentile, replica-level aggregation, scaling recommendations).
Safety / edge cases to be aware of
Invalid memRequest strings produce parse errors; analyzer returns error on parse failure.
Very small observed average values might produce extremely small targets — code enforces a 4Mi minimum.
The sample output uses placeholder apiVersion/kind strings; when you later wire controller-runtime you should regenerate types and CRD manifests accordingly.
Suggested next steps (pick one)
Make the sample CLI accept arguments (pod name, namespace, request, observed average) — quick and useful for demos.
Integrate Prometheus:
Implement a function to query average memory usage (7d avg or p95) and feed that into AnalyzePodMemory.
Wire a Reconciler and Kubernetes client:
Add controller-runtime, generate deep-copy code and CRD manifests, implement a Reconciler that creates CostRecommendation CRs in-cluster.
Add cloud pricing integration:
Small adapter to convert resource savings into $/month using a chosen cloud's pricing API.
Improve heuristic:
Use percentile queries, consider per-deployment aggregation, factor in replicas.



#### What this project does 
Imagine a container (pod) in Kubernetes asks the cluster for 2 GiB of memory, but in practice it only ever uses 200 MiB. That means you’re paying for memory you don’t need.
This project contains a small program (an “analyzer”) that:
Takes a pod’s memory request (what it asked for) and the observed average memory usage (what it actually used),
Checks whether the pod is over-provisioned (right now it flags it when usage is less than 30% of the requested memory),
If it’s over-provisioned, it creates a human-readable suggestion — a CostRecommendation — telling you to reduce the pod memory request (e.g., “Reduce to 300Mi”).
The tool outputs that suggestion as a YAML file (so it’s easy to read or apply later).
It’s safe: it does not automatically change anything in your cluster — it only suggests changes.
What it’s not (important)
It is NOT a full Kubernetes operator yet. It will not automatically apply changes or create Kubernetes CRs in-cluster.
It does not calculate real dollar savings (the savings field is a placeholder).
It doesn’t automatically fetch real metrics from Prometheus yet — that can be added.
