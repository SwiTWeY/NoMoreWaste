package bdd

import (
	"database/sql"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func ListTournees(db *sql.DB) ([]models.Tournee, error) {
	rows, err := db.Query(`
		SELECT id, reference, vehicule_id, chauffeur_id, date_prevue, date_realisee, statut, COALESCE(commentaire, '')
		FROM tournee
		ORDER BY date_prevue DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournees := []models.Tournee{}
	for rows.Next() {
		var t models.Tournee
		if err := rows.Scan(&t.ID, &t.Reference, &t.VehiculeID, &t.ChauffeurID, &t.DatePrevue, &t.DateRealisee, &t.Statut, &t.Commentaire); err != nil {
			return nil, err
		}
		tournees = append(tournees, t)
	}
	return tournees, rows.Err()
}

func GetTournee(db *sql.DB, id int) (models.Tournee, error) {
	var t models.Tournee
	err := db.QueryRow(`
		SELECT id, reference, vehicule_id, chauffeur_id, date_prevue, date_realisee, statut, COALESCE(commentaire, '')
		FROM tournee
		WHERE id = $1`, id).
		Scan(&t.ID, &t.Reference, &t.VehiculeID, &t.ChauffeurID, &t.DatePrevue, &t.DateRealisee, &t.Statut, &t.Commentaire)
	return t, err
}

func ListArrets(db *sql.DB, tourneeID int) ([]models.Arret, error) {
	rows, err := db.Query(`
		SELECT a.id, a.beneficiaire_id, b.nom, b.ville, a.ordre_passage,
		       COALESCE(to_char(a.heure_prevue, 'HH24:MI'), ''), a.livre
		FROM arret a
		JOIN beneficiaire b ON b.id = a.beneficiaire_id
		WHERE a.tournee_id = $1
		ORDER BY a.ordre_passage`, tourneeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	arrets := []models.Arret{}
	for rows.Next() {
		var a models.Arret
		if err := rows.Scan(&a.ID, &a.BeneficiaireID, &a.BeneficiaireNom, &a.BeneficiaireVille, &a.OrdrePassage, &a.HeurePrevue, &a.Livre); err != nil {
			return nil, err
		}
		arrets = append(arrets, a)
	}
	return arrets, rows.Err()
}

func ListLignesTournee(db *sql.DB, tourneeID int) ([]models.LigneTournee, error) {
	rows, err := db.Query(`
		SELECT tp.produit_id, p.code_barre, p.libelle, tp.quantite
		FROM tournee_produit tp
		JOIN produit p ON p.id = tp.produit_id
		WHERE tp.tournee_id = $1
		ORDER BY p.libelle`, tourneeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lignes := []models.LigneTournee{}
	for rows.Next() {
		var l models.LigneTournee
		if err := rows.Scan(&l.ProduitID, &l.CodeBarre, &l.Libelle, &l.Quantite); err != nil {
			return nil, err
		}
		lignes = append(lignes, l)
	}
	return lignes, rows.Err()
}

func CreateTournee(db *sql.DB, t models.Tournee) (models.Tournee, error) {
	err := db.QueryRow(`
		INSERT INTO tournee (reference, vehicule_id, chauffeur_id, date_prevue, statut, commentaire)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		t.Reference, t.VehiculeID, t.ChauffeurID, t.DatePrevue, t.Statut, t.Commentaire).
		Scan(&t.ID)
	return t, err
}
