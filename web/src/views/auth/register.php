<h2 class="mb-3"><?= t('inscription') ?></h2>
<?php if ($erreur): ?>
    <div class="alert alert-danger"><?= htmlspecialchars($erreur) ?></div>
<?php endif; ?>
<form method="post" action="/register" class="col-md-6">
    <div class="mb-3"><label class="form-label"><?= t('prenom') ?></label><input class="form-control" name="prenom" required></div>
    <div class="mb-3"><label class="form-label"><?= t('nom') ?></label><input class="form-control" name="nom" required></div>
    <div class="mb-3"><label class="form-label"><?= t('email') ?></label><input class="form-control" type="email" name="email" required></div>
    <div class="mb-3"><label class="form-label"><?= t('telephone') ?></label><input class="form-control" name="telephone"></div>
    <div class="mb-3"><label class="form-label"><?= t('mot_de_passe') ?></label><input class="form-control" type="password" name="mot_de_passe" required></div>
    <button type="submit" class="btn btn-success"><?= t('creer_compte') ?></button>
</form>
