# NOTES.md — NO MORE WASTE

Journal des décisions d'architecture. À tenir à jour à chaque session.

- **Rendu** : samedi 29 août 2026, 23h59
- **Stack imposée** : API Go + front PHP ou JavaScript
- **Soutenance** : 5 min de pitch commercial + live coding (ajout de champs, filtres, tris)
- **Principe directeur** : code simple, répétitif, navigable. Un seul patron dupliqué sur tous les modules. Pas de généricité maligne — ce qui compte est de pouvoir modifier n'importe quoi en 3 minutes sous pression.

---

## Règle de travail

- Tout ce qui touche à la **structure des données** se décide avant de coder.
- Tout ce qui touche à l'**apparence ou au confort** se décide en codant.
- Nommage strictement cohérent de la base au front : `birth_date` en SQL → `BirthDate` en Go → `birth_date` en JSON. Aucune traduction mentale en live.

---

## Décisions actées

### D1 — Utilisateur central avec tables satellites

Une seule table `utilisateur` porte **l'identité et l'authentification**. Les rôles sont
définis par l'**existence d'une ligne satellite**, pas par une colonne `role`.

- ligne dans `adhesion` → adhérent
- ligne dans `profil_benevole` → bénévole
- ligne dans `commercant` → commerçant
- les trois sont cumulables sans conflit

Le seul rôle qui ne se déduit d'aucune donnée métier est le personnel de l'association
(accès back-office) : booléen `est_personnel` sur `utilisateur`.

**Pourquoi** : évite une table à 25 colonnes majoritairement NULL, évite une table de rôles
à maintenir, et permet le cumul (un bénévole peut être adhérent).

### D2 — Adhésion en 1-N (historique conservé)

Une ligne d'adhésion par période, pas d'écrasement.

**Conséquences** :
- Le rappel automatique de renouvellement = requête sur l'adhésion la plus récente dont
  `date_fin` tombe dans X jours. Un cron, une requête.
- « Est-il adhérent ? » n'est plus trivial : il faut la ligne la plus récente ET non expirée.
  → Prévoir une fonction unique `EstAdherentActif()`. Ne jamais recopier la logique de date
  à plusieurs endroits.

### D3 — Commerçant en satellite, pas en entité autonome

`utilisateur` porte le contact (le gérant, celui qui reçoit les mails).
`commercant` porte l'entreprise : SIRET, raison sociale, adresse, horaires.

