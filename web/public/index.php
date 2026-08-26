<?php

require __DIR__ . '/../src/lib/Session.php';
require __DIR__ . '/../src/lib/ApiClient.php';
require __DIR__ . '/../src/lib/I18n.php';

$config = require __DIR__ . '/../config/config.php';

if (php_sapi_name() === 'cli-server') {
    $fichier = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
    if ($fichier !== '/' && file_exists(__DIR__ . $fichier)) {
        return false;
    }
}

Session::demarrer();
I18n::init();

function rendre(string $vue, array $vars = [], string $layout = 'layouts/back'): void
{
    extract($vars);
    ob_start();
    require __DIR__ . '/../src/views/' . $vue . '.php';
    $contenu = ob_get_clean();
    require __DIR__ . '/../src/views/' . $layout . '.php';
}

function exigerPersonnel(): void
{
    if (!Session::estConnecte()) {
        header('Location: /login');
        exit;
    }
    if (!Session::estPersonnel()) {
        http_response_code(403);
        echo 'Acces reserve au personnel.';
        exit;
    }
}

$chemin = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);

switch ($chemin) {
    case '/':
        if (!Session::estConnecte()) {
            rendre('accueil-public', ['titre' => 'NO MORE WASTE'], 'layouts/public');
            break;
        }
        if (!Session::estPersonnel()) {
            header('Location: /espace');
            exit;
        }
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('accueil', ['titre' => 'Tableau de bord', 'message' => null, 'stats' => $api->get('/stats')['donnees']]);
        break;
    case '/login':
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $api = new ApiClient($config['api_url']);
            $reponse = $api->post('/login', [
                'email' => $_POST['email'] ?? '',
                'mot_de_passe' => $_POST['mot_de_passe'] ?? '',
            ]);
            if ($reponse['statut'] === 200) {
                Session::connecter($reponse['donnees']['token'], $reponse['donnees']['utilisateur']);
                $_SESSION['est_benevole'] = $reponse['donnees']['est_benevole'] ?? false;
                    header('Location: ' . (Session::estPersonnel() ? '/' : '/espace'));
                exit;
            }
            rendre('auth/login', ['titre' => 'Connexion', 'erreur' => 'Identifiants invalides']);
            break;
        }
        rendre('auth/login', ['titre' => 'Connexion', 'erreur' => null],'layouts/public');
        break;
    case '/espace/agenda':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('espace/agenda', ['titre' => 'Mon agenda', 'evenements' => $api->get('/mon-agenda')['donnees']], 'layouts/front');
        break;
    case '/espace/benevole':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('espace/benevole', ['titre' => 'Espace bénévole', 'creneaux' => $api->get('/creneaux')['donnees'], 'message' => null], 'layouts/front');
        break;

    case '/espace/proposer-animation':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        $creneauId = (int) ($_POST['creneau_id'] ?? 0);
        $reponse = $api->post('/creneaux/' . $creneauId . '/affectation', []);
        $msg = $reponse['statut'] === 201 ? 'Proposition enregistrée !' : ($reponse['donnees']['error'] ?? 'Proposition impossible');
        rendre('espace/benevole', ['titre' => 'Espace bénévole', 'creneaux' => $api->get('/creneaux')['donnees'], 'message' => $msg], 'layouts/front');
        break;
    case '/register':
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $api = new ApiClient($config['api_url']);
            $reponse = $api->post('/register', [
                'prenom'       => $_POST['prenom'] ?? '',
                'nom'          => $_POST['nom'] ?? '',
                'email'        => $_POST['email'] ?? '',
                'telephone'    => $_POST['telephone'] ?? '',
                'mot_de_passe' => $_POST['mot_de_passe'] ?? '',
            ]);
            if ($reponse['statut'] === 201) { header('Location: /login'); exit; }
            rendre('auth/register', ['titre' => 'Inscription', 'erreur' => ($reponse['donnees']['error'] ?? 'Inscription impossible')], 'layouts/public');
            break;
        }
        rendre('auth/register', ['titre' => 'Inscription', 'erreur' => null], 'layouts/public');
        break;

    case '/espace':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('espace/accueil', ['titre' => 'Mon espace', 'message' => null, 'evenements' => $api->get('/mon-agenda')['donnees']], 'layouts/front');
        break;

    case '/espace/services':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('espace/services', ['titre' => 'Services', 'creneaux' => $api->get('/creneaux')['donnees'], 'message' => null], 'layouts/front');
        break;

    case '/espace/inscription':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        $creneauId = (int) ($_POST['creneau_id'] ?? 0);
        $reponse = $api->post('/creneaux/' . $creneauId . '/inscription', []);
        $msg = $reponse['statut'] === 201 ? 'Inscription enregistrée !' : ($reponse['donnees']['error'] ?? 'Inscription impossible');
        rendre('espace/services', ['titre' => 'Services', 'creneaux' => $api->get('/creneaux')['donnees'], 'message' => $msg], 'layouts/front');
        break;

    case '/espace/abonnement':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        $message = null; $type_message = 'info';
        if (($_GET['paiement'] ?? '') === 'ok' && !empty($_GET['session_id'])) {
            $reponse = $api->post('/paiement/confirmer', ['session_id' => $_GET['session_id']]);
            if ($reponse['statut'] === 200) {
                $message = 'Merci ! Votre cotisation est réglée, votre compte est actif.';
                $type_message = 'success';
            } else {
                $message = $reponse['donnees']['error'] ?? 'La confirmation a échoué.';
                $type_message = 'danger';
            }
        } elseif (($_GET['paiement'] ?? '') === 'annule') {
            $message = 'Paiement annulé.';
            $type_message = 'warning';
        }
        $adhesion = $api->get('/mon-adhesion')['donnees'];
        rendre('espace/abonnement', ['titre' => 'Ma cotisation', 'message' => $message, 'type_message' => $type_message, 'adhesion' => $adhesion], 'layouts/front');
        break;

    case '/espace/paiement':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        $api = new ApiClient($config['api_url'], Session::token());
        $reponse = $api->post('/paiement/checkout', []);
        if ($reponse['statut'] === 200 && !empty($reponse['donnees']['url'])) {
            header('Location: ' . $reponse['donnees']['url']);
            exit;
        }
        rendre('espace/abonnement', ['titre' => 'Ma cotisation', 'message' => $reponse['donnees']['error'] ?? 'Paiement indisponible', 'type_message' => 'danger'], 'layouts/front');
        break;
    case '/adhesion-statut':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $api->post('/adhesions/' . (int) ($_POST['id'] ?? 0) . '/statut', ['statut' => $_POST['statut'] ?? '']);
        header('Location: /adhesions');
        exit;
    case '/logout':
        Session::deconnecter();
        header('Location: /login');
        exit;
    case '/stock':
            exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $reponse = $api->get('/stock');
        rendre('stock/index', ['titre' => 'Stock', 'stock' => $reponse['donnees']]);
        break;
    case '/benevoles':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $reponse = $api->get('/benevoles');
        rendre('benevoles/index', ['titre' => 'Bénévoles', 'benevoles' => $reponse['donnees']]);
        break;
    case '/rappels':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $reponse = $api->post('/admin/rappels', []);
        $n = $reponse['donnees']['rappels_envoyes'] ?? 0;
        rendre('accueil', ['titre' => 'Tableau de bord', 'message' => "Rappels envoyés : $n", 'stats' => $api->get('/stats')['donnees']]);
        break;
    case '/adhesions':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('adhesions/index', ['titre' => 'Adhésions', 'adhesions' => $api->get('/adhesions')['donnees']]);
        break;
    case '/espace/devenir-benevole':
        if (!Session::estConnecte()) { header('Location: /login'); exit; }
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $api = new ApiClient($config['api_url'], Session::token());
            $reponse = $api->post('/benevoles/candidature', ['disponibilites' => $_POST['disponibilites'] ?? '']);
            $msg = $reponse['statut'] === 201 ? 'Candidature envoyée ! Elle sera étudiée par l\'équipe.' : ($reponse['donnees']['error'] ?? 'Envoi impossible');
            rendre('espace/devenir-benevole', ['titre' => 'Devenir bénévole', 'message' => $msg], 'layouts/front');
            break;
        }
        rendre('espace/devenir-benevole', ['titre' => 'Devenir bénévole', 'message' => null], 'layouts/front');
        break;
    case '/benevole-statut':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $api->post('/benevoles/' . (int) ($_POST['id'] ?? 0) . '/statut', ['statut' => $_POST['statut'] ?? '']);
        header('Location: /benevoles');
        exit;
    case '/utilisateur-ban':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $api->post('/utilisateurs/' . (int) ($_POST['id'] ?? 0) . '/ban', ['actif' => (($_POST['actif'] ?? '1') === '1')]);
        header('Location: /utilisateurs');
        exit;
    case '/collectes':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('collectes/index', ['titre' => 'Collectes', 'collectes' => $api->get('/collectes')['donnees']]);
        break;

    case '/tournees':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('tournees/index', ['titre' => 'Tournées', 'tournees' => $api->get('/tournees')['donnees']]);
        break;
    case '/planning.xlsx':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $f = $api->telecharger('/planning.xlsx');
        header('Content-Type: ' . $f['type']);
        header('Content-Disposition: attachment; filename="planning.xlsx"');
        echo $f['corps'];
        exit;

    case '/tournee-pdf':
        exigerPersonnel();
        $id = (int) ($_GET['id'] ?? 0);
        $api = new ApiClient($config['api_url'], Session::token());
        $f = $api->telecharger('/tournees/' . $id . '/pdf');
        header('Content-Type: ' . $f['type']);
        header('Content-Disposition: attachment; filename="tournee.pdf"');
        echo $f['corps'];
        exit;
    case '/utilisateurs':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        rendre('utilisateurs/index', ['titre' => 'Utilisateurs', 'utilisateurs' => $api->get('/utilisateurs')['donnees']]);
        break;
        case '/produits/nouveau':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        $categories = $api->get('/categories-produits')['donnees'];
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $reponse = $api->post('/produits', [
                'code_barre' => $_POST['code_barre'] ?? '',
                'libelle'    => $_POST['libelle'] ?? '',
                'categorie'  => $_POST['categorie'] ?? '',
                'unite'      => $_POST['unite'] ?? 'piece',
            ]);
            if ($reponse['statut'] === 201) {
                header('Location: /stock');
                exit;
            }
            rendre('produits/nouveau', ['titre' => 'Nouveau produit', 'erreur' => ($reponse['donnees']['error'] ?? 'Création impossible'), 'categories' => $categories]);
            break;
        }
        rendre('produits/nouveau', ['titre' => 'Nouveau produit', 'erreur' => null, 'categories' => $categories]);
        break;
    case '/tournees/nouveau':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $reponse = $api->post('/tournees', [
                'reference'   => $_POST['reference'] ?? '',
                'date_prevue' => date('c', strtotime($_POST['date_prevue'] ?? 'now')),
                'statut'      => $_POST['statut'] ?? 'planifiee',
            ]);
            if ($reponse['statut'] === 201) { header('Location: /tournees'); exit; }
            rendre('tournees/nouveau', ['titre' => 'Nouvelle tournée', 'erreur' => ($reponse['donnees']['error'] ?? 'Création impossible')]);
            break;
        }
        rendre('tournees/nouveau', ['titre' => 'Nouvelle tournée', 'erreur' => null]);
        break;

    case '/collectes/nouveau':
        exigerPersonnel();
        $api = new ApiClient($config['api_url'], Session::token());
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $reponse = $api->post('/collectes', [
                'commercant_id' => (int) ($_POST['commercant_id'] ?? 0),
                'date_prevue'   => date('c', strtotime($_POST['date_prevue'] ?? 'now')),
                'statut'        => $_POST['statut'] ?? 'planifiee',
            ]);
            if ($reponse['statut'] === 201) { header('Location: /collectes'); exit; }
            rendre('collectes/nouveau', ['titre' => 'Nouvelle collecte', 'erreur' => ($reponse['donnees']['error'] ?? 'Création impossible')]);
            break;
        }
        rendre('collectes/nouveau', ['titre' => 'Nouvelle collecte', 'erreur' => null]);
        break;
    default:
        http_response_code(404);
        echo 'Page introuvable';
}
