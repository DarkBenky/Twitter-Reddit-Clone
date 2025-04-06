<template>
    <div>
        <NavBar :user="getUserWithId($store.state.userId)"></NavBar>
        <div class="single-post-container">
            <div v-if="loadingPost" class="loading">
                Loading post...
            </div>

            <div v-else-if="error" class="error">
                {{ error }}
            </div>

            <div v-else-if="post" class="post">
                <!-- Post Header with User Info -->
                <div class="post-header">
                    <UserProfile :user="getUserFromId(post.userID)" />
                    <small v-if="post" class="post-date">{{ formatDate(post.created_at) }}</small>
                    <h2 v-if="category" class="category"> {{ category }}</h2>
                </div>

                <!-- Post Content -->
                <div class="post-content">
                    <!-- Edit mode -->
                    <div v-if="isEditing" class="edit-mode">
                        <textarea v-model="editedContent" class="edit-textarea"></textarea>
                        <div class="edit-buttons">
                            <button class="save-edit" @click="saveEdit">Save</button>
                            <button class="cancel-edit" @click="cancelEdit">Cancel</button>
                        </div>
                    </div>

                    <!-- View mode -->
                    <div v-else>
                        <span v-html="formatLinks(post.content_text)"></span>
                        <div class="edit-delete-options" v-if="post.userID === $store.state.userId">
                            <button class="toggle-delete" @click="deletePost(post.idPost)">
                                Delete
                            </button>
                            <button class="toggle-edit" @click="startEdit">
                                Edit
                            </button>
                            <button class="post-save-button"
                                :class="{ 'post-saved': postSaved, 'post-not-saved': !postSaved }" @click="savePost">
                                {{ postSaved ? 'Saved' : 'Save Post' }}
                            </button>
                        </div>
                    </div>

                    <div class="image-container" v-if="post.imageURL">
                        <img :src="post.imageURL" alt="Post Image" @error="handleImageError" class="post-image" />
                    </div>

                    <div class="like-dislike">
                        <button @click="like" v-if="liked" class="liked">
                            <span>Likes: {{ likes }}</span>
                        </button>
                        <button @click="like" v-else class="not-liked">
                            <span>Likes: {{ likes }}</span>
                        </button>
                        <button @click="dislike" v-if="disliked" class="disliked">
                            <span>Dislikes: {{ dislikes }}</span>
                        </button>
                        <button @click="dislike" v-else class="not-disliked">
                            <span>Dislikes: {{ dislikes }}</span>
                        </button>
                    </div>
                </div>

                <!-- Comments Section -->
                <div class="comments-section">
                    <h3>Comments</h3>

                    <!-- Add Comment Section -->
                    <div class="add-comment-section">
                        <textarea v-model="newComment" placeholder="Add a comment..."></textarea>
                        <button @click="addComment">Add Comment</button>
                    </div>

                    <div v-if="loadingComments" class="loading">
                        Loading comments...
                    </div>
                    <div v-else>
                        <ul v-if="comments" class="comments-list">
                            <li v-for="comment in comments" :key="comment.idComment" class="comment">
                                <div class="comment-header">
                                    <UserProfile :user="getUserWithId(comment.idUser)" compact />
                                </div>
                                <div class="comment-content">
                                    <span v-html="formatLinks(comment.content_text)"></span>
                                    <small class="comment-date">{{ formatDate(comment.created_at) }}</small>
                                </div>
                            </li>
                        </ul>
                        <p v-else class="no-comments">No comments yet</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import UserProfile from './UserProfile.vue';
import NavBar from './NavBar.vue';

