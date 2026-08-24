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
		SELECT id, nom, prenom, email, mot_de_passe, COALESCE(telephone, ''), est_personnel, actif, langue_pref, created_at
		FROM utilisateur
		WHERE email = $1`, email).
		Scan(&u.ID, &u.Nom, &u.Prenom, &u.Email, &hash, &u.Telephone, &u.EstPersonnel, &u.Actif, &u.LanguePref, &u.CreatedAt)
	return u, hash, err
}

func ListUtilisateurs(db *sql.DB) ([]models.Utilisateur, error) {
	rows, err := db.Query(`
		SELECT id, nom, prenom, email, COALESCE(telephone, ''), est_personnel, actif, langue_pref, created_at
		FROM utilisateur
		ORDER BY nom, prenom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.Utilisateur{}
	for rows.Next() {
		var u models.Utilisateur
		if err := rows.Scan(&u.ID, &u.Nom, &u.Prenom, &u.Email, &u.Telephone, &u.EstPersonnel, &u.Actif, &u.LanguePref, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func BannirUtilisateur(db *sql.DB, id int, actif bool) error {
	_, err := db.Exec(`UPDATE utilisateur SET actif = $1 WHERE id = $2`, actif, id)
	return err
}

func Stats(db *sql.DB) (models.Stats, error) {
	var s models.Stats
	err := db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM utilisateur),
			(SELECT COUNT(*) FROM adhesion),
			(SELECT COUNT(*) FROM profil_benevole),
			(SELECT COUNT(*) FROM produit),
			(SELECT COUNT(*) FROM collecte),
			(SELECT COUNT(*) FROM tournee)`).
		Scan(&s.Utilisateurs, &s.Adhesions, &s.Benevoles, &s.Produits, &s.Collectes, &s.Tournees)
	return s, err
}
