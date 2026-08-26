<?php $actif = !empty($adhesion['actif']); ?>
<h2 class="mb-3">Ma cotisation</h2>

<?php if (!empty($message)): ?>
    <div class="alert alert-<?= htmlspecialchars($type_message ?? 'info') ?>"><?= htmlspecialchars($message) ?></div>
<?php endif; ?>

<?php if ($actif): ?>
    <div class="card shadow-sm col-md-7 border-success">
        <div class="card-body">
            <h5 class="card-title text-success">✅ Cotisation à jour</h5>
            <p class="card-text">Votre adhésion est valable jusqu'au
                <strong><?= htmlspecialchars(date('d/m/Y', strtotime($adhesion['date_fin']))) ?></strong>.
                Vous avez accès à tous les services de l'association.</p>
            <form method="post" action="/espace/paiement" style="margin:0">
                <button type="submit" class="btn btn-outline-success">Renouveler ma cotisation</button>
            </form>
        </div>
    </div>
<?php else: ?>
    <div class="card shadow-sm col-md-7">
        <div class="card-body">
            <h5 class="card-title">Cotisation annuelle</h5>
            <p class="card-text">Votre cotisation d'adhérent vous donne accès à tous les services de l'association (cours, ateliers, échanges de services…) pendant 1 an.</p>
            <p class="display-6 text-success mb-3">12 €<span class="fs-6 text-muted"> / an</span></p>
            <form method="post" action="/espace/paiement" style="margin:0">
                <button type="submit" class="btn btn-success btn-lg">Payer ma cotisation</button>
            </form>
            <p class="text-muted small mt-3 mb-0">Paiement sécurisé via Stripe. Carte de test : <code>4242 4242 4242 4242</code>, date future, CVC au hasard.</p>
        </div>
    </div>
<?php endif; ?>
