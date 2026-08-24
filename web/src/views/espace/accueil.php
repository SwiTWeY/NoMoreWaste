<?php $u = Session::utilisateur(); ?>
<h2 class="mb-3"><?= t('bonjour') ?> <?= htmlspecialchars($u['prenom']) ?> 👋</h2>
<?php if (!empty($message)): ?>
    <div class="alert alert-success"><?= htmlspecialchars($message) ?></div>
<?php endif; ?>

<div class="row g-3 mb-4">
    <div class="col-md-4">
        <a href="/espace/services" class="text-decoration-none">
            <div class="card shadow-sm h-100"><div class="card-body text-center">
                <h5 class="card-title"><?= t('services') ?></h5>
                <p class="text-muted mb-0"><?= t('services_desc') ?></p>
            </div></div>
        </a>
    </div>
    <div class="col-md-4">
        <a href="/espace/agenda" class="text-decoration-none">
            <div class="card shadow-sm h-100"><div class="card-body text-center">
                <h5 class="card-title"><?= t('mon_agenda') ?></h5>
                <p class="text-muted mb-0"><?= t('agenda_desc') ?></p>
            </div></div>
        </a>
    </div>
    <div class="col-md-4">
        <a href="/espace/devenir-benevole" class="text-decoration-none">
            <div class="card shadow-sm h-100"><div class="card-body text-center">
                <h5 class="card-title"><?= t('devenir_benevole') ?></h5>
                <p class="text-muted mb-0"><?= t('benevole_desc') ?></p>
            </div></div>
        </a>
    </div>
</div>

<h4 class="mb-3"><?= t('prochains_evenements') ?></h4>
<?php if (empty($evenements)): ?>
    <div class="alert alert-info"><?= t('aucun_evenement') ?> <a href="/espace/services"><?= t('inscrivez_service') ?></a>.</div>
<?php else: ?>
<table class="table table-striped table-bordered align-middle">
    <thead>
        <tr><th><?= t('th_type') ?></th><th><?= t('th_date') ?></th><th><?= t('th_heure') ?></th><th><?= t('th_intitule') ?></th><th><?= t('th_lieu') ?></th></tr>
    </thead>
    <tbody>
        <?php foreach (array_slice($evenements, 0, 5) as $e): ?>
            <tr>
                <td>
                    <?php if ($e['type'] === 'service'): ?>
                        <span class="badge bg-primary"><?= t('type_service') ?></span>
                    <?php else: ?>
                        <span class="badge bg-warning text-dark"><?= t('type_collecte') ?></span>
                    <?php endif; ?>
                </td>
                <td><?= htmlspecialchars(substr($e['date'], 0, 10)) ?></td>
                <td><?= htmlspecialchars($e['heure']) ?></td>
                <td><?= htmlspecialchars($e['libelle']) ?></td>
                <td><?= htmlspecialchars($e['lieu']) ?></td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
<p><a href="/espace/agenda"><?= t('voir_tout_agenda') ?></a></p>
<?php endif; ?>
