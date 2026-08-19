# Architecture des fichiers — NO MORE WASTE

## Vue d'ensemble

```
NoMoreWaste/
├── NOTES.md                  Journal des decisions (D1 a D14)
├── README.md                 Installation et lancement
├── docker-compose.yml        PostgreSQL en local (optionnel)
│
├── db/
│   ├── schema.sql            Les 20 tables + la vue v_stock
│   └── seed.sql              Jeu de donnees de demo
│
├── scripts/
│   ├── reset-db.sh           Recree la base + injecte le seed
│   ├── init-structure.sh     Creation de l'arborescence (a lancer une fois)
│   └── deploy.sh             Script de packaging exige par le sujet
│
├── docs/
│   ├── mcd.md                Le schema de donnees
│   └── api.md                Liste des endpoints
│
├── api/                      ===== BACKEND GO =====
└── web/                      ===== FRONTEND PHP =====
```

---

## API Go — `api/`

```
api/
├── go.mod / go.sum
├── .env / .env.example
├── main.go              Point d'entree : config, connexion DB, routes, serveur HTTP
│                        (+ goroutine des rappels de renouvellement)
│
├── config/
│   ├── config.go        Lecture des variables d'environnement
│   └── db.go            Connexion PostgreSQL (database/sql + lib/pq)
│
├── models/              Une struct par table
│   ├── utilisateur.go
│   ├── commercant.go
│   ├── adhesion.go
│   ├── benevole.go
│   ├── produit.go
│   ├── collecte.go
│   ├── tournee.go
│   ├── beneficiaire.go
│   └── service.go
│
├── bdd/                 Les requetes SQL, une par module
│   ├── utilisateur.go
│   ├── commercant.go
│   └── ...
│
├── handlers/            Les endpoints HTTP, un par module
│   ├── auth.go
│   ├── commercant.go
│   └── ...
│
├── middleware/
│   ├── auth.go          Verification du token JWT + roles
│   └── logger.go
│
├── mailer/mailer.go     Envoi des emails (rappels), templates //go:embed
│
├── export/
│   ├── pdf.go           Recapitulatif de tournee (fpdf)
│   └── excel.go         Planning des benevoles (excelize)
│
└── utils/response.go    Helpers JSON : succes, erreur
```

### Le chemin d'une requete

```
handler  ->  bdd  ->  PostgreSQL
   |          |
   |          +-- construit et execute le SQL
   |
   +-- decode la requete, appelle bdd, renvoie du JSON
```

**Trois couches, jamais plus.** Le handler ne fait pas de SQL. Le paquet bdd ne
connaît pas le HTTP. C'est ce qui rend le code navigable : quand on demande
d'ajouter un champ, on sait exactement quels trois fichiers ouvrir.

---

## Front PHP — `web/`

```
web/
├── public/                   SEUL dossier expose par le serveur web
│   ├── index.php             Routeur : toutes les URLs passent par ici
│   ├── .htaccess             Reecriture d'URL + pages d'erreur
│   └── assets/
│       ├── css/style.css
│       ├── js/app.js
│       └── img/
│
├── config/config.php         URL de l'API Go, langue par defaut
│
└── src/
    ├── lib/
    │   ├── ApiClient.php     Appelle l'API Go en HTTP (cURL)
    │   ├── Session.php       Gestion de la session et du token
    │   └── I18n.php          Chargement des traductions
    │
    ├── lang/
    │   ├── fr.php            Libelles de l'interface en francais
    │   └── en.php            Libelles de l'interface en anglais
    │
    └── views/
        ├── layouts/
        │   ├── back.php      Gabarit du back-office
        │   └── front.php     Gabarit du front-office
        ├── partials/
        │   ├── nav.php
        │   └── flash.php
        ├── auth/
        ├── commercants/      index.php, form.php, show.php
        ├── collectes/
        ├── stocks/
        ├── tournees/
        ├── benevoles/
        └── services/
```

### Regle absolue

**PHP n'a pas les identifiants de PostgreSQL.** Il ne fait que des appels HTTP
vers l'API Go, decode le JSON, et affiche.

Si tu es tente d'ecrire une requete SQL depuis PHP « juste pour ce cas-la »,
c'est le moment ou l'architecture casse. Le sujet impose l'API en Go.

### Pourquoi `public/` est le seul dossier expose

Le serveur web pointe sur `web/public/`. Tout le reste (`src/`, `config/`) est
au-dessus de la racine web, donc inaccessible depuis un navigateur. Les
identifiants et la logique ne sont jamais servis en clair.

---

## Le patron a repliquer sur chaque module

C'est le point central de la strategie : **un seul chemin, duplique six fois**.

Pour le module commerçants :

| Couche | Fichier |
|---|---|
| Table | `db/schema.sql` |
| Struct Go | `api/models/commercant.go` |
| SQL | `api/bdd/commercant.go` |
| Endpoints | `api/handlers/commercant.go` |
| Route front | `web/public/index.php` |
| Affichage | `web/src/views/commercants/index.php` |

Les cinq autres modules suivent **exactement** la même liste.

### Ajouter un champ en live coding

1. `ALTER TABLE` dans `schema.sql`
2. Ajouter le champ dans la struct (`models/`)
3. L'ajouter au SELECT et au INSERT (`bdd/`)
4. L'afficher dans la vue (`views/`)

Quatre fichiers, toujours les mêmes, toujours dans cet ordre.
**A repeter au chronometre en semaine 3.**

---

## Nommage : la regle qui evite les erreurs sous pression

Le meme concept garde le meme nom a travers toutes les couches :

| SQL | Go (struct) | JSON | PHP |
|---|---|---|---|
| `raison_sociale` | `RaisonSociale` | `raison_sociale` | `$c['raison_sociale']` |
| `date_limite` | `DateLimite` | `date_limite` | `$p['date_limite']` |

En Go, les tags le forcent :

```go
type Commercant struct {
    ID            int    `json:"id"            db:"id"`
    RaisonSociale string `json:"raison_sociale" db:"raison_sociale"`
}
```

**Aucune traduction mentale.** C'est ce qui fait la difference quand on code
devant un examinateur.