export default {
    name: 'SinglePost',

    components: {
        UserProfile,
        NavBar
    },

    data() {
        return {
            // like and dislike variables
            likes: 0,
            dislikes: 0,
            liked: false,
            disliked: false,
            postSaved: false,
            category: null,
            post: null,
            user: null,
            comments: [],
            users: [],
            loadingPost: true,
            loadingComments: false,
            error: null,
            baseUrl: "http://localhost:5533",
            newComment: "",
            isEditing: false,
            editedContent: ""
        };
    },

    async created() {
        const postId = this.$route.params.id;
        await this.fetchPost(postId);
        await this.fetchComments(postId);
        await this.fetchUsers();
        await this.getLikesDislikes();
        await this.checkPostSaved();
        await this.getUserLikeDislike();
    },

    methods: {
        async getLikesDislikes() {
            try {
                const response = await axios.get(`${this.baseUrl}/likesDislikes`, {
                    params: {
                        idPost: this.post.idPost
                    }
                });
                this.likes = response.data.likes || 0;
                this.dislikes = response.data.dislikes || 0;
            } catch (error) {
                console.error("Error:", error);
                this.likes = 0;
                this.dislikes = 0;
                this.liked = false;
                this.disliked = false;
            }
        },

        async savePost() {
            try {
                const response = await axios.post(`${this.baseUrl}/savePost`, {
                    postID: String(this.post.idPost),
                    userID: String(this.$store.state.userId)
                });

                if (response.status === 200) {
                    this.postSaved = !this.postSaved;
                    console.log("Post saved status toggled successfully");
                }
            } catch (error) {
                console.error("Error saving post:", error);
            }
        },

        async checkPostSaved() {
            try {
                const response = await axios.get(`${this.baseUrl}/checkPostSaved`, {
                    params: {  
                        postID: String(this.post.idPost),
                        userID: String(this.$store.state.userId)
                    }
                });

                this.postSaved = response.data.saved;
            } catch (error) {
                console.error("Error checking if post is saved:", error);
            }
        },

        async like() {
            try {
                if (this.$store.state.userId === -1) {
                    alert("Please login to like posts");
                    return;
                }

                const response = await axios.get(`${this.baseUrl}/like`, {
                    params: {
                        postId: this.post.idPost,
                        userId: this.$store.state.userId
                    }
                });

                if (response.status === 200) {
                    await this.getLikesDislikes();
                    await this.getUserLikeDislike();
                }
            } catch (error) {
                console.error("Error liking post:", error);
            }
        },

        async dislike() {
            try {
                if (this.$store.state.userId === -1) {
                    alert("Please login to dislike posts");
                    return;
                }

                const response = await axios.get(`${this.baseUrl}/dislike`, {
                    params: {
                        postId: this.post.idPost,
                        userId: this.$store.state.userId
                    }
                });

                if (response.status === 200) {
                    await this.getLikesDislikes();
                    await this.getUserLikeDislike();
                }
            } catch (error) {
                console.error("Error disliking post:", error);
            }
        },

        async getUserLikeDislike() {
            try {
                if (this.$store.state.userId === -1) return;

                const response = await axios.get(`${this.baseUrl}/userLikeDislike`, {
                    params: {
                        idPost: this.post.idPost,
                        userID: this.$store.state.userId
                    }
                });

                const userLikeDislike = response.data;
                if (userLikeDislike.like === 1) {
                    this.liked = true;
                    this.disliked = false;
                } if (userLikeDislike.like === -1) {
                    this.liked = false;
                    this.disliked = true;
                }
            } catch (error) {
                this.liked = false;
                this.disliked = false;
                console.error("Error getting user like/dislike status:", error);
            }
        },


        async deletePost(postID) {
            try {
                const response = await axios.delete(`${this.baseUrl}/deletePost`, {
                    data: { postID: String(postID) }
                });

                if (response.status === 200) {
                    this.$router.push('/');
                }
            } catch (error) {
                console.error("Error deleting post:", error);
                alert("Failed to delete post. Please try again.");
            }
        },

        editPost() {
            this.$router.push(`/edit/${this.post.idPost}`);
        },
        async addComment() {
            // Ensure there is content in the new comment field
            if (!this.newComment.trim()) {
                alert("Comment content cannot be empty.");
                return;
            }

            try {
                // Make a POST request to the backend with the new comment data
                const response = await axios.post(`${this.baseUrl}/addComment`, {
                    postID: String(this.post.idPost),
                    userID: String(this.$store.state.userId),
                    contentText: this.newComment
                });

                // If the request was successful, fetch the updated list of comments
                if (response.status === 200) {
                    this.newComment = ""; // Clear the comment input field
                    await this.fetchComments(this.post.idPost); // Refresh the comments
                }
            } catch (error) {
                console.error("Error adding comment:", error);
            }
        },

        getUserFromId(id) {
            return this.users.find(user => user.idUser === id);
        },

        async fetchPost(postId) {
            try {
                const response = await axios.get(`${this.baseUrl}/post`, {
                    params: {
                        id: postId
                    }
                })
                this.post = await response.data;
                // this.fetchUser(this.post.userID);
                this.user = this.getUserFromId(this.post.userID);
                this.category = this.post.category;
            } catch (error) {
                this.error = 'Failed to load post: ' + error.message;
                console.error('Error fetching post:', error);
            } finally {
                this.loadingPost = false;
            }
        },

        async fetchUser(userId) {
            try {
                const response = await axios.get(`${this.baseUrl}/users/${userId}`);
                this.user = response.data;
            } catch (error) {
                console.error('Error fetching user:', error);
            }
        },

        async fetchUsers() {
            try {
                const response = await axios.get(`${this.baseUrl}/users`);
                this.users = response.data;
            } catch (error) {
                console.error('Error fetching users:', error);
            }
        },

        async fetchComments(postId) {
            this.loadingComments = true;
            try {
                const response = await axios.get(`${this.baseUrl}/comments`, {
                    params: { idPost: postId }
                });
                this.comments = response.data;
            } catch (error) {
                console.error('Error fetching comments:', error);
            } finally {
                this.loadingComments = false;
            }
        },

        getUserWithId(id) {
            return this.users.find(user => user.idUser === id);
        },

        formatDate(dateString) {
            try {
                const date = new Date(dateString);
                return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
            } catch (err) {
                return dateString;
            }
        },

        startEdit() {
            this.isEditing = true;
            this.editedContent = this.post.content_text;
        },

        cancelEdit() {
            this.isEditing = false;
            this.editedContent = "";
        },

        async saveEdit() {
            try {
                const response = await axios.put(`${this.baseUrl}/editPost`, {
                    postID: String(this.post.idPost),
                    contentText: this.editedContent
                });

                if (response.status === 200) {
                    this.post.content_text = this.editedContent;
                    this.isEditing = false;
                }
            } catch (error) {
                console.error("Error saving edit:", error);
            }
        },

        formatLinks(text) {
            if (!text) return '';
            const urlRegex = /(https?:\/\/[^\s]+)/g;
            return text.replace(urlRegex, url =>
                `<a href="${url}" 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    class="content-link">${url}</a>`
            );
        }
    }
};
</script>

