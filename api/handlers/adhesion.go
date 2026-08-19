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

type AdhesionHandler struct {
	DB *sql.DB
}

func (h AdhesionHandler) List(w http.ResponseWriter, r *http.Request) {
	adhesions, err := bdd.ListAdhesions(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, adhesions)
}

func (h AdhesionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	a, err := bdd.GetAdhesion(h.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "adhesion introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, a)
}

func (h AdhesionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a models.Adhesion
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if a.UtilisateurID == 0 || a.DateDebut.IsZero() || a.DateFin.IsZero() {
		utils.Error(w, http.StatusBadRequest, "utilisateur_id, date_debut et date_fin sont requis")
		return
	}
	if a.StatutPaiement == "" {
		a.StatutPaiement = "en_attente"
	}
	created, err := bdd.CreateAdhesion(h.DB, a)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, created)
}
