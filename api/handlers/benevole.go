package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/middleware"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type BenevoleHandler struct {
	DB *sql.DB
}

var statutsBenevole = map[string]bool{
	"candidat": true,
	"valide":   true,
	"refuse":   true,
	"inactif":  true,
}

func (h BenevoleHandler) List(w http.ResponseWriter, r *http.Request) {
	benevoles, err := bdd.ListBenevoles(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, benevoles)
}

func (h BenevoleHandler) ChangerStatut(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	var req struct {
		Statut string `json:"statut"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if !statutsBenevole[req.Statut] {
		utils.Error(w, http.StatusBadRequest, "statut invalide (candidat, valide, refuse, inactif)")
		return
	}
	if err := bdd.ChangerStatutBenevole(h.DB, id, req.Statut); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"statut": req.Statut})
}

func (h BenevoleHandler) Competences(w http.ResponseWriter, r *http.Request) {
	competences, err := bdd.ListCompetences(h.DB)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, competences)
}

func (h BenevoleHandler) AjouterCompetence(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id invalide")
		return
	}
	var req struct {
		CompetenceID int `json:"competence_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if req.CompetenceID == 0 {
		utils.Error(w, http.StatusBadRequest, "competence_id est requis")
		return
	}
	if err := bdd.AjouterCompetenceBenevole(h.DB, id, req.CompetenceID); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "competence ajoutee"})
}

func (h BenevoleHandler) ProposerAnimation(w http.ResponseWriter, r *http.Request) {
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
	profilID, statut, err := bdd.GetProfilBenevoleParUtilisateur(h.DB, claims.UtilisateurID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusForbidden, "reserve aux benevoles")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if statut != "valide" {
		utils.Error(w, http.StatusForbidden, "votre candidature de benevole n'est pas encore validee")
		return
	}
	if err := bdd.CreerAffectation(h.DB, creneauID, profilID); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "proposition enregistree"})
}

func (h BenevoleHandler) Postuler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsDepuis(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "non authentifie")
		return
	}
	var req struct {
		Disponibilites string `json:"disponibilites"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if err := bdd.CreerCandidatureBenevole(h.DB, claims.UtilisateurID, req.Disponibilites); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "candidature enregistree"})
}
