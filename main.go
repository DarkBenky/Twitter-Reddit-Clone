package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	testing         = true
	testingLogin    = "test"
	testingPassword = "test"
)

type User struct {
	IDUser      int    `json:"idUser"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type Post struct {
	IDPost      int    `json:"idPost"`
	ContentText string `json:"content_text"`
	CreatedAt   string `json:"created_at"`
	UserID      int    `json:"userID"`
	CategoryID  int    `json:"categoryID"`
	ImageURL    string `json:"imageURL"`
}

type PostWithCategory struct {
	IDPost      int    `json:"idPost"`
	ContentText string `json:"content_text"`
	CreatedAt   string `json:"created_at"`
	UserID      int    `json:"userID"`
	CategoryID  int    `json:"categoryID"`
	Category    string `json:"category"`
	ImageURL    string `json:"imageURL"`
}

type SubscribeUser struct {
	IDSubscription int `json:"idSubscription"`
	SubscriberID   int `json:"subscriberID"`
	SubscribedToID int `json:"subscribedToID"`
}

type Item struct {
	IDItem      int    `json:"idItem"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Price       string `json:"price"`
	UserID      int    `json:"userID"`
	CreatedAt   string `json:"created_at"`
}

type Comment struct {
	IDComment   int    `json:"idComment"`
	IDPost      int    `json:"idPost"`
	IDUser      int    `json:"idUser"`
	ContentText string `json:"content_text"`
	CreatedAt   string `json:"created_at"`
}

type Category struct {
	IDCategory  int    `json:"idCategory"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LikeDislike struct {
	IDLikeDislike int `json:"idLikeDislike"`
	IDPost        int `json:"idPost"`
	IDUser        int `json:"idUser"`
	Like          int `json:"like"`
}

type SavedPost struct {
	IDSavedPost int `json:"idSavedPost"`
	IDPost      int `json:"idPost"`
	IDUser      int `json:"idUser"`
}

type Message struct {
	IDMessage  int    `json:"idMessage"`
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

func GetAllPosts(c echo.Context) error {
	// Keep struct definition
	type Post struct {
		IDPost      int    `json:"idPost"`
		ContentText string `json:"content_text"`
		CreatedAt   string `json:"created_at"`
		UserID      int    `json:"userID"`
		Category    string `json:"category"`
		ImageURL    string `json:"imageURL"`
	}

	offset := c.QueryParam("offset")
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Printf("Invalid offset: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid offset parameter"})
	}

	// Reorder SELECT to match struct field order
	query := `
        SELECT 
            p.idPost,
            p.content_text,
            p.created_at,
            p.userID,
            COALESCE(c.name, '') as category,
            COALESCE(p.imageURL, '') as imageURL
        FROM posts p
        LEFT JOIN categories c ON p.categoryID = c.idCategory
        ORDER BY p.created_at DESC 
        LIMIT 10 OFFSET ?`

	rows, err := db.Query(query, offsetInt)
	if err != nil {
		log.Printf("Query error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query posts"})
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		// Scan order now matches SELECT order
		if err := rows.Scan(
			&post.IDPost,
			&post.ContentText,
			&post.CreatedAt,
			&post.UserID,
			&post.Category,
			&post.ImageURL,
		); err != nil {
			log.Printf("Scan error: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan post data"})
		}
		posts = append(posts, post)
	}

	return c.JSON(http.StatusOK, posts)
}

func GetAllPostsForCategory(c echo.Context) error {
	category := c.QueryParam("category")

	if category == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Category parameter is required"})
	}

	type Post struct {
		IDPost      int    `json:"idPost"`
		ContentText string `json:"content_text"`
		CreatedAt   string `json:"created_at"`
		UserID      int    `json:"userID"`
		Category    string `json:"category"`
		ImageURL    string `json:"imageURL"`
	}

	offset := c.QueryParam("offset")
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Printf("Invalid offset: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid offset parameter"})
	}

	query := `
        SELECT 
            p.idPost,
            p.content_text,
            p.created_at,
            p.userID,
            COALESCE(c.name, '') as category,
            COALESCE(p.imageURL, '') as imageURL
        FROM posts p
        LEFT JOIN categories c ON p.categoryID = c.idCategory
        WHERE c.name = ?
        ORDER BY p.created_at DESC 
        LIMIT 10 OFFSET ?`

	rows, err := db.Query(query, category, offsetInt)
	if err != nil {
		log.Printf("Query error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query posts"})
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.IDPost,
			&post.ContentText,
			&post.CreatedAt,
			&post.UserID,
			&post.Category,
			&post.ImageURL,
		); err != nil {
			log.Printf("Scan error: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan post data"})
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "No posts found for the given category"})
	}

	return c.JSON(http.StatusOK, posts)
}

func GetUsers() []User {
	// Get all users from the database
	query := `SELECT idUser, username, displayName, email FROM users`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.IDUser, &user.Username, &user.DisplayName, &user.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return users
}

func GetAllCategories(c echo.Context) error {
	categories := GetCategories()
	return c.JSON(http.StatusOK, categories)
}

func GetCategories() []Category {
	// Get all categories from the database
	query := `SELECT idCategory, name, description FROM categories`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.IDCategory, &category.Name, &category.Description); err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return categories
}

func GetPosts() []Post {
	// Get all posts from the database
	query := `SELECT idPost, content_text, created_at, userID FROM posts`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.IDPost, &post.ContentText, &post.CreatedAt, &post.UserID); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return posts
}

func GetPostByUserID(c echo.Context) error {
	userID := c.QueryParam("id")

	// Get all posts from the database with category information
	query := `
        SELECT 
            p.idPost, 
            p.content_text, 
            p.created_at, 
            p.userID, 
            p.imageURL,
            p.categoryID,
            COALESCE(c.name, '') as category_name
        FROM posts p
        LEFT JOIN categories c ON p.categoryID = c.idCategory
        WHERE p.userID = ?
        ORDER BY p.created_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to query posts",
		})
	}
	defer rows.Close()

	var posts []PostWithCategory
	for rows.Next() {
		var post PostWithCategory
		if err := rows.Scan(
			&post.IDPost,
			&post.ContentText,
			&post.CreatedAt,
			&post.UserID,
			&post.ImageURL,
			&post.CategoryID,
			&post.Category,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan post data",
			})
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error iterating over rows",
		})
	}

	return c.JSON(http.StatusOK, posts)
}

