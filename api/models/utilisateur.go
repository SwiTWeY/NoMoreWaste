package models

import "time"

type Utilisateur struct {
	ID           int       `json:"id"`
	Nom          string    `json:"nom"`
	Prenom       string    `json:"prenom"`
	Email        string    `json:"email"`
	Telephone    string    `json:"telephone"`
	EstPersonnel bool      `json:"est_personnel"`
	Actif        bool      `json:"actif"`
	LanguePref   string    `json:"langue_pref"`
	CreatedAt    time.Time `json:"created_at"`
}

type Stats struct {
	Utilisateurs int `json:"utilisateurs"`
	Adhesions    int `json:"adhesions"`
	Benevoles    int `json:"benevoles"`
	Produits     int `json:"produits"`
	Collectes    int `json:"collectes"`
	Tournees     int `json:"tournees"`
}
