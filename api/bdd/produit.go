package bdd

import (
	"database/sql"

	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

func ListProduits(db *sql.DB) ([]models.Produit, error) {
	rows, err := db.Query(`
		SELECT id, code_barre, libelle, COALESCE(categorie, ''), unite, date_limite, created_at
		FROM produit
		ORDER BY libelle`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	produits := []models.Produit{}
	for rows.Next() {
		var p models.Produit
		if err := rows.Scan(&p.ID, &p.CodeBarre, &p.Libelle, &p.Categorie, &p.Unite, &p.DateLimite, &p.CreatedAt); err != nil {
			return nil, err
		}
		produits = append(produits, p)
	}
	return produits, rows.Err()
}

func GetProduit(db *sql.DB, id int) (models.Produit, error) {
	var p models.Produit
	err := db.QueryRow(`
		SELECT id, code_barre, libelle, COALESCE(categorie, ''), unite, date_limite, created_at
		FROM produit
		WHERE id = $1`, id).
		Scan(&p.ID, &p.CodeBarre, &p.Libelle, &p.Categorie, &p.Unite, &p.DateLimite, &p.CreatedAt)
	return p, err
}

func GetProduitParCodeBarre(db *sql.DB, codeBarre string) (models.Produit, error) {
	var p models.Produit
	err := db.QueryRow(`
		SELECT id, code_barre, libelle, COALESCE(categorie, ''), unite, date_limite, created_at
		FROM produit
		WHERE code_barre = $1`, codeBarre).
		Scan(&p.ID, &p.CodeBarre, &p.Libelle, &p.Categorie, &p.Unite, &p.DateLimite, &p.CreatedAt)
	return p, err
}

func CreateProduit(db *sql.DB, p models.Produit) (models.Produit, error) {
	err := db.QueryRow(`
		INSERT INTO produit (code_barre, libelle, categorie, unite, date_limite)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		p.CodeBarre, p.Libelle, p.Categorie, p.Unite, p.DateLimite).
		Scan(&p.ID, &p.CreatedAt)
	return p, err
}

func ListStock(db *sql.DB) ([]models.StockItem, error) {
	rows, err := db.Query(`
		SELECT id, code_barre, libelle, COALESCE(categorie, ''), date_limite, quantite_stock
		FROM v_stock
		ORDER BY date_limite NULLS LAST`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stock := []models.StockItem{}
	for rows.Next() {
		var s models.StockItem
		if err := rows.Scan(&s.ID, &s.CodeBarre, &s.Libelle, &s.Categorie, &s.DateLimite, &s.QuantiteStock); err != nil {
			return nil, err
		}
		stock = append(stock, s)
	}
	return stock, rows.Err()
}

func GetStockParCodeBarre(db *sql.DB, codeBarre string) (models.StockItem, error) {
	var s models.StockItem
	err := db.QueryRow(`
		SELECT id, code_barre, libelle, COALESCE(categorie, ''), date_limite, quantite_stock
		FROM v_stock
		WHERE code_barre = $1`, codeBarre).
		Scan(&s.ID, &s.CodeBarre, &s.Libelle, &s.Categorie, &s.DateLimite, &s.QuantiteStock)
	return s, err
}

func ListCategories(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT categorie
		FROM produit
		WHERE categorie IS NOT NULL AND categorie <> ''
		ORDER BY categorie`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []string{}
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}
