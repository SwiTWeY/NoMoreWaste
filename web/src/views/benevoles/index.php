<h2 class="mb-3">Bénévoles</h2>
<table class="table table-striped table-bordered align-middle">
    <thead>
        <tr>
            <th>Nom</th>
            <th>Email</th>
            <th>Statut</th>
            <th>Compétences</th>
            <th>Action</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($benevoles as $b): ?>
            <tr>
                <td><?= htmlspecialchars($b['prenom'] . ' ' . $b['nom']) ?></td>
                <td><?= htmlspecialchars($b['email']) ?></td>
                <td>
                    <?php if ($b['statut_candidature'] === 'valide'): ?>
                        <span class="badge bg-success">validé</span>
                    <?php elseif ($b['statut_candidature'] === 'candidat'): ?>
                        <span class="badge bg-warning text-dark">candidat</span>
                    <?php elseif ($b['statut_candidature'] === 'refuse'): ?>
                        <span class="badge bg-danger">refusé</span>
                    <?php else: ?>
                        <span class="badge bg-secondary"><?= htmlspecialchars($b['statut_candidature']) ?></span>
                    <?php endif; ?>
                </td>
                <td><?= htmlspecialchars(implode(', ', $b['competences'])) ?></td>
                <td>
                    <form method="post" action="/benevole-statut" class="d-inline">
                        <input type="hidden" name="id" value="<?= (int) $b['id'] ?>">
                        <input type="hidden" name="statut" value="valide">
                        <button type="submit" class="btn btn-success btn-sm">Valider</button>
                    </form>
                    <form method="post" action="/benevole-statut" class="d-inline">
                        <input type="hidden" name="id" value="<?= (int) $b['id'] ?>">
                        <input type="hidden" name="statut" value="refuse">
                        <button type="submit" class="btn btn-danger btn-sm">Refuser</button>
                    </form>
                </td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
