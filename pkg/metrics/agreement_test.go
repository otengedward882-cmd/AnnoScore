package metrics

import (
	"math"
	"testing"
)

func TestCalculateCohenKappa(t *testing.T) {
	r1 := []string{"Positive", "Negative", "Neutral", "Positive"}
	r2 := []string{"Positive", "Negative", "Neutral", "Positive"}

	po, kappa, err := CalculateCohenKappa(r1, r2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if po != 1.0 {
		t.Errorf("expected observed agreement 1.0, got %f", po)
	}

	if kappa != 1.0 {
		t.Errorf("expected kappa 1.0, got %f", kappa)
	}
}

func TestCohenKappaPartialAgreement(t *testing.T) {
	r1 := []string{"Yes", "Yes", "No", "No", "Yes"}
	r2 := []string{"Yes", "No", "No", "No", "Yes"}

	po, kappa, err := CalculateCohenKappa(r1, r2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if po != 0.8 {
		t.Errorf("expected po 0.8, got %f", po)
	}

	if kappa <= 0.0 || kappa >= 1.0 {
		t.Errorf("expected 0 < kappa < 1, got %f", kappa)
	}
}

func TestFleissKappa(t *testing.T) {
	// 5 items, 3 categories, rated by 4 evaluators
	matrix := [][]int{
		{0, 0, 4},
		{1, 3, 0},
		{0, 4, 0},
		{4, 0, 0},
		{1, 2, 1},
	}

	kappa, err := CalculateFleissKappa(matrix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if math.IsNaN(kappa) || kappa < 0 {
		t.Errorf("expected valid positive kappa, got %f", kappa)
	}
}
