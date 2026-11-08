package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// AnnotationRecord represents a single evaluator rating.
type AnnotationRecord struct {
	ItemID      string
	AnnotatorID string
	Label       string
	Confidence  string
}

// AnnotationDataset holds records indexed by item and annotator.
type AnnotationDataset struct {
	Records           []AnnotationRecord
	ItemRatings       map[string]map[string]string // ItemID -> AnnotatorID -> Label
	Annotators        []string
	Categories        []string
	TotalItems        int
}

// LoadCSV parses an annotation CSV dataset.
func LoadCSV(filePath string) (*AnnotationDataset, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dataset file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Map headers
	itemIdx, annotatorIdx, labelIdx := -1, -1, -1
	for idx, h := range headers {
		lower := strings.ToLower(strings.TrimSpace(h))
		if strings.Contains(lower, "item") || strings.Contains(lower, "prompt") || strings.Contains(lower, "id") && itemIdx == -1 {
			itemIdx = idx
		} else if strings.Contains(lower, "annotator") || strings.Contains(lower, "rater") || strings.Contains(lower, "user") {
			annotatorIdx = idx
		} else if strings.Contains(lower, "label") || strings.Contains(lower, "rating") || strings.Contains(lower, "choice") {
			labelIdx = idx
		}
	}

	if itemIdx == -1 || annotatorIdx == -1 || labelIdx == -1 {
		// Fallback to standard columns 0, 1, 2
		itemIdx, annotatorIdx, labelIdx = 0, 1, 2
	}

	dataset := &AnnotationDataset{
		ItemRatings: make(map[string]map[string]string),
	}

	annotatorSet := make(map[string]bool)
	categorySet := make(map[string]bool)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}

		if len(row) <= itemIdx || len(row) <= annotatorIdx || len(row) <= labelIdx {
			continue
		}

		itemID := strings.TrimSpace(row[itemIdx])
		annotatorID := strings.TrimSpace(row[annotatorIdx])
		label := strings.TrimSpace(row[labelIdx])

		if itemID == "" || annotatorID == "" || label == "" {
			continue
		}

		record := AnnotationRecord{
			ItemID:      itemID,
			AnnotatorID: annotatorID,
			Label:       label,
		}
		dataset.Records = append(dataset.Records, record)

		if _, exists := dataset.ItemRatings[itemID]; !exists {
			dataset.ItemRatings[itemID] = make(map[string]string)
		}
		dataset.ItemRatings[itemID][annotatorID] = label
		annotatorSet[annotatorID] = true
		categorySet[label] = true
	}

	for a := range annotatorSet {
		dataset.Annotators = append(dataset.Annotators, a)
	}
	for c := range categorySet {
		dataset.Categories = append(dataset.Categories, c)
	}
	dataset.TotalItems = len(dataset.ItemRatings)

	return dataset, nil
}
