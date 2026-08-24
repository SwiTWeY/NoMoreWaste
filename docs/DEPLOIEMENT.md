# Déploiement — NO MORE WASTE

Architecture conteneurisée :

```
Navigateur ─▶ Nginx (port 80)
                ├─ /assets/*        → fichiers statiques (CSS, polices, images)
                ├─ réécriture d'URL → index.php (front controller)
                └─ *.php            → php-fpm (front PHP)
                                        └─ appels API ─▶ api (Go, port 8086)
                                                            └─▶ PostgreSQL (db)
```

4 conteneurs : **nginx**, **php** (php-fpm), **api** (Go), **db** (PostgreSQL 16).

## 1. Prérequis
- Docker + Docker Compose (`docker --version`, `docker compose version`).

## 2. Configuration
```bash
cp deploy.env.example .env
# éditer .env : DB_PASSWORD, JWT_SECRET, et SMTP_USER / SMTP_PASS (rappels email)
```
Le `.env` n'est pas versionné (secrets).

## 3. Lancement
```bash
docker compose up -d --build
```
- Site (front + back-office) : **http://localhost:8088**
- API (debug) : **http://localhost:8087/health**

La base est **initialisée automatiquement** au premier démarrage (`db/schema.sql` puis `db/seed.sql`).

Comptes de démo (mot de passe `password123`) :
- back-office : `sylvie.moreau@nomorewaste.org`
- adhérent : `nathalie.leroy@email.fr`

## 4. Commandes utiles
```bash
docker compose logs -f            # suivre les logs
docker compose ps                 # état des conteneurs
docker compose exec db psql -U postgres nomorewaste   # accès base
docker compose down               # arrêter
docker compose down -v            # arrêter + effacer la base (repart de zéro)
```

## 5. Fonctionnalités serveur web (Nginx) — `docker/nginx/default.conf`
- **Réécriture d'URL** : `try_files $uri /index.php?$query_string` (URLs propres → front controller).
- **Codes d'erreur personnalisés** : `error_page 404 → /erreur-404.html`, `500/502/503/504 → /erreur-50x.html`.
- **Fichiers statiques** servis directement (cache 7 jours), fichiers cachés bloqués.

## 6. Packaging (livrable)
```bash
./scripts/package.sh 1.0.0
# produit dist/nomorewaste-1.0.0.tar.gz (sans secrets)
```

## 7. Mise en ligne (VPS) — Phase 2
1. Un serveur (VPS) avec Docker installé.
2. Copier le livrable (ou `git clone`), créer le `.env`.
3. Enregistrement DNS : `nomorewaste.upcycle-connect.fr` (type A) → IP du serveur.
4. `docker compose up -d --build`.
5. HTTPS : ajouter un reverse-proxy TLS (Let's Encrypt via Caddy ou Certbot) devant Nginx, ou terminer le TLS dans Nginx.