func GetCategoryByID(c echo.Context) error {
	categoryID := c.QueryParam("categoryId")

	// Get category from the database
	query := `SELECT idCategory, name, description FROM categories WHERE idCategory = ?`
	row := db.QueryRow(query, categoryID)

	var category Category
	if err := row.Scan(&category.IDCategory, &category.Name, &category.Description); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan category data"})
	}

	// Return category as JSON response
	return c.JSON(http.StatusOK, category)
}

func GetAllCommentsToPost(c echo.Context) error {
	postID := c.QueryParam("idPost")

	// Get all comments from the database
	query := `SELECT idComment, idPost, idUser, content_text, created_at FROM comments WHERE idPost = ?`
	rows, err := db.Query(query, postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query comments"})
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.IDComment, &comment.IDPost, &comment.IDUser, &comment.ContentText, &comment.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan comment data"})
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error iterating over rows"})
	}

	// Return comments as JSON response
	return c.JSON(http.StatusOK, comments)
}

func GetUserByID(c echo.Context) error {
	userID := c.QueryParam("id")

	// Get user from the database
	query := `SELECT idUser, username, displayName, email FROM users WHERE idUser = ?`
	row := db.QueryRow(query, userID)

	var user User
	if err := row.Scan(&user.IDUser, &user.Username, &user.DisplayName, &user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan user data"})
	}

	// Return user as JSON response
	return c.JSON(http.StatusOK, user)
}

func GetAllUsers(c echo.Context) error {
	// Get all users from the database
	query := `SELECT idUser, username, displayName, email FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query users"})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.IDUser, &user.Username, &user.DisplayName, &user.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan user data"})
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error iterating over rows"})
	}

	// Return users as JSON response
	return c.JSON(http.StatusOK, users)
}

func NumberOfSubscribers(c echo.Context) error {
    userID := c.QueryParam("userID")

    if userID == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "error": "User ID is required",
        })
    }

    query := `SELECT COUNT(*) FROM subscriptions WHERE subscribedToID = ?`
    var count int
    err := db.QueryRow(query, userID).Scan(&count)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "error": "Failed to count subscribers",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "numberOfSubscribers": count,
    })
}

func NumberOfSubscribeTo(c echo.Context) error {
    userID := c.QueryParam("userID")

    if userID == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "error": "User ID is required",
        })
    }

    query := `SELECT COUNT(*) FROM subscriptions WHERE subscriberID = ?`
    var count int
    err := db.QueryRow(query, userID).Scan(&count)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "error": "Failed to count subscriptions",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "numberOfSubscriptions": count,
    })
}

func AddPost(c echo.Context) error {
	type PostRequest struct {
		ContentText string `json:"content_text"`
		ImageURL    string `json:"imageURL"` // Added
		UserID      string `json:"userID"`
		CategoryID  string `json:"categoryID"`
	}

	postReq := new(PostRequest)
	if err := c.Bind(postReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request data"})
	}

	query := `INSERT INTO posts (content_text, imageURL, created_at, userID, categoryID)
              VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query,
		postReq.ContentText,
		postReq.ImageURL,
		time.Now().Format(time.RFC3339),
		postReq.UserID,
		postReq.CategoryID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Post created"})
}

