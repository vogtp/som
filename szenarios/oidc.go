package szenarios

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/pkg/stater/user"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

// OIDCSzenario check a OIDC Relying Party
type OIDCSzenario struct {
	*szenario.Base
	Issuer   string
	ClientID string
	Port     int
}

// Execute the szenario
func (s *OIDCSzenario) Execute(engine szenario.Engine) (err error) {
	if s.Port < 1 {
		s.Port = 4444
	}
	go func() {
		if err := s.startRP(context.Background(), engine); err != nil {
			fmt.Printf("starting rp: %w", err)
		}
	}()
	engine.Step("Loading",
		chromedp.Navigate(fmt.Sprintf("http://localhost:%v", s.Port)),
	)
	// document.querySelector("#idToken1_0")
	// <input id="idToken1_0" name="callback_0" type="submit" role="button" index="0" value="unibas" class="btn btn-lg btn-block btn-uppercase btn-primary">
	err = engine.StepTimeout("click unibas button", 10*time.Second,
		chromedp.WaitReady("idToken1_0", chromedp.ByID),
		chromedp.Click("#idToken1_0", chromedp.ByID),
	)
	if err != nil {
		return fmt.Errorf("cannot click unibas button %w", err)
	}
	if err := AdfsLogin(engine, s); err != nil {
		return fmt.Errorf("ADFS login: %w", err)
	}
	var body string
	engine.Step("Check if email is correct",
		
		chromedp.WaitReady("email", chromedp.ByID),
		engine.Body(engine.Strings(&body)),
	)
	if !strings.EqualFold(strings.TrimSpace(body), s.User().Email()) {
		return fmt.Errorf("mail not in ID Token: body %q", body)
	}
	return nil
}

func (s OIDCSzenario) startRP(ctx context.Context, engine szenario.Engine) error {
	oidcUserApp, err := user.Store.Get(s.ClientID)
	if err != nil {
		return fmt.Errorf("loading user for OIDC secret: %w", err)
	}
	ClientSecret := oidcUserApp.Password()

	port := s.Port
	callbackPath := "/auth/callback"
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	options := []rp.Option{
		//rp.WithSigningAlgsFromDiscovery(),
	}
	RedirectURI := fmt.Sprintf("http://localhost:%v%s", port, callbackPath)
	scopes := []string{"openid", "email"}
	//responseMode := "code token id_token"

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "OIDC test webserver") })
	relyingParty, err := rp.NewRelyingPartyOIDC(ctx, s.Issuer, s.ClientID, ClientSecret, RedirectURI, scopes, options...)
	if err != nil {
		return fmt.Errorf("creating oidc relying party: %w", err)
	}
	urlParam := []rp.URLParamOpt{
		//rp.WithResponseModeURLParam(oidc.ResponseMode(responseMode))},
	}
	http.Handle("/", rp.AuthURLHandler(
		uuid.NewString,
		relyingParty,
		urlParam...,
	))
	callbackHandler := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rpReg rp.RelyingParty, info *oidc.UserInfo) {
		c, err := rp.VerifyIDToken[*oidc.IDTokenClaims](r.Context(), tokens.IDToken, rpReg.IDTokenVerifier())
		if err != nil {
			http.Error(w, fmt.Sprintf("Verify IDToken: %v", err), http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "<html><body><div id='email'>%s</div></body></html>", c.GetUserInfo().Email)
	}

	http.DefaultServeMux.Handle(callbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(callbackHandler), relyingParty))
	go func() {
		<-ctx.Done()
		sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(sdCtx); err != nil {
			slog.Error("cannot shutdown the webserver", "err", err)
		}
	}()
	engine.Log().Info("Starting webserver for OIDC", "port", s.Port, "ClientID", s.ClientID)
	return srv.ListenAndServe()
}
