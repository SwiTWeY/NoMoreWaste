package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/handlers"
	"github.com/SwiTWeY/NoMoreWaste/api/middleware"
	"github.com/SwiTWeY/NoMoreWaste/api/rappel"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := config.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	adhesion := handlers.AdhesionHandler{DB: db}
	utilisateur := handlers.UtilisateurHandler{DB: db}
	admin := handlers.AdminHandler{DB: db, Cfg: cfg}
	authH := handlers.AuthHandler{DB: db, Cfg: cfg}
	produit := handlers.ProduitHandler{DB: db}
	collecte := handlers.CollecteHandler{DB: db}
	benevole := handlers.BenevoleHandler{DB: db}
	service := handlers.ServiceHandler{DB: db}
	tournee := handlers.TourneeHandler{DB: db}
	paiement := handlers.PaiementHandler{DB: db, Cfg: cfg}

	auth := middleware.Auth(cfg.JWTSecret)
	connecte := func(h http.HandlerFunc) http.Handler {
		return auth(h)
	}
	perso := func(h http.HandlerFunc) http.Handler {
		return auth(middleware.Personnel(h))
	}

	mux := http.NewServeMux()

	// --- Public ---
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			utils.Error(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /register", authH.Register)
	mux.HandleFunc("POST /login", authH.Login)
	mux.HandleFunc("GET /services", service.ListServices)
	mux.HandleFunc("GET /creneaux", service.ListCreneaux)

	// --- Adherent connecte (front-office) ---
	mux.Handle("POST /creneaux/{id}/inscription", connecte(service.Inscrire))
	mux.Handle("GET /mon-agenda", connecte(service.MonAgenda))
	mux.Handle("GET /creneaux-disponibles", connecte(service.ListCreneauxDisponibles))
	mux.Handle("POST /creneaux/{id}/affectation", connecte(benevole.ProposerAnimation))
	mux.Handle("POST /benevoles/candidature", connecte(benevole.Postuler))
	mux.Handle("POST /paiement/checkout", connecte(paiement.CreerCheckout))
	mux.Handle("POST /paiement/confirmer", connecte(paiement.ConfirmerPaiement))
	mux.Handle("GET /mon-adhesion", connecte(paiement.MonAdhesion))

	// --- Back-office (personnel uniquement) ---
	mux.Handle("GET /utilisateurs", perso(utilisateur.List))
	mux.Handle("POST /utilisateurs/{id}/ban", perso(utilisateur.Bannir))
	mux.Handle("GET /utilisateurs/{id}/adherent-actif", perso(utilisateur.AdherentActif))

	mux.Handle("GET /adhesions", perso(adhesion.List))
	mux.Handle("GET /adhesions/{id}", perso(adhesion.Get))
	mux.Handle("POST /adhesions", perso(adhesion.Create))
	mux.Handle("POST /adhesions/{id}/statut", perso(adhesion.ChangerStatut))
	mux.Handle("POST /admin/rappels", perso(admin.DeclencherRappels))
	mux.Handle("GET /stats", perso(admin.Stats))

	mux.Handle("GET /produits", perso(produit.List))
	mux.Handle("GET /produits/{id}", perso(produit.Get))
	mux.Handle("GET /produits/code-barre/{code}", perso(produit.GetParCodeBarre))
	mux.Handle("POST /produits", perso(produit.Create))
	mux.Handle("GET /stock", perso(produit.Stock))
	mux.Handle("GET /stock/code-barre/{code}", perso(produit.StockParCodeBarre))
	mux.Handle("GET /categories-produits", perso(produit.Categories))

	mux.Handle("GET /collectes", perso(collecte.List))
	mux.Handle("GET /collectes/{id}", perso(collecte.Get))
	mux.Handle("POST /collectes", perso(collecte.Create))
	mux.Handle("GET /collectes/{id}/produits", perso(collecte.Lignes))
	mux.Handle("POST /collectes/{id}/produits", perso(collecte.AjouterProduit))

	mux.Handle("GET /benevoles", perso(benevole.List))
	mux.Handle("POST /benevoles/{id}/statut", perso(benevole.ChangerStatut))
	mux.Handle("GET /competences", perso(benevole.Competences))
	mux.Handle("POST /benevoles/{id}/competences", perso(benevole.AjouterCompetence))

	mux.Handle("GET /planning.xlsx", perso(service.ExportPlanning))

	mux.Handle("GET /tournees", perso(tournee.List))
	mux.Handle("POST /tournees", perso(tournee.Create))
	mux.Handle("GET /tournees/{id}", perso(tournee.Get))
	mux.Handle("GET /tournees/{id}/pdf", perso(tournee.ExportPDF))

	handler := middleware.Logger(mux)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rappel.DemarrerScheduler(ctx, db, cfg)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("arret du serveur...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("arret force: %v", err)
	}
	log.Println("serveur arrete proprement")
}