func AddComment(c echo.Context) error {
	type CommentRequest struct {
		PostID      string `json:"postID"`
		UserID      string `json:"userID"`
		ContentText string `json:"contentText"`
	}

	comment := new(CommentRequest)
	if err := c.Bind(comment); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request data",
		})
	}

	// Validate input
	if comment.PostID == "" || comment.UserID == "" || comment.ContentText == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "PostID, userID, and content text are required",
		})
	}

	// Insert comment into the database
	query := `INSERT INTO comments (idPost, idUser, content_text, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, comment.PostID, comment.UserID, comment.ContentText, time.Now().Format(time.RFC3339))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to insert comment: " + err.Error(),
		})
	}

	// Check if the insert was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to confirm comment insertion",
		})
	}

	// Return success response
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Comment added successfully",
	})
}

func EditPost(c echo.Context) error {
	// Create struct for request body
	type EditRequest struct {
		PostID      string `json:"postID"`
		ContentText string `json:"contentText"`
		ImageURL    string `json:"imageURL"`
		CategoryID  string `json:"categoryID"`
	}

	// Parse request body
	var req EditRequest
	if err := c.Bind(&req); err != nil {
		// Add debug logging
		fmt.Printf("Request binding error: %v\n", err)
		fmt.Printf("Request body: %+v\n", c.Request().Body)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format: " + err.Error(),
		})
	}

	// Log received data
	fmt.Printf("Received request: %+v\n", req)

	// Validate input
	if req.PostID == "" || req.ContentText == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Post ID and content text are required",
		})
	}

	fmt.Printf("Executing query with values: content=%s, image=%s, category=%s, postID=%s\n",
		req.ContentText, req.ImageURL, req.CategoryID, req.PostID)

	query := `UPDATE posts SET content_text = ?, imageURL = ?, categoryID = ? WHERE idPost = ?`
	result, err := db.Exec(query, req.ContentText, req.ImageURL, req.CategoryID, req.PostID)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to update post: " + err.Error(),
		})
	}

	// Verify update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to confirm update",
		})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Post not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post updated successfully",
	})
}

// In main.go, update the DeletePost function:
func DeletePost(c echo.Context) error {
	// Create struct for request body
	type DeleteRequest struct {
		PostID string `json:"postID"`
	}

	// Parse request body
	var req DeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	// Validate postID
	if req.PostID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Post ID is required",
		})
	}

	// Delete post from database
	query := `DELETE FROM posts WHERE idPost = ?`
	result, err := db.Exec(query, req.PostID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to delete post: " + err.Error(),
		})
	}

	// Verify deletion
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to confirm deletion",
		})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Post not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post deleted successfully",
	})
}

func GetPostById(c echo.Context) error {
	type Post struct {
		IDPost      int    `json:"idPost"`
		ContentText string `json:"content_text"`
		CreatedAt   string `json:"created_at"`
		UserID      int    `json:"userID"`
		Category    string `json:"category"`
		ImageURL    string `json:"imageURL"` // Added ImageURL field
	}

	postID := c.QueryParam("id")
	if postID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Post ID is required"})
	}

	query := `
        SELECT p.idPost, p.content_text, p.created_at, p.userID, 
               COALESCE(c.name, '') as category_name,
               COALESCE(p.imageURL, '') as imageURL
        FROM posts p 
        LEFT JOIN categories c ON p.categoryID = c.idCategory 
        WHERE p.idPost = ?`

	var post Post
	err := db.QueryRow(query, postID).Scan(
		&post.IDPost,
		&post.ContentText,
		&post.CreatedAt,
		&post.UserID,
		&post.Category,
		&post.ImageURL, // Added ImageURL scan
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Post not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan post data"})
	}

	return c.JSON(http.StatusOK, post)
}

func AddCategory(c echo.Context) error {
	// Define request structure
	type CategoryRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// Bind request body
	category := new(CategoryRequest)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	// Validate input
	if category.Name == "" || category.Description == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Name and description are required",
		})
	}

	// Insert category
	query := `INSERT INTO categories (name, description) VALUES (?, ?)`
	result, err := db.Exec(query, category.Name, category.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to insert category",
		})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to get category ID",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":    "Category added successfully",
		"categoryId": id,
	})
}

func UpdatePassword(c echo.Context) error {
	type PasswordRequest struct {
		UserID   int    `json:"userID"`
		Password string `json:"password"`
	}

	passwordReq := new(PasswordRequest)
	if err := c.Bind(passwordReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request data"})
	}

	fmt.Println(passwordReq)

	query := `UPDATE users SET password = ? WHERE idUser = ?`
	_, err := db.Exec(query, passwordReq.Password, passwordReq.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Password updated"})
}

func addPostToSavedPosts(postID, userID int) error {
	// Check if the post is already saved by the user
	isSaved, err := checkIfPostIsSaved(postID, userID)
	if err != nil {
		return err
	}

	if isSaved {
		// remove post from saved posts
		query := `DELETE FROM saved_posts WHERE idPost = ? AND idUser = ?`
		_, err := db.Exec(query, postID, userID)
		if err != nil {
			return err
		}
		return nil
	} else {
		// Add post to saved posts
		query := `INSERT INTO saved_posts (idPost, idUser) VALUES (?, ?)`
		_, err := db.Exec(query, postID, userID)
		if err != nil {
			return err
		}
		return nil
	}
}

func AddPostToSavedPosts(c echo.Context) error {
	type SavePostRequest struct {
		PostID string `json:"postID"`
		UserID string `json:"userID"`
	}

	// Bind the request body to the struct
	request := new(SavePostRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	// Check if the required fields are provided
	if request.PostID == "" || request.UserID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Post ID and User ID are required",
		})
	}

	// Convert string IDs to integers
	postIDInt, err := strconv.Atoi(request.PostID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid Post ID format",
		})
	}

	userIDInt, err := strconv.Atoi(request.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid User ID format",
		})
	}

	err = addPostToSavedPosts(postIDInt, userIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to add post to saved posts",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post added to saved posts successfully",
	})
}

func checkIfPostIsSaved(postID, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM saved_posts WHERE idPost = ? AND idUser = ?`
	var count int
	err := db.QueryRow(query, postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckIfPostIsSaved(c echo.Context) error {
	// Get parameters from query instead of binding from request body
	postId := c.QueryParam("postID")
	userId := c.QueryParam("userID")

	if postId == "" || userId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Post ID and User ID are required",
		})
	}

	// Convert string IDs to integers
	postIDInt, err := strconv.Atoi(postId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid Post ID format",
		})
	}

	userIDInt, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid User ID format",
		})
	}

	isSaved, err := checkIfPostIsSaved(postIDInt, userIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to check if post is saved",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"saved": isSaved,
	})
}

