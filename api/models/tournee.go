package models

import "time"

type Tournee struct {
	ID           int        `json:"id"`
	Reference    string     `json:"reference"`
	VehiculeID   *int       `json:"vehicule_id"`
	ChauffeurID  *int       `json:"chauffeur_id"`
	DatePrevue   time.Time  `json:"date_prevue"`
	DateRealisee *time.Time `json:"date_realisee"`
	Statut       string     `json:"statut"`
	Commentaire  string     `json:"commentaire"`
}

type Arret struct {
	ID                int    `json:"id"`
	BeneficiaireID    int    `json:"beneficiaire_id"`
	BeneficiaireNom   string `json:"beneficiaire_nom"`
	BeneficiaireVille string `json:"beneficiaire_ville"`
	OrdrePassage      int    `json:"ordre_passage"`
	HeurePrevue       string `json:"heure_prevue"`
	Livre             bool   `json:"livre"`
}

type LigneTournee struct {
	ProduitID int    `json:"produit_id"`
	CodeBarre string `json:"code_barre"`
	Libelle   string `json:"libelle"`
	Quantite  int    `json:"quantite"`
}
