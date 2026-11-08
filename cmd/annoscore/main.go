package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/otengedward882-cmd/AnnoScore/pkg/metrics"
	"github.com/otengedward882-cmd/AnnoScore/pkg/parser"
	"github.com/otengedward882-cmd/AnnoScore/pkg/report"
)

const version = "1.0.0"

func main() {
	filePath := flag.String("file", "", "Path to annotation CSV dataset (e.g. data.csv)")
	format := flag.String("format", "text", "Output format: text, markdown, or json")
	showVersion := flag.Bool("version", false, "Display AnnoScore version")
	threshold := flag.Float64("min-kappa", 0.60, "Minimum required Cohen's Kappa for passing exit code (default: 0.60)")

	flag.Parse()

	if *showVersion {
		fmt.Printf("AnnoScore version %s (Go)\n", version)
		os.Exit(0)
	}

	if *filePath == "" {
		if flag.NArg() > 0 {
			*filePath = flag.Arg(0)
		} else {
			fmt.Println("Error: Missing dataset file path.")
			fmt.Println("Usage: annoscore --file dataset.csv [--format text|markdown|json]")
			os.Exit(2)
		}
	}

	dataset, err := parser.LoadCSV(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading dataset: %v\n", err)
		os.Exit(1)
	}

	if len(dataset.Annotators) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Dataset must contain at least 2 distinct annotators to calculate agreement.\n")
		os.Exit(1)
	}

	// Pairwise evaluation between the primary two annotators (or overall consensus)
	rater1 := dataset.Annotators[0]
	rater2 := dataset.Annotators[1]

	var r1Labels []string
	var r2Labels []string
	var disagreedItems []string

	for itemID, ratings := range dataset.ItemRatings {
		l1, ok1 := ratings[rater1]
		l2, ok2 := ratings[rater2]
		if ok1 && ok2 {
			r1Labels = append(r1Labels, l1)
			r2Labels = append(r2Labels, l2)
			if l1 != l2 {
				disagreedItems = append(disagreedItems, itemID)
			}
		}
	}

	po, kappa, err := metrics.CalculateCohenKappa(r1Labels, r2Labels)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Calculation error: %v\n", err)
		os.Exit(1)
	}

	// Compute per-annotator agreement against majority vote
	annotatorScores := make(map[string]float64)
	for _, a := range dataset.Annotators {
		matches := 0
		total := 0
		for _, ratings := range dataset.ItemRatings {
			if userRating, ok := ratings[a]; ok {
				// Count majority in ratings
				counts := make(map[string]int)
				for _, r := range ratings {
					counts[r]++
				}
				majorityLabel := ""
				maxCount := 0
				for lbl, cnt := range counts {
					if cnt > maxCount {
						maxCount = cnt
						majorityLabel = lbl
					}
				}
				if userRating == majorityLabel {
					matches++
				}
				total++
			}
		}
		if total > 0 {
			annotatorScores[a] = float64(matches) / float64(total)
		}
	}

	res := &metrics.AgreementResult{
		TotalItems:         len(r1Labels),
		TotalAnnotators:    len(dataset.Annotators),
		ObservedAgreement:  po,
		CohenKappa:         kappa,
		Interpretation:     metrics.InterpretKappa(kappa),
		AnnotatorAgreement: annotatorScores,
		DisagreedItemIDs:   disagreedItems,
	}

	switch *format {
	case "markdown":
		fmt.Print(report.GenerateMarkdownReport(res))
	case "json":
		jsonStr, _ := report.GenerateJSONReport(res)
		fmt.Println(jsonStr)
	default:
		fmt.Print(report.GenerateTerminalReport(res))
	}

	if kappa < *threshold {
		os.Exit(1)
	}
	os.Exit(0)
}
