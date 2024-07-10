package server

import (
	"First-TgBot-On-GO/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketCLient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketCLient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}
	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		log.Printf("400_%d", 1)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		log.Printf("400_%d", 2)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestToken)
	if err != nil {
		log.Printf("400_%d", 3)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authResp, err := s.pocketCLient.Authorize(r.Context(), requestToken)
	if err != nil {
		log.Printf("500_%d", 1)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.tokenRepository.Save(chatID, authResp.AccessToken, repository.AccessToken)
	if err != nil {
		log.Printf("500_%d", 2)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("chat_id: %d, request_token: %s, access_token: %s", chatID, requestToken, authResp.AccessToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