func GetUsersSavedPosts(c echo.Context) error {
	userID := c.QueryParam("userID")

	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "User ID is required",
		})
	}

	// Updated query to join with posts table
	query := `
        SELECT p.idPost, p.content_text, p.imageURL, p.created_at, p.userID, p.categoryID 
        FROM saved_posts sp
        JOIN posts p ON sp.idPost = p.idPost
        WHERE sp.idUser = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to query saved posts",
		})
	}
	defer rows.Close()

	var savedPosts []PostWithCategory
	for rows.Next() {
		var post PostWithCategory
		if err := rows.Scan(&post.IDPost, &post.ContentText, &post.ImageURL, &post.CreatedAt, &post.UserID, &post.CategoryID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan saved post data",
			})
		}
		savedPosts = append(savedPosts, post)
	}

	// Get category name for each post
	for i, post := range savedPosts {
		query := `SELECT name FROM categories WHERE idCategory = ?`
		row := db.QueryRow(query, post.CategoryID)
		var categoryName string
		err := row.Scan(&categoryName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan category name",
			})
		}
		savedPosts[i].Category = categoryName
	}

	return c.JSON(http.StatusOK, savedPosts)
}

func CheckIfUserSubscribed(c echo.Context) error {
	subscribedToID := c.QueryParam("subscribedToID")
	subscriberID := c.QueryParam("subscriberID")

	if subscribedToID == "" || subscriberID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "User ID and Subscriber ID are required",
		})
	}

	query := `SELECT COUNT(*) FROM subscriptions WHERE subscriberID = ? AND subscribedToID = ?`
	var count int
	err := db.QueryRow(query, subscriberID, subscribedToID).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to check subscription status",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"subscribed": count > 0,
	})
}

func SubscribeORUnsubscribe(c echo.Context) error {
	subscribedToID := c.QueryParam("subscribedToID")
	subscriberID := c.QueryParam("subscriberID")

	if subscribedToID == "" || subscriberID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "User ID and Subscriber ID are required",
		})
	}

	// Check if the subscription already exists
	query := `SELECT COUNT(*) FROM subscriptions WHERE subscriberID = ? AND subscribedToID = ?`
	var count int
	err := db.QueryRow(query, subscriberID, subscribedToID).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to check subscription status",
		})
	}
	if count > 0 {
		// If it exists, unsubscribe
		query = `DELETE FROM subscriptions WHERE subscriberID = ? AND subscribedToID = ?`
		_, err := db.Exec(query, subscriberID, subscribedToID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to unsubscribe",
			})
		}
	} else {
		// If it doesn't exist, subscribe
		query = `INSERT INTO subscriptions (subscriberID, subscribedToID) VALUES (?, ?)`
		_, err := db.Exec(query, subscriberID, subscribedToID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to subscribe",
			})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Subscription status changed successfully",
	})
}

func GetListOfSubscribers(c echo.Context) error {
	userID := c.QueryParam("userID")

	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "User ID is required",
		})
	}

	query := `SELECT subscribedToID FROM subscriptions WHERE subscriberID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to query subscribers",
		})
	}

	defer rows.Close()
	var subscribers []int
	for rows.Next() {
		var subscriberID int
		if err := rows.Scan(&subscriberID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan subscriber data",
			})
		}
		subscribers = append(subscribers, subscriberID)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error iterating over rows",
		})
	}
	// Return the list of subscribers
	return c.JSON(http.StatusOK, subscribers)
}

