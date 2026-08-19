package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type UtilisateurHandler struct {
	DB *sql.DB
}

func (h UtilisateurHandler) AdherentActif(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	actif, err := bdd.EstAdherentActif(h.DB, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]bool{"adherent_actif": actif})
}
