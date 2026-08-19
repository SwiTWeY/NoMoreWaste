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

type CollecteHandler struct {
	DB *sql.DB
}

func (h CollecteHandler) List(w http.ResponseWriter, r *http.Request) {
	collectes, err := bdd.ListCollectes(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, collectes)
}

func (h CollecteHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	c, err := bdd.GetCollecte(h.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "collecte introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, c)
}

func (h CollecteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Collecte
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if (c.CommercantID == nil) == (c.DonateurID == nil) {
		utils.Error(w, http.StatusBadRequest, "renseigner exactement une source : commercant_id OU donateur_id")
		return
	}
	if c.DatePrevue.IsZero() {
		utils.Error(w, http.StatusBadRequest, "date_prevue est requise")
		return
	}
	if c.Statut == "" {
		c.Statut = "planifiee"
	}
	cree, err := bdd.CreateCollecte(h.DB, c)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, cree)
}

func (h CollecteHandler) Lignes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	lignes, err := bdd.ListLignesCollecte(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, lignes)
}

type ligneRequete struct {
	CodeBarre string `json:"code_barre"`
	Quantite  int    `json:"quantite"`
}

func (h CollecteHandler) AjouterProduit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}

	var req ligneRequete
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if req.Quantite <= 0 {
		utils.Error(w, http.StatusBadRequest, "quantite doit etre positive")
		return
	}

	produit, err := bdd.GetProduitParCodeBarre(h.DB, req.CodeBarre)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "code-barres inconnu")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := bdd.AjouterProduitCollecte(h.DB, id, produit.ID, req.Quantite); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	lignes, err := bdd.ListLignesCollecte(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, lignes)
}
