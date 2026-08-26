package bdd

import (
	"database/sql"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func ListServices(db *sql.DB, langue string) ([]models.Service, error) {
	rows, err := db.Query(`
		SELECT s.id, s.code, st.libelle, COALESCE(st.description, ''), s.actif
		FROM service s
		JOIN langue l ON l.code = $1
		JOIN service_traduction st ON st.service_id = s.id AND st.langue_id = l.id
		WHERE s.actif = TRUE
		ORDER BY st.libelle`, langue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []models.Service{}
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Code, &s.Libelle, &s.Description, &s.Actif); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, rows.Err()
}

func ListCreneaux(db *sql.DB, langue string) ([]models.Creneau, error) {
	rows, err := db.Query(`
		SELECT c.id, c.service_id, st.libelle, c.date_creneau,
		       to_char(c.heure_debut, 'HH24:MI'), to_char(c.heure_fin, 'HH24:MI'),
		       COALESCE(c.lieu, ''), c.capacite_max, c.statut
		FROM creneau c
		JOIN service s ON s.id = c.service_id
		JOIN langue l ON l.code = $1
		JOIN service_traduction st ON st.service_id = s.id AND st.langue_id = l.id
		ORDER BY c.date_creneau, c.heure_debut`, langue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	creneaux := []models.Creneau{}
	for rows.Next() {
		var c models.Creneau
		if err := rows.Scan(&c.ID, &c.ServiceID, &c.ServiceLibelle, &c.DateCreneau,
			&c.HeureDebut, &c.HeureFin, &c.Lieu, &c.CapaciteMax, &c.Statut); err != nil {
			return nil, err
		}
		creneaux = append(creneaux, c)
	}
	return creneaux, rows.Err()
}

func ListCreneauxDisponibles(db *sql.DB, langue string, utilisateurID int) ([]models.Creneau, error) {
	rows, err := db.Query(`
		SELECT c.id, c.service_id, st.libelle, c.date_creneau,
		       to_char(c.heure_debut, 'HH24:MI'), to_char(c.heure_fin, 'HH24:MI'),
		       COALESCE(c.lieu, ''), c.capacite_max, c.statut
		FROM creneau c
		JOIN service s ON s.id = c.service_id
		JOIN langue l ON l.code = $1
		JOIN service_traduction st ON st.service_id = s.id AND st.langue_id = l.id
		WHERE NOT EXISTS (
			SELECT 1 FROM inscription i
			WHERE i.creneau_id = c.id AND i.utilisateur_id = $2
		)
		ORDER BY c.date_creneau, c.heure_debut`, langue, utilisateurID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	creneaux := []models.Creneau{}
	for rows.Next() {
		var c models.Creneau
		if err := rows.Scan(&c.ID, &c.ServiceID, &c.ServiceLibelle, &c.DateCreneau,
			&c.HeureDebut, &c.HeureFin, &c.Lieu, &c.CapaciteMax, &c.Statut); err != nil {
			return nil, err
		}
		creneaux = append(creneaux, c)
	}
	return creneaux, rows.Err()
}

func CreerInscription(db *sql.DB, creneauID, utilisateurID int) error {
	_, err := db.Exec(`
		INSERT INTO inscription (creneau_id, utilisateur_id)
		VALUES ($1, $2)
		ON CONFLICT (creneau_id, utilisateur_id) DO NOTHING`,
		creneauID, utilisateurID)
	return err
}

func MonAgenda(db *sql.DB, utilisateurID int) ([]models.EvenementAgenda, error) {
	rows, err := db.Query(`
		SELECT 'service' AS type, c.date_creneau AS d, to_char(c.heure_debut, 'HH24:MI') AS h,
		       st.libelle, COALESCE(c.lieu, ''), i.statut
		FROM inscription i
		JOIN creneau c ON c.id = i.creneau_id
		JOIN service s ON s.id = c.service_id
		JOIN langue l ON l.code = 'fr'
		JOIN service_traduction st ON st.service_id = s.id AND st.langue_id = l.id
		WHERE i.utilisateur_id = $1
		UNION ALL
		SELECT 'collecte', co.date_prevue::date, to_char(co.date_prevue, 'HH24:MI'),
		       'Collecte', COALESCE(co.adresse_collecte, ''), co.statut
		FROM collecte co
		WHERE co.donateur_id = $1
		   OR co.commercant_id IN (SELECT id FROM commercant WHERE utilisateur_id = $1)
		ORDER BY d, h`, utilisateurID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evenements := []models.EvenementAgenda{}
	for rows.Next() {
		var e models.EvenementAgenda
		if err := rows.Scan(&e.Type, &e.Date, &e.Heure, &e.Libelle, &e.Lieu, &e.Statut); err != nil {
			return nil, err
		}
		evenements = append(evenements, e)
	}
	return evenements, rows.Err()
}
