# NO MORE WASTE

Plateforme de gestion de l'association NO MORE WASTE (lutte contre le gaspillage).

- **API** : Go (net/http, pgx), dossier `api/`
- **Front** : PHP, dossier `web/`
- **Base** : PostgreSQL, dossier `db/`

## Prerequis

- Go 1.22+
- PostgreSQL 16+
- PHP 8.3+

## Installation

```bash
# 1. Base de donnees
createdb nomorewaste
bash scripts/reset-db.sh          # schema + jeu de donnees de demo

# 2. API
cd api
cp .env.example .env              # ajuster DATABASE_URL si besoin
go mod tidy
go run ./cmd/server               # API sur http://localhost:8080
```

Verifier : `curl http://localhost:8080/health` -> `{"status":"ok"}`

## Comptes de demo (mot de passe : password123)

- back-office : `sylvie.moreau@nomorewaste.org`
- commercant  : `claire.lefevre@boulangerie.fr`
- benevole    : `julien.martin@email.fr`
- adherent    : `nathalie.leroy@email.fr`

Voir `NOTES.md` (decisions d'architecture) et `ARCHITECTURE.md` (structure des fichiers).
