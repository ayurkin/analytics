package http

import "net/http"

func (s *Server) ValidateAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		accessToken, err := r.Cookie("access")
		if err != nil {
			http.Error(w, "failed to obtain access token from cookie", http.StatusUnauthorized)
		}
		ctx := r.Context()
		isAuthorized, err := s.auth.Validate(ctx, accessToken.Value)
		if err != nil {
			s.logger.Errorf("validate token failed: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if !isAuthorized {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
