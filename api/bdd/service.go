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
