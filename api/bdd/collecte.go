package bdd

import (
	"database/sql"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func ListCollectes(db *sql.DB) ([]models.Collecte, error) {
	rows, err := db.Query(`
		SELECT id, commercant_id, donateur_id, vehicule_id, chauffeur_id,
		       COALESCE(adresse_collecte, ''), date_prevue, date_realisee, statut, COALESCE(commentaire, ''), created_at
		FROM collecte
		ORDER BY date_prevue DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collectes := []models.Collecte{}
	for rows.Next() {
		var c models.Collecte
		if err := rows.Scan(&c.ID, &c.CommercantID, &c.DonateurID, &c.VehiculeID, &c.ChauffeurID,
			&c.AdresseCollecte, &c.DatePrevue, &c.DateRealisee, &c.Statut, &c.Commentaire, &c.CreatedAt); err != nil {
			return nil, err
		}
		collectes = append(collectes, c)
	}
	return collectes, rows.Err()
}

func GetCollecte(db *sql.DB, id int) (models.Collecte, error) {
	var c models.Collecte
	err := db.QueryRow(`
		SELECT id, commercant_id, donateur_id, vehicule_id, chauffeur_id,
		       COALESCE(adresse_collecte, ''), date_prevue, date_realisee, statut, COALESCE(commentaire, ''), created_at
		FROM collecte
		WHERE id = $1`, id).
		Scan(&c.ID, &c.CommercantID, &c.DonateurID, &c.VehiculeID, &c.ChauffeurID,
			&c.AdresseCollecte, &c.DatePrevue, &c.DateRealisee, &c.Statut, &c.Commentaire, &c.CreatedAt)
	return c, err
}

func CreateCollecte(db *sql.DB, c models.Collecte) (models.Collecte, error) {
	err := db.QueryRow(`
		INSERT INTO collecte (commercant_id, donateur_id, vehicule_id, chauffeur_id, adresse_collecte, date_prevue, statut, commentaire)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`,
		c.CommercantID, c.DonateurID, c.VehiculeID, c.ChauffeurID, c.AdresseCollecte, c.DatePrevue, c.Statut, c.Commentaire).
		Scan(&c.ID, &c.CreatedAt)
	return c, err
}

func AjouterProduitCollecte(db *sql.DB, collecteID, produitID, quantite int) error {
	_, err := db.Exec(`
		INSERT INTO collecte_produit (collecte_id, produit_id, quantite)
		VALUES ($1, $2, $3)
		ON CONFLICT (collecte_id, produit_id)
		DO UPDATE SET quantite = collecte_produit.quantite + EXCLUDED.quantite`,
		collecteID, produitID, quantite)
	return err
}

func ListLignesCollecte(db *sql.DB, collecteID int) ([]models.LigneCollecte, error) {
	rows, err := db.Query(`
		SELECT cp.produit_id, p.code_barre, p.libelle, cp.quantite
		FROM collecte_produit cp
		JOIN produit p ON p.id = cp.produit_id
		WHERE cp.collecte_id = $1
		ORDER BY p.libelle`, collecteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lignes := []models.LigneCollecte{}
	for rows.Next() {
		var l models.LigneCollecte
		if err := rows.Scan(&l.ProduitID, &l.CodeBarre, &l.Libelle, &l.Quantite); err != nil {
			return nil, err
		}
		lignes = append(lignes, l)
	}
	return lignes, rows.Err()
}
