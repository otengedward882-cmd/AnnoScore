# 🎯 AnnoScore

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Domain](https://img.shields.io/badge/focus-AI%20Data%20Annotation%20%26%20RLHF-blue.svg)](#why-annoscore)
[![Zero Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen.svg)](#prerequisites)

> **High-Performance Inter-Annotator Agreement (IAA) & Quality Scoring Engine for AI Datasets.**  
> Written in pure Go with zero external dependencies.

```text
=================================================================
  🎯 AnnoScore — AI Dataset Inter-Annotator Agreement Report
=================================================================

• Total Evaluated Items:     8
• Active Annotators:         3
• Observed Agreement Rate:   87.50%
• Cohen's Kappa (κ):         0.783
• Quality Classification:    Substantial Agreement (Production Quality)

-----------------------------------------------------------------
  Individual Annotator Consistency vs Consensus:
-----------------------------------------------------------------
  [✅ RELIABLE] evaluator_alice  Consensus Match: 100.0%
  [✅ RELIABLE] evaluator_carol  Consensus Match: 100.0%
  [⚠️  OUTLIER ] evaluator_bob    Consensus Match: 75.0%

-----------------------------------------------------------------
  High-Ambiguity Items Requiring Guideline Clarification (2):
-----------------------------------------------------------------
  • Item ID: prompt_104 (Annotator Disagreement Detected)
  • Item ID: prompt_107 (Annotator Disagreement Detected)

=================================================================
```

---

## ⚡ Why AnnoScore?

In remote AI data annotation, RLHF (Reinforcement Learning from Human Feedback), and model safety evaluation, **data quality is everything**.

When remote evaluators rate chatbot safety or prompt helpfulness, companies need to verify:
1. **Are the guidelines clear?** High disagreement indicates ambiguous rubric instructions.
2. **Are raters aligned?** Low Inter-Annotator Agreement (IAA) means the data cannot be trusted to train or fine-tune models.
3. **Who are the outlier evaluators?** Identify annotators whose ratings consistently deviate from consensus.

**AnnoScore** ingests CSV/JSONL annotation datasets, computes statistical consensus metrics (**Cohen's Kappa**, **Fleiss' Kappa**, **Observed Agreement**), detects inconsistent raters, and isolates ambiguous items that need guideline refinement.

---

## ✨ Key Features

- **Statistical Rigor**:
  - **Cohen's Kappa ($\kappa$)**: Corrects for chance agreement in pairwise evaluator ratings.
  - **Fleiss' Kappa**: Multi-rater generalized agreement across arbitrary categories.
  - **Landis & Koch Interpretation**: Automated benchmark grading from *Near Perfect* to *Poor*.
- **Outlier Rater Identification**: Flags individual annotators with sub-threshold consensus correlation.
- **Ambiguity Detection**: Highlights exact prompt item IDs where raters split, pinpointing flaws in labeling guidelines.
- **Zero Third-Party Dependencies**: Pure Go standard library. Compiles to a single lightweight, high-speed binary.
- **Multiple Output Formats**: Human-readable terminal reports, Markdown summary tables for documentation, and structured JSON for automated pipelines.

---

## 🚀 Quick Start

### 1. Installation
Clone the repository and build with the Go toolchain:
```bash
git clone https://github.com/otengedward882-cmd/AnnoScore.git
cd AnnoScore
go build -o annoscore cmd/annoscore/main.go
```

### 2. Run Quality Audit on a Dataset
```bash
./annoscore --file examples/ai_safety_eval.csv
```

### 3. Generate Markdown Report for Documentation
```bash
./annoscore --file examples/ai_safety_eval.csv --format markdown > report.md
```

### 4. Export JSON for Automated CI Pipelines
```bash
./annoscore --file examples/ai_safety_eval.csv --format json
```

---

## 📖 CLI Flags

| Flag | Default | Description |
|---|---|---|
| `--file <PATH>` | `""` | Path to the annotation CSV dataset |
| `--format <FMT>` | `text` | Output format (`text`, `markdown`, `json`) |
| `--min-kappa <K>` | `0.60` | Minimum passing Cohen's Kappa score (exits `1` if below) |
| `--version` | `false` | Display version information |

---

## 🧪 Running Unit Tests

Execute the test suite with Go's built-in testing tool:
```bash
go test ./...
```

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.
