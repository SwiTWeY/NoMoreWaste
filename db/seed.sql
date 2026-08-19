-- =============================================================
-- NO MORE WASTE — jeu de donnees de demonstration
-- A executer APRES schema.sql
--
-- Mot de passe de tous les comptes : password123
-- Hash bcrypt (cout 10) : $2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG
--
-- Toutes les dates sont RELATIVES a CURRENT_DATE.
-- Le seed reste donc coherent quel que soit le jour ou on le lance.
-- =============================================================

TRUNCATE TABLE
    affectation, inscription, creneau, service_traduction, service, langue,
    tournee_produit, collecte_produit, arret, tournee, collecte,
    produit, beneficiaire, vehicule,
    benevole_competence, competence, profil_benevole, commercant, adhesion, utilisateur
RESTART IDENTITY CASCADE;


-- =============================================================
-- LANGUES
-- =============================================================

INSERT INTO langue (code, libelle) VALUES
    ('fr', 'Francais'),
    ('en', 'English');


-- =============================================================
-- COMPETENCES
-- =============================================================

INSERT INTO competence (code, libelle) VALUES
    ('CHAUFFEUR',   'Chauffeur'),
    ('CUISINIER',   'Cuisinier'),
    ('PLOMBIER',    'Plombier'),
    ('ELECTRICIEN', 'Electricien'),
    ('BRICOLEUR',   'Bricoleur'),
    ('MANUTENTION', 'Manutention'),
    ('ACCUEIL',     'Accueil et administratif');


-- =============================================================
-- UTILISATEURS
--   1-2   : personnel du siege (back-office)
--   3-7   : gerants de commerces
--   8-14  : benevoles
--   15-18 : adherents particuliers
-- =============================================================

INSERT INTO utilisateur (nom, prenom, email, mot_de_passe, telephone, est_personnel, langue_pref) VALUES
    ('Moreau',   'Sylvie',   'sylvie.moreau@nomorewaste.org',  '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0142380192', TRUE,  'fr'),
    ('Diallo',   'Amadou',   'amadou.diallo@nomorewaste.org',  '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0142380193', TRUE,  'fr'),

    ('Lefevre',  'Claire',   'claire.lefevre@boulangerie.fr',  '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0145227781', FALSE, 'fr'),
    ('Nguyen',   'Tuan',     'tuan.nguyen@primeurs-est.fr',    '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0143901220', FALSE, 'fr'),
    ('Bonnet',   'Philippe', 'p.bonnet@supermarche-nation.fr', '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0143678845', FALSE, 'fr'),
    ('Rossi',    'Giulia',   'giulia.rossi@traiteur-italia.fr','$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0148055512', FALSE, 'fr'),
    ('O''Brien', 'Sean',     'sean.obrien@dublin-grocer.ie',   '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0035312345', FALSE, 'en'),

    ('Martin',   'Julien',   'julien.martin@email.fr',         '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345601', FALSE, 'fr'),
    ('Benali',   'Karim',    'karim.benali@email.fr',          '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345602', FALSE, 'fr'),
    ('Dupont',   'Marie',    'marie.dupont@email.fr',          '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345603', FALSE, 'fr'),
    ('Fontaine', 'Lucas',    'lucas.fontaine@email.fr',        '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345604', FALSE, 'fr'),
    ('Chen',     'Wei',      'wei.chen@email.fr',              '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345605', FALSE, 'fr'),
    ('Garcia',   'Elena',    'elena.garcia@email.fr',          '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345606', FALSE, 'fr'),
    ('Petit',    'Thomas',   'thomas.petit@email.fr',          '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345607', FALSE, 'fr'),

    ('Leroy',    'Nathalie', 'nathalie.leroy@email.fr',        '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345608', FALSE, 'fr'),
    ('Ferreira', 'Joao',     'joao.ferreira@email.pt',         '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0351912345', FALSE, 'en'),
    ('Roux',     'Camille',  'camille.roux@email.fr',          '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345609', FALSE, 'fr'),
    ('Mercier',  'Antoine',  'antoine.mercier@email.fr',       '$2b$10$zYRma7QSIgD8eqoNWQYll.LUViLJwlhtyeq6KZ4AuaqBE4U6LeDHG', '0612345610', FALSE, 'fr');


