package models

import "time"

type Utilisateur struct {
	ID           int       `json:"id"`
	Nom          string    `json:"nom"`
	Prenom       string    `json:"prenom"`
	Email        string    `json:"email"`
	Telephone    string    `json:"telephone"`
	EstPersonnel bool      `json:"est_personnel"`
	LanguePref   string    `json:"langue_pref"`
	CreatedAt    time.Time `json:"created_at"`
}
