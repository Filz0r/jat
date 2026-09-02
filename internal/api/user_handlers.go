package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// @Summary Create user
// @Description Creates the first user. If no admin exists and the service is not initialized, the new user becomes the first admin.
// @Tags users
// @Accept json
// @Produce json
// @Param request body userCreateRequest true "User creation payload"
// @Success 201 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /users [post]
func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := userCreateRequest{}
		err := decoder.Decode(&user)
		if err != nil {
			s.respondWithError(w, 400, "error parsing json", err)
			return
		}
		if user.Username == "" || user.Password == "" || user.Email == "" {
			s.respondWithError(w, 400, "invalid payload", nil)
			return
		}
		dbUser, err := s.services.CreateUser(
			user.Username,
			user.Email,
			user.Password,
		)
		if err != nil {
			s.respondWithError(w, 400, "error creating user", err)
			return
		}
		err = s.services.CreateInitialApplicationStatus(dbUser.ID)
		if err != nil {
			s.respondWithError(w, 400, "error creating initial application status", err)
			return
		}
		// Create first admin if it doesn't exist AND the server isn't at a initialized state
		if s.services.GetFirstAdmin() == uuid.Nil && !s.services.IsInitialized() {
			set := s.services.SetFirstAdmin(dbUser.ID)
			if !set {
				s.respondWithError(w, 400, "error setting first admin", nil)
				return
			}
			err := s.services.ChangeUserAdminStatus(dbUser.ID, true)
			if err != nil {
				s.respondWithError(w, 400, "error changing user admin status", err)
				return
			}
		}
		s.respondWithJSON(w, 201, apiResponse{
			Ok:      true,
			Message: "User created",
		})
	}
}

// @Summary List all users
// @Description Admin-only endpoint that returns every user.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apiResponse{data=[]adminUserResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /users [get]
func (s *Server) handleGetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := userAdminFromContext(r.Context())
		if !ok || !isAdmin {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		users, err := s.services.GetAllUsers()
		if err != nil {
			s.respondWithError(w, 400, "could not fetch users", err)
			return
		}
		response := make([]adminUserResponse, 0, len(users))
		for _, user := range users {
			temp := newAdminUserResponse(user)
			response = append(response, temp)
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    response,
			Message: fmt.Sprintf("Found %d users", len(response)),
		})

	}
}

// @Summary	Get current user
// @Description	Returns the currently authenticated user. This is the only user endpoint that web clients can call on page load without knowing their own UUID.
// @Tags users
// @Accept json
// @Produce	json
// @Security BearerAuth
// @Success	200	{object} apiResponse{data=meUserResponse}
// @Failure	401	{object} apiResponse
// @Router /users/me [get]
func (s *Server) handleGetCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		user, err := s.services.GetUserByID(userID)
		if err != nil {
			s.respondWithError(w, 404, "user not found", err)
			return
		}
		response := newMeUserResponse(user)
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "User found",
			Data:    response,
		})
	}
}

// @Summary	Get a user
// @Description	Returns a single user. Users can read their own record; admins can read any record.
// @Tags users
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userID path string true "User UUID"
// @Success 200 {object} apiResponse{data=adminUserResponse}
// @Failure 401 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /users/{userID} [get]
func (s *Server) handleGetSingleUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		param := r.PathValue("userID")
		paramUUID, err := uuid.Parse(param)
		if err != nil {
			s.respondWithError(w, 404, "user not found", nil)
			return
		}
		// only admins can check other users
		isAdmin := s.services.IsUserAdmin(userID)
		if paramUUID != userID && !isAdmin {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		user, err := s.services.GetUserByID(paramUUID)
		if err != nil {
			s.respondWithError(w, 404, "user not found", err)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "User found",
			Data:    newAdminUserResponse(user),
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Update current user
// @Description Updates the authenticated user's profile.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userCreateRequest true "Updated user fields"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /users [put]
func (s *Server) handleUserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := userCreateRequest{}
		err := decoder.Decode(&user)
		if err != nil {
			s.respondWithError(w, 400, "invalid body", err)
			return
		}
		userID, _ := userIDFromContext(r.Context())
		_, err = s.services.UpdateUser(
			userID,
			user.Username,
			user.Password,
			user.Email,
		)
		if err != nil {
			s.respondWithError(w, 400, "could not update user", err)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "User updated",
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Toggle admin status
// @Description Admin-only endpoint to grant (GET) or revoke (DELETE) admin privileges for a user.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userID path string true "User UUID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /admin/users/{userID} [get]
func (s *Server) handleMakeUserAdmin(give bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("userID")
		paramUUID, err := uuid.Parse(param)
		if err != nil {
			s.respondWithError(w, 404, "user not found", nil)
			return
		}
		err = s.services.ChangeUserAdminStatus(paramUUID, give)
		if err != nil {
			s.respondWithError(w, 400, "could not change admin status for this user", err)
			return
		}
		message := "User with Id: " + param + " is now admin"
		if !give {
			message = "User with Id: " + param + " is no longer admin"
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: message,
		})
	}
}