-- =============================================================
-- COMMERCANTS
-- =============================================================

INSERT INTO commercant (utilisateur_id, raison_sociale, siret, type_commerce, adresse, ville, code_postal, horaires) VALUES
    (3, 'Boulangerie Lefevre',     '81234567800011', 'Boulangerie',  '18 rue de Charonne',        'Paris',   '75011', 'Lun-Sam 6h-20h'),
    (4, 'Primeurs de l''Est',      '81234567800022', 'Primeur',      '45 avenue Philippe Auguste','Paris',   '75011', 'Mar-Dim 8h-19h30'),
    (5, 'Supermarche Nation',      '81234567800033', 'Supermarche',  '3 place de la Nation',      'Paris',   '75012', 'Lun-Dim 8h30-21h'),
    (6, 'Traiteur Italia',         '81234567800044', 'Traiteur',     '27 rue du Faubourg',        'Nantes',  '44000', 'Mar-Sam 10h-22h'),
    (7, 'Dublin Corner Grocer',    '81234567800055', 'Epicerie',     '12 Camden Street',          'Dublin',  'D02',   'Mon-Sun 7h-23h');


-- =============================================================
-- ADHESIONS
--   Volontairement heterogenes pour que le rappel de renouvellement
--   ait de la matiere a montrer en demo :
--     - id 3 : expire dans 12 jours  -> DOIT declencher un rappel
--     - id 5 : expire dans 25 jours  -> a la limite
--     - id 7 : DEJA EXPIREE          -> a relancer
--     - id 4 : historique 2 lignes   -> montre le 1-N (D2)
-- =============================================================

INSERT INTO adhesion (utilisateur_id, date_debut, date_fin, montant, statut_paiement) VALUES
    -- Commercants
    (3, CURRENT_DATE - INTERVAL '353 days', CURRENT_DATE + INTERVAL '12 days',  120.00, 'paye'),
    (4, CURRENT_DATE - INTERVAL '730 days', CURRENT_DATE - INTERVAL '365 days', 120.00, 'paye'),
    (4, CURRENT_DATE - INTERVAL '365 days', CURRENT_DATE + INTERVAL '65 days',  120.00, 'paye'),
    (5, CURRENT_DATE - INTERVAL '340 days', CURRENT_DATE + INTERVAL '25 days',  250.00, 'paye'),
    (6, CURRENT_DATE - INTERVAL '200 days', CURRENT_DATE + INTERVAL '165 days', 120.00, 'paye'),
    (7, CURRENT_DATE - INTERVAL '400 days', CURRENT_DATE - INTERVAL '35 days',  120.00, 'paye'),

    -- Adherents particuliers
    (15, CURRENT_DATE - INTERVAL '100 days', CURRENT_DATE + INTERVAL '265 days', 25.00, 'paye'),
    (16, CURRENT_DATE - INTERVAL '50 days',  CURRENT_DATE + INTERVAL '315 days', 25.00, 'paye'),
    (17, CURRENT_DATE - INTERVAL '355 days', CURRENT_DATE + INTERVAL '10 days',  25.00, 'paye'),
    (18, CURRENT_DATE - INTERVAL '10 days',  CURRENT_DATE + INTERVAL '355 days', 25.00, 'en_attente'),

    -- Benevoles egalement adherents (montre le cumul de roles, D1)
    (8,  CURRENT_DATE - INTERVAL '120 days', CURRENT_DATE + INTERVAL '245 days', 25.00, 'paye'),
    (10, CURRENT_DATE - INTERVAL '80 days',  CURRENT_DATE + INTERVAL '285 days', 25.00, 'paye');


-- =============================================================
-- PROFILS BENEVOLES
-- =============================================================

