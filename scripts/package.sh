#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$ROOT_DIR"

VERSION="${1:-1.0.0}"
NOM="nomorewaste-${VERSION}"
DIST="$ROOT_DIR/dist"

echo "==> [1/4] Verification du build Go (vet + build)"
( cd api && go vet ./... && go build ./... )

echo "==> [2/4] Build des images Docker"
docker compose build

echo "==> [3/4] Preparation de l'archive $NOM"
rm -rf "${DIST:?}/$NOM"
mkdir -p "$DIST/$NOM"
cp -r api web db docker scripts docker-compose.yml docker-compose.prod.yml deploy.env.example "$DIST/$NOM/"
for f in README.md NOTES.md ARCHITECTURE.md; do
    [ -f "$f" ] && cp "$f" "$DIST/$NOM/"
done

# Retirer les secrets et artefacts de compilation
find "$DIST/$NOM" -type f -name '*.env' ! -name '*.example' -delete
rm -f "$DIST/$NOM/api/api" "$DIST/$NOM/api/server" "$DIST/$NOM/api/cron/cron"

echo "==> [4/4] Compression"
tar -czf "$DIST/${NOM}.tar.gz" -C "$DIST" "$NOM"
rm -rf "$DIST/$NOM"

echo ""
echo "==> Livrable pret : dist/${NOM}.tar.gz"
echo "    Sur le serveur cible :"
echo "      tar xzf ${NOM}.tar.gz && cd ${NOM}"
echo "      cp deploy.env.example .env   # puis remplir SMTP / JWT / DB_PASSWORD"
echo "      docker compose up -d --build"