**Pourquoi** : un commerçant se connecte (demande de collecte, suivi d'adhésion). Une entité
autonome imposerait un second système d'authentification et rendrait `adhesion` polymorphe.

**Important** : `commercant` ne se relie PAS à `adhesion`. L'adhésion pointe déjà sur
`utilisateur`. Deux chemins vers la même chose = doublon.

**Limite connue** : un changement de gérant impose de transférer le `commercant` vers un autre
`utilisateur`. Hors périmètre.

### D4 — Compétences en table N-N

Le sujet cite chauffeurs, cuisiniers, plombiers, et un bénévole peut cumuler.
→ table `competence` + table de liaison `benevole_competence`.

**Pourquoi** : une colonne texte avec des valeurs séparées par des virgules rendrait impossible
la requête « qui sait conduire ? », qui est exactement ce dont le module d'affectation a besoin.

### D5 — Collecte et tournée : deux tables distinctes

**Pourquoi** (le vrai argument) :
- Les entités liées ne sont pas les mêmes de chaque côté. Une collecte pointe vers une source
  (commerçant). Une tournée pointe vers des bénéficiaires. Fusionner imposerait une clé
  polymorphe.
- L'impact sur le stock est inverse : une collecte incrémente, une tournée décrémente.
- Le PDF récapitulatif du sujet est attaché aux livraisons, donc aux tournées uniquement.

Ce qui est factorisé malgré tout : `vehicule` et le chauffeur (clé vers `utilisateur`).

### D6 — Tournée à arrêts multiples, produits au niveau de la tournée

Une tournée comprend plusieurs `arret` (un par bénéficiaire, avec ordre de passage).
Mais les produits sont rattachés à la **tournée**, pas à l'arrêt.

**Pourquoi** : traçabilité au niveau de la cargaison, pas du destinataire. Évite un troisième
niveau d'imbrication (`tournee` → `arret` → `produits`) qui aurait rendu les formulaires
et le PDF nettement plus risqués en live coding.

**Bénéfice** : `collecte_produit` et `tournee_produit` deviennent structurellement identiques.
Un seul patron, réutilisé.

**Coût accepté** : on ne sait pas quel produit est allé chez quel bénéficiaire.

### D7 — Bénéficiaire en table autonome avec champ `type`

Pas de satellite autour d'`utilisateur`.

**Pourquoi** : le critère est « est-ce que cette personne se connecte au site ? ». Un particulier
en détresse ou une association bénéficiaire ne se connecte jamais. Créer des `utilisateur`
sans mot de passe polluerait la table et casserait la logique de rôles.

Un champ `type` (`association` / `particulier`) suffit : les relations sont identiques des deux
côtés, seule l'étiquette change. C'est le cas où un discriminant est légitime — contrairement à
collecte/tournée.

### D8 — Stock calculé, pas stocké

Pas de colonne `quantite_stock`. Le stock = somme des collectes − somme des tournées.

**Alternative si trop lent ou trop pénible à afficher** : ajouter une colonne mise à jour à
chaque mouvement. Plus simple à lire, mais risque de désynchronisation. À trancher en codant.

### D9 — i18n dès le schéma

Deux niveaux à ne pas confondre :
1. **Interface** (boutons, labels) → fichiers de traduction `fr.json`, `en.json`. Ajoutable tard.
2. **Données métier** (nom d'un service, d'une catégorie) → **table de traduction en base**.
   Structurel. À faire maintenant, sinon réécriture du schéma en semaine 3.

Cible : FR + EN. Deux langues bien faites valent mieux que quatre bâclées.

---

### D10 — Services : catalogue et créneaux séparés

`service` = le catalogue (quasi statique, alimenté par le back-office).
`creneau` = les occurrences datées d'un service (date, horaires, lieu, capacité).

**Pourquoi** : avec une table unique, la description du service serait recopiée sur chaque
séance. Trente séances = trente copies à corriger.

**Bénéfice pour l'export Excel** : le planning quotidien se requête sur `creneau` filtré par
date, avec une jointure sur `service` pour le libellé. Une requête simple.

### D11 — Inscription et affectation : deux tables

`inscription` = les adhérents qui s'inscrivent pour recevoir le service (front-office, annulable).
`affectation` = les bénévoles affectés pour rendre le service (back-office, validé sur compétence).

**Pourquoi** : deux flux, deux logiques métier, deux jeux de permissions. Une table unique avec
un champ `role` obligerait à filtrer sur `role = 'animateur'` dans chaque requête — dans
l'export Excel, dans le planning, dans les notifications. Source d'erreur en live coding.

Plusieurs animateurs possibles par créneau (souplesse voulue : un déménagement demande
trois personnes).

### D12 — Périmètre de l'i18n en base

Traduction en base **uniquement** pour `service` et `competence`.
Tout le reste (libellés d'interface, boutons, messages) passe par les fichiers `fr.json` / `en.json`.

**Pourquoi** : multiplier les tables `_traduction` n'apporte aucun bénéfice visible en démo.
`service` ne porte aucun texte — seulement un code (`COURS_CUISINE`) et un booléen `actif`.
Les libellés vivent dans `service_traduction`.

### D13 — Service lié à une compétence

`service.competence_id` : un cours de cuisine requiert la compétence cuisinier.

**Bénéfice** : à l'affectation d'un bénévole sur un créneau, ne proposer que ceux qui ont la
bonne compétence. Deux jointures, et c'est un bon moment de démo.

### D14 — Inscription liée à `utilisateur`, pas à `adhesion`

Le sujet réserve les services aux adhérents, mais c'est un **contrôle à l'inscription**
(`EstAdherentActif()`), pas une clé étrangère. Lier à une adhésion précise poserait problème
au renouvellement.

---

## État : les 6 modules sont modélisés

21 tables. Chacune est plate et sans piège ; neuf sont de simples tables de liaison ou de
référence.

| Module | Tables |
|---|---|
| Adhésions commerçants | `utilisateur`, `commercant`, `adhesion` |
| Collectes | `collecte`, `collecte_produit`, `vehicule` |
| Stocks | `produit` (+ calcul sur les mouvements) |
| Tournées | `tournee`, `arret`, `beneficiaire`, `tournee_produit` |
| Bénévoles | `profil_benevole`, `competence`, `benevole_competence` |
| Services | `service`, `service_traduction`, `langue`, `creneau`, `inscription`, `affectation` |

---

### D15 — Front en PHP (templates serveur), pas React — DECISION FERMEE

**Critère de décision** : la techno la mieux maîtrisée aujourd'hui, sans documentation.
C'est le seul critère qui compte pour un oral noté sur le live coding.

**Pourquoi PHP** :
- Ajouter un champ = éditer 4 fichiers (SQL, struct Go, requête, template). Recharger la page.
  En React : state + composant + fetch + rebuild. Si le build casse devant le jury, l'oral est perdu.
- L'exigence de configuration serveur du sujet (réécriture d'URL, pages d'erreur) n'a de sens
  qu'en server-side. Avec une SPA, c'est artificiel.
- L'i18n est plus simple : un tableau de traduction chargé selon la langue, pas de librairie.
- Aucun besoin temps réel dans ce projet — l'asynchrone n'apporte rien ici.

**Ce qu'on accepte** : interface moins fluide, rechargement à chaque action. Bon échange sur
un projet noté sur la maîtrise, pas sur l'UX.

**Règle absolue** : PHP n'a PAS les identifiants PostgreSQL. Il appelle l'API Go en HTTP,
décode le JSON, affiche. Le jour où une requête SQL est écrite depuis PHP « juste pour ce
cas-là », l'architecture est cassée et ça se voit.

**NE PLUS ROUVRIR CE DÉBAT.** Trois revirements le 9 août ont coûté du temps sans produire
une ligne de code.

### D16 — Layout Go plat, façon UpcycleConnect

`config/`, `models/`, `bdd/`, `handlers/`, `middleware/`, `mailer/`, un seul `main.go`.

**Pourquoi** : structure déjà maîtrisée sur un projet précédent. Un code navigué d'instinct
vaut mieux qu'une architecture théoriquement plus propre où on hésite devant l'examinateur.

**Contrainte obligatoire** : **un fichier par module** dans chaque dossier
(`handlers/commercant.go`, `bdd/commercant.go`…), jamais de fichier fourre-tout.
Avec 20 tables, il faut pouvoir ouvrir directement le bon fichier.
Objectif : chaque fichier sous ~200 lignes.

### D17 — Driver : `database/sql` + `lib/pq`

Pas `pgx`. Même raison que D16 : syntaxe déjà connue. Aucun problème de performance sur ce
projet ne justifie d'apprendre un nouveau driver.

### D18 — Cron en binaire séparé, pas en goroutine

Les rappels de renouvellement tournent dans un binaire lancé par crontab.

**Pourquoi** : une goroutine qui plante, plante en silence. Un binaire séparé se lance à la
demande devant le jury — on montre le mail partir. Plus simple à débugger aussi.

---

## Points ouverts

- ~~Collecte chez un particulier~~ → **RÉSOLU** : ajout d'un `utilisateur_id` nullable sur
  `collecte`, avec une contrainte `CHECK` garantissant une seule source renseignée.
- **i18n** : fr + en pour le rendu. `it` et `pt` en semaine 3 seulement s'il reste du temps.
  La table `langue` est déjà prévue pour.
- **Disponibilités des bénévoles** : actuellement un champ `text` sur `profil_benevole`. Si les
  plannings doivent croiser disponibilités et créneaux, il faudra une vraie table. À trancher au
  module services.
- **Contrôle de capacité véhicule** : `vehicule.capacite_kg` existe. Vérifier que le poids de la
  cargaison ne dépasse pas la capacité serait un bon moment de démo. Confort, pas structure —
  à ajouter en semaine 3 s'il reste du temps.

---

## Non fait / à corriger

- Le repo s'appelle `NoMoreWast` — il manque le `e` de `Waste`. À renommer avant que ça se
  propage dans les URLs, les imports Go et le rapport.
