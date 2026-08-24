package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/export"
	"github.com/SwiTWeY/NoMoreWaste/api/middleware"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type ServiceHandler struct {
	DB *sql.DB
}

func langueDe(r *http.Request) string {
	l := r.URL.Query().Get("lang")
	if l == "" {
		return "fr"
	}
	return l
}

func (h ServiceHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := bdd.ListServices(h.DB, langueDe(r))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, services)
}

func (h ServiceHandler) ListCreneaux(w http.ResponseWriter, r *http.Request) {
	creneaux, err := bdd.ListCreneaux(h.DB, langueDe(r))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, creneaux)
}

func (h ServiceHandler) ExportPlanning(w http.ResponseWriter, r *http.Request) {
	creneaux, err := bdd.ListCreneaux(h.DB, langueDe(r))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	data, err := export.PlanningXLSX(creneaux)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="planning.xlsx"`)
	w.Write(data)
}

func (h ServiceHandler) Inscrire(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsDepuis(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "non authentifie")
		return
	}
	creneauID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	actif, err := bdd.EstAdherentActif(h.DB, claims.UtilisateurID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !actif {
		utils.Error(w, http.StatusForbidden, "reserve aux adherents a jour de cotisation")
		return
	}
	if err := bdd.CreerInscription(h.DB, creneauID, claims.UtilisateurID); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "inscription enregistree"})
}

func (h ServiceHandler) MonAgenda(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsDepuis(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "non authentifie")
		return
	}
	evenements, err := bdd.MonAgenda(h.DB, claims.UtilisateurID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, evenements)
}