func main() {

	// Open a connection to the SQLite database
	database, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	db = database

	// Create tables
	createUsersTable(database)
	createSavedPostsTable(database)
	createPostsTable(database)
	createCommentsTable(database)
	createCategoriesTable(database)
	createLikesDislikesTable(database)
	createMessagesTable(database)
	createSubscriptionsTable(database)

	// // Generate random users, posts, and comments
	// n := 20 // Number of random entries to generate

	// // fmt.Printf("Generating %d random users...\n", n)
	// insertRandomUsers(n)

	// InsertRandomMessages(n)
	// InsertRandomCategories(n / 10)

	// n = n / 5
	// fmt.Printf("Generating %d random posts...\n", n)
	// insertRandomPosts(n)

	// n = n / 2
	// fmt.Printf("Generating %d random comments...\n", n)
	// insertRandomComments(n)

	// InsertTestUser()

	// Start the server
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080", "http://138.68.76.63:8080"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.GET("/posts", GetAllPosts)
	e.GET("/posts/category", GetAllPostsForCategory)
	e.GET("/comments", GetAllCommentsToPost)
	e.GET("/user", GetUserByID)
	e.GET("/users", GetAllUsers)
	e.GET("/posts/user", GetPostByUserID)
	e.GET("/post", GetPostById)
	e.POST("/addPost", AddPost)
	e.POST("/addComment", AddComment)
	e.DELETE("/deletePost", DeletePost)
	e.PUT("/editPost", EditPost)
	e.POST("/login", Login)
	e.PUT("/userEdit", UpdateUser)
	e.GET("/likesDislikes", getLikesDislikes)
	e.GET("/userLikeDislike", getUserLikeDislikeForPost)
	e.GET("/like", like)
	e.GET("/dislike", dislike)
	e.GET("/messages", getMessages)
	e.POST("/sendMessage", SendMessages)
	e.GET("/conversations", GetUserConversations)
	e.POST("/register", Register)
	e.GET("/categories", GetAllCategories)
	e.GET("/category", GetCategoryByID)
	e.POST("/addCategory", AddCategory)
	e.POST("/updatePassword", UpdatePassword)
	e.POST("/savePost", AddPostToSavedPosts)
	e.GET("/checkPostSaved", CheckIfPostIsSaved)
	e.GET("/savedPosts", GetUsersSavedPosts)
	e.GET("/checkSubscription", CheckIfUserSubscribed)
	e.GET("/subscribe", SubscribeORUnsubscribe)
	e.GET("/listOfSubscribers", GetListOfSubscribers)
	e.GET("/numberOfSubscribers", NumberOfSubscribers)
	e.GET("/numberOfSubscribeTo", NumberOfSubscribeTo)
	e.Logger.Fatal(e.Start(":5533"))

}

func Register(c echo.Context) error {
	// Create user struct to bind request body
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if user.Username == "" || user.DisplayName == "" || user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "All fields are required",
		})
	}

	// Debug log incoming request
	fmt.Printf("Register attempt - Username: %s, Email: %s\n", user.Username, user.Email)

	// Check if user already exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?",
		user.Username, user.Email).Scan(&exists)
	if err != nil {
		fmt.Printf("Database error checking user existence: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}

	if exists > 0 {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "User already exists",
		})
	}

	// Insert user into database
	query := `INSERT INTO users (username, displayName, email, password) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.DisplayName, user.Email, user.Password)
	if err != nil {
		fmt.Printf("Error inserting new user: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to insert user",
		})
	}

	id, _ := result.LastInsertId()
	fmt.Printf("Successfully registered new user with ID: %d\n", id)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User registered successfully",
	})
}

func InsertTestUser() {
	_, err := db.Exec(`INSERT INTO users (username, displayName, email, password) VALUES (?, ?, ?, ?)`,
		testingLogin, testingLogin, testingLogin+"@gmail.com", testingPassword)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted test user.\n")

	// insert random post and comment for test user
	_, err = db.Exec(`INSERT INTO posts (content_text, created_at, userID) VALUES (?, ?, ?)`,
		faker.Sentence(), time.Now().Format(time.RFC3339), 1)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO comments (idPost, idUser, content_text, created_at) VALUES (?, ?, ?, ?)`,
		1, 1, faker.Sentence(), time.Now().Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted random post and comment for test user.\n")
}

type UpdateUserRequest struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

func UpdateUser(c echo.Context) error {
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request payload",
		})
	}

	// Update the user in the database
	stmt, err := db.Prepare("UPDATE users SET username = ?, displayName = ?, email = ? WHERE idUser = ?")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}
	defer stmt.Close()

	res, err := stmt.Exec(req.Username, req.DisplayName, req.Email, req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to update user",
		})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "User not found",
		})
	}

	// Retrieve the updated user
	var updatedUser User
	err = db.QueryRow("SELECT idUser, username, displayName, email FROM users WHERE idUser = ?", req.ID).
		Scan(&updatedUser.IDUser, &updatedUser.Username, &updatedUser.DisplayName, &updatedUser.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve updated user",
		})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getLikesDislikes(c echo.Context) error {
	postId := c.QueryParam("idPost")
	query := `SELECT like FROM likes_dislikes WHERE idPost = ?`
	rows, err := db.Query(query, postId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to query likes/dislikes",
		})
	}
	defer rows.Close()

	var LikeCount int
	var DislikeCount int

	// print row count
	// fmt.Println(rows.Columns())

	for rows.Next() {
		var like int
		if err := rows.Scan(&like); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan like/dislike data",
			})
		}
		// fmt.Println(like, "like")
		if like == 1 {
			LikeCount++
		} else if like == -1 {
			DislikeCount++
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"likes":    LikeCount,
		"dislikes": DislikeCount,
	})
}

