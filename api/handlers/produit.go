package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/models"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type ProduitHandler struct {
	DB *sql.DB
}

func (h ProduitHandler) List(w http.ResponseWriter, r *http.Request) {
	produits, err := bdd.ListProduits(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, produits)
}

func (h ProduitHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	p, err := bdd.GetProduit(h.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "produit introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, p)
}

func (h ProduitHandler) GetParCodeBarre(w http.ResponseWriter, r *http.Request) {
	p, err := bdd.GetProduitParCodeBarre(h.DB, r.PathValue("code"))
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "produit introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, p)
}

func (h ProduitHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Produit
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if p.CodeBarre == "" || p.Libelle == "" {
		utils.Error(w, http.StatusBadRequest, "code_barre et libelle sont requis")
		return
	}
	if p.Unite == "" {
		p.Unite = "piece"
	}
	cree, err := bdd.CreateProduit(h.DB, p)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, cree)
}

func (h ProduitHandler) Stock(w http.ResponseWriter, r *http.Request) {
	stock, err := bdd.ListStock(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, stock)
}

func (h ProduitHandler) StockParCodeBarre(w http.ResponseWriter, r *http.Request) {
	s, err := bdd.GetStockParCodeBarre(h.DB, r.PathValue("code"))
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "produit introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, s)
}
