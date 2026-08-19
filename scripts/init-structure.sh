#!/usr/bin/env bash
# =============================================================
# NO MORE WASTE — creation de l'arborescence du projet
#
# A lancer UNE FOIS a la racine du repo NoMoreWaste :
#   bash init-structure.sh
#
# Cree tous les dossiers et les fichiers .gitkeep necessaires.
# =============================================================

set -euo pipefail

echo "==> Creation de l'arborescence..."

# ---------- Racine ----------
mkdir -p db scripts docs

# ---------- API Go ----------
mkdir -p api/cmd/server              # point d'entree du serveur HTTP
mkdir -p api/cmd/cron                # rappels d'adhesion (tache planifiee)
mkdir -p api/internal/config         # lecture des variables d'environnement
mkdir -p api/internal/database       # connexion PostgreSQL
mkdir -p api/internal/models         # structs Go = tables SQL
mkdir -p api/internal/repository     # requetes SQL
mkdir -p api/internal/handlers       # endpoints HTTP
mkdir -p api/internal/middleware     # auth, logs, CORS
mkdir -p api/internal/export         # generation PDF et Excel
mkdir -p api/internal/mailer         # envoi des emails
mkdir -p api/pkg/response            # helpers JSON (succes / erreur)

# ---------- Front PHP ----------
mkdir -p web/public/assets/css
mkdir -p web/public/assets/js
mkdir -p web/public/assets/img
mkdir -p web/src/lib                 # client API, session, helpers
mkdir -p web/src/lang                # fichiers de traduction i18n
mkdir -p web/src/views/layouts
mkdir -p web/src/views/partials
mkdir -p web/src/views/auth
mkdir -p web/src/views/commercants
mkdir -p web/src/views/collectes
mkdir -p web/src/views/stocks
mkdir -p web/src/views/tournees
mkdir -p web/src/views/benevoles
mkdir -p web/src/views/services

# ---------- .gitkeep pour que Git suive les dossiers vides ----------
find api web db scripts docs -type d -empty -exec touch {}/.gitkeep \;

echo "==> Arborescence creee."
echo ""
echo "Prochaines etapes :"
echo "  1. Deplacer schema.sql et seed.sql dans db/"
echo "  2. Deplacer reset-db.sh dans scripts/ puis chmod +x"
echo "  3. cd api && go mod init github.com/<ton-user>/nomorewaste/api"
