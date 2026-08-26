package bdd

import (
	"database/sql"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func ListAdhesions(db *sql.DB) ([]models.Adhesion, error) {
	rows, err := db.Query(`
		SELECT a.id, a.utilisateur_id, u.nom, u.prenom, a.date_debut, a.date_fin,
		       a.montant, a.statut_paiement, a.rappel_envoye_le, a.created_at
		FROM adhesion a
		JOIN utilisateur u ON u.id = a.utilisateur_id
		ORDER BY a.date_fin DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adhesions := []models.Adhesion{}
	for rows.Next() {
		var a models.Adhesion
		if err := rows.Scan(&a.ID, &a.UtilisateurID, &a.Nom, &a.Prenom, &a.DateDebut, &a.DateFin,
			&a.Montant, &a.StatutPaiement, &a.RappelEnvoyeLe, &a.CreatedAt); err != nil {
			return nil, err
		}
		adhesions = append(adhesions, a)
	}
	return adhesions, rows.Err()
}

func GetAdhesion(db *sql.DB, id int) (models.Adhesion, error) {
	var a models.Adhesion
	err := db.QueryRow(`
		SELECT a.id, a.utilisateur_id, u.nom, u.prenom, a.date_debut, a.date_fin,
		       a.montant, a.statut_paiement, a.rappel_envoye_le, a.created_at
		FROM adhesion a
		JOIN utilisateur u ON u.id = a.utilisateur_id
		WHERE a.id = $1`, id).
		Scan(&a.ID, &a.UtilisateurID, &a.Nom, &a.Prenom, &a.DateDebut, &a.DateFin,
			&a.Montant, &a.StatutPaiement, &a.RappelEnvoyeLe, &a.CreatedAt)
	return a, err
}

func CreateAdhesion(db *sql.DB, a models.Adhesion) (models.Adhesion, error) {
	err := db.QueryRow(`
		INSERT INTO adhesion (utilisateur_id, date_debut, date_fin, montant, statut_paiement)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		a.UtilisateurID, a.DateDebut, a.DateFin, a.Montant, a.StatutPaiement).
		Scan(&a.ID, &a.CreatedAt)
	return a, err
}

func ChangerStatutAdhesion(db *sql.DB, id int, statut string) error {
	_, err := db.Exec(`UPDATE adhesion SET statut_paiement = $1 WHERE id = $2`, statut, id)
	return err
}

func AdhesionsARappeler(db *sql.DB, joursAvant int) ([]models.RappelAdhesion, error) {
	rows, err := db.Query(`
		SELECT a.id, u.email, u.nom, u.prenom, a.date_fin, u.langue_pref
		FROM adhesion a
		JOIN utilisateur u ON u.id = a.utilisateur_id
		WHERE a.rappel_envoye_le IS NULL
		  AND a.date_fin BETWEEN CURRENT_DATE AND CURRENT_DATE + ($1 * INTERVAL '1 day')
		ORDER BY a.date_fin`, joursAvant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rappels := []models.RappelAdhesion{}
	for rows.Next() {
		var r models.RappelAdhesion
		if err := rows.Scan(&r.AdhesionID, &r.Email, &r.Nom, &r.Prenom, &r.DateFin, &r.LanguePref); err != nil {
			return nil, err
		}
		rappels = append(rappels, r)
	}
	return rappels, rows.Err()
}

func MarquerRappelEnvoye(db *sql.DB, adhesionID int) error {
	_, err := db.Exec(`UPDATE adhesion SET rappel_envoye_le = NOW() WHERE id = $1`, adhesionID)
	return err
}

func CreerAdhesionPayee(db *sql.DB, utilisateurID int, montant float64, sessionID string) error {
	_, err := db.Exec(`
		INSERT INTO adhesion (utilisateur_id, date_debut, date_fin, montant, statut_paiement, stripe_session_id)
		VALUES ($1, CURRENT_DATE, CURRENT_DATE + INTERVAL '1 year', $2, 'paye', $3)
		ON CONFLICT (stripe_session_id) DO NOTHING`,
		utilisateurID, montant, sessionID)
	return err
}
