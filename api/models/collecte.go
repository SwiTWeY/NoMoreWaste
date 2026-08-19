package models

import "time"

type Collecte struct {
	ID              int        `json:"id"`
	CommercantID    *int       `json:"commercant_id"`
	DonateurID      *int       `json:"donateur_id"`
	VehiculeID      *int       `json:"vehicule_id"`
	ChauffeurID     *int       `json:"chauffeur_id"`
	AdresseCollecte string     `json:"adresse_collecte"`
	DatePrevue      time.Time  `json:"date_prevue"`
	DateRealisee    *time.Time `json:"date_realisee"`
	Statut          string     `json:"statut"`
	Commentaire     string     `json:"commentaire"`
	CreatedAt       time.Time  `json:"created_at"`
}

type LigneCollecte struct {
	ProduitID int    `json:"produit_id"`
	CodeBarre string `json:"code_barre"`
	Libelle   string `json:"libelle"`
	Quantite  int    `json:"quantite"`
}
