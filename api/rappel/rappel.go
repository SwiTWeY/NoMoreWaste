package rappel

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/mailer"
	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

const limiteConcurrence = 3

func envoyerUn(db *sql.DB, cfg config.Config, r models.RappelAdhesion) error {
	if cfg.SMTPUser == "" {
		log.Printf("[rappel a blanc] adhesion %d -> %s (SMTP non configure)", r.AdhesionID, r.Email)
	} else {
		if err := mailer.EnvoyerRappel(cfg, r); err != nil {
			return err
		}
	}
	return bdd.MarquerRappelEnvoye(db, r.AdhesionID)
}

func EnvoyerRappels(ctx context.Context, db *sql.DB, cfg config.Config) (int, error) {
	joursAvant, err := strconv.Atoi(cfg.RappelJours)
	if err != nil || joursAvant <= 0 {
		joursAvant = 30
	}

	rappels, err := bdd.AdhesionsARappeler(db, joursAvant)
	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, limiteConcurrence)
	var mu sync.Mutex
	envoyes := 0

	for _, r := range rappels {
		select {
		case <-ctx.Done():
			wg.Wait()
			return envoyes, ctx.Err()
		default:
		}

		wg.Add(1)
		go func(r models.RappelAdhesion) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()
			log.Printf("[rappel] envoi -> %s", r.Email)

			if err := envoyerUn(db, cfg, r); err != nil {
				log.Printf("rappel adhesion %d: %v", r.AdhesionID, err)
				return
			}

			mu.Lock()
			envoyes++
			mu.Unlock()
		}(r)
	}

	wg.Wait()
	return envoyes, nil
}
func DemarrerScheduler(ctx context.Context, db *sql.DB, cfg config.Config) {
	intervalle, err := time.ParseDuration(cfg.SchedulerIntervalle)
	if err != nil || intervalle <= 0 {
		intervalle = 24 * time.Hour
	}

	go func() {
		ticker := time.NewTicker(intervalle)
		defer ticker.Stop()

		log.Printf("scheduler demarre (intervalle %s)", intervalle)
		for {
			select {
			case <-ctx.Done():
				log.Println("scheduler: arret")
				return
			case <-ticker.C:
				n, err := EnvoyerRappels(ctx, db, cfg)
				if err != nil {
					log.Printf("scheduler: %v", err)
					continue
				}
				log.Printf("scheduler: %d rappels envoyes", n)
			}
		}
	}()
}
