<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title><?= htmlspecialchars($titre ?? 'NO MORE WASTE') ?></title>
</head>
<body>
    <header>
        <h1>NO MORE WASTE — back-office</h1>
    </header>
    <main>
        <?= $contenu ?>
    </main>
</body>
</html>