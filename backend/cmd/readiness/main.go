// The readiness command verifies repository evidence required for production.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"aeolyzer/internal/releasegate"
)

func main() {
	root := flag.String("root", ".", "repository root to verify")
	jsonOutput := flag.Bool("json", false, "emit the report as JSON")
	flag.Parse()

	report, err := releasegate.Check(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "readiness check failed: %v\n", err)
		os.Exit(2)
	}

	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			fmt.Fprintf(os.Stderr, "encode readiness report: %v\n", err)
			os.Exit(2)
		}
	} else {
		printReport(report)
	}

	if !report.Ready() {
		os.Exit(1)
	}
}

func printReport(report releasegate.Report) {
	if report.Ready() {
		fmt.Println("production readiness: PASS")
		return
	}

	fmt.Printf("production readiness: FAIL (%d blockers)\n", len(report.Findings))
	for _, finding := range report.Findings {
		location := finding.Area
		if finding.Path != "" {
			location += " " + finding.Path
		}
		fmt.Printf("- [%s] %s: %s\n", finding.Code, location, finding.Message)
	}
}
