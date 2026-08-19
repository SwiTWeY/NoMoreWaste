package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/SwiTWeY/NoMoreWaste/api/auth"
	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/models"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type AuthHandler struct {
	DB  *sql.DB
	Cfg config.Config
}

type inscriptionRequete struct {
	Nom        string `json:"nom"`
	Prenom     string `json:"prenom"`
	Email      string `json:"email"`
	MotDePasse string `json:"mot_de_passe"`
	Telephone  string `json:"telephone"`
	LanguePref string `json:"langue_pref"`
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req inscriptionRequete
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}
	if req.Nom == "" || req.Prenom == "" || req.Email == "" || req.MotDePasse == "" {
		utils.Error(w, http.StatusBadRequest, "nom, prenom, email et mot_de_passe sont requis")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	u := models.Utilisateur{
		Nom:          req.Nom,
		Prenom:       req.Prenom,
		Email:        req.Email,
		Telephone:    req.Telephone,
		EstPersonnel: false,
		LanguePref:   req.LanguePref,
	}
	if u.LanguePref == "" {
		u.LanguePref = "fr"
	}

	cree, err := bdd.CreateUtilisateur(h.DB, u, string(hash))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, cree)
}

type connexionRequete struct {
	Email      string `json:"email"`
	MotDePasse string `json:"mot_de_passe"`
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req connexionRequete
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "corps de requete invalide")
		return
	}

	u, hash, err := bdd.GetUtilisateurParEmail(h.DB, req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error(w, http.StatusUnauthorized, "identifiants invalides")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.MotDePasse)); err != nil {
		utils.Error(w, http.StatusUnauthorized, "identifiants invalides")
		return
	}

	token, err := auth.GenererToken(h.Cfg.JWTSecret, u.ID, u.EstPersonnel)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]any{
		"token":       token,
		"utilisateur": u,
	})
}
