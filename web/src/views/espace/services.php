<h2 class="mb-3"><?= t('services_disponibles') ?></h2>
<?php if (!empty($message)): ?>
    <div class="alert alert-success"><?= htmlspecialchars($message) ?></div>
<?php endif; ?>

<input type="text" id="recherche" class="form-control mb-3" placeholder="<?= htmlspecialchars(t('recherche_placeholder')) ?>">

<table class="table table-striped table-bordered align-middle" id="tableServices">
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
                    <form method="post" action="/espace/inscription" style="margin:0">
                        <input type="hidden" name="creneau_id" value="<?= (int) $c['id'] ?>">
                        <button type="submit" class="btn btn-primary btn-sm"><?= t('s_inscrire') ?></button>
                    </form>
                </td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>

<script>
document.getElementById('recherche').addEventListener('input', function () {
    var q = this.value.toLowerCase();
    document.querySelectorAll('#tableServices tbody tr').forEach(function (ligne) {
        ligne.style.display = ligne.textContent.toLowerCase().includes(q) ? '' : 'none';
    });
});
</script>