<style scoped>
.image-container {
    width: 100%;
    max-width: 600px;
    margin: 1rem auto;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s ease;
}

.image-container:hover {
    transform: scale(1.02);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.post-image {
    width: 100%;
    height: auto;
    max-height: 500px;
    object-fit: cover;
    display: block;
    transition: opacity 0.3s ease;
}

.post-image:hover {
    opacity: 0.95;
}

.liked {
    border: 1px solid #666;
    background-color: #28a745;
    color: white;
    border-radius: 3px;
}

.not-liked {
    border: 1px solid #666;
    background-color: #fff;
    color: #333;
    border-radius: 3px;
}

.disliked {
    border: 1px solid #666;
    background-color: #dc3545;
    color: white;
    border-radius: 3px;
}

.not-disliked {
    border: 1px solid #666;
    background-color: #fff;
    color: #333;
    border-radius: 3px;
}

.add-comment-section {
    display: flex;
    flex-direction: column;
    margin-bottom: 1.5rem;
}

.add-comment-section textarea {
    min-height: 80px;
    padding: 10px;
    margin: 15px;
    border-radius: 4px;
    border: 1px solid #ddd;
    resize: vertical;
    font-size: 1rem;
    font-family: inherit;
    margin-bottom: 0.5rem;
    transition: border-color 0.2s ease-in-out;
}

.add-comment-section textarea:focus {
    border-color: #999;
    outline: none;
    padding: 10px;
    margin: 15px;
}

.add-comment-section button {
    align-self: flex-end;
    padding: 0.5rem 1rem;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.2s ease-in-out;
}

.add-comment-section button:hover {
    background-color: #0056b3;
}

.add-comment-section button:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
}

