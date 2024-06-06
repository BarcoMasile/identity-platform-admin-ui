// Copyright 2024 Canonical Ltd.
// SPDX-License-Identifier: AGPL-3.0

package authentication

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"

	"github.com/canonical/identity-platform-admin-ui/internal/logging"
	"github.com/canonical/identity-platform-admin-ui/internal/monitoring"
	"github.com/canonical/identity-platform-admin-ui/internal/tracing"
)

type JWKSTokenVerifier struct {
	verifier providerVerifierInterface
	logger   logging.LoggerInterface
	tracer   tracing.TracingInterface
	monitor  monitoring.MonitorInterface
}

func verifyJWT(ctx context.Context, rawJwt string, verifier providerVerifierInterface) (*Principal, error) {
	i, err := verifier.Verify(ctx, rawJwt)
	if err != nil {
		return nil, err
	}

	return NewPrincipalFromClaims(i)
}

func (j *JWKSTokenVerifier) VerifyAccessToken(ctx context.Context, rawAccessToken string) (*Principal, error) {
	_, span := j.tracer.Start(ctx, "authentication.JWKSTokenVerifier.VerifyAccessToken")
	defer span.End()

	return verifyJWT(ctx, rawAccessToken, j.verifier)
}

func (j *JWKSTokenVerifier) VerifyIDToken(ctx context.Context, rawIDToken string) (*Principal, error) {
	_, span := j.tracer.Start(ctx, "authentication.JWKSTokenVerifier.VerifyIDToken")
	defer span.End()

	return verifyJWT(ctx, rawIDToken, j.verifier)
}

func NewJWKSTokenVerifier(provider ProviderInterface, clientID string, tracer trace.Tracer, logger logging.LoggerInterface, monitor monitoring.MonitorInterface) *JWKSTokenVerifier {
	j := new(JWKSTokenVerifier)
	j.tracer = tracer
	j.logger = logger
	j.monitor = monitor

	j.verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	return j
}

type UserinfoTokenVerifier struct {
	clientID string
	provider ProviderInterface
	verifier providerVerifierInterface

	logger  logging.LoggerInterface
	tracer  tracing.TracingInterface
	monitor monitoring.MonitorInterface
}

type claims struct {
	Audience []string `json:"aud"`
}

func (u *UserinfoTokenVerifier) VerifyAccessToken(ctx context.Context, rawAccessToken string) (*Principal, error) {
	_, span := u.tracer.Start(ctx, "authentication.UserinfoTokenVerifier.VerifyAccessToken")
	defer span.End()

	info, err := u.provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: rawAccessToken}))
	if err != nil {
		return nil, err
	}

	err = u.validateAdditionalClaims(info)
	if err != nil {
		return nil, err
	}

	return NewPrincipalFromClaims(info)
}

func (u *UserinfoTokenVerifier) validateAdditionalClaims(userinfo ReadableClaims) error {
	claimsToCheck := claims{}

	err := userinfo.Claims(&claimsToCheck)
	if err != nil {
		return err
	}

	for _, audience := range claimsToCheck.Audience {
		if audience == u.clientID {
			return nil
		}
	}

	return fmt.Errorf("access token audiece doesn't match expected value")
}

func (u *UserinfoTokenVerifier) VerifyIDToken(ctx context.Context, rawIDToken string) (*Principal, error) {
	_, span := u.tracer.Start(ctx, "authentication.UserinfoTokenVerifier.VerifyIDToken")
	defer span.End()

	return verifyJWT(ctx, rawIDToken, u.verifier)
}

func NewUserinfoTokenVerifier(provider ProviderInterface, clientID string, tracer trace.Tracer, logger logging.LoggerInterface, monitor monitoring.MonitorInterface) *UserinfoTokenVerifier {
	u := new(UserinfoTokenVerifier)
	u.tracer = tracer
	u.logger = logger
	u.monitor = monitor

	u.clientID = clientID
	u.provider = provider
	u.verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	return u
}
