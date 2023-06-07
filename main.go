//go:generate go run ./variants/main.go -input ./variants/input.tsv -output ./variants/generated/zz_generated.variants.go

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/bertinatto/testgrid/internal/crawler"
	"github.com/bertinatto/testgrid/internal/report"
)

func main() {
	prFlag := flag.String("pr", "", "pull request in the format 'org/repo#prID'")
	outputFlag := flag.String("output", "report.html", "specify the output file for the report (default: report.html)")
	flag.Parse()

	if *outputFlag == "" {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "ERROR: Output file cannot be empty.\n")
		os.Exit(1)
	}

	// Extract organization, repository, and pull request ID
	re := regexp.MustCompile(`^(\w+)/(\w+)#(\d+)$`)
	matches := re.FindStringSubmatch(*prFlag)
	if len(matches) < 3 {
		fmt.Fprintf(os.Stderr, "ERROR Invalid input format %q. Expected: 'org/repo#pr\n", *prFlag)
		flag.PrintDefaults()
		os.Exit(1)
	}

	org := matches[1]
	repo := matches[2]
	prID, err := strconv.Atoi(matches[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to convert pull request ID to integer: %v", err)
		os.Exit(1)
	}

	jobs := crawler.New(org, repo, prID).Do()
	report := report.New(org, repo, prID)
	err = report.Create(jobs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create report: %v", err)
		os.Exit(1)
	}

	if err := report.WriteToFile(*outputFlag); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to write report to file: %v", err)
		os.Exit(1)
	}
}
