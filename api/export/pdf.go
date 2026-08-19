package export

import (
	"bytes"
	"fmt"

	"github.com/go-pdf/fpdf"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func TourneePDF(t models.Tournee, arrets []models.Arret, lignes []models.LigneTournee) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "NO MORE WASTE - Recapitulatif de tournee")
	pdf.Ln(14)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 7, "Reference : "+t.Reference)
	pdf.Ln(6)
	pdf.Cell(0, 7, "Date prevue : "+t.DatePrevue.Format("02/01/2006 15:04"))
	pdf.Ln(6)
	pdf.Cell(0, 7, "Statut : "+t.Statut)
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 8, "Arrets / beneficiaires")
	pdf.Ln(9)
	pdf.SetFont("Arial", "", 10)
	for _, a := range arrets {
		etat := "non livre"
		if a.Livre {
			etat = "livre"
		}
		pdf.Cell(0, 6, fmt.Sprintf("%d. %s (%s) - %s - %s", a.OrdrePassage, a.BeneficiaireNom, a.BeneficiaireVille, a.HeurePrevue, etat))
		pdf.Ln(6)
	}
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 8, "Produits distribues")
	pdf.Ln(9)
	pdf.SetFont("Arial", "", 10)
	for _, l := range lignes {
		pdf.Cell(0, 6, fmt.Sprintf("%s (%s) : %d", l.Libelle, l.CodeBarre, l.Quantite))
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