INSERT INTO profil_benevole (utilisateur_id, statut_candidature, date_candidature, disponibilites) VALUES
    (8,  'valide',   CURRENT_DATE - INTERVAL '150 days', 'Lundi et mercredi matin, samedi toute la journee'),
    (9,  'valide',   CURRENT_DATE - INTERVAL '120 days', 'Mardi et jeudi apres-midi'),
    (10, 'valide',   CURRENT_DATE - INTERVAL '95 days',  'Weekends uniquement'),
    (11, 'valide',   CURRENT_DATE - INTERVAL '60 days',  'Tous les soirs apres 18h'),
    (12, 'valide',   CURRENT_DATE - INTERVAL '45 days',  'Mercredi journee complete'),
    (13, 'candidat', CURRENT_DATE - INTERVAL '8 days',   'Flexible, en recherche d''emploi'),
    (14, 'candidat', CURRENT_DATE - INTERVAL '3 days',   'Vendredi et samedi');


-- =============================================================
-- COMPETENCES DES BENEVOLES
--   Julien (1) cumule chauffeur + manutention
--   Karim  (2) est le cuisinier principal
-- =============================================================

INSERT INTO benevole_competence (profil_benevole_id, competence_id) VALUES
    (1, 1), (1, 6),
    (2, 2), (2, 7),
    (3, 1), (3, 2),
    (4, 3), (4, 5),
    (5, 4), (5, 5),
    (6, 6),
    (7, 7), (7, 2);


-- =============================================================
-- VEHICULES
-- =============================================================

INSERT INTO vehicule (immatriculation, modele, capacite_kg) VALUES
    ('AB-123-CD', 'Renault Master frigorifique', 1200),
    ('EF-456-GH', 'Citroen Jumper',               900),
    ('IJ-789-KL', 'Peugeot Partner',              450),
    ('MN-012-OP', 'Iveco Daily frigorifique',    1800);


-- =============================================================
-- BENEFICIAIRES
-- =============================================================

INSERT INTO beneficiaire (type, nom, contact, telephone, adresse, ville, code_postal) VALUES
    ('association', 'Restos du Coeur - Paris 11',  'Mme Dubois',   '0143551200', '8 rue Keller',            'Paris',    '75011'),
    ('association', 'Secours Populaire - Nation',  'M. Lambert',   '0143551201', '22 boulevard Voltaire',   'Paris',    '75011'),
    ('association', 'Foyer Solidarite Est',        'Mme Kaplan',   '0143551202', '5 rue de Montreuil',      'Paris',    '75011'),
    ('association', 'Epicerie Solidaire Nantes',   'M. Guerin',    '0240551203', '14 rue Crebillon',        'Nantes',   '44000'),
    ('particulier',  'Famille Sow',                 NULL,           '0612000101', '31 rue Sedaine',          'Paris',    '75011'),
    ('particulier',  'Mme Vasseur',                 NULL,           '0612000102', '9 passage Josset',        'Paris',    '75011'),
    ('particulier',  'M. Traore',                   NULL,           '0612000103', '47 rue Saint-Maur',       'Paris',    '75011'),
    ('particulier',  'Famille Andrade',             NULL,           '0612000104', '2 rue du Calvaire',       'Nantes',   '44000');


-- =============================================================
-- PRODUITS
--   Codes barres EAN-13 plausibles.
--   Certaines dates limites sont proches ou depassees : utile pour
--   demontrer un filtre "produits a distribuer en priorite".
-- =============================================================

