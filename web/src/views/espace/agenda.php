<h2 class="mb-3"><?= t('mon_agenda') ?></h2>
<?php if (empty($evenements)): ?>
    <div class="alert alert-info"><?= t('aucun_evenement_agenda') ?></div>
<?php else: ?>
<table class="table table-striped table-bordered align-middle">
    <thead>
        <tr>
            <th><?= t('th_type') ?></th>
            <th><?= t('th_date') ?></th>
            <th><?= t('th_heure') ?></th>
            <th><?= t('th_intitule') ?></th>
            <th><?= t('th_lieu') ?></th>
            <th><?= t('th_statut') ?></th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($evenements as $e): ?>
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
                <td><?= htmlspecialchars($e['statut']) ?></td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
<?php endif; ?>
