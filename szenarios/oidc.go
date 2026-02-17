package szenarios

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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
	Issuer    string
	ClientID  string
	Port      int
	rpHttpSrv *http.Server
}

// Execute the szenario
func (s *OIDCSzenario) Execute(engine szenario.Engine) (err error) {
	_, stop, err := s.InitOIDCRelyingParty(engine)
	defer stop()
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
		chromedp.WaitReady("displayName", chromedp.ByID),
		engine.Body(engine.Strings(&body)),
	)
	//engine.WaitForEver()
	// if !strings.EqualFold(strings.TrimSpace(body), s.User().Email()) {
	// 	return fmt.Errorf("mail not in ID Token: body %q", body)
	// }
	return nil
}

func (s *OIDCSzenario) InitOIDCRelyingParty(engine szenario.Engine) (context.Context, context.CancelFunc, error) {
	ctx, stop := context.WithTimeout(context.Background(), 1*time.Minute)

	if err := s.initRelyingPartyHttpsrv(ctx); err != nil {
		return ctx, stop, fmt.Errorf("init relying party http servce: %w", err)
	}
	go s.startRelyingPartyHttpsrv(engine)
	return ctx, stop, nil
}

func (s *OIDCSzenario) initRelyingPartyHttpsrv(ctx context.Context) error {
	oidcUserApp, err := user.Store.Get(s.ClientID)
	if err != nil {
		return fmt.Errorf("loading user for OIDC secret: %w", err)
	}
	ClientSecret := oidcUserApp.Password()
	mux := http.NewServeMux()
	callbackPath := "/auth/callback"
	s.rpHttpSrv = &http.Server{
		Addr:              fmt.Sprintf(":%v", s.Port),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Handler:           mux,
	}
	options := []rp.Option{
		//rp.WithSigningAlgsFromDiscovery(),
	}
	RedirectURI := fmt.Sprintf("http://localhost:%v%s", s.Port, callbackPath)
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
	mux.Handle("/", rp.AuthURLHandler(
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

	mux.Handle(callbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(callbackHandler), relyingParty))
	go func() {
		<-ctx.Done()
		sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.rpHttpSrv.Shutdown(sdCtx); err != nil {
			slog.Error("cannot shutdown the webserver", "err", err)
		}
	}()
	return nil
}

func (s OIDCSzenario) startRelyingPartyHttpsrv(engine szenario.Engine) {
	engine.Log().Info("Starting webserver for OIDC", "port", s.Port, "ClientID", s.ClientID)
	if err := s.rpHttpSrv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return
		}
		engine.Log().Warn("Could not close OIDC http server", "err", err)
	}
}