func getUserLikeDislikeForPost(c echo.Context) error {
	postId := c.QueryParam("idPost")
	userId := c.QueryParam("userID")

	// fmt.Println("Received request for user like/dislike")
	// fmt.Printf("Post ID: %s, User ID: %s\n", postId, userId)

	query := `SELECT idLikeDislike, idPost, idUser, like FROM likes_dislikes WHERE idPost = ? AND idUser = ?`
	row := db.QueryRow(query, postId, userId)

	var likeDislike LikeDislike
	if err := row.Scan(&likeDislike.IDLikeDislike, &likeDislike.IDPost, &likeDislike.IDUser, &likeDislike.Like); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"notLikedDisliked": "User has not liked/disliked this post",
		})
	}

	return c.JSON(http.StatusOK, likeDislike)
}

func like(c echo.Context) error {
	postId := c.QueryParam("postId")
	userId := c.QueryParam("userId")

	// Add debug logging
	// fmt.Printf("Received like request - postId: %s, userId: %s\n", postId, userId)

	// Convert string IDs to integers for database query
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid post ID format",
		})
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid user ID format",
		})
	}

	// Check if post exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE idPost = ?)", postIdInt).Scan(&exists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Post not found",
		})
	}

	// Check if user exists
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE idUser = ?)", userIdInt).Scan(&exists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "User not found",
		})
	}

	// Check existing like/dislike
	query := `SELECT idLikeDislike, like FROM likes_dislikes WHERE idPost = ? AND idUser = ?`
	row := db.QueryRow(query, postIdInt, userIdInt)

	var likeDislike LikeDislike
	err = row.Scan(&likeDislike.IDLikeDislike, &likeDislike.Like)

	if err == sql.ErrNoRows {
		// No existing like/dislike, insert new like
		query = `INSERT INTO likes_dislikes (idPost, idUser, like) VALUES (?, ?, 1)`
		_, err := db.Exec(query, postIdInt, userIdInt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to add like",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Like added",
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}

	if likeDislike.Like == 1 {
		// Already liked, remove it
		query = `DELETE FROM likes_dislikes WHERE idPost = ? AND idUser = ?`
		_, err := db.Exec(query, postIdInt, userIdInt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to remove like",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Like removed",
		})
	} else {
		// Was disliked, update to like
		query = `UPDATE likes_dislikes SET like = 1 WHERE idPost = ? AND idUser = ?`
		_, err := db.Exec(query, postIdInt, userIdInt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to update to like",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Changed from dislike to like",
		})
	}
}

func dislike(c echo.Context) error {
	postId := c.QueryParam("postId")
	userId := c.QueryParam("userId")

	// Check if post exists
	query := `SELECT idPost FROM posts WHERE idPost = ?`
	row := db.QueryRow(query, postId)
	var post Post
	if err := row.Scan(&post.IDPost); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Post not found",
		})
	}

	// Check if user exists
	query = `SELECT idUser FROM users WHERE idUser = ?`
	row = db.QueryRow(query, userId)
	var user User
	if err := row.Scan(&user.IDUser); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "User not found",
		})
	}

	// Check existing like/dislike
	query = `SELECT idLikeDislike, like FROM likes_dislikes WHERE idPost = ? AND idUser = ?`
	row = db.QueryRow(query, postId, userId)

	var likeDislike LikeDislike
	err := row.Scan(&likeDislike.IDLikeDislike, &likeDislike.Like)

	if err == sql.ErrNoRows {
		// No existing like/dislike, insert new dislike
		query = `INSERT INTO likes_dislikes (idPost, idUser, like) VALUES (?, ?, -1)`
		_, err := db.Exec(query, postId, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to add dislike",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Dislike added",
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database error",
		})
	}

	if likeDislike.Like == -1 {
		// Already disliked, remove it
		query = `DELETE FROM likes_dislikes WHERE idPost = ? AND idUser = ?`
		_, err := db.Exec(query, postId, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to remove dislike",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Dislike removed",
		})
	} else {
		// Was liked, update to dislike
		query = `UPDATE likes_dislikes SET like = -1 WHERE idPost = ? AND idUser = ?`
		_, err := db.Exec(query, postId, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to update to dislike",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Changed from like to dislike",
		})
	}
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Username and password are required",
		})
	}

	// Check if user exists
	query := `SELECT idUser, username, displayName, email FROM users WHERE (username = ? OR email = ?) AND password = ?`
	row := db.QueryRow(query, req.Username, req.Username, req.Password)

	var user User
	if err := row.Scan(&user.IDUser, &user.Username, &user.DisplayName, &user.Email); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Invalid username or password",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// Create Users table
func createUsersTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"idUser" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"username" TEXT UNIQUE,
		"displayName" TEXT UNIQUE,
		"email" TEXT UNIQUE,
		"password" TEXT
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Users table created")
}

