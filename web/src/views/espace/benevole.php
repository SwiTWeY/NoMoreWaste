<h2 class="mb-3"><?= t('espace_benevole') ?></h2>
<p class="text-muted"><?= t('benevole_intro') ?></p>
<?php if (!empty($message)): ?>
    <div class="alert alert-info"><?= htmlspecialchars($message) ?></div>
<?php endif; ?>
<table class="table table-striped table-bordered align-middle">
    <thead>
        <tr>
            <th><?= t('th_service') ?></th>
            <th><?= t('th_date') ?></th>
            <th><?= t('th_horaire') ?></th>
            <th><?= t('th_lieu') ?></th>
            <th></th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($creneaux as $c): ?>
            <tr>
                <td><?= htmlspecialchars($c['service_libelle']) ?></td>
                <td><?= htmlspecialchars(substr($c['date_creneau'], 0, 10)) ?></td>
                <td><?= htmlspecialchars($c['heure_debut'] . ' - ' . $c['heure_fin']) ?></td>
                <td><?= htmlspecialchars($c['lieu']) ?></td>
                <td>
                    <form method="post" action="/espace/proposer-animation" style="margin:0">
                        <input type="hidden" name="creneau_id" value="<?= (int) $c['id'] ?>">
                        <button type="submit" class="btn btn-primary btn-sm"><?= t('proposer_animer') ?></button>
                    </form>
                </td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
