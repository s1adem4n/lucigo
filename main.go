package main

import (
	"context"
	"fmt"
	"lucigo/pkg/auth"
	"lucigo/pkg/db"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

func setAuthContextMiddleware(a *auth.Auth, queries *db.Queries, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, ok := s.Values["token"].(string)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		session, err := a.GetSession(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := queries.GetUser(r.Context(), session.GetUserID())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func main() {
	godotenv.Load()

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	conn, err := db.OpenDB()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	queries := db.New(conn)
	authDB := db.NewAuthDatabase(queries)

	providerMap := make(auth.ProviderMap)
	providerMap["github"] = auth.NewGithubOAuth2Provider(
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		"http://localhost:8080/oauth2/callback/github",
	)

	a := auth.NewAuth(authDB)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /oauth2/login/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		p, exists := providerMap.Get(provider)
		if !exists {
			http.Error(w, "provider not found", http.StatusNotFound)
			return
		}

		authURL := p.GetAuthURL()
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	mux.HandleFunc("GET /oauth2/callback/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider, exists := providerMap.Get(r.PathValue("provider"))
		if !exists {
			http.Error(w, "provider not found", http.StatusNotFound)
			return
		}

		token, err := auth.GenerateSessionToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = a.RegisterOrLoginOAuth2(r.Context(), token, provider, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Values["token"] = token
		if err := s.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*db.User)

		name := "stranger"
		if ok {
			name = user.ID
		}
		fmt.Fprintf(w, "Hello, %s!", name)
	})
	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, ok := s.Values["token"].(string)
		if !ok {
			http.Error(w, "no token found", http.StatusNotFound)
			return
		}

		if err := a.DeleteSession(r.Context(), token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Values["token"] = ""
		if err := s.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe(":8080", setAuthContextMiddleware(a, queries, mux))
}
