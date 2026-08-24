<h2 class="mb-3">Utilisateurs</h2>
<table class="table table-striped table-bordered align-middle">
    <thead>
        <tr>
            <th>Nom</th>
            <th>Email</th>
            <th>Téléphone</th>
            <th>Rôle</th>
            <th>Statut</th>
            <th>Action</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($utilisateurs as $u): ?>
            <tr>
                <td><?= htmlspecialchars($u['prenom'] . ' ' . $u['nom']) ?></td>
                <td><?= htmlspecialchars($u['email']) ?></td>
                <td><?= htmlspecialchars($u['telephone']) ?></td>
                <td><?= $u['est_personnel'] ? 'Personnel' : 'Adhérent' ?></td>
                <td>
                    <?php if ($u['actif']): ?>
                        <span class="badge bg-success">actif</span>
                    <?php else: ?>
                        <span class="badge bg-danger">banni</span>
                    <?php endif; ?>
                </td>
                <td>
                    <form method="post" action="/utilisateur-ban" class="d-inline">
                        <input type="hidden" name="id" value="<?= (int) $u['id'] ?>">
                        <input type="hidden" name="actif" value="<?= $u['actif'] ? '0' : '1' ?>">
                        <button type="submit" class="btn btn-sm <?= $u['actif'] ? 'btn-danger' : 'btn-success' ?>">
                            <?= $u['actif'] ? 'Bannir' : 'Réactiver' ?>
                        </button>
                    </form>
                </td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
