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
	IDLikeDislike int  `json:"idLikeDislike"`
	IDPost        int  `json:"idPost"`
	IDUser        int  `json:"idUser"`
	Like          int `json:"like"`
}

type Message struct {
	IDMessage  int    `json:"idMessage"`
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

func GetAllPosts(c echo.Context) error {
	// Get all posts from the database
	query := `SELECT idPost, content_text, created_at, userID FROM posts`
	rows, err := db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query posts"})
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.IDPost, &post.ContentText, &post.CreatedAt, &post.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan post data"})
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error iterating over rows"})
	}

	// Return posts as JSON response
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

	// Get all posts from the database
	query := `SELECT idPost, content_text, created_at, userID FROM posts WHERE userID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query posts"})
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.IDPost, &post.ContentText, &post.CreatedAt, &post.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan post data"})
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error iterating over rows"})
	}

	// Return posts as JSON response
	return c.JSON(http.StatusOK, posts)
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

func AddPost(c echo.Context) error {
	type PostRequest struct {
		UserID      string `json:"userID"`
		ContentText string `json:"contentText"`
	}

	post := new(PostRequest)
	if err := c.Bind(post); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request data",
		})
	}

	if post.UserID == "" || post.ContentText == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "UserID and content text are required",
		})
	}

	query := `INSERT INTO posts (userID, content_text, created_at) VALUES (?, ?, ?)`
	result, err := db.Exec(query, post.UserID, post.ContentText, time.Now().Format(time.RFC3339))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to insert post: " + err.Error(),
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to confirm post insertion",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post added successfully",
	})
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

// In main.go, update the EditPost function:
func EditPost(c echo.Context) error {
	// Create struct for request body
	type EditRequest struct {
		PostID      string `json:"postID"`
		ContentText string `json:"contentText"`
	}

	// Parse request body
	var req EditRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	// Validate input
	if req.PostID == "" || req.ContentText == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Post ID and content text are required",
		})
	}

	// Update post in database
	query := `UPDATE posts SET content_text = ? WHERE idPost = ?`
	result, err := db.Exec(query, req.ContentText, req.PostID)
	if err != nil {
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
	postID := c.QueryParam("id")

	// Get post from the database
	query := `SELECT idPost, content_text, created_at, userID FROM posts WHERE idPost = ?`
	row := db.QueryRow(query, postID)

	var post Post
	if err := row.Scan(&post.IDPost, &post.ContentText, &post.CreatedAt, &post.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to scan post data"})
	}

	// Return post as JSON response
	return c.JSON(http.StatusOK, post)
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
	createPostsTable(database)
	createCommentsTable(database)
	createCategoriesTable(database)
	createLikesDislikesTable(database)
	createMessagesTable(database)

	// Generate random users, posts, and comments
	// n := 100 // Number of random entries to generate

	// InsertRandomMessages(n)

	// fmt.Printf("Generating %d random users...\n", n)
	// insertRandomUsers(n)

	// n = n / 2
	// fmt.Printf("Generating %d random posts...\n", n)
	// insertRandomPosts(n)

	// n = n / 2
	// fmt.Printf("Generating %d random comments...\n", n)
	// insertRandomComments(n)

	// InsertTestUser()

	// Start the server
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.GET("/posts", GetAllPosts)
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
	e.Logger.Fatal(e.Start(":5050"))

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

// Create Posts table
func createPostsTable(db *sql.DB) {
	fmt.Println("Creating posts table")
	createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
        "idPost" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
        "content_text" TEXT,
        "created_at" TEXT,
        "userID" INTEGER,
        "categoryID" INTEGER,
        FOREIGN KEY ("userID") REFERENCES users(idUser),
        FOREIGN KEY ("categoryID") REFERENCES categories(idCategory)
    );`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
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

	for _, user := range Users {
		numberOfPosts := rand.Intn(n) + 1
		for i := 0; i < numberOfPosts; i++ {
			content := faker.Sentence()
			createdAt := time.Now().Format(time.RFC3339)
			_, err := db.Exec(`INSERT INTO posts (content_text, created_at, userID) VALUES (?, ?, ?)`,
				content, createdAt, user.IDUser)
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
