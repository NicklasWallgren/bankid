package bankid

import (
	"context"
	"fmt"

	"github.com/NicklasWallgren/bankid/configuration"
	"gopkg.in/go-playground/validator.v9"
)

// BankID contains the validator and configuration context.
type BankID struct {
	validator     *validator.Validate
	configuration *configuration.Configuration
	client        Client
}

// New returns a new instance of 'BankID'.
func New(configuration *configuration.Configuration) *BankID {
	return &BankID{validator: newValidator(), configuration: configuration}
}

// Authenticate - Initiates an authentication order.
//
// Use the collect method to query the status of the order.
// If the request is successful, the orderRef and autoStartToken is returned.
func (b BankID) Authenticate(context context.Context, payload *AuthenticationPayload) (*AuthenticateResponse, error) {
	request := newAuthenticationRequest(payload)
	response, err := b.call(context, request)
	if err != nil {
		return nil, err
	}

	authenticateResponse := (response).(*AuthenticateResponse)
	return authenticateResponse, nil
}

// Sign - Initiates an sign order.
//
// Use the collect method to query the status of the order.
// If the request is successful, the orderRef and autoStartToken is returned.
func (b BankID) Sign(context context.Context, payload *SignPayload) (*SignResponse, error) {
	request := newSignRequest(payload)
	response, err := b.call(context, request)
	if err != nil {
		return nil, err
	}

	signResponse := (response).(*SignResponse)
	return signResponse, nil
}

// Collect - Collects the result of a sign or auth order suing the orderRef as reference.
//
// RP should keep calling collect every two seconds as long as status indicates pending.
// RP must abort if status indicates failed. The user identity is returned when complete.
func (b BankID) Collect(context context.Context, payload *CollectPayload) (*CollectResponse, error) {
	request := newCollectRequest(payload)
	response, err := b.call(context, request)
	if err != nil {
		return nil, err
	}

	collectResponse := (response).(*CollectResponse)
	return collectResponse, nil
}

// Cancel - Cancels an ongoing sign or auth order.
//
// This is typically used if the user cancels the order in your service or app.
func (b BankID) Cancel(context context.Context, payload *CancelPayload) (*CancelResponse, error) {
	request := newCancelRequest(payload)
	response, err := b.call(context, request)
	if err != nil {
		return nil, err
	}

	cancelResponse := (response).(*CancelResponse)
	return cancelResponse, nil
}

// call validates the prerequisites of the requests and invokes the REST API method.
func (b *BankID) call(context context.Context, request Request) (Response, error) {
	// Validate the integrity of the Payload
	if err := b.validator.Struct(request.Payload()); err != nil {
		return nil, fmt.Errorf("payload validation error %w", err)
	}

	if err := b.initialize(); err != nil {
		return nil, err
	}

	return b.client.call(context, request, b)
}

// initialize prepares the client in head of a request.
func (b *BankID) initialize() error {
	// Check whether the client has been initialized
	if b.client != nil {
		return nil
	}

	// Lazy initialization
	client, err := newClient(b.configuration)
	if err != nil {
		return err
	}

	b.client = client

	return nil
}
