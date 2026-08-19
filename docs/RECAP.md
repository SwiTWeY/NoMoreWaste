# Récap NO MORE WASTE — Backend (au 17 août)

## 0. À corriger d'abord (build cassé)
- `cron/main.go` : supprimer la ligne parasite `admin := handlers.AdminHandler{...}`.
- `main.go` : ajouter `admin := handlers.AdminHandler{DB: db, Cfg: cfg}` près des autres handlers.

---

## 1. Le projet
Plateforme de gestion de l'association NO MORE WASTE (anti-gaspillage). **API en Go**, front **PHP** (à venir), base **PostgreSQL**. Un module métier = un patron répété (models → bdd → handlers).

**Fait à ce jour** : socle backend + module **adhésion** complet + système de **rappel de renouvellement** (mailer, concurrence, cron, endpoint admin). Il manque le scheduler à ticker.

---

## 2. L'architecture en couches

```
Requête HTTP
   │
   ▼
main.go (route)  ──►  handler (HTTP)  ──►  bdd (SQL)  ──►  PostgreSQL
                         │
                         └─ décode la requête, appelle bdd, renvoie du JSON (via utils)
```

- **Paquet (package)** : un dossier = un paquet Go. Chaque fichier commence par `package X`.
- **models/** : les structs = la forme des données.
- **bdd/** : les requêtes SQL. **Seule** couche qui parle à PostgreSQL.
- **handlers/** : le HTTP. Décode, appelle `bdd`, répond. **Jamais** de SQL ici.
- **config/**, **utils/**, **middleware/**, **mailer/**, **rappel/** : briques transverses.
- **Règle d'or** : `handler → bdd → base`. Trois couches, jamais plus.

---

## 3. Le socle

| Fichier | Rôle |
|---|---|
| `config/config.go` | struct `Config` + `Load()` : lit les variables d'environnement (.env) |
| `config/db.go` | `Connect(cfg)` → ouvre le pool PostgreSQL (`database/sql` + `lib/pq`), fait un `Ping` |
| `utils/response.go` | `JSON()` / `Error()` : helpers de réponse JSON |
| `middleware/logger.go` | `Logger` : logge `méthode chemin durée` à chaque requête |
| `main.go` | point d'entrée : charge .env, connecte la base, déclare les routes, lance le serveur (port 8086) |

---

## 4. Module adhésion

### `models/adhesion.go`
- `Adhesion` : reflète la table `adhesion` (id, utilisateur_id, dates, montant, statut, rappel_envoye_le, created_at).
- `RappelAdhesion` : read-model (jointure adhésion+utilisateur) pour écrire le mail.

### `bdd/adhesion.go`
- `ListAdhesions` → toutes les adhésions (`Query` + boucle).
- `GetAdhesion(id)` → une adhésion (`QueryRow`), `sql.ErrNoRows` si absente.
- `CreateAdhesion(a)` → insertion (`INSERT ... RETURNING id, created_at`).
- `AdhesionsARappeler(joursAvant)` → jointure : adhésions expirant sous N jours ET `rappel_envoye_le IS NULL`.
- `MarquerRappelEnvoye(id)` → `UPDATE ... SET rappel_envoye_le = NOW()` (le "tampon").

### `bdd/utilisateur.go`
- `EstAdherentActif(id)` → `SELECT EXISTS(...)` : booléen, adhésion payée couvrant aujourd'hui (décision D2).

### `handlers/`
- `adhesion.go` : `List` / `Get` (404 si absent) / `Create` (201).
- `utilisateur.go` : `AdherentActif` → `{"adherent_actif": true/false}`.
- `admin.go` : `DeclencherRappels` → lance le rappel via HTTP.

### Routes (`main.go`)
```
GET  /health
GET  /adhesions
GET  /adhesions/{id}
POST /adhesions
GET  /utilisateurs/{id}/adherent-actif
POST /admin/rappels
```

---

## 5. Le système de rappel de renouvellement

### Le cycle
```
AdhesionsARappeler (rappel_envoye_le IS NULL)  →  envoyer le mail  →  MarquerRappelEnvoye (tampon)
```
Grâce au tampon, chaque adhésion n'est rappelée **qu'une fois** (idempotence : 2e passage = 0).

### Les fichiers
- `mailer/rappel_renouvellement.html` : le template du mail (`{{.Prenom}}`…), fichier séparé.
- `mailer/mailer.go` : `EnvoyerRappel(cfg, r)` — charge le template via `//go:embed`, le remplit, construit le message, envoie par SMTP.
- `rappel/rappel.go` :
  - `envoyerUn(db, cfg, r)` : traite UN rappel (envoi ou mode "à blanc" si SMTP non configuré, puis tampon).
  - `EnvoyerRappels(ctx, db, cfg)` : orchestration **concurrente** (goroutines + WaitGroup + canal-sémaphore + mutex + context). Renvoie le nombre envoyés.

### Les 3 déclencheurs (1 seule implémentation : `EnvoyerRappels`)
1. **`cron/main.go`** : binaire one-shot (crontab). ✅
2. **`POST /admin/rappels`** : à la demande via HTTP. ✅
3. **Scheduler goroutine** (ticker dans `main.go`) : à faire.

---

## 6. Glossaire — définitions

### Go — bases
- **struct** : un type qui regroupe des champs nommés (≈ un objet). `type Adhesion struct { ... }`.
- **pointeur (`*T`)** : l'adresse d'une valeur, pas une copie. `*sql.DB` = pointeur vers le pool. `&x` = "adresse de x".
- **`*time.Time` (pointeur nullable)** : permet la valeur `nil` = NULL en base (ex. `rappel_envoye_le`).
- **exporté / non exporté** : nom qui commence par une **majuscule** = visible des autres paquets ; minuscule = privé au paquet.
- **`:=` vs `=`** : `:=` déclare + assigne (nouvelle variable) ; `=` assigne à une variable existante.
- **erreur** : renvoyée en dernier ; on teste toujours `if err != nil`.
- **`defer`** : exécute une instruction à la sortie de la fonction (nettoyage garanti).
- **`errors.Is(err, cible)`** : teste si une erreur est d'un type précis (ex. `sql.ErrNoRows`).
- **tags `json:"..."`** : indiquent le nom du champ en JSON.

### Go — HTTP
- **handler** : fonction qui traite une requête, signature `(w http.ResponseWriter, r *http.Request)`.
- **`w` (ResponseWriter)** : là où on écrit la réponse.
- **`r` (*Request)** : la requête entrante (URL, en-têtes, corps).
- **receiver** `func (h T) Nom(...)` : attache une méthode à un type `T`.
- **`r.PathValue("id")`** : lit un paramètre d'URL `{id}` (Go 1.22).
- **`json.NewDecoder(r.Body).Decode(&x)`** : lit le JSON de la requête dans une struct.
- **codes HTTP** : 200 OK, 201 Created, 400 Bad Request, 404 Not Found, 500 Internal Error.

### SQL / base
- **`database/sql`** : la lib standard Go pour les bases SQL. **`lib/pq`** : le driver PostgreSQL.
- **pool de connexions (`*sql.DB`)** : stock de connexions réutilisables, partagé.
- **`db.Query`** : lit **plusieurs** lignes. **`db.QueryRow`** : lit **une** ligne. **`db.Exec`** : modifie sans rien lire (INSERT/UPDATE/DELETE).
- **`rows.Scan(&...)`** : recopie les colonnes de la ligne dans les champs (ordre = ordre du SELECT).
- **`$1, $2...`** : paramètres — la valeur est passée séparément (protège de l'injection SQL).
- **`RETURNING`** : renvoie des colonnes après un INSERT (ex. l'`id` généré).
- **`JOIN ... ON`** : combine deux tables reliées par une clé.
- **`SELECT EXISTS(...)`** : renvoie un booléen (au moins une ligne ?).
- **`CURRENT_DATE` / `NOW()`** : date du jour / instant présent, côté PostgreSQL.

### Go — templates / embed
- **`//go:embed fichier`** : directive qui inclut un fichier dans le binaire à la compilation.
- **`html/template`** : moteur de template qui remplace `{{.Champ}}` par des valeurs (et échappe le HTML).
- **`bytes.Buffer`** : tampon de texte en mémoire.
- **`time.Format("02/01/2006")`** : format de date ; Go utilise une **date de référence** (02=jour, 01=mois, 2006=année).

### Go — concurrence (le morceau "oral")
- **goroutine** (`go func(){...}()`) : une fonction lancée **en parallèle**, sans attendre.
- **`sync.WaitGroup`** : compteur pour **attendre** que toutes les goroutines finissent. `Add(1)` / `Done()` / `Wait()`.
- **canal (channel)** : tuyau typé pour communiquer entre goroutines. `make(chan T, n)` = bufferisé (capacité n).
- **sémaphore par canal** : un canal bufferisé de taille N utilisé pour **limiter** le nombre de goroutines actives (au plus N en même temps).
- **`struct{}`** : struct vide (0 octet) ; sert de "jeton" quand seule compte la quantité.
- **`sync.Mutex`** : verrou. `Lock()` / `Unlock()` protègent une donnée partagée (ici le compteur) d'une **data race**.
- **`context.Context`** : porte l'**annulation**. `ctx.Done()` = canal qui se ferme à l'annulation ; permet un arrêt propre.
- **`select { case <-ch: ... default: }`** : attend plusieurs canaux ; `default` le rend non bloquant.

---

## 7. Lancer / tester

```bash
# (Re)charger la base + le jeu de démo
cd ~/Rattrapage/nomorewaste/NoMoreWaste
PGPASSWORD=zak bash scripts/reset-db.sh

# Lancer l'API
cd api
go run .

# Tester (autre terminal)
curl -s http://localhost:8086/adhesions | head -c 300; echo
curl -s http://localhost:8086/adhesions/1; echo
curl -s http://localhost:8086/utilisateurs/15/adherent-actif; echo
curl -s -X POST http://localhost:8086/admin/rappels; echo

# Le cron (rappels)
go run ./cron        # 1er passage : 3 ; 2e passage : 0 (idempotent)
```

---

## 8. Ce qui reste (backend)
- Scheduler à ticker (déclencheur n°3) + arrêt propre par `context`.
- Auth JWT/RBAC (register/login/middleware, `est_personnel`).
- 5 modules : collectes, stocks (+ code-barres), tournées (+ PDF), bénévoles (+ compétences), services (+ Excel).
- i18n fr/en.
- Puis : front PHP, déploiement Nginx + packaging.
