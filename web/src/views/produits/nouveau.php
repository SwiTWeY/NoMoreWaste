<h2>Nouveau produit</h2>
<?php if ($erreur): ?>
    <p style="color:red"><?= htmlspecialchars($erreur) ?></p>
<?php endif; ?>
<form method="post" action="/produits/nouveau">
    <p><label>Code-barres<br><input class="form-control" name="code_barre" required></label></p>
    <p><label>Libellé<br><input class="form-control" name="libelle" required></label></p>
    <p><label>Catégorie<br>
        <select class="form-select" name="categorie">
            <?php foreach ($categories as $cat): ?>
                <option value="<?= htmlspecialchars($cat) ?>"><?= htmlspecialchars($cat) ?></option>
            <?php endforeach; ?>
        </select>
    </label></p>
    <p><label>Unité<br><input class="form-control" name="unite" value="piece"></label></p>
    <button type="submit" class="btn btn-primary btn-sm">Créer</button>
</form>
