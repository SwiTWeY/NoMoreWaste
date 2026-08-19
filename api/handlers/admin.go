package handlers

import (
	"database/sql"
	"net/http"

	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/rappel"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type AdminHandler struct {
	DB  *sql.DB
	Cfg config.Config
}

func (h AdminHandler) DeclencherRappels(w http.ResponseWriter, r *http.Request) {
	n, err := rappel.EnvoyerRappels(r.Context(), h.DB, h.Cfg)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]int{"rappels_envoyes": n})
}
