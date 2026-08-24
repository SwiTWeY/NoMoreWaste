<h2>Adhésions</h2>
<table class="table table-striped table-bordered align-middle">
    <thead>
        <tr>
            <th>Adhérent</th>
            <th>Début</th>
            <th>Fin</th>
            <th>Montant</th>
            <th>Statut</th>
            <th>Changer le statut</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($adhesions as $a): ?>
            <tr>
                <td><?= htmlspecialchars($a['prenom'] . ' ' . $a['nom']) ?></td>
                <td><?= htmlspecialchars(substr($a['date_debut'], 0, 10)) ?></td>
                <td><?= htmlspecialchars(substr($a['date_fin'], 0, 10)) ?></td>
                <td><?= htmlspecialchars($a['montant']) ?> €</td>
                <td><?= htmlspecialchars($a['statut_paiement']) ?></td>
                <td>
                    <form method="post" action="/adhesion-statut" style="margin:0">
                        <input class="form-control" type="hidden" name="id" value="<?= (int) $a['id'] ?>">
                        <select class="form-select" name="statut">
                            <option value="paye">payé</option>
                            <option value="en_attente">en attente</option>
                            <option value="annule">annulé</option>
                        </select>
                        <button type="submit" class="btn btn-primary btn-sm">OK</button>
                    </form>
                </td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>
