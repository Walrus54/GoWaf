package web

import (
	"html/template"
	"net/http"
	"waf/internal/waf"
)

type Handler struct {
	WAF *waf.WAF
}

func (h *Handler) AdminPanel(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *Handler) GetRules(w http.ResponseWriter, r *http.Request) {
	rules := h.WAF.GetRules()
	respondJSON(w, rules)
}

func (h *Handler) AddRule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Pattern string `json:"pattern"`
	}
	if err := parseJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.WAF.AddRule(req.Name, req.Pattern); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, map[string]string{"status": "ok"})
}

func (h *Handler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.WAF.DeleteRule(id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, map[string]string{"status": "ok"})
}

func (h *Handler) UpdateMode(w http.ResponseWriter, r *http.Request) {
	mode := r.PathValue("mode")
	h.WAF.SetMode(mode)
	respondJSON(w, map[string]string{"status": "ok"})
}