INSERT INTO produit (code_barre, libelle, categorie, unite, date_limite) VALUES
    ('3560070139682', 'Baguette tradition',           'Boulangerie', 'piece', CURRENT_DATE + INTERVAL '1 day'),
    ('3560070139699', 'Pain de campagne 500g',        'Boulangerie', 'piece', CURRENT_DATE + INTERVAL '2 days'),
    ('3560070139705', 'Croissants (lot de 6)',        'Boulangerie', 'lot',   CURRENT_DATE + INTERVAL '1 day'),
    ('3560070139712', 'Brioche tressee',              'Boulangerie', 'piece', CURRENT_DATE + INTERVAL '3 days'),

    ('3245412765431', 'Pommes Golden',                'Fruits',      'kg',    CURRENT_DATE + INTERVAL '8 days'),
    ('3245412765448', 'Bananes',                      'Fruits',      'kg',    CURRENT_DATE + INTERVAL '4 days'),
    ('3245412765455', 'Oranges a jus',                'Fruits',      'kg',    CURRENT_DATE + INTERVAL '6 days'),
    ('3245412765462', 'Poires Williams',              'Fruits',      'kg',    CURRENT_DATE + INTERVAL '3 days'),
    ('3245412765479', 'Fraises barquette 250g',       'Fruits',      'piece', CURRENT_DATE + INTERVAL '1 day'),

    ('3011360004512', 'Carottes',                     'Legumes',     'kg',    CURRENT_DATE + INTERVAL '12 days'),
    ('3011360004529', 'Pommes de terre',              'Legumes',     'kg',    CURRENT_DATE + INTERVAL '25 days'),
    ('3011360004536', 'Courgettes',                   'Legumes',     'kg',    CURRENT_DATE + INTERVAL '5 days'),
    ('3011360004543', 'Tomates grappe',               'Legumes',     'kg',    CURRENT_DATE + INTERVAL '4 days'),
    ('3011360004550', 'Salade batavia',               'Legumes',     'piece', CURRENT_DATE + INTERVAL '2 days'),
    ('3011360004567', 'Oignons jaunes',               'Legumes',     'kg',    CURRENT_DATE + INTERVAL '30 days'),

    ('3033490004521', 'Lait demi-ecreme 1L',          'Cremerie',    'piece', CURRENT_DATE + INTERVAL '7 days'),
    ('3033490004538', 'Yaourts nature (pack de 8)',   'Cremerie',    'lot',   CURRENT_DATE + INTERVAL '5 days'),
    ('3033490004545', 'Beurre doux 250g',             'Cremerie',    'piece', CURRENT_DATE + INTERVAL '14 days'),
    ('3033490004552', 'Emmental rape 200g',           'Cremerie',    'piece', CURRENT_DATE + INTERVAL '9 days'),
    ('3033490004569', 'Fromage blanc 500g',           'Cremerie',    'piece', CURRENT_DATE + INTERVAL '4 days'),

    ('8001234567890', 'Pates penne 500g',             'Epicerie',    'piece', CURRENT_DATE + INTERVAL '400 days'),
    ('8001234567906', 'Riz long grain 1kg',           'Epicerie',    'piece', CURRENT_DATE + INTERVAL '500 days'),
    ('8001234567913', 'Sauce tomate basilic 400g',    'Epicerie',    'piece', CURRENT_DATE + INTERVAL '300 days'),
    ('8001234567920', 'Huile d''olive 75cl',          'Epicerie',    'piece', CURRENT_DATE + INTERVAL '250 days'),
    ('8001234567937', 'Lentilles vertes 500g',        'Epicerie',    'piece', CURRENT_DATE + INTERVAL '450 days'),
    ('8001234567944', 'Conserve haricots verts',      'Epicerie',    'piece', CURRENT_DATE + INTERVAL '600 days'),
    ('8001234567951', 'Farine T55 1kg',               'Epicerie',    'piece', CURRENT_DATE + INTERVAL '200 days'),
    ('8001234567968', 'Sucre en poudre 1kg',          'Epicerie',    'piece', CURRENT_DATE + INTERVAL '700 days'),

    ('3760091725431', 'Jus d''orange 1L',             'Boissons',    'piece', CURRENT_DATE + INTERVAL '60 days'),
    ('3760091725448', 'Eau minerale (pack de 6)',     'Boissons',    'lot',   CURRENT_DATE + INTERVAL '365 days'),

    ('3564700223451', 'Plat prepare lasagnes',        'Traiteur',    'piece', CURRENT_DATE + INTERVAL '2 days'),
    ('3564700223468', 'Quiche lorraine',              'Traiteur',    'piece', CURRENT_DATE + INTERVAL '1 day'),
    ('3564700223475', 'Salade composee 300g',         'Traiteur',    'piece', CURRENT_DATE),
    ('3564700223482', 'Soupe de legumes 1L',          'Traiteur',    'piece', CURRENT_DATE + INTERVAL '3 days'),

    ('3168930100119', 'Oeufs (boite de 12)',          'Cremerie',    'lot',   CURRENT_DATE + INTERVAL '11 days'),
    ('3168930100126', 'Jambon blanc 4 tranches',      'Charcuterie', 'piece', CURRENT_DATE + INTERVAL '3 days'),
    ('3168930100133', 'Filet de poulet 500g',         'Boucherie',   'piece', CURRENT_DATE + INTERVAL '2 days'),
    ('3168930100140', 'Steaks haches x4',             'Boucherie',   'lot',   CURRENT_DATE + INTERVAL '1 day'),
    ('3168930100157', 'Saumon fume 200g',             'Poissonnerie','piece', CURRENT_DATE + INTERVAL '5 days'),
    ('3168930100164', 'Cabillaud surgele 400g',       'Poissonnerie','piece', CURRENT_DATE + INTERVAL '180 days');


