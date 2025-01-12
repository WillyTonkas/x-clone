package controllers

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("userid")
	parentIDStr := r.FormValue("parent")
	quoteIDStr := r.FormValue("quote")
	body := r.FormValue("body")
	fmt.Println(body)

	if body == constants.EMPTY {
		http.Error(w, "Body cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid userid", http.StatusBadRequest)
		return
	}

	var parentID, quoteID *uint
	if parentIDStr != constants.EMPTY {
		parsedParentID, parentErr := strconv.ParseUint(parentIDStr, 10, 32)
		if parentErr != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		tempParentID := uint(parsedParentID)
		parentID = &tempParentID
	}

	if quoteIDStr != constants.EMPTY {
		parsedQuoteID, parsedErr := strconv.ParseUint(quoteIDStr, 10, 32)
		if parsedErr != nil {
			http.Error(w, "Invalid quote ID", http.StatusBadRequest)
			return
		}
		tempQuoteID := uint(parsedQuoteID)
		quoteID = &tempQuoteID
	}

	err = user.CreatePost(db, uint(userID), parentID, quoteID, body)
	if err != nil {
		if err.Error() == constants.ERRNOUSER {
			http.Error(w, constants.ERRNOUSER, http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Post created successfully"))
	if err != nil {
		return
	}
}

func ViewSpecificPost(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	pathID := r.PathValue("postid")
	if pathID == constants.EMPTY {
		http.Error(w, "Missing 'postid' parameter", http.StatusBadRequest)
		return
	}

	postID, errorParse := strconv.ParseUint(pathID, 10, 32)
	if errorParse != nil {
		http.Error(w, "Invalid postid", http.StatusBadRequest)
		return
	}

	post, getPostError := user.GetPostByID(db, uint(postID))

	if getPostError != nil {
		http.Error(w, "Failed to get post", http.StatusNotFound)
		return
	}

	response, marshalErr := json.Marshal(post)
	if marshalErr != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

var ViewSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/{postid}",
	HandlerFunction: ViewSpecificPost,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}
