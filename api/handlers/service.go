package handlers

import (
	"database/sql"
	"net/http"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/export"
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
