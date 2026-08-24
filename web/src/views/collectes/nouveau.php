<h2>Nouvelle collecte</h2>
<?php if ($erreur): ?>
    <p style="color:red"><?= htmlspecialchars($erreur) ?></p>
<?php endif; ?>
<form method="post" action="/collectes/nouveau">
    <p><label>Commerçant (id)<br><input class="form-control" type="number" name="commercant_id" required></label></p>
    <p><label>Date prévue<br><input class="form-control" type="datetime-local" name="date_prevue" required></label></p>
    <p><label>Statut<br>
        <select class="form-select" name="statut">
            <option value="planifiee">planifiee</option>
            <option value="en_cours">en_cours</option>
            <option value="terminee">terminee</option>
            <option value="annulee">annulee</option>
        </select>
    </label></p>
    <button type="submit" class="btn btn-primary btn-sm">Créer</button>
</form>
