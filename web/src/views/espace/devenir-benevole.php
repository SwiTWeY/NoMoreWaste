<h2 class="mb-3"><?= t('devenir_benevole') ?></h2>
<?php if (!empty($message)): ?>
    <div class="alert alert-success"><?= htmlspecialchars($message) ?></div>
<?php endif; ?>
<p><?= t('candidature_intro') ?></p>
<form method="post" action="/espace/devenir-benevole" class="col-md-8">
    <div class="mb-3">
        <label class="form-label"><?= t('candidature_label') ?></label>
        <textarea class="form-control" name="disponibilites" rows="4" required></textarea>
    </div>
    <button type="submit" class="btn btn-primary"><?= t('envoyer_candidature') ?></button>
</form>