.single-post-container {
    max-width: 800px;
    margin: 2rem auto;
    padding: 1rem;
}

.post {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
}

.post-header {
    padding: 1rem;
    border-bottom: 1px solid #eee;
}

.post-content {
    padding: 1.5rem;
}

.post-date {
    display: block;
    color: #666;
    margin-top: 0.5rem;
}

.comments-section {
    padding: 1.5rem;
    background: #f9f9f9;
}

.comments-section h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    color: #333;
}

.comment {
    padding: 1rem;
    border-bottom: 1px solid #eee;
    background: white;
    margin-bottom: 0.5rem;
    border-radius: 4px;
}

.comment-header {
    margin-bottom: 0.5rem;
}

.comment-content {
    padding-left: 1rem;
}

.comment:last-child {
    border-bottom: none;
}

.loading {
    text-align: center;
    padding: 2rem;
    color: #666;
}

.error {
    color: #d32f2f;
    background-color: #ffebee;
    padding: 1rem;
    border-radius: 4px;
    margin: 1rem 0;
}

.no-comments {
    text-align: center;
    color: #666;
    font-style: italic;
    padding: 1rem;
}

.edit-delete-options {
    margin-top: 1rem;
    display: flex;
    gap: 10px;
}

.toggle-delete,
.toggle-edit {
    padding: 0.3em 0.8em;
    background-color: #e0e0e0;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
}

.toggle-delete:hover {
    background-color: #f44336;
    color: white;
}

.toggle-edit:hover {
    background-color: #4CAF50;
    color: white;
}

.edit-mode {
    display: flex;
    flex-direction: column;
}

.edit-textarea {
    min-height: 100px;
    padding: 10px;
    margin-bottom: 1rem;
    border-radius: 4px;
    border: 1px solid #ddd;
    resize: vertical;
    font-size: 1rem;
    font-family: inherit;
    transition: border-color 0.2s ease-in-out;
}

.edit-textarea:focus {
    border-color: #999;
    outline: none;
}

.edit-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}

.save-edit,
.cancel-edit {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.2s ease-in-out;
}

.save-edit {
    background-color: #4CAF50;
    color: white;
}

.save-edit:hover {
    background-color: #45a049;
}

.cancel-edit {
    background-color: #e0e0e0;
}

.cancel-edit:hover {
    background-color: #d4d4d4;
}

.category {
    display: inline-block;
    background-color: #007bff;
    color: white;
    padding: 0.3rem 0.8rem;
    border-radius: 16px;
    font-size: 0.9rem;
    font-weight: 500;
    margin: 0.5rem 0;
    transition: background-color 0.2s ease;
}

.category:hover {
    background-color: #007bff6c;
}

:deep(.content-link) {
    color: #0066cc;
    text-decoration: underline;
    cursor: pointer;
    word-break: break-word;
}

:deep(.content-link:hover) {
    color: #004499;
    text-decoration: underline;
}

.comment :deep(.content-link) {
    font-size: 0.9em;
}

.post-save-button {
    padding: 0.3em 0.8em;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 500;
}

.post-saved {
    background-color: #28a745;
    color: white;
}

.post-saved:hover {
    background-color: #218838;
}

.post-not-saved {
    background-color: #e0e0e0;
    color: #333;
}

.post-not-saved:hover {
    background-color: #d4d4d4;
}

/* You can add a bookmark icon with pseudo-element if desired */
.post-saved::before {
    content: "★";
    /* Filled star */
    font-size: 1rem;
}

.post-not-saved::before {
    content: "☆";
    /* Empty star */
    font-size: 1rem;
}
</style>