package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/export"
	"github.com/SwiTWeY/NoMoreWaste/api/models"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type TourneeHandler struct {
	DB *sql.DB
}

func (h TourneeHandler) List(w http.ResponseWriter, r *http.Request) {
	tournees, err := bdd.ListTournees(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, tournees)
}

func (h TourneeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	t, err := bdd.GetTournee(h.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "tournee introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	arrets, err := bdd.ListArrets(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	lignes, err := bdd.ListLignesTournee(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]any{
		"tournee":  t,
		"arrets":   arrets,
		"produits": lignes,
	})
}

func (h TourneeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t models.Tournee
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if t.Reference == "" || t.DatePrevue.IsZero() {
		utils.Error(w, http.StatusBadRequest, "reference et date_prevue sont requis")
		return
	}
	if t.Statut == "" {
		t.Statut = "planifiee"
	}
	cree, err := bdd.CreateTournee(h.DB, t)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, cree)
}

func (h TourneeHandler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	t, err := bdd.GetTournee(h.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusNotFound, "tournee introuvable")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	arrets, err := bdd.ListArrets(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	lignes, err := bdd.ListLignesTournee(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	data, err := export.TourneePDF(t, arrets, lignes)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="tournee-`+t.Reference+`.pdf"`)
	w.Write(data)
}
