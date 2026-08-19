package bdd

import (
	"database/sql"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func EstAdherentActif(db *sql.DB, utilisateurID int) (bool, error) {
	var actif bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM adhesion
			WHERE utilisateur_id = $1
			  AND statut_paiement = 'paye'
			  AND date_debut <= CURRENT_DATE
			  AND date_fin   >= CURRENT_DATE
		)`, utilisateurID).Scan(&actif)
	return actif, err
}

func CreateUtilisateur(db *sql.DB, u models.Utilisateur, motDePasseHash string) (models.Utilisateur, error) {
	err := db.QueryRow(`
		INSERT INTO utilisateur (nom, prenom, email, mot_de_passe, telephone, est_personnel, langue_pref)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`,
		u.Nom, u.Prenom, u.Email, motDePasseHash, u.Telephone, u.EstPersonnel, u.LanguePref).
		Scan(&u.ID, &u.CreatedAt)
	return u, err
}

func GetUtilisateurParEmail(db *sql.DB, email string) (models.Utilisateur, string, error) {
	var u models.Utilisateur
	var hash string
	err := db.QueryRow(`
		SELECT id, nom, prenom, email, mot_de_passe, COALESCE(telephone, ''), est_personnel, langue_pref, created_at
		FROM utilisateur
		WHERE email = $1`, email).
		Scan(&u.ID, &u.Nom, &u.Prenom, &u.Email, &hash, &u.Telephone, &u.EstPersonnel, &u.LanguePref, &u.CreatedAt)
	return u, hash, err
}
