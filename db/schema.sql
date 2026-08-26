-- =============================================================
-- NO MORE WASTE — schema PostgreSQL
-- Genere a partir du MCD valide (voir NOTES.md, decisions D1 a D14)
-- =============================================================

DROP TABLE IF EXISTS affectation CASCADE;
DROP TABLE IF EXISTS inscription CASCADE;
DROP TABLE IF EXISTS creneau CASCADE;
DROP TABLE IF EXISTS service_traduction CASCADE;
DROP TABLE IF EXISTS service CASCADE;
DROP TABLE IF EXISTS langue CASCADE;
DROP TABLE IF EXISTS tournee_produit CASCADE;
DROP TABLE IF EXISTS collecte_produit CASCADE;
DROP TABLE IF EXISTS arret CASCADE;
DROP TABLE IF EXISTS tournee CASCADE;
DROP TABLE IF EXISTS collecte CASCADE;
DROP TABLE IF EXISTS produit CASCADE;
DROP TABLE IF EXISTS beneficiaire CASCADE;
DROP TABLE IF EXISTS vehicule CASCADE;
DROP TABLE IF EXISTS benevole_competence CASCADE;
DROP TABLE IF EXISTS competence CASCADE;
DROP TABLE IF EXISTS profil_benevole CASCADE;
DROP TABLE IF EXISTS commercant CASCADE;
DROP TABLE IF EXISTS adhesion CASCADE;
DROP TABLE IF EXISTS utilisateur CASCADE;


-- =============================================================
-- BLOC 1 : IDENTITE
-- =============================================================

