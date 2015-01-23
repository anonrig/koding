package main

import (
	"encoding/json"
	"socialapi/workers/payment/paymentemail"
	"socialapi/workers/payment/paymentwebhook/webhookmodels"
	"socialapi/workers/payment/stripe"
)

type subscriptionActionType func(*webhookmodels.StripeSubscription) error

func stripeSubscriptionCreated(raw []byte) error {
	actions := []subscriptionActionType{
		sendSubscriptionCreatedEmail,
	}

	return _stripeSubscription(raw, actions)
}

func stripeSubscriptionDeleted(raw []byte) error {
	actions := []subscriptionActionType{
		stripe.SubscriptionDeletedWebhook,
		sendSubscriptionDeletedEmail,
	}

	return _stripeSubscription(raw, actions)
}

func _stripeSubscription(raw []byte, actions []subscriptionActionType) error {
	var req *webhookmodels.StripeSubscription

	err := json.Unmarshal(raw, &req)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := action(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendSubscriptionCreatedEmail(req *webhookmodels.StripeSubscription) error {
	email, err := getEmailForCustomer(req.CustomerId)
	if err != nil {
		return err
	}

	opts := &paymentemail.Options{
		PlanName: req.Plan.Name,
	}

	return paymentemail.Send(paymentemail.SubscriptionCreated, email, opts)
}

func sendSubscriptionDeletedEmail(req *webhookmodels.StripeSubscription) error {
	email, err := getEmailForCustomer(req.CustomerId)
	if err != nil {
		return err
	}

	opts := &paymentemail.Options{
		PlanName: req.Plan.Name,
	}

	return paymentemail.Send(paymentemail.SubscriptionDeleted, email, opts)
}

//----------------------------------------------------------
// TODO: move to stripe package
//----------------------------------------------------------

type stripeCard struct {
	Id      string `json:"id"`
	ExpYear string `json:"exp_year"`
	Last4   string `json:"last4"`
	Brand   string `json:"brand"`
}

type stripeChargeRefundWebhookReq struct {
	Card       *stripeCard `json:"card"`
	Currency   string      `json:"currency"`
	Amount     float64     `json:"amount"`
	CustomerId string      `json:"customer"`
}
