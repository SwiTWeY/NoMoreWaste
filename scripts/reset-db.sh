#!/usr/bin/env bash
set -euo pipefail

DB_NAME="${DB_NAME:-nomorewaste}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5433}"
export PGPASSWORD="${DB_PASS:-${PGPASSWORD:-}}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

PSQL="psql -h $DB_HOST -p $DB_PORT -U $DB_USER -v ON_ERROR_STOP=1"

echo "==> Cluster    : $DB_HOST:$DB_PORT (user $DB_USER)"

if ! $PSQL -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" | grep -q 1; then
    echo "==> Creation de la base $DB_NAME"
    $PSQL -d postgres -c "CREATE DATABASE $DB_NAME"
fi

echo "==> Schema     : $ROOT_DIR/db/schema.sql"
$PSQL -d "$DB_NAME" -f "$ROOT_DIR/db/schema.sql" -q

echo "==> Seed       : $ROOT_DIR/db/seed.sql"
$PSQL -d "$DB_NAME" -f "$ROOT_DIR/db/seed.sql"

echo ""
echo "==> Base reinitialisee."
echo "    Comptes de demo, mot de passe : password123"
echo "      back-office : sylvie.moreau@nomorewaste.org"
echo "      commercant  : claire.lefevre@boulangerie.fr"
echo "      benevole    : julien.martin@email.fr"
echo "      adherent    : nathalie.leroy@email.fr"
