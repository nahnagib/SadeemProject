package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"trining/models"
	"trining/utils"

	"github.com/Masterminds/squirrel"
)

// Vendor Columns
var (
	Vendor_columns = []string{
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
)

func IndexVendorHandler(w http.ResponseWriter, r *http.Request) {
	var vendors []models.Vendor

	query, args, err := QB.Select(strings.Join(Vendor_columns, ", ")).
		From("vendors").
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Select(&vendors, query, args...); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, vendors)
}

func ShowVendorHandler(w http.ResponseWriter, r *http.Request) {
	var vendors models.Vendor
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(Vendor_columns, ", ")).
		From("vendors").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Get(&vendors, query, args...); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, vendors)
}

func StoreVendorHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request body to get vendor details
	var newVendor models.Vendor
	if err := json.NewDecoder(r.Body).Decode(&newVendor); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the input data (you can customize this as needed)
	if newVendor.Name == "" || newVendor.Description == "" {
		utils.HandleError(w, http.StatusBadRequest, "Vendor name and description are required")
		return
	}

	// Insert the new vendor into the database
	query, args, err := QB.Insert("vendors").
		Columns("name", "image", "description", "created_at", "updated_at").
		Values(newVendor.Name, newVendor.Image, newVendor.Description, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Execute the query and retrieve the new vendor's ID
	var vendorID int
	if err := db.QueryRow(query, args...).Scan(&vendorID); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch the newly created vendor (to return it in the response)
	var createdVendor models.Vendor
	query, args, err = QB.Select(strings.Join(Vendor_columns, ", ")).
		From("vendors").
		Where("id = ?", vendorID).
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Get(&createdVendor, query, args...); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the newly created vendor in the response
	utils.SendJSONResponse(w, http.StatusCreated, createdVendor)
}

func UpdateVendorHandler(w http.ResponseWriter, r *http.Request) {
	var vendor models.Vendor
	id := r.PathValue("id")
	query, args, err := QB.Select(Vendor_columns...).
		From("vendors").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Get(&vendor, query, args...); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// update user
	if r.FormValue("name") != "" {
		vendor.Name = r.FormValue("name")
	}
	if r.FormValue("description") != "" {
		vendor.Description = r.FormValue("description")
	}

	var oldImg *string
	var newImg *string
	// Handle image file upload
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		utils.HandleError(w, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	} else if err == nil {
		defer file.Close()
		if vendor.Image != nil {
			oldImg = vendor.Image
		}
		imageName, err := utils.SaveImageFile(file, "vendors", fileHeader.Filename)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Error saving image file: "+err.Error())
			return
		}
		vendor.Image = &imageName
		newImg = &imageName
	}
	if vendor.Image != nil {
		*vendor.Image = strings.TrimPrefix(*vendor.Image, utils.Domain+"/")
	}

	query, args, err = QB.
		Update("vendors").
		Set("img", vendor.Image).
		Set("name", vendor.Name).
		Set("description", vendor.Description).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": vendor.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).
		ToSql()
	if err != nil {
		utils.DeleteImageFile(*newImg)
		utils.HandleError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&vendor); err != nil {
		utils.DeleteImageFile(*newImg)
		utils.HandleError(w, http.StatusInternalServerError, "Error creating Vendor"+err.Error())
		return
	}

	if oldImg != nil {
		if err := utils.DeleteImageFile(*oldImg); err != nil {
			log.Println(err)
		}
	}

	utils.SendJSONResponse(w, http.StatusOK, vendor)
}

func DeleteVendorandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// Use QB to build the delete query
	query, args, err := QB.Delete("vendors").
		Where("id = ?", id).
		Suffix("RETURNING img").
		ToSql()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	var img *string
	if err := db.QueryRow(query, args...).Scan(&img); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error deleting Vendor: "+err.Error())
		return
	}
	// If the user has an image, delete it
	if img != nil {
		if err := utils.DeleteImageFile(*img); err != nil {
			log.Println(err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}