package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/otengedward882-cmd/AnnoScore/pkg/metrics"
)

// GenerateTerminalReport renders an ASCII formatted summary for the command line.
func GenerateTerminalReport(res *metrics.AgreementResult) string {
	var b strings.Builder

	b.WriteString("=================================================================\n")
	b.WriteString("  🎯 AnnoScore — AI Dataset Inter-Annotator Agreement Report\n")
	b.WriteString("=================================================================\n\n")

	b.WriteString(fmt.Sprintf("• Total Evaluated Items:     %d\n", res.TotalItems))
	b.WriteString(fmt.Sprintf("• Active Annotators:         %d\n", res.TotalAnnotators))
	b.WriteString(fmt.Sprintf("• Observed Agreement Rate:   %.2f%%\n", res.ObservedAgreement*100))
	b.WriteString(fmt.Sprintf("• Cohen's Kappa (κ):         %.3f\n", res.CohenKappa))
	b.WriteString(fmt.Sprintf("• Quality Classification:    %s\n\n", res.Interpretation))

	b.WriteString("-----------------------------------------------------------------\n")
	b.WriteString("  Individual Annotator Consistency vs Consensus:\n")
	b.WriteString("-----------------------------------------------------------------\n")
	for annotator, rate := range res.AnnotatorAgreement {
		badge := "✅ RELIABLE"
		if rate < 0.70 {
			badge = "⚠️  OUTLIER "
		}
		b.WriteString(fmt.Sprintf("  [%s] %-16s Consensus Match: %.1f%%\n", badge, annotator, rate*100))
	}

	if len(res.DisagreedItemIDs) > 0 {
		b.WriteString("\n-----------------------------------------------------------------\n")
		b.WriteString(fmt.Sprintf("  High-Ambiguity Items Requiring Guideline Clarification (%d):\n", len(res.DisagreedItemIDs)))
		b.WriteString("-----------------------------------------------------------------\n")
		limit := 5
		if len(res.DisagreedItemIDs) < limit {
			limit = len(res.DisagreedItemIDs)
		}
		for i := 0; i < limit; i++ {
			b.WriteString(fmt.Sprintf("  • Item ID: %s (Annotator Disagreement Detected)\n", res.DisagreedItemIDs[i]))
		}
		if len(res.DisagreedItemIDs) > limit {
			b.WriteString(fmt.Sprintf("  ... and %d more items.\n", len(res.DisagreedItemIDs)-limit))
		}
	}

	b.WriteString("\n=================================================================\n")
	return b.String()
}

// GenerateMarkdownReport outputs a clean table for PRs or documentation.
func GenerateMarkdownReport(res *metrics.AgreementResult) string {
	var b strings.Builder

	b.WriteString("### 🎯 AnnoScore AI Dataset Agreement Report\n\n")
	b.WriteString(fmt.Sprintf("**Status:** %s\n\n", res.Interpretation))
	b.WriteString("| Metric | Value |\n")
	b.WriteString("|:---|:---|\n")
	b.WriteString(fmt.Sprintf("| **Observed Agreement** | `%.2f%%` |\n", res.ObservedAgreement*100))
	b.WriteString(fmt.Sprintf("| **Cohen's Kappa (κ)** | `%.3f` |\n", res.CohenKappa))
	b.WriteString(fmt.Sprintf("| **Evaluated Items** | `%d` |\n", res.TotalItems))
	b.WriteString(fmt.Sprintf("| **Active Annotators** | `%d` |\n\n", res.TotalAnnotators))

	b.WriteString("#### Annotator Consensus Correlation\n\n")
	b.WriteString("| Annotator ID | Consensus Agreement | Status |\n")
	b.WriteString("|:---|:---|:---|\n")
	for annotator, rate := range res.AnnotatorAgreement {
		status := "✅ Passing"
		if rate < 0.70 {
			status = "⚠️ Outlier / Needs Review"
		}
		b.WriteString(fmt.Sprintf("| `%s` | `%.1f%%` | %s |\n", annotator, rate*100, status))
	}

	return b.String()
}

// GenerateJSONReport serializes results to JSON.
func GenerateJSONReport(res *metrics.AgreementResult) (string, error) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
