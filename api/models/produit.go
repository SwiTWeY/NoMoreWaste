package models

import "time"

type Produit struct {
	ID         int        `json:"id"`
	CodeBarre  string     `json:"code_barre"`
	Libelle    string     `json:"libelle"`
	Categorie  string     `json:"categorie"`
	Unite      string     `json:"unite"`
	DateLimite *time.Time `json:"date_limite"`
	CreatedAt  time.Time  `json:"created_at"`
}

type StockItem struct {
	ID            int        `json:"id"`
	CodeBarre     string     `json:"code_barre"`
	Libelle       string     `json:"libelle"`
	Categorie     string     `json:"categorie"`
	DateLimite    *time.Time `json:"date_limite"`
	QuantiteStock int        `json:"quantite_stock"`
}
