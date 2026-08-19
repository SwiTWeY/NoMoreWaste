<?php

require __DIR__ . '/../src/lib/Session.php';
require __DIR__ . '/../src/lib/ApiClient.php';

$config = require __DIR__ . '/../config/config.php';

Session::demarrer();

function rendre(string $vue, array $vars = []): void
{
    extract($vars);
    ob_start();
    require __DIR__ . '/../src/views/' . $vue . '.php';
    $contenu = ob_get_clean();
    require __DIR__ . '/../src/views/layouts/back.php';
}

$chemin = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);

switch ($chemin) {
    case '/':
        rendre('accueil', ['titre' => 'Accueil']);
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
                header('Location: /');
                exit;
            }
            rendre('auth/login', ['titre' => 'Connexion', 'erreur' => 'Identifiants invalides']);
            break;
        }
        rendre('auth/login', ['titre' => 'Connexion', 'erreur' => null]);
        break;

    case '/logout':
        Session::deconnecter();
        header('Location: /login');
        exit;
    case '/stock':
        if (!Session::estConnecte()) {
            header('Location: /login');
            exit;
        }
        $api = new ApiClient($config['api_url'], Session::token());
        $reponse = $api->get('/stock');
        rendre('stock/index', ['titre' => 'Stock', 'stock' => $reponse['donnees']]);
        break;
    default:
        http_response_code(404);
        echo 'Page introuvable';
}
