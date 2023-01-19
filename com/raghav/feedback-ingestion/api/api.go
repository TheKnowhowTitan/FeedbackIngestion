package api

import (
	"encoding/json"
	"net/http"

	"com/raghav/feedback-ingestion/feedback"
	"com/raghav/feedback-ingestion/poller"
	"com/raghav/feedback-ingestion/service"
)

type API struct {
	feedbackStore feedback.Store
	poller        *poller.Poller
	serviceReg    *service.Registry
}

func NewAPI(feedbackStore feedback.Store, poller *poller.Poller, serviceReg *service.Registry) *API {
	return &API{
		feedbackStore: feedbackStore,
		poller:        poller,
		serviceReg:    serviceReg,
	}
}

func (a *API) GetFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	feedback, err := a.feedbackStore.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(feedback)
}

func (a *API) AddFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback feedback.Feedback
	json.NewDecoder(r.Body).Decode(&feedback)

	err := a.feedbackStore.Add(feedback)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) RegisterFeedbackSource(w http.ResponseWriter, r *http.Request) {
	var req service.RegistrationRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := a.serviceReg.Add(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) StartServer() {
	http.HandleFunc("/feedback", a.GetFeedback)
	http.HandleFunc("/feedback", a.AddFeedback)
	http.HandleFunc("/register", a.RegisterFeedbackSource)

	http.ListenAndServe(":8080", nil)
}
