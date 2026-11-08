package metrics

import (
	"fmt"
	"math"
)

// AgreementResult holds statistical agreement metrics for annotation datasets.
type AgreementResult struct {
	TotalItems          int                `json:"total_items"`
	TotalAnnotators     int                `json:"total_annotators"`
	ObservedAgreement   float64            `json:"observed_agreement_pct"`
	CohenKappa          float64            `json:"cohen_kappa"`
	FleissKappa         float64            `json:"fleiss_kappa"`
	Interpretation      string             `json:"interpretation"`
	AnnotatorAgreement  map[string]float64 `json:"annotator_agreement_pct"`
	DisagreedItemIDs    []string           `json:"disagreed_item_ids"`
}

// PairwiseAgreement computes observed agreement and Cohen's Kappa between two annotators.
func CalculateCohenKappa(rater1Labels, rater2Labels []string) (observed float64, kappa float64, err error) {
	n := len(rater1Labels)
	if n == 0 || n != len(rater2Labels) {
		return 0, 0, fmt.Errorf("label slices must be non-empty and of equal length")
	}

	// Distinct categories
	categories := make(map[string]bool)
	agreedCount := 0
	r1Dist := make(map[string]int)
	r2Dist := make(map[string]int)

	for i := 0; i < n; i++ {
		l1 := rater1Labels[i]
		l2 := rater2Labels[i]
		categories[l1] = true
		categories[l2] = true
		r1Dist[l1]++
		r2Dist[l2]++

		if l1 == l2 {
			agreedCount++
		}
	}

	// Observed agreement Po
	po := float64(agreedCount) / float64(n)

	// Expected agreement Pe by chance
	pe := 0.0
	for cat := range categories {
		p1 := float64(r1Dist[cat]) / float64(n)
		p2 := float64(r2Dist[cat]) / float64(n)
		pe += p1 * p2
	}

	// Cohen's Kappa = (Po - Pe) / (1 - Pe)
	if math.Abs(1.0-pe) < 1e-9 {
		kappa = 1.0
	} else {
		kappa = (po - pe) / (1.0 - pe)
	}

	return po, kappa, nil
}

// InterpretKappa returns the standard Landis & Koch (1977) benchmark classification.
func InterpretKappa(k float64) string {
	switch {
	case k >= 0.81:
		return "Near Perfect Agreement (Consensus High)"
	case k >= 0.61:
		return "Substantial Agreement (Production Quality)"
	case k >= 0.41:
		return "Moderate Agreement (Guideline Clarification Advised)"
	case k >= 0.21:
		return "Fair Agreement (High Ambiguity)"
	case k > 0.0:
		return "Slight Agreement (Needs Review)"
	default:
		return "Poor Agreement (No Consensus / Random)"
	}
}

// CalculateFleissKappa computes generalized multi-rater agreement across N items and k categories.
func CalculateFleissKappa(matrix [][]int) (float64, error) {
	n := len(matrix)
	if n == 0 {
		return 0, fmt.Errorf("matrix cannot be empty")
	}

	k := len(matrix[0])
	if k == 0 {
		return 0, fmt.Errorf("matrix must have at least one category")
	}

	// Number of ratings per subject
	raters := 0
	for _, val := range matrix[0] {
		raters += val
	}
	if raters <= 1 {
		return 0, fmt.Errorf("need at least 2 ratings per item")
	}

	// Calculate P_i for each subject
	pSubject := make([]float64, n)
	for i := 0; i < n; i++ {
		sumSquares := 0
		for j := 0; j < k; j++ {
			sumSquares += matrix[i][j] * matrix[i][j]
		}
		pSubject[i] = float64(sumSquares-raters) / float64(raters*(raters-1))
	}

	// Mean observed agreement P_bar
	pBar := 0.0
	for _, p := range pSubject {
		pBar += p
	}
	pBar /= float64(n)

	// Proportion of all assignments to j-th category (p_j)
	pCat := make([]float64, k)
	for j := 0; j < k; j++ {
		sumCat := 0
		for i := 0; i < n; i++ {
			sumCat += matrix[i][j]
		}
		pCat[j] = float64(sumCat) / float64(n*raters)
	}

	// Expected agreement P_bar_e
	pBarE := 0.0
	for j := 0; j < k; j++ {
		pBarE += pCat[j] * pCat[j]
	}

	if math.Abs(1.0-pBarE) < 1e-9 {
		return 1.0, nil
	}

	kappa := (pBar - pBarE) / (1.0 - pBarE)
	return kappa, nil
}
