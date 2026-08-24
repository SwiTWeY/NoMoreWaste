<h2 class="mb-4">Tableau de bord</h2>
<?php if (!empty($message)): ?>
    <div class="alert alert-success"><?= htmlspecialchars($message) ?></div>
<?php endif; ?>

<div class="row g-3 mb-4">
    <?php
    $cartes = [
        ['Utilisateurs', $stats['utilisateurs'] ?? 0, '/utilisateurs'],
        ['Adhésions',    $stats['adhesions'] ?? 0,    '/adhesions'],
        ['Bénévoles',    $stats['benevoles'] ?? 0,    '/benevoles'],
        ['Produits',     $stats['produits'] ?? 0,     '/stock'],
        ['Collectes',    $stats['collectes'] ?? 0,    '/collectes'],
        ['Tournées',     $stats['tournees'] ?? 0,     '/tournees'],
    ];
    foreach ($cartes as $carte): ?>
        <div class="col-6 col-md-4 col-lg-2">
            <a href="<?= $carte[2] ?>" class="text-decoration-none">
                <div class="card text-center shadow-sm h-100">
                    <div class="card-body">
                        <div class="display-6 text-success"><?= (int) $carte[1] ?></div>
                        <div class="text-muted"><?= htmlspecialchars($carte[0]) ?></div>
                    </div>
                </div>
            </a>
        </div>
    <?php endforeach; ?>
</div>

<div class="card shadow-sm">
    <div class="card-body">
        <h5 class="card-title mb-3">Actions rapides</h5>
        <form method="post" action="/rappels" class="d-inline">
            <button type="submit" class="btn btn-warning">Envoyer les rappels de renouvellement</button>
        </form>
        <a href="/planning.xlsx" class="btn btn-success">Télécharger le planning (Excel)</a>
    </div>
</div>