-- =============================================================
-- COLLECTES
--   Melange de collectes terminees (alimentent le stock),
--   en cours, et planifiees.
-- =============================================================

-- Collectes chez des commercants
INSERT INTO collecte (commercant_id, donateur_id, vehicule_id, chauffeur_id, adresse_collecte, date_prevue, date_realisee, statut, commentaire) VALUES
    (1, NULL, 1, 8,  NULL, CURRENT_DATE - INTERVAL '6 days'  + TIME '07:00', CURRENT_DATE - INTERVAL '6 days'  + TIME '07:35', 'terminee',  'Invendus de la veille'),
    (2, NULL, 2, 10, NULL, CURRENT_DATE - INTERVAL '5 days'  + TIME '09:00', CURRENT_DATE - INTERVAL '5 days'  + TIME '09:50', 'terminee',  NULL),
    (3, NULL, 4, 8,  NULL, CURRENT_DATE - INTERVAL '4 days'  + TIME '18:00', CURRENT_DATE - INTERVAL '4 days'  + TIME '19:10', 'terminee',  'Grosse collecte fin de journee'),
    (1, NULL, 1, 10, NULL, CURRENT_DATE - INTERVAL '3 days'  + TIME '07:00', CURRENT_DATE - INTERVAL '3 days'  + TIME '07:30', 'terminee',  NULL),
    (4, NULL, 2, 8,  NULL, CURRENT_DATE - INTERVAL '2 days'  + TIME '14:00', CURRENT_DATE - INTERVAL '2 days'  + TIME '14:45', 'terminee',  NULL),
    (3, NULL, 4, 10, NULL, CURRENT_DATE - INTERVAL '1 day'   + TIME '18:00', CURRENT_DATE - INTERVAL '1 day'   + TIME '18:55', 'terminee',  NULL),
    (2, NULL, 2, 8,  NULL, CURRENT_DATE + TIME '09:00',                      NULL,                                              'en_cours',  'En cours ce matin'),
    (1, NULL, 1, 10, NULL, CURRENT_DATE + INTERVAL '1 day'   + TIME '07:00', NULL,                                              'planifiee', NULL),
    (5, NULL, 3, 8,  NULL, CURRENT_DATE + INTERVAL '2 days'  + TIME '10:00', NULL,                                              'planifiee', 'Premiere collecte a Dublin'),
    (3, NULL, 4, 10, NULL, CURRENT_DATE + INTERVAL '3 days'  + TIME '18:00', NULL,                                              'planifiee', NULL);

