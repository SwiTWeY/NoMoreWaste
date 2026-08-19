<h2>Bienvenue</h2>
<?php if (Session::estConnecte()): ?>
    <?php $u = Session::utilisateur(); ?>
    <p>Connecté : <strong><?= htmlspecialchars($u['prenom'] . ' ' . $u['nom']) ?></strong>
       (<?= $u['est_personnel'] ? 'personnel' : 'adherent' ?>)</p>
    <p><a href="/logout">Déconnexion</a></p>
<?php else: ?>
    <p><a href="/login">Se connecter</a></p>
<?php endif; ?>