func createSavedPostsTable(db *sql.DB) {
	fmt.Println("Creating saved_posts table")
	createTableSQL := `CREATE TABLE IF NOT EXISTS saved_posts (
		"idSavedPost" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"idPost" INTEGER,
		"idUser" INTEGER,
		FOREIGN KEY(idPost) REFERENCES posts(idPost),
		FOREIGN KEY(idUser) REFERENCES users(idUser)
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Saved posts table created")
}

// Create Posts table
func createPostsTable(db *sql.DB) {
	fmt.Println("Creating posts table")
	createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
        "idPost" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "content_text" TEXT,
        "imageURL" TEXT,
        "created_at" TEXT,
        "userID" INTEGER,
        "categoryID" INTEGER,
        FOREIGN KEY ("userID") REFERENCES users(idUser),
        FOREIGN KEY ("categoryID") REFERENCES categories(idCategory)
    );`
	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
	fmt.Println("Posts table created")
}

// Create Comments table
func createCommentsTable(db *sql.DB) {
	fmt.Println("Creating comments table")
	createTableSQL := `CREATE TABLE IF NOT EXISTS comments (
		"idComment" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"idPost" INTEGER,
		"idUser" INTEGER,
		"content_text" TEXT,
		"created_at" TEXT,
		FOREIGN KEY(idPost) REFERENCES posts(idPost),
		FOREIGN KEY(idUser) REFERENCES users(idUser)
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Comments table created")
}

func createCategoriesTable(db *sql.DB) {
	fmt.Println("Creating categories table")
	createTableSQL := `CREATE TABLE IF NOT EXISTS categories (
		"idCategory" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"description" TEXT
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Categories table created")
}

func createLikesDislikesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS likes_dislikes (
		"idLikeDislike" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"idPost" INTEGER,
		"idUser" INTEGER,
		"like" INTEGER,
		FOREIGN KEY(idPost) REFERENCES posts(idPost),
		FOREIGN KEY(idUser) REFERENCES users(idUser)
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Likes/Dislikes table created")
}

func createMessagesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"idMessage" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"senderID" INTEGER,
		"receiverID" INTEGER,
		"content" TEXT,
		"created_at" TEXT,
		FOREIGN KEY(senderID) REFERENCES users(idUser),
		FOREIGN KEY(receiverID) REFERENCES users(idUser)
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Messages table created")
}

func createSubscriptionsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS subscriptions (
		"idSubscription" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"subscriberID" INTEGER,
		"subscribedToID" INTEGER,
		FOREIGN KEY(subscriberID) REFERENCES users(idUser),
		FOREIGN KEY(subscribedToID) REFERENCES users(idUser)
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	fmt.Println("Subscriptions table created")
}

func getMessages(c echo.Context) error {
	senderId := c.QueryParam("senderID")
	receiverId := c.QueryParam("receiverID")

	query := `SELECT idMessage, senderID, receiverID, content, created_at FROM messages WHERE senderID = ? AND receiverID = ? OR senderID = ? AND receiverID = ? ORDER BY created_at`
	rows, err := db.Query(query, senderId, receiverId, receiverId, senderId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to query messages",
		})
	}

	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.IDMessage, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan message data",
			})
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error iterating over rows",
		})
	}

	return c.JSON(http.StatusOK, messages)
}

