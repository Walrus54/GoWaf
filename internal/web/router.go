package web

import (
	"encoding/json"
	"net/http"
	"waf/internal/waf"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	// Статические файлы
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Веб-интерфейс
	r.Get("/admin", h.AdminPanel)

	// API
	r.Route("/api", func(r chi.Router) {
		r.Get("/rules", h.GetRules)
		r.Post("/rules", h.AddRule)
		r.Delete("/rules/{id}", h.DeleteRule)
		r.Post("/mode/{mode}", h.UpdateMode)
	})

	return r
}

func StartAdminInterface(waf *waf.WAF) error {
	handler := NewHandler(waf)
	router := NewRouter(handler)
	return http.ListenAndServe(waf.Config.AdminPort, router)
}

func NewHandler(waf *waf.WAF) *Handler {
	return &Handler{WAF: waf}
}

// Вспомогательные функции
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func parseJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
