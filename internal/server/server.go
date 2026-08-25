package server

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/freeze", handleFreeze)
	mux.HandleFunc("/api/compare", handleCompare)
	mux.HandleFunc("/api/scale", handleScale)
	mux.HandleFunc("/api/riser", handleRiser)
	mux.HandleFunc("/api/validate", handleValidate)
	mux.HandleFunc("/api/rules", handleRules)
	mux.HandleFunc("/api/modulus", handleModulus)
	mux.HandleFunc("/api/superheat", handleSuperheat)
	mux.HandleFunc("/api/info", handleInfo)
	mux.HandleFunc("/api/endpoints", handleEndpoints)
	return mux
}

func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, New())
}