func SendMessages(c echo.Context) error {
	type MessageRequest struct {
		SenderID   string `json:"senderID"`
		ReceiverID string `json:"receiverID"`
		Content    string `json:"content"`
	}

	message := new(MessageRequest)
	if err := c.Bind(message); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request data",
		})
	}

	if message.SenderID == "" || message.ReceiverID == "" || message.Content == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "SenderID, receiverID, and content are required",
		})
	}

	query := `INSERT INTO messages (senderID, receiverID, content, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, message.SenderID, message.ReceiverID, message.Content, time.Now().Format(time.RFC3339))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to insert message: " + err.Error(),
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to confirm message insertion",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Message sent successfully",
	})
}

func GetUserConversations(c echo.Context) error {
	userID := c.QueryParam("userID")

	type Conversation struct {
		SenderID     int    `json:"senderID"`
		SenderName   string `json:"senderName"`
		ReceiverID   int    `json:"receiverID"`
		ReceiverName string `json:"receiverName"`
		LastMessage  string `json:"lastMessage"`
	}

	// Get distinct conversations with last message
	query := `
        WITH RankedMessages AS (
            SELECT 
                m.*,
                ROW_NUMBER() OVER (
                    PARTITION BY 
                        CASE 
                            WHEN senderID < receiverID 
                            THEN senderID || '-' || receiverID 
                            ELSE receiverID || '-' || senderID 
                        END 
                    ORDER BY created_at DESC
                ) as rn
            FROM messages m
            WHERE senderID = ? OR receiverID = ?
        )
        SELECT 
            m.senderID,
            s.displayName as senderName,
            m.receiverID,
            r.displayName as receiverName,
            m.content as lastMessage
        FROM RankedMessages m
        JOIN users s ON m.senderID = s.idUser
        JOIN users r ON m.receiverID = r.idUser
        WHERE rn = 1
        ORDER BY m.created_at DESC`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to query conversations: " + err.Error(),
		})
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		if err := rows.Scan(
			&conv.SenderID,
			&conv.SenderName,
			&conv.ReceiverID,
			&conv.ReceiverName,
			&conv.LastMessage,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to scan conversation data: " + err.Error(),
			})
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error iterating over rows: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, conversations)
}

// insertRandomUsers generates and inserts synthetic user data into the database for testing purposes.
// It creates 'n' users with randomly generated usernames, display names, and emails using the faker library.
//
// Parameters:
//   - n: The number of random users to generate and insert
//
// The function will:
//   - Generate unique random usernames, display names and emails for each user
//   - Insert the generated data into the 'users' table
//   - Print confirmation message after successful insertion
//
// Example usage:
//
//	insertRandomUsers(100) // Generates 100 random users
//
// Note: This function will fail fast with log.Fatal if any database errors occur
func insertRandomUsers(n int) {
	for i := 0; i < n; i++ {
		username := faker.Username()
		displayName := faker.Username()
		email := faker.Email()
		password := faker.Password()
		_, err := db.Exec(`INSERT INTO users (username, displayName, email, password) VALUES (?, ?, ?, ?)`,
			username, displayName, email, password)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Inserted %d random users.\n", n)
}

func InsertRandomMessages(n int) {
	Users := GetUsers()

	for i := 0; i < n; i++ {
		senderID := Users[rand.Intn(len(Users))].IDUser
		receiverID := Users[rand.Intn(len(Users))].IDUser
		content := faker.Sentence()
		createdAt := time.Now().Format(time.RFC3339)
		_, err := db.Exec(`INSERT INTO messages (senderID, receiverID, content, created_at) VALUES (?, ?, ?, ?)`,
			senderID, receiverID, content, createdAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Inserted %d random messages.\n", n)
}

func InsertRandomCategories(n int) {
	for i := 0; i < n; i++ {
		name := faker.Word()
		description := faker.Sentence()
		_, err := db.Exec(`INSERT INTO categories (name, description) VALUES (?, ?)`,
			name, description)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Inserted %d random categories.\n", n)
}

// insertRandomPosts generates and inserts synthetic post data for existing users in the database.
// For each user, it creates a random number of posts (1 to n) with randomly generated content.
//
// Parameters:
//   - n: The maximum number of posts to generate per user
//
// The function will:
//   - Retrieve all existing users from the database
//   - For each user, generate between 1 and n random posts
//   - Set the creation timestamp to the current time
//   - Insert the generated posts into the 'posts' table
//   - Print confirmation message after successful insertion
//
// Example usage:
//
//	insertRandomPosts(5) // Generates up to 5 posts per user
//
// Note:
//   - Requires existing users in the database
//   - Will fail fast with log.Fatal if any database errors occur
func insertRandomPosts(n int) {
	Users := GetUsers()
	Categories := GetCategories()

	for _, user := range Users {
		numberOfPosts := rand.Intn(n) + 1
		for i := 0; i < numberOfPosts; i++ {
			content := faker.Sentence()
			createdAt := time.Now().Format(time.RFC3339)
			categoryID := Categories[rand.Intn(len(Categories))].IDCategory
			testImage := "https://picsum.photos/200/300"
			_, err := db.Exec(`INSERT INTO posts (content_text, created_at, userID, categoryID, imageURL ) VALUES (?, ?, ?, ?, ?)`,
				content, createdAt, user.IDUser, categoryID, testImage)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("Inserted random posts for %d users.\n", len(Users))
}

// insertRandomComments generates and inserts synthetic comment data for existing posts in the database.
// For each post, it creates a random number of comments (1 to n) with randomly generated content
// and assigns them to random users.
//
// Parameters:
//   - n: The maximum number of comments to generate per post
//
// The function will:
//   - Retrieve all existing posts and users from the database
//   - For each post, generate between 1 and n random comments
//   - Randomly assign each comment to an existing user
//   - Set the creation timestamp to the current time
//   - Insert the generated comments into the 'comments' table
//   - Print confirmation message after successful insertion
//
// Example usage:
//
//	insertRandomComments(10) // Generates up to 10 comments per post
//
// Note:
//   - Requires existing posts and users in the database
//   - Will fail fast with log.Fatal if any database errors occur
//   - Uses random user selection, so comment distribution may not be uniform
func insertRandomComments(n int) {
	Posts := GetPosts()
	Users := GetUsers()

	for _, post := range Posts {
		numberOfComments := rand.Intn(n) + 1

		for i := 0; i < numberOfComments; i++ {
			content := faker.Sentence()
			createdAt := time.Now().Format(time.RFC3339)
			_, err := db.Exec(`INSERT INTO comments (idPost, idUser, content_text, created_at) VALUES (?, ?, ?, ?)`,
				post.IDPost, Users[rand.Intn(len(Users))].IDUser, content, createdAt)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("Inserted random comments for %d posts.\n", len(Posts))
}
