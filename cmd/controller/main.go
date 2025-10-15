package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourname/k8s-cost-optimizer/controllers"
)

func main() {
	var prometheusAddr string
	var runSample bool
	flag.StringVar(&prometheusAddr, "prometheus-addr", "http://localhost:9090", "Prometheus base URL")
	flag.BoolVar(&runSample, "sample", false, "Run a sample analyzer invocation and print YAML recommendation")
	flag.Parse()

	fmt.Println("k8s-cost-optimizer scaffold starting")

	if runSample {
		a := controllers.NewAnalyzer()
		// hypothetical pod: request 2Gi, observed average 200Mi
		avgBytes := 200.0 * 1024.0 * 1024.0
		rec, err := a.AnalyzePodMemory("production", "my-app-abc123", "2Gi", avgBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "analyze error: %v\n", err)
			os.Exit(1)
		}
		if rec == nil {
			fmt.Println("No recommendation")
			return
		}
		out, err := a.EmitYAML(rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "emit yaml error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(out))
		return
	}

	promClient, err := controllers.NewPrometheusClient(prometheusAddr)
	if err != nil {
		fmt.Printf("warning: could not create prometheus client: %v\n", err)
	} else {
		_ = promClient
	}

	// start a simple ticker to demonstrate background work (replace with controller-runtime later)
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Analyzer tick — TODO: implement analysis loop and reconcile CostRecommendation CRs")
	}
}
