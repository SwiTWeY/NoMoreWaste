package bdd

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func ListBenevoles(db *sql.DB) ([]models.Benevole, error) {
	rows, err := db.Query(`
		SELECT pb.id, pb.utilisateur_id, u.nom, u.prenom, u.email,
		       pb.statut_candidature, pb.date_candidature, COALESCE(pb.disponibilites, ''),
		       COALESCE(array_agg(c.libelle) FILTER (WHERE c.id IS NOT NULL), '{}') AS competences
		FROM profil_benevole pb
		JOIN utilisateur u ON u.id = pb.utilisateur_id
		LEFT JOIN benevole_competence bc ON bc.profil_benevole_id = pb.id
		LEFT JOIN competence c ON c.id = bc.competence_id
		GROUP BY pb.id, u.nom, u.prenom, u.email
		ORDER BY pb.date_candidature DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	benevoles := []models.Benevole{}
	for rows.Next() {
		var b models.Benevole
		if err := rows.Scan(&b.ID, &b.UtilisateurID, &b.Nom, &b.Prenom, &b.Email,
			&b.StatutCandidature, &b.DateCandidature, &b.Disponibilites, pq.Array(&b.Competences)); err != nil {
			return nil, err
		}
		benevoles = append(benevoles, b)
	}
	return benevoles, rows.Err()
}

func ChangerStatutBenevole(db *sql.DB, id int, statut string) error {
	_, err := db.Exec(`UPDATE profil_benevole SET statut_candidature = $1 WHERE id = $2`, statut, id)
	return err
}

func ListCompetences(db *sql.DB) ([]models.Competence, error) {
	rows, err := db.Query(`SELECT id, code, libelle FROM competence ORDER BY libelle`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	competences := []models.Competence{}
	for rows.Next() {
		var c models.Competence
		if err := rows.Scan(&c.ID, &c.Code, &c.Libelle); err != nil {
			return nil, err
		}
		competences = append(competences, c)
	}
	return competences, rows.Err()
}

func AjouterCompetenceBenevole(db *sql.DB, benevoleID, competenceID int) error {
	_, err := db.Exec(`
		INSERT INTO benevole_competence (profil_benevole_id, competence_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`, benevoleID, competenceID)
	return err
}

func GetProfilBenevoleParUtilisateur(db *sql.DB, utilisateurID int) (int, string, error) {
	var id int
	var statut string
	err := db.QueryRow(`
		SELECT id, statut_candidature
		FROM profil_benevole
		WHERE utilisateur_id = $1`, utilisateurID).Scan(&id, &statut)
	return id, statut, err
}

func CreerAffectation(db *sql.DB, creneauID, profilBenevoleID int) error {
	_, err := db.Exec(`
		INSERT INTO affectation (creneau_id, profil_benevole_id, statut)
		VALUES ($1, $2, 'proposee')
		ON CONFLICT (creneau_id, profil_benevole_id) DO NOTHING`,
		creneauID, profilBenevoleID)
	return err
}

func CreerCandidatureBenevole(db *sql.DB, utilisateurID int, disponibilites string) error {
	_, err := db.Exec(`
		INSERT INTO profil_benevole (utilisateur_id, statut_candidature, disponibilites)
		VALUES ($1, 'candidat', $2)
		ON CONFLICT (utilisateur_id) DO NOTHING`,
		utilisateurID, disponibilites)
	return err
}
