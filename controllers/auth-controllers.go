package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"trining/models"
	"trining/utils"

	"github.com/google/uuid"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		ID:       uuid.New(),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if user.Password == "" {
		utils.HandleError(w, http.StatusBadRequest, "Password is required")
		return
	}

	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		utils.HandleError(w, http.StatusBadRequest, "Invalid file")
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := utils.SaveImageFile(file, "users", fileHeader.Filename)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Error saving image")
		}
		user.Img = &imageName
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user.Password = hashedPassword

	query, args, err := QB.
		Insert("users").
		Columns("id", "img", "name", "phone", "email", "password").
		Values(user.ID, user.Img, user.Name, user.Phone, user.Email, user.Password).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error generate query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&user); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, user)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		utils.HandleError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	var user models.User
	query := "SELECT id, img, name, phone, email, password FROM users WHERE email=$1"
	if err := db.Get(&user, query, email); err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		utils.HandleError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"img":   user.Img,
	})
}
