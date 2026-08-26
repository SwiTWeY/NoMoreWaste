package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"

	"github.com/SwiTWeY/NoMoreWaste/api/bdd"
	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/middleware"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type PaiementHandler struct {
	DB  *sql.DB
	Cfg config.Config
}

const montantCotisationCents = 1200

func (h PaiementHandler) CreerCheckout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsDepuis(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "non authentifie")
		return
	}

	stripe.Key = h.Cfg.StripeSecretKey

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Cotisation annuelle NO MORE WASTE"),
					},
					UnitAmount: stripe.Int64(montantCotisationCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(h.Cfg.AppBaseURL + "/espace/abonnement?paiement=ok&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(h.Cfg.AppBaseURL + "/espace/abonnement?paiement=annule"),
	}
	params.AddMetadata("utilisateur_id", strconv.Itoa(claims.UtilisateurID))

	s, err := session.New(params)
	if err != nil {
		utils.Error(w, http.StatusBadGateway, "stripe: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"url": s.URL})
}

func (h PaiementHandler) ConfirmerPaiement(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsDepuis(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "non authentifie")
		return
	}

	var corps struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&corps); err != nil || corps.SessionID == "" {
		utils.Error(w, http.StatusBadRequest, "session_id manquant")
		return
	}

	stripe.Key = h.Cfg.StripeSecretKey

	s, err := session.Get(corps.SessionID, nil)
	if err != nil {
		utils.Error(w, http.StatusBadGateway, "stripe: "+err.Error())
		return
	}

	if s.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		utils.Error(w, http.StatusPaymentRequired, "paiement non effectue")
		return
	}
	if s.Metadata["utilisateur_id"] != strconv.Itoa(claims.UtilisateurID) {
		utils.Error(w, http.StatusForbidden, "session d'un autre utilisateur")
		return
	}

	montant := float64(s.AmountTotal) / 100
	if err := bdd.CreerAdhesionPayee(h.DB, claims.UtilisateurID, montant, s.ID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "creation adhesion")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"statut": "adhesion activee"})
}
