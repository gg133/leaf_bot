package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/yalagtyarzh/leaf_bot/pocket"
	"github.com/yalagtyarzh/leaf_bot/repository"
)

type AuthServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
}

func NewAuthServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectUrl string) *AuthServer {
	return &AuthServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectUrl:     redirectUrl,
	}
}

func (a *AuthServer) Start() error {
	a.server = &http.Server{
		Handler: a,
		Addr:    ":80",
	}

	return a.server.ListenAndServe()
}

func (a *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := a.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResp, err := a.pocketClient.Auth(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.tokenRepository.Save(chatID, authResp.AccessToken, repository.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("\nchat_id: %d\nrequest_token: %s\naccess_token: %s\n", chatID, requestToken, authResp.AccessToken)

	w.Header().Add("Location", a.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}
