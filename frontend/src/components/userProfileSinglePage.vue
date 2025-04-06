<template>
    <div class="page-container">
        <NavBar :user="user"></NavBar>
        
        <!-- Main content grid -->
        <div class="main-content">
            <!-- Left sidebar with user info -->
            <div class="profile-section">
                <div v-if="loading" class="loading-indicator">Loading...</div>
                <div v-else-if="error" class="error-message">{{ error }}</div>
                <div v-else class="profile-card">
                    <div class="profile-header">
                        <div class="avatar">{{ user.displayName?.charAt(0) || '?' }}</div>
                        <h2>{{ user.displayName }}</h2>
                    </div>
                    
                    <div class="user-info" v-if="!isEditing && !isEditingPassword">
                        <div class="info-item">
                            <span class="label">Username:</span>
                            <span class="value">{{ user.username }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">Email:</span>
                            <span class="value">{{ user.email }}</span>
                        </div>
                        
                        <div class="action-buttons">
                            <button class="edit-btn" @click="toggleEdit">Edit Profile</button>
                            <button class="password-btn" @click="toggleEditPassword">Change Password</button>
                        </div>
                    </div>

                    <!-- Edit forms remain the same -->
                    <div class="edit-user-info" v-else-if="isEditing">
                        <h2>Edit Profile</h2>
                        <form @submit.prevent="updateProfile">
                            <div class="form-group">
                                <label for="username">Username:</label>
                                <input
                                    type="text"
                                    id="username"
                                    v-model="editableUser.username"
                                    required
                                />
                            </div>
                            <div class="form-group">
                                <label for="displayName">Display Name:</label>
                                <input
                                    type="text"
                                    id="displayName"
                                    v-model="editableUser.displayName"
                                    required
                                />
                            </div>
                            <div class="form-group">
                                <label for="email">Email:</label>
                                <input
                                    type="email"
                                    id="email"
                                    v-model="editableUser.email"
                                    required
                                />
                            </div>
                            <button type="submit" :disabled="updating">
                                {{ updating ? 'Updating...' : 'Save Changes' }}
                            </button>
                            <button type="button" @click="toggleEdit" :disabled="updating">Cancel</button>
                            <p v-if="updateError" class="error">{{ updateError }}</p>
                        </form>
                    </div>

                    <div class="edit-user-info" v-else-if="isEditingPassword">
                        <h2>Edit Profile</h2>
                        <form @submit.prevent="updatePassword">
                            <div class="form-group">
                                <label for="password">Password:</label>
                                <input
                                    type="password"
                                    id="password"
                                    v-model="editableUser.password"
                                    required
                                />
                            </div>
                            <button type="submit" :disabled="updating">
                                {{ updating ? 'Updating...' : 'Save Changes' }}
                            </button>
                            <button type="button" @click="toggleEditPassword" :disabled="updating">Cancel</button>
                            <p v-if="updateError" class="error">{{ updateError }}</p>
                        </form>
                    </div>
                </div>
            </div>

            <!-- Right content area -->
            <div class="content-section">
                <!-- Posts Section -->
                <section class="posts-section">
                    <h2>My Posts</h2>
                    <div v-if="loadingPosts" class="loading-indicator">Loading...</div>
                    <div v-else-if="postsError" class="error-message">{{ postsError }}</div>
                    <div v-else class="posts-grid">
                        <PostView 
                            v-for="post in userPosts" 
                            :key="post.idPost" 
                            :post="post" 
                            :user="user" 
                            :users="users"
                            @post-deleted="removePost"
                            class="post-item"
                        />
                    </div>
                </section>

                <!-- Saved Posts Section -->
                <section class="saved-posts-section">
                    <h2>Saved Posts</h2>
                    <div v-if="loadingPostsSaved" class="loading-indicator">Loading...</div>
                    <div v-else-if="postsErrorSaved" class="error-message">{{ postsError }}</div>
                    <div v-else class="posts-grid">
                        <PostView 
                            v-for="post in savedPosts" 
                            :key="post.idPost" 
                            :post="post" 
                            :user="user" 
                            :users="users"
                            @post-deleted="removePost"
                            class="post-item"
                        />
                    </div>
                </section>

                <!-- Subscriptions Section -->
                <section class="subscriptions-section">
                    <h2>Subscriptions</h2>
                    <div class="subscriptions-grid">
                        <UserProfile 
                            v-for="idUser in subscriptions" 
                            :key="idUser"
                            :user="users.find(user => user.idUser === idUser)" 
                            :compact="false" 
                            :expanded="false" 
                            :enableSubscribe="false"
                            class="subscription-item"
                        />
                    </div>
                </section>
            </div>
        </div>
    </div>
</template>

<script>

import axios from 'axios'
import NavBar from './NavBar.vue';
import PostView from './Post.vue';
import UserProfile from './UserProfile.vue';

export default {
    name: 'UserProfileSinglePage',
    components: {
        NavBar,
        PostView,
        UserProfile
    },
    data() {
        return {
            user: {},
            editableUser: {},
            loading: true,
            error: null,
            userPosts: [],
            loadingPosts: false,
            loadingPostsSaved: false,
            postsError: null,
            postsErrorSaved: null,
            users: [],
            baseUrl: "http://localhost:5533",
            isEditing: false,
            updating: false,
            updateError: null,
            isEditingPassword: false,
            savedPosts: [],
            subscriptions: [],
        }
    },

    methods: {
        async GetListOfSubscriptions() {
            try {
                const response = await axios.get(`${this.baseUrl}/listOfSubscribers`, {
                    params: { userID: this.$store.state.userId }
                })
                this.subscriptions = response.data
                console.log(this.subscriptions, "subscriptions")
            } catch (error) {
                console.error('Error fetching subscriptions:', error)
                this.postsError = 'Failed to load subscriptions'
            } finally {
                this.loadingPosts = false
            }
        },

        async GetSavedPosts() {
            try {
                const response = await axios.get(`${this.baseUrl}/savedPosts`, {
                    params: { userID: this.$store.state.userId }
                })
                this.savedPosts = response.data
                console.log(this.savedPosts, "saved Posts")
            } catch (error) {
                console.error('Error fetching saved posts:', error)
                this.postsErrorSaved = 'Failed to load saved posts'
            } finally {
                this.loadingPostsSaved = false
            }
        },

        toggleEdit() {
            this.isEditing = !this.isEditing
            if (this.isEditing) {
                // Create a copy of the user data for editing
                this.editableUser = { ...this.user }
            }
        },
        toggleEditPassword() {
            this.isEditingPassword = !this.isEditingPassword
            if (this.isEditing) {
                // Create a copy of the user data for editing
                this.editableUser = { ...this.user }
            }
        },

        async updateProfile() {
            this.updating = true
            this.updateError = null
            try {
                const response = await axios.put(`${this.baseUrl}/userEdit`, {
                    id: this.user.idUser,
                    username: this.editableUser.username,
                    displayName: this.editableUser.displayName,
                    email: this.editableUser.email,
                })
                // Update local user data and Vuex store
                this.user = response.data
                this.$store.commit('setCurrentUser', response.data)
                this.isEditing = false
            } catch (error) {
                console.error('Error updating profile:', error)
                this.updateError = error.response?.data?.error || 'Failed to update profile. Please try again.'
            } finally {
                this.updating = false
            }
        },
        async updatePassword() {
            this.updating = true
            this.updateError = null
            try {
                const response = await axios.post(`${this.baseUrl}/updatePassword`, {
                    userID: this.user.idUser,
                    password: this.editableUser.password
                })
                // Update local user data and Vuex store
                this.user = response.data
                this.$store.commit('setCurrentUser', response.data)
                this.isEditing = false
            } catch (error) {
                console.error('Error updating profile:', error)
                this.updateError = error.response?.data?.error || 'Failed to update profile. Please try again.'
            } finally {
                this.updating = false
            }
        },

        removePost(postId) {
            this.userPosts = this.userPosts.filter(post => post.idPost !== postId);
        }
    },
    async created() {
        if (this.$store.state.userId == -1) {
            this.$router.push({ path: '/' });
        }
        try {
            // Fetch user data from the API
            const response = await axios.get(`${this.baseUrl}/user`, {
                params: { id: this.$store.state.userId }
            })
            this.user = response.data
            await this.GetSavedPosts()
        } catch (error) {
            console.error('Error fetching user data:', error)
            this.error = 'Failed to load user data'
        } finally {
            this.loading = false
        }

        try {
            // Fetch user posts from the API
            const response = await axios.get(`${this.baseUrl}/posts/user`, {
                params: { id: this.$store.state.userId }
            })
            this.userPosts = response.data
            console.log(this.userPosts)
        } catch (error) {
            console.error('Error fetching user posts:', error)
            this.postsError = 'Failed to load user posts'
        } finally {
            this.loadingPosts = false
        }

        try {
            const response = await axios.get(`${this.baseUrl}/users`);
            this.users = response.data;
            console.log(this.users)
        } catch (error) {
            console.error('Error fetching users:', error);
        }
        this.GetListOfSubscriptions()
    },
}
</script>

<style scoped>
.user-profile {
    max-width: 800px;
    margin: 0 auto;
    padding: 1rem;
    background-color: #f9f9f9;
    border-radius: 8px;
    border: 1px solid #ddd;
}

.user-profile-nav {
    margin-bottom: 1em;
    padding: 0.5em 1em;
    background-color: #fff;
    border-bottom: 1px solid #eee;
    border-radius: 8px 8px 0 0;
}

.user-info, .edit-user-info {
    padding: 1em;
    background-color: #f9f9f9;
}

.user-info p {
    margin: 0.5em 0;
    color: #333;
}

.edit-user-info .form-group {
    margin-bottom: 1rem;
}

.edit-user-info label {
    display: block;
    margin-bottom: 0.5rem;
    color: #555;
}

.edit-user-info input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
}