-- Collectes chez des particuliers donateurs (produits proches de la date limite)
INSERT INTO collecte (commercant_id, donateur_id, vehicule_id, chauffeur_id, adresse_collecte, date_prevue, date_realisee, statut, commentaire) VALUES
    (NULL, 15, 3, 8,  '12 rue Oberkampf, 75011 Paris',  CURRENT_DATE - INTERVAL '4 days' + TIME '11:00', CURRENT_DATE - INTERVAL '4 days' + TIME '11:20', 'terminee',  'Depart en vacances'),
    (NULL, 17, 3, 10, '5 rue Popincourt, 75011 Paris',  CURRENT_DATE - INTERVAL '2 days' + TIME '16:00', CURRENT_DATE - INTERVAL '2 days' + TIME '16:15', 'terminee',  NULL),
    (NULL, 18, 3, 8,  '88 avenue Parmentier, 75011 Paris', CURRENT_DATE + INTERVAL '1 day' + TIME '15:00', NULL,                                           'planifiee', 'Demande via le front-office');


INSERT INTO collecte_produit (collecte_id, produit_id, quantite) VALUES
    -- Collecte 1 : boulangerie
    (1, 1, 40), (1, 2, 15), (1, 3, 12), (1, 4, 8),
    -- Collecte 2 : primeur
    (2, 5, 30), (2, 6, 25), (2, 10, 40), (2, 12, 18), (2, 13, 22), (2, 14, 15),
    -- Collecte 3 : supermarche (grosse)
    (3, 16, 60), (3, 17, 30), (3, 18, 25), (3, 21, 50), (3, 22, 40),
    (3, 23, 35), (3, 26, 45), (3, 35, 20), (3, 36, 15), (3, 37, 12),
    -- Collecte 4 : boulangerie
    (4, 1, 35), (4, 2, 10), (4, 3, 18),
    -- Collecte 5 : traiteur
    (5, 31, 14), (5, 32, 10), (5, 33, 8), (5, 34, 12),
    -- Collecte 6 : supermarche
    (6, 7, 28), (6, 11, 55), (6, 19, 20), (6, 24, 18), (6, 27, 30),
    (6, 29, 24), (6, 30, 20), (6, 38, 16), (6, 39, 10),
    -- Collectes 11 et 12 : particuliers donateurs (petits volumes)
    (11, 21, 3), (11, 22, 2), (11, 26, 4), (11, 28, 1),
    (12, 16, 2), (12, 18, 1), (12, 23, 3), (12, 25, 2);


-- =============================================================
-- TOURNEES ET ARRETS
-- =============================================================

INSERT INTO tournee (reference, vehicule_id, chauffeur_id, date_prevue, date_realisee, statut, commentaire) VALUES
    ('TRN-2026-0001', 1, 8,  CURRENT_DATE - INTERVAL '5 days' + TIME '10:00', CURRENT_DATE - INTERVAL '5 days' + TIME '13:20', 'terminee',  'Tournee Paris Est'),
    ('TRN-2026-0002', 2, 10, CURRENT_DATE - INTERVAL '3 days' + TIME '10:00', CURRENT_DATE - INTERVAL '3 days' + TIME '14:05', 'terminee',  'Tournee Paris 11e'),
    ('TRN-2026-0003', 4, 8,  CURRENT_DATE - INTERVAL '1 day'  + TIME '09:30', CURRENT_DATE - INTERVAL '1 day'  + TIME '13:45', 'terminee',  NULL),
    ('TRN-2026-0004', 1, 10, CURRENT_DATE + TIME '10:00',                     NULL,                                             'en_cours',  'Tournee du jour'),
    ('TRN-2026-0005', 2, 8,  CURRENT_DATE + INTERVAL '1 day' + TIME '10:00',  NULL,                                             'planifiee', NULL),
    ('TRN-2026-0006', 3, 10, CURRENT_DATE + INTERVAL '2 days' + TIME '11:00', NULL,                                             'planifiee', 'Tournee Nantes');


INSERT INTO arret (tournee_id, beneficiaire_id, ordre_passage, heure_prevue, livre) VALUES
    (1, 1, 1, '10:30', TRUE), (1, 2, 2, '11:15', TRUE), (1, 5, 3, '12:00', TRUE),
    (2, 3, 1, '10:30', TRUE), (2, 6, 2, '11:30', TRUE), (2, 7, 3, '12:15', TRUE), (2, 1, 4, '13:00', TRUE),
    (3, 2, 1, '10:00', TRUE), (3, 1, 2, '11:00', TRUE), (3, 3, 3, '12:00', TRUE),
    (4, 1, 1, '10:30', TRUE), (4, 5, 2, '11:15', FALSE), (4, 6, 3, '12:00', FALSE),
    (5, 2, 1, '10:30', FALSE), (5, 3, 2, '11:30', FALSE), (5, 7, 3, '12:30', FALSE),
    (6, 4, 1, '11:30', FALSE), (6, 8, 2, '12:30', FALSE);


