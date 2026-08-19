package models

import "time"

type Benevole struct {
	ID                int       `json:"id"`
	UtilisateurID     int       `json:"utilisateur_id"`
	Nom               string    `json:"nom"`
	Prenom            string    `json:"prenom"`
	Email             string    `json:"email"`
	StatutCandidature string    `json:"statut_candidature"`
	DateCandidature   time.Time `json:"date_candidature"`
	Disponibilites    string    `json:"disponibilites"`
	Competences       []string  `json:"competences"`
}

type Competence struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Libelle string `json:"libelle"`
}
