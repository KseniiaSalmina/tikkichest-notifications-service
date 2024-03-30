package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/api/models"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/api/validation"
)

// @Summary Notifications mode on
// @Tags notifications
// @Description turn on notifications for selected user
// @Accept json
// @Param user body models.User true "user information: id and telegram username"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /notifications/{userID} [post]
func (s *Server) notificationsOn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validation.User(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.storage.SaveUsername(r.Context(), user.ID, user.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Notifications mode off
// @Tags notifications
// @Description turn off notifications for selected user
// @Accept json
// @Param userID path int true "profile id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /notifications/{userID} [delete]
func (s *Server) notificationsOff(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = validation.UserID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.storage.DeleteUsername(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Change username
// @Tags notifications
// @Description change username for user with notifications mode on
// @Accept json
// @Param user body models.User true "user information: id and telegram username"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /notifications/{userID} [patch]
func (s *Server) changeUsername(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validation.User(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.storage.ChangeUsername(r.Context(), user.ID, user.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