CREATE TABLE utilisateur (
    id            SERIAL PRIMARY KEY,
    nom           VARCHAR(100) NOT NULL,
    prenom        VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    mot_de_passe  VARCHAR(255) NOT NULL,
    telephone     VARCHAR(20),
    est_personnel BOOLEAN      NOT NULL DEFAULT FALSE,
    actif         BOOLEAN      NOT NULL DEFAULT TRUE,
    langue_pref   VARCHAR(5)   NOT NULL DEFAULT 'fr',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_utilisateur_email ON utilisateur (email);


CREATE TABLE adhesion (
    id              SERIAL PRIMARY KEY,
    utilisateur_id  INTEGER      NOT NULL REFERENCES utilisateur (id) ON DELETE CASCADE,
    date_debut      DATE         NOT NULL,
    date_fin        DATE         NOT NULL,
    montant         NUMERIC(8,2) NOT NULL DEFAULT 0,
    statut_paiement VARCHAR(20)  NOT NULL DEFAULT 'en_attente',
    rappel_envoye_le TIMESTAMPTZ,
    stripe_session_id VARCHAR(255) UNIQUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_adhesion_dates   CHECK (date_fin > date_debut),
    CONSTRAINT chk_adhesion_paiement CHECK (statut_paiement IN ('en_attente','paye','annule'))
);

CREATE INDEX idx_adhesion_utilisateur ON adhesion (utilisateur_id);
CREATE INDEX idx_adhesion_date_fin    ON adhesion (date_fin);


CREATE TABLE commercant (
    id             SERIAL PRIMARY KEY,
    utilisateur_id INTEGER      NOT NULL UNIQUE REFERENCES utilisateur (id) ON DELETE CASCADE,
    raison_sociale VARCHAR(150) NOT NULL,
    siret          VARCHAR(14)  NOT NULL UNIQUE,
    type_commerce  VARCHAR(50),
    adresse        VARCHAR(255) NOT NULL,
    ville          VARCHAR(100) NOT NULL,
    code_postal    VARCHAR(10)  NOT NULL,
    horaires       TEXT,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_commercant_ville ON commercant (ville);


CREATE TABLE profil_benevole (
    id                 SERIAL PRIMARY KEY,
    utilisateur_id     INTEGER     NOT NULL UNIQUE REFERENCES utilisateur (id) ON DELETE CASCADE,
    statut_candidature VARCHAR(20) NOT NULL DEFAULT 'candidat',
    date_candidature   DATE        NOT NULL DEFAULT CURRENT_DATE,
    disponibilites     TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_benevole_statut
        CHECK (statut_candidature IN ('candidat','valide','refuse','inactif'))
);


CREATE TABLE competence (
    id      SERIAL PRIMARY KEY,
    code    VARCHAR(50) NOT NULL UNIQUE,
    libelle VARCHAR(100) NOT NULL
);


CREATE TABLE benevole_competence (
    profil_benevole_id INTEGER NOT NULL REFERENCES profil_benevole (id) ON DELETE CASCADE,
    competence_id      INTEGER NOT NULL REFERENCES competence (id)      ON DELETE CASCADE,

    PRIMARY KEY (profil_benevole_id, competence_id)
);


-- =============================================================
-- BLOC 2 : LOGISTIQUE
-- =============================================================

CREATE TABLE vehicule (
    id              SERIAL PRIMARY KEY,
    immatriculation VARCHAR(20)  NOT NULL UNIQUE,
    modele          VARCHAR(100) NOT NULL,
    capacite_kg     INTEGER      NOT NULL DEFAULT 0,
    actif           BOOLEAN      NOT NULL DEFAULT TRUE
);


CREATE TABLE produit (
    id          SERIAL PRIMARY KEY,
    code_barre  VARCHAR(50)  NOT NULL UNIQUE,
    libelle     VARCHAR(150) NOT NULL,
    categorie   VARCHAR(50),
    unite       VARCHAR(20)  NOT NULL DEFAULT 'piece',
    date_limite DATE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Index critique : le sujet exige de retrouver un produit "tres rapidement"
CREATE INDEX idx_produit_code_barre ON produit (code_barre);
CREATE INDEX idx_produit_libelle    ON produit (libelle);


CREATE TABLE beneficiaire (
    id          SERIAL PRIMARY KEY,
    type        VARCHAR(20)  NOT NULL,
    nom         VARCHAR(150) NOT NULL,
    contact     VARCHAR(150),
    telephone   VARCHAR(20),
    adresse     VARCHAR(255) NOT NULL,
    ville       VARCHAR(100) NOT NULL,
    code_postal VARCHAR(10),
    actif       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_beneficiaire_type CHECK (type IN ('association','particulier'))
);


-- La collecte se fait soit chez un commercant, soit chez un particulier donateur
-- (le sujet mentionne les deux). Le CHECK garantit qu'une seule source est renseignee.
CREATE TABLE collecte (
    id            SERIAL PRIMARY KEY,
    commercant_id INTEGER     REFERENCES commercant (id)  ON DELETE SET NULL,
    donateur_id   INTEGER     REFERENCES utilisateur (id) ON DELETE SET NULL,
    vehicule_id   INTEGER     REFERENCES vehicule (id)    ON DELETE SET NULL,
    chauffeur_id  INTEGER     REFERENCES utilisateur (id) ON DELETE SET NULL,
    adresse_collecte VARCHAR(255),
    date_prevue   TIMESTAMPTZ NOT NULL,
    date_realisee TIMESTAMPTZ,
    statut        VARCHAR(20) NOT NULL DEFAULT 'planifiee',
    commentaire   TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_collecte_statut
        CHECK (statut IN ('planifiee','en_cours','terminee','annulee')),

    CONSTRAINT chk_collecte_source
        CHECK (
            (commercant_id IS NOT NULL AND donateur_id IS NULL)
         OR (commercant_id IS NULL     AND donateur_id IS NOT NULL)
        )
);

CREATE INDEX idx_collecte_date   ON collecte (date_prevue);
CREATE INDEX idx_collecte_statut ON collecte (statut);


CREATE TABLE collecte_produit (
    collecte_id INTEGER NOT NULL REFERENCES collecte (id) ON DELETE CASCADE,
    produit_id  INTEGER NOT NULL REFERENCES produit (id)  ON DELETE RESTRICT,
    quantite    INTEGER NOT NULL,

    PRIMARY KEY (collecte_id, produit_id),
    CONSTRAINT chk_collecte_produit_qte CHECK (quantite > 0)
);


CREATE TABLE tournee (
    id            SERIAL PRIMARY KEY,
    reference     VARCHAR(30) NOT NULL UNIQUE,
    vehicule_id   INTEGER     REFERENCES vehicule (id)    ON DELETE SET NULL,
    chauffeur_id  INTEGER     REFERENCES utilisateur (id) ON DELETE SET NULL,
    date_prevue   TIMESTAMPTZ NOT NULL,
    date_realisee TIMESTAMPTZ,
    statut        VARCHAR(20) NOT NULL DEFAULT 'planifiee',
    commentaire   TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_tournee_statut
        CHECK (statut IN ('planifiee','en_cours','terminee','annulee'))
);

CREATE INDEX idx_tournee_date   ON tournee (date_prevue);
CREATE INDEX idx_tournee_statut ON tournee (statut);


CREATE TABLE arret (
    id              SERIAL PRIMARY KEY,
    tournee_id      INTEGER NOT NULL REFERENCES tournee (id)      ON DELETE CASCADE,
    beneficiaire_id INTEGER NOT NULL REFERENCES beneficiaire (id) ON DELETE RESTRICT,
    ordre_passage   INTEGER NOT NULL,
    heure_prevue    TIME,
    livre           BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT uq_arret_ordre UNIQUE (tournee_id, ordre_passage)
);

CREATE INDEX idx_arret_tournee ON arret (tournee_id);


CREATE TABLE tournee_produit (
    tournee_id INTEGER NOT NULL REFERENCES tournee (id) ON DELETE CASCADE,
    produit_id INTEGER NOT NULL REFERENCES produit (id) ON DELETE RESTRICT,
    quantite   INTEGER NOT NULL,

    PRIMARY KEY (tournee_id, produit_id),
    CONSTRAINT chk_tournee_produit_qte CHECK (quantite > 0)
);


-- =============================================================
-- BLOC 3 : SERVICES ET I18N
-- =============================================================

CREATE TABLE langue (
    id      SERIAL PRIMARY KEY,
    code    VARCHAR(5)  NOT NULL UNIQUE,
    libelle VARCHAR(50) NOT NULL,
    actif   BOOLEAN     NOT NULL DEFAULT TRUE
);


CREATE TABLE service (
    id            SERIAL PRIMARY KEY,
    code          VARCHAR(50) NOT NULL UNIQUE,
    competence_id INTEGER     REFERENCES competence (id) ON DELETE SET NULL,
    actif         BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE service_traduction (
    id          SERIAL PRIMARY KEY,
    service_id  INTEGER      NOT NULL REFERENCES service (id) ON DELETE CASCADE,
    langue_id   INTEGER      NOT NULL REFERENCES langue (id)  ON DELETE CASCADE,
    libelle     VARCHAR(150) NOT NULL,
    description TEXT,

    CONSTRAINT uq_service_traduction UNIQUE (service_id, langue_id)
);


CREATE TABLE creneau (
    id           SERIAL PRIMARY KEY,
    service_id   INTEGER     NOT NULL REFERENCES service (id) ON DELETE CASCADE,
    date_creneau DATE        NOT NULL,
    heure_debut  TIME        NOT NULL,
    heure_fin    TIME        NOT NULL,
    lieu         VARCHAR(255),
    capacite_max INTEGER     NOT NULL DEFAULT 10,
    statut       VARCHAR(20) NOT NULL DEFAULT 'ouvert',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_creneau_heures   CHECK (heure_fin > heure_debut),
    CONSTRAINT chk_creneau_capacite CHECK (capacite_max > 0),
    CONSTRAINT chk_creneau_statut   CHECK (statut IN ('ouvert','complet','annule','termine'))
);

-- Index critique : l'export Excel filtre le planning par date
CREATE INDEX idx_creneau_date    ON creneau (date_creneau);
CREATE INDEX idx_creneau_service ON creneau (service_id);


CREATE TABLE inscription (
    id               SERIAL PRIMARY KEY,
    creneau_id       INTEGER     NOT NULL REFERENCES creneau (id)     ON DELETE CASCADE,
    utilisateur_id   INTEGER     NOT NULL REFERENCES utilisateur (id) ON DELETE CASCADE,
    date_inscription TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    statut           VARCHAR(20) NOT NULL DEFAULT 'confirmee',

    CONSTRAINT uq_inscription     UNIQUE (creneau_id, utilisateur_id),
    CONSTRAINT chk_inscription_st CHECK (statut IN ('confirmee','annulee','presente','absente'))
);

CREATE INDEX idx_inscription_creneau ON inscription (creneau_id);


CREATE TABLE affectation (
    id                 SERIAL PRIMARY KEY,
    creneau_id         INTEGER     NOT NULL REFERENCES creneau (id)         ON DELETE CASCADE,
    profil_benevole_id INTEGER     NOT NULL REFERENCES profil_benevole (id) ON DELETE CASCADE,
    date_affectation   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    statut             VARCHAR(20) NOT NULL DEFAULT 'proposee',

    CONSTRAINT uq_affectation     UNIQUE (creneau_id, profil_benevole_id),
    CONSTRAINT chk_affectation_st CHECK (statut IN ('proposee','acceptee','refusee','realisee'))
);

CREATE INDEX idx_affectation_creneau  ON affectation (creneau_id);
CREATE INDEX idx_affectation_benevole ON affectation (profil_benevole_id);


-- =============================================================
-- VUE : stock courant (D8 — stock calcule, pas stocke)
-- =============================================================

CREATE OR REPLACE VIEW v_stock AS
SELECT
    p.id,
    p.code_barre,
    p.libelle,
    p.categorie,
    p.date_limite,
    COALESCE(entrees.total, 0) - COALESCE(sorties.total, 0) AS quantite_stock
FROM produit p
LEFT JOIN (
    SELECT produit_id, SUM(quantite) AS total
    FROM collecte_produit
    GROUP BY produit_id
) AS entrees ON entrees.produit_id = p.id
LEFT JOIN (
    SELECT produit_id, SUM(quantite) AS total
    FROM tournee_produit
    GROUP BY produit_id
) AS sorties ON sorties.produit_id = p.id;
