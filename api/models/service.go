package models

import "time"

type Service struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Libelle     string `json:"libelle"`
	Description string `json:"description"`
	Actif       bool   `json:"actif"`
}

type Creneau struct {
	ID             int       `json:"id"`
	ServiceID      int       `json:"service_id"`
	ServiceLibelle string    `json:"service_libelle"`
	DateCreneau    time.Time `json:"date_creneau"`
	HeureDebut     string    `json:"heure_debut"`
	HeureFin       string    `json:"heure_fin"`
	Lieu           string    `json:"lieu"`
	CapaciteMax    int       `json:"capacite_max"`
	Statut         string    `json:"statut"`
}

type EvenementAgenda struct {
	Type    string    `json:"type"`
	Date    time.Time `json:"date"`
	Heure   string    `json:"heure"`
	Libelle string    `json:"libelle"`
	Lieu    string    `json:"lieu"`
	Statut  string    `json:"statut"`
}
