// internal/handler/sectionDetail.go
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"myapp/internal/model"
	"myapp/internal/service"
)

type SectionDetailHandler struct {
	service *service.SectionDetailService
}

func NewSectionDetailHandler(service *service.SectionDetailService) *SectionDetailHandler {
	return &SectionDetailHandler{
		service: service,
	}
}

func (h *SectionDetailHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sectionDetail model.SectionDetail
	if err := json.NewDecoder(r.Body).Decode(&sectionDetail); err != nil {
		respondWithErrorSD(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.Create(r.Context(), &sectionDetail); err != nil {
		respondWithErrorSD(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONSD(w, http.StatusCreated, sectionDetail)
}

func (h *SectionDetailHandler) List(w http.ResponseWriter, r *http.Request) {
	sectionDetails, err := h.service.FindAll(r.Context())
	if err != nil {
		respondWithErrorSD(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONSD(w, http.StatusOK, sectionDetails)
}

func (h *SectionDetailHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sectionDetail, err := h.service.Find(r.Context(), id)
	if err != nil {
		respondWithErrorSD(w, http.StatusNotFound, "Section Detail not found")
		return
	}

	respondWithJSONSD(w, http.StatusOK, sectionDetail)
}

func (h *SectionDetailHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var sectionDetail model.SectionDetail
	if err := json.NewDecoder(r.Body).Decode(&sectionDetail); err != nil {
		respondWithErrorSD(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	sectionDetail.ID, _ = primitive.ObjectIDFromHex(id)
	if err := h.service.Update(r.Context(), &sectionDetail); err != nil {
		respondWithErrorSD(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONSD(w, http.StatusOK, sectionDetail)
}

func (h *SectionDetailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		respondWithErrorSD(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONSD(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithErrorSD(w http.ResponseWriter, code int, message string) {
	respondWithJSONSD(w, code, map[string]string{"error": message})
}

func respondWithJSONSD(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *SectionDetailHandler) FindBySectionType(w http.ResponseWriter, r *http.Request) {
	sectionType := chi.URLParam(r, "sectionType")
	log.Println("Request received")
	sectionDetails, err := h.service.FindBySectionType(r.Context(), sectionType)
	if err != nil {
		respondWithErrorSD(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONSD(w, http.StatusOK, sectionDetails)
}