INSERT INTO tournee_produit (tournee_id, produit_id, quantite) VALUES
    (1, 1, 25), (1, 5, 15), (1, 10, 20), (1, 16, 20), (1, 21, 25),
    (2, 2, 12), (2, 6, 15), (2, 12, 10), (2, 17, 15), (2, 22, 20), (2, 26, 20),
    (3, 3, 20), (3, 13, 12), (3, 18, 12), (3, 23, 18), (3, 35, 10), (3, 36, 8),
    (4, 1, 20), (4, 11, 25), (4, 27, 15), (4, 29, 12);


-- =============================================================
-- SERVICES ET TRADUCTIONS
-- =============================================================

INSERT INTO service (code, competence_id) VALUES
    ('CONSEILS_ANTIGASPI', 7),
    ('COURS_CUISINE',      2),
    ('PARTAGE_VEHICULE',   1),
    ('REPARATION',         5),
    ('PLOMBERIE',          3),
    ('ELECTRICITE',        4),
    ('GARDIENNAGE',        7);


INSERT INTO service_traduction (service_id, langue_id, libelle, description) VALUES
    (1, 1, 'Conseils anti-gaspi',      'Ateliers pratiques pour reduire le gaspillage alimentaire au quotidien.'),
    (1, 2, 'Anti-waste advice',        'Practical workshops to reduce daily food waste.'),
    (2, 1, 'Cours de cuisine',         'Apprendre a cuisiner les restes et les produits de saison.'),
    (2, 2, 'Cooking class',            'Learn to cook leftovers and seasonal produce.'),
    (3, 1, 'Partage de vehicules',     'Mise en relation entre adherents pour partager un vehicule.'),
    (3, 2, 'Vehicle sharing',          'Connecting members to share a vehicle.'),
    (4, 1, 'Service de reparation',    'Reparation de petit electromenager et de mobilier.'),
    (4, 2, 'Repair service',           'Repair of small appliances and furniture.'),
    (5, 1, 'Plomberie',                'Interventions de plomberie realisees par des benevoles qualifies.'),
    (5, 2, 'Plumbing',                 'Plumbing work carried out by qualified volunteers.'),
    (6, 1, 'Electricite',              'Depannage electrique par des benevoles habilites.'),
    (6, 2, 'Electrical work',          'Electrical repairs by certified volunteers.'),
    (7, 1, 'Gardiennage',              'Garde ponctuelle de logement ou d''animaux entre adherents.'),
    (7, 2, 'House sitting',            'Occasional home or pet sitting between members.');


-- =============================================================
-- CRENEAUX
--   Repartis sur la semaine ecoulee et la semaine a venir,
--   pour que l'export Excel du planning ait toujours des donnees.
-- =============================================================

