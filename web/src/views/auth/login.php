<h2><?= t('connexion') ?></h2>
<?php if ($erreur): ?>
    <div class="alert alert-danger"><?= htmlspecialchars($erreur) ?></div>
<?php endif; ?>
<form method="post" action="/login" class="col-md-6">
    <div class="mb-3"><label class="form-label"><?= t('email') ?></label><input class="form-control" type="email" name="email" required></div>
    <div class="mb-3"><label class="form-label"><?= t('mot_de_passe') ?></label><input class="form-control" type="password" name="mot_de_passe" required></div>
    <button type="submit" class="btn btn-primary"><?= t('se_connecter') ?></button>
</form>
