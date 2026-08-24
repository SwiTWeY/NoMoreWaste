<h2>Collectes</h2>
<table class="table table-striped table-bordered align-middle">
    <p><a class="btn btn-success mb-3" href="/collectes/nouveau">Nouvelle collectes</a></p>
    <thead>
        <tr>
            <th>ID</th>
            <th>Source</th>
            <th>Date prévue</th>
            <th>Statut</th>
            <th>Adresse</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($collectes as $c): ?>
            <tr>
                <td><?= (int) $c['id'] ?></td>
                <td>
                    <?php if (!empty($c['commercant_id'])): ?>
                        Commerçant #<?= (int) $c['commercant_id'] ?>
                    <?php else: ?>
                        Donateur #<?= (int) $c['donateur_id'] ?>
                    <?php endif; ?>
                </td>
                <td><?= htmlspecialchars(substr($c['date_prevue'], 0, 16)) ?></td>
                <td><?= htmlspecialchars($c['statut']) ?></td>
                <td><?= htmlspecialchars($c['adresse_collecte']) ?></td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
