package export

import (
	"bytes"

	"github.com/xuri/excelize/v2"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func PlanningXLSX(creneaux []models.Creneau) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	feuille := "Planning"
	f.SetSheetName("Sheet1", feuille)

	entetes := []string{"Date", "Debut", "Fin", "Service", "Lieu", "Capacite", "Statut"}
	for i, h := range entetes {
		cellule, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(feuille, cellule, h)
	}

	for i, c := range creneaux {
		ligne := i + 2
		valeurs := []interface{}{
			c.DateCreneau.Format("02/01/2006"),
			c.HeureDebut,
			c.HeureFin,
			c.ServiceLibelle,
			c.Lieu,
			c.CapaciteMax,
			c.Statut,
		}
		for col, v := range valeurs {
			cellule, _ := excelize.CoordinatesToCellName(col+1, ligne)
			f.SetCellValue(feuille, cellule, v)
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
