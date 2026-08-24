<!DOCTYPE html>
<html lang="<?= htmlspecialchars(I18n::langue()) ?>">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title><?= htmlspecialchars($titre ?? 'NO MORE WASTE') ?></title>
    <link rel="stylesheet" href="/assets/css/bootstrap.min.css">
    <link rel="stylesheet" href="/assets/css/custom.css">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container-fluid">
            <a class="navbar-brand" href="/espace"><img src="/assets/img/logo.svg" alt="" height="32" class="me-2">NO MORE WASTE</a>
            <div class="navbar-nav">
                <a class="nav-link" href="/espace"><?= t('mon_espace') ?></a>
                <a class="nav-link" href="/espace/services"><?= t('services') ?></a>
                <a class="nav-link" href="/espace/agenda"><?= t('mon_agenda') ?></a>
                <?php if (Session::estBenevole()): ?>
                    <a class="nav-link" href="/espace/benevole"><?= t('espace_benevole') ?></a>
                <?php endif; ?>
                <a class="nav-link" href="/logout"><?= t('deconnexion') ?></a>
                <span class="navbar-text ms-2"><a class="text-white" href="?lang=fr">FR</a> | <a class="text-white" href="?lang=en">EN</a></span>
            </div>
        </div>
    </nav>
    <main class="container my-4">
        <?= $contenu ?>
    </main>
</body>
</html>
