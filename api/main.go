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
	mux := http.NewServeMux()
	protege := middleware.Auth(cfg.JWTSecret)
	produit := handlers.ProduitHandler{DB: db}
	collecte := handlers.CollecteHandler{DB: db}
	benevole := handlers.BenevoleHandler{DB: db}
	service := handlers.ServiceHandler{DB: db}
	tournee := handlers.TourneeHandler{DB: db}
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			utils.Error(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /adhesions", adhesion.List)
	mux.HandleFunc("GET /adhesions/{id}", adhesion.Get)
	mux.HandleFunc("POST /adhesions", adhesion.Create)
	mux.Handle("POST /admin/rappels", protege(middleware.Personnel(http.HandlerFunc(admin.DeclencherRappels))))
	mux.HandleFunc("GET /utilisateurs/{id}/adherent-actif", utilisateur.AdherentActif)
	mux.HandleFunc("POST /register", authH.Register)
	mux.HandleFunc("POST /login", authH.Login)
	mux.HandleFunc("GET /produits", produit.List)
	mux.HandleFunc("GET /produits/{id}", produit.Get)
	mux.HandleFunc("GET /produits/code-barre/{code}", produit.GetParCodeBarre)
	mux.HandleFunc("POST /produits", produit.Create)
	mux.HandleFunc("GET /stock", produit.Stock)
	mux.HandleFunc("GET /stock/code-barre/{code}", produit.StockParCodeBarre)
	mux.HandleFunc("GET /collectes", collecte.List)
	mux.HandleFunc("GET /collectes/{id}", collecte.Get)
	mux.HandleFunc("POST /collectes", collecte.Create)
	mux.HandleFunc("GET /collectes/{id}/produits", collecte.Lignes)
	mux.HandleFunc("POST /collectes/{id}/produits", collecte.AjouterProduit)
	mux.HandleFunc("GET /benevoles", benevole.List)
	mux.HandleFunc("POST /benevoles/{id}/statut", benevole.ChangerStatut)
	mux.HandleFunc("GET /competences", benevole.Competences)
	mux.HandleFunc("POST /benevoles/{id}/competences", benevole.AjouterCompetence)
	mux.HandleFunc("GET /services", service.ListServices)
	mux.HandleFunc("GET /creneaux", service.ListCreneaux)
	mux.HandleFunc("GET /planning.xlsx", service.ExportPlanning)
	mux.HandleFunc("GET /tournees", tournee.List)
	mux.HandleFunc("GET /tournees/{id}", tournee.Get)
	mux.HandleFunc("GET /tournees/{id}/pdf", tournee.ExportPDF)
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
