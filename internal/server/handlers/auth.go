package handlers

import (
	"context"
	"encoding/json"
	"github.com/AliRizaAynaci/rlaas/internal/database"
	"github.com/AliRizaAynaci/rlaas/internal/service"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
	"time"
)

// oauthCfg reads the env vars we put in .env
var oauthCfg = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthCfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in query", http.StatusBadRequest)
		return
	}

	tok, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "OAuth exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauthCfg.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?alt=json")
	if err != nil {
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dbSvc, _ := database.GetInstance()
	user, err := service.CreateUser(ctx, dbSvc.Pool(), userInfo.ID, userInfo.Email)
	if err != nil {
		http.Error(w, "CreateUser error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Failed to sign JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // for production
	})

	http.Redirect(w, r, "http://localhost:3000/dashboard", http.StatusSeeOther)
}