INSERT INTO creneau (service_id, date_creneau, heure_debut, heure_fin, lieu, capacite_max, statut) VALUES
    (2, CURRENT_DATE - INTERVAL '4 days', '14:00', '17:00', 'Cuisine du siege - Paris 11e',  12, 'termine'),
    (1, CURRENT_DATE - INTERVAL '2 days', '10:00', '12:00', 'Salle associative - Paris 11e', 20, 'termine'),

    (2, CURRENT_DATE,                     '14:00', '17:00', 'Cuisine du siege - Paris 11e',  12, 'ouvert'),
    (4, CURRENT_DATE,                     '09:00', '12:00', 'Atelier - Paris 12e',            8, 'ouvert'),

    (1, CURRENT_DATE + INTERVAL '1 day',  '10:00', '12:00', 'Salle associative - Paris 11e', 20, 'ouvert'),
    (3, CURRENT_DATE + INTERVAL '1 day',  '08:00', '18:00', 'Parking du siege',               4, 'ouvert'),
    (5, CURRENT_DATE + INTERVAL '2 days', '09:00', '13:00', 'A domicile',                     3, 'ouvert'),
    (2, CURRENT_DATE + INTERVAL '3 days', '18:00', '21:00', 'Cuisine du siege - Paris 11e',  12, 'ouvert'),
    (6, CURRENT_DATE + INTERVAL '4 days', '14:00', '18:00', 'A domicile',                     3, 'ouvert'),
    (4, CURRENT_DATE + INTERVAL '5 days', '09:00', '12:00', 'Atelier - Paris 12e',            8, 'ouvert'),
    (7, CURRENT_DATE + INTERVAL '6 days', '08:00', '20:00', 'A domicile',                     5, 'ouvert'),
    (2, CURRENT_DATE + INTERVAL '7 days', '14:00', '17:00', 'Cuisine du siege - Nantes',     10, 'ouvert');


-- =============================================================
-- AFFECTATIONS (benevoles qui animent)
--   Plusieurs benevoles possibles par creneau (D11).
-- =============================================================

INSERT INTO affectation (creneau_id, profil_benevole_id, statut) VALUES
    (1, 2, 'realisee'), (1, 7, 'realisee'),
    (2, 6, 'realisee'),
    (3, 2, 'acceptee'), (3, 3, 'acceptee'),
    (4, 4, 'acceptee'), (4, 5, 'proposee'),
    (5, 6, 'acceptee'),
    (6, 1, 'acceptee'), (6, 3, 'acceptee'),
    (7, 4, 'acceptee'),
    (8, 2, 'proposee'), (8, 7, 'proposee'),
    (9, 5, 'acceptee'),
    (10, 4, 'proposee'), (10, 5, 'acceptee'),
    (11, 6, 'proposee'),
    (12, 2, 'proposee');


-- =============================================================
-- INSCRIPTIONS (adherents qui participent)
-- =============================================================

INSERT INTO inscription (creneau_id, utilisateur_id, statut) VALUES
    (1, 15, 'presente'), (1, 16, 'presente'), (1, 17, 'absente'), (1, 8, 'presente'),
    (2, 15, 'presente'), (2, 18, 'presente'), (2, 10, 'presente'),
    (3, 15, 'confirmee'), (3, 16, 'confirmee'), (3, 17, 'confirmee'), (3, 10, 'confirmee'),
    (4, 16, 'confirmee'), (4, 18, 'confirmee'),
    (5, 15, 'confirmee'), (5, 17, 'confirmee'), (5, 8, 'confirmee'),
    (6, 16, 'confirmee'),
    (7, 17, 'confirmee'),
    (8, 15, 'confirmee'), (8, 16, 'annulee'), (8, 18, 'confirmee'),
    (9, 15, 'confirmee'),
    (10, 17, 'confirmee'), (10, 18, 'confirmee'),
    (11, 16, 'confirmee'),
    (12, 15, 'confirmee');


-- =============================================================
-- CONTROLE
-- =============================================================

SELECT 'utilisateurs' AS table_name, COUNT(*) FROM utilisateur
UNION ALL SELECT 'commercants',   COUNT(*) FROM commercant
UNION ALL SELECT 'adhesions',     COUNT(*) FROM adhesion
UNION ALL SELECT 'benevoles',     COUNT(*) FROM profil_benevole
UNION ALL SELECT 'produits',      COUNT(*) FROM produit
UNION ALL SELECT 'collectes',     COUNT(*) FROM collecte
UNION ALL SELECT 'tournees',      COUNT(*) FROM tournee
UNION ALL SELECT 'beneficiaires', COUNT(*) FROM beneficiaire
UNION ALL SELECT 'creneaux',      COUNT(*) FROM creneau
UNION ALL SELECT 'inscriptions',  COUNT(*) FROM inscription
UNION ALL SELECT 'affectations',  COUNT(*) FROM affectation;
