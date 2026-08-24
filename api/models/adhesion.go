package models

import "time"

type Adhesion struct {
	ID             int        `json:"id"`
	UtilisateurID  int        `json:"utilisateur_id"`
	Nom            string     `json:"nom"`
	Prenom         string     `json:"prenom"`
	DateDebut      time.Time  `json:"date_debut"`
	DateFin        time.Time  `json:"date_fin"`
	Montant        float64    `json:"montant"`
	StatutPaiement string     `json:"statut_paiement"`
	RappelEnvoyeLe *time.Time `json:"rappel_envoye_le"`
	CreatedAt      time.Time  `json:"created_at"`
}

type RappelAdhesion struct {
	AdhesionID int
	Email      string
	Nom        string
	Prenom     string
	DateFin    time.Time
	LanguePref string
}