.edit-user-info button {
    margin-top: 1rem;
    padding: 0.5rem 1rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    margin-right: 0.5rem;
}

.edit-user-info button:disabled {
    background-color: #a5d6a7;
    cursor: not-allowed;
}

.edit-user-info .error {
    margin-top: 1rem;
    color: #d32f2f;
    text-align: center;
}

h2 {
    margin: 0;
    color: #333;
}

.loading {
    padding: 1em;
    text-align: center;
    color: #666;
    font-style: italic;
}

.error {
    padding: 1em;
    text-align: center;
    color: #d32f2f;
    background-color: #ffebee;
    border-radius: 4px;
}

button {
    padding: 0.5rem 1rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

button:hover {
    background-color: #45a049;
}

.page-container {
    min-height: 100vh;
    background-color: #f5f8fa;
}

.main-content {
    max-width: 1200px;
    margin: 20px auto;
    padding: 0 20px;
    display: grid;
    grid-template-columns: 300px 1fr;
    gap: 20px;
}

.profile-section {
    position: sticky;
    top: 20px;
    height: fit-content;
}

.profile-card {
    background: white;
    border-radius: 15px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    padding: 20px;
}

.profile-header {
    text-align: center;
    margin-bottom: 20px;
}

.avatar {
    width: 80px;
    height: 80px;
    background: #1da1f2;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 2em;
    margin: 0 auto 15px;
}

.info-item {
    display: flex;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid #eee;
}

.label {
    color: #536471;
    font-weight: 500;
}

.action-buttons {
    display: flex;
    gap: 10px;
    margin-top: 20px;
}

.edit-btn, .password-btn {
    flex: 1;
    padding: 8px;
    border-radius: 20px;
    border: none;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.2s ease;
}

.edit-btn {
    background: #1da1f2;
    color: white;
}

.password-btn {
    background: #e8f5fe;
    color: #1da1f2;
}

.content-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.posts-section, .saved-posts-section, .subscriptions-section {
    background: white;
    border-radius: 15px;
    padding: 20px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.posts-grid {
    display: grid;
    gap: 15px;
    margin-top: 15px;
}

.subscriptions-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
    margin-top: 15px;
}

.loading-indicator {
    text-align: center;
    padding: 20px;
    color: #536471;
}

.error-message {
    color: #dc3545;
    background: #ffebee;
    padding: 10px;
    border-radius: 8px;
    margin: 10px 0;
}

@media (max-width: 768px) {
    .main-content {
        grid-template-columns: 1fr;
    }

    .profile-section {
        position: static;
    }
}
</style>