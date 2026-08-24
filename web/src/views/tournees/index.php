<h2>Tournées</h2>
<table class="table table-striped table-bordered align-middle">
    <p><a class="btn btn-success mb-3" href="/tournees/nouveau">+ Nouvelle tournée</a></p>
    <thead>
        <tr>
            <th>Référence</th>
            <th>Date prévue</th>
            <th>Statut</th>
            <th>PDF</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($tournees as $t): ?>
            <tr>
                <td><?= htmlspecialchars($t['reference']) ?></td>
                <td><?= htmlspecialchars(substr($t['date_prevue'], 0, 16)) ?></td>
                <td><?= htmlspecialchars($t['statut']) ?></td>
                <td><a href="/tournee-pdf?id=<?= (int) $t['id'] ?>">PDF</a></td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
