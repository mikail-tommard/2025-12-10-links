package report

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
)

type PDFReportGenerator struct {
	title string
	dateFormat string
}

func NewReportGenerator(title string, dateFormat string) *PDFReportGenerator {
	return &PDFReportGenerator{
		title: title,
		dateFormat: dateFormat,
	}
}

func (g *PDFReportGenerator) GenerateReport(batches []*domain.LinkBatch) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, g.title)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "Generated at: " + time.Now().Format(g.dateFormat))
	pdf.Ln(10)

	for _, batch := range batches {
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(0, 6, fmt.Sprintf("Batch ID: %d", batch.ID))
    pdf.Ln(6)

    pdf.SetFont("Arial", "", 10)
    pdf.Cell(0, 5, fmt.Sprintf("Status: %s", batch.Status))
    pdf.Ln(5)
    pdf.Cell(0, 5, fmt.Sprintf("Created at: %s", batch.CreatedAt.Format(g.dateFormat)))
    pdf.Ln(5)
    pdf.Cell(0, 5, fmt.Sprintf("Updated at: %s", batch.UpdatedAt.Format(g.dateFormat)))
    pdf.Ln(8)

    pdf.SetFont("Arial", "B", 10)
    pdf.CellFormat(80, 6, "URL", "1", 0, "", false, 0, "")
    pdf.CellFormat(30, 6, "STATUS", "1", 0, "", false, 0, "")
    pdf.CellFormat(80, 6, "ERROR", "1", 0, "", false, 0, "")
    pdf.Ln(6)

    pdf.SetFont("Arial", "", 9)
    for _, res := range batch.Results {
        pdf.CellFormat(80, 5, res.Link.URL, "1", 0, "", false, 0, "")
        pdf.CellFormat(30, 5, string(res.Status), "1", 0, "", false, 0, "")
        pdf.CellFormat(80, 5, res.Error, "1", 0, "", false, 0, "")
        pdf.Ln(5)
    }

    pdf.Ln(10)
	}
	
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}