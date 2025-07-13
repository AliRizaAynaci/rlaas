package auth

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv" // << added
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/AliRizaAynaci/rlaas/internal/user"
)

/* ---------------------------------------------------------------------- */
/*  Lazy-loaded Google OAuth2 config                                      */
/* ---------------------------------------------------------------------- */

var (
	googleCfg *oauth2.Config
	initOnce  sync.Once
)

// cfg() is called wherever we need the OAuth2 config.
// It makes sure the .env file is loaded and the config is built exactly once.
func cfg() *oauth2.Config {
	initOnce.Do(func() {
		_ = godotenv.Load() // silently ignore if .env is missing

		googleCfg = &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"), // e.g. http://localhost:8080/auth/google/callback
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			}, Endpoint: google.Endpoint,
		}
	})
	return googleCfg
}

/* ---------------------------------------------------------------------- */
/*  Handlers                                                              */
/* ---------------------------------------------------------------------- */

// GET /auth/google/login
func Login(c *fiber.Ctx) error {
	url := cfg().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// GET /auth/google/callback
func Callback(svc *user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing code")
		}

		ctx := context.Background()
		tok, err := cfg().Exchange(ctx, code)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "oauth exchange: "+err.Error())
		}

		client := cfg().Client(ctx, tok)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?alt=json")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "userinfo fetch: "+err.Error())
		}
		defer resp.Body.Close()

		var info struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "decode: "+err.Error())
		}

		// Find the user by GoogleID or create a new record
		u, err := svc.FindOrCreate(info.ID, info.Email, info.Name, info.Picture)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "user create: "+err.Error())
		}

		// Issue JWT
		claims := jwt.MapClaims{
			"user_id": u.ID,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "jwt sign: "+err.Error())
		}

		// Set cookie
		c.Cookie(&fiber.Cookie{
			Name:     "session_token",
			Value:    signed,
			Path:     "/",
			HTTPOnly: true,
			SameSite: "None",
			Secure:   true,
		})

		frontendRedirect := os.Getenv("FRONTEND_REDIRECT_URL")
		if frontendRedirect == "" {
			frontendRedirect = "http://localhost:3000/dashboard" // fallback
		}
		return c.Redirect(frontendRedirect, fiber.StatusSeeOther)
	}
}
