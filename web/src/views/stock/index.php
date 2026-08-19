<h2>Stock</h2>
<table border="1" cellpadding="6">
    <thead>
        <tr>
            <th>Code-barres</th>
            <th>Libellé</th>
            <th>Catégorie</th>
            <th>Quantité</th>
            <th>Date Limite</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($stock as $item): ?>
            <tr>
                <td><?= htmlspecialchars($item['code_barre']) ?></td>
                <td><?= htmlspecialchars($item['libelle']) ?></td>
                <td><?= htmlspecialchars($item['categorie']) ?></td>
                <td><?= (int) $item['quantite_stock'] ?></td>
                <td><?= htmlspecialchars($item['date_limite'] ?? '') ?></td>
            </tr>
        <?php endforeach; ?>
    </tbody>
</table>