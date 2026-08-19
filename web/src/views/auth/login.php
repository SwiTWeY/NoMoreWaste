<h2>Connexion</h2>
<?php if ($erreur): ?>
    <p style="color:red"><?= htmlspecialchars($erreur) ?></p>
<?php endif; ?>
<form method="post" action="/login">
    <p><label>Email<br><input type="email" name="email" required></label></p>
    <p><label>Mot de passe<br><input type="password" name="mot_de_passe" required></label></p>
    <button type="submit">Se connecter</button>
</form>