<template>
    <div class="page-container">
        <NavBar :user="currentUser"></NavBar>
        
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
                        
                        <!-- Subscribe button section -->
                        <div v-if="$store.state.userId !== -1 && $store.state.userId !== user.idUser" class="subscribe-container">
                            <button 
                                :class="['subscribe-button', { subscribed: isSubscribed }]" 
                                @click="subscribeUnsubscribe"
                                :disabled="subscribing"
                            >
                                {{ subscribing ? 'Processing...' : (isSubscribed ? 'Unsubscribe' : 'Subscribe') }}
                            </button>
                        </div>

                        <!-- Styled stats cards -->
                        <div class="user-stats">
                            <div class="stat-card">
                                <div class="stat-number">{{userPosts.length}}</div>
                                <div class="stat-label">Posts</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-number">{{numberOfSubscribers}}</div>
                                <div class="stat-label">Subscribers</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-number">{{numberOfSubscribeTo}}</div>
                                <div class="stat-label">Following</div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="user-info">
                        <div class="info-item">
                            <span class="label">Username:</span>
                            <span class="value">{{ user.username }}</span>
                        </div>

                        <div class="action-buttons" v-if="$store.state.userId !== -1">
                            <button class="message-btn" @click="goToMessages">Message</button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Right content area -->
            <div class="content-section">
                <!-- Posts Section -->
                <section class="posts-section">
                    <h2>Posts</h2>
                    <div v-if="loadingPosts" class="loading-indicator">Loading...</div>
                    <div v-else-if="postsError" class="error-message">{{ postsError }}</div>
                    <div v-else-if="userPosts.length === 0" class="no-posts-message">
                        No posts to display
                    </div>
                    <div v-else class="posts-grid">
                        <PostView 
                            v-for="post in userPosts" 
                            :key="post.idPost" 
                            :post="post" 
                            :user="user" 
                            :users="users"
                            class="post-item"
                        />
                    </div>
                </section>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import NavBar from './NavBar.vue'
import PostView from './Post.vue'

export default {
    name: 'UsersProfile',
    components: {
        NavBar,
        PostView
    },
    data() {
        return {
            user: {},
            currentUser: {},
            loading: true,
            error: null,
            userPosts: [],
            loadingPosts: true,
            postsError: null,
            users: [],
            baseUrl: "http://138.68.76.63:5533",
            isSubscribed: false,
            subscribing: false,
            subscriptionCheckInterval: null,
            numberOfSubscribers: 0,
            numberOfSubscribeTo: 0
        }
    },
    methods: {
        async fetchNumberOfSubscribers() {
            const userId = this.$route.params.id
            try {
                const response = await axios.get(`${this.baseUrl}/numberOfSubscribers`, {
                    params: { userID: userId }
                })
                this.numberOfSubscribers = response.data.numberOfSubscribers
            } catch (error) {
                console.error('Error fetching number of subscribers:', error)
            }
        },
        async fetchNumberOfSubscribeTo() {
            const userId = this.$route.params.id
            try {
                const response = await axios.get(`${this.baseUrl}/numberOfSubscribeTo`, {
                    params: { userID: userId }
                })
                this.numberOfSubscribeTo = response.data.numberOfSubscriptions
            } catch (error) {
                console.error('Error fetching number of subscriptions:', error)
            }
        },

        async fetchUserData() {
            try {
                // Get the user ID from the route parameter
                const userId = this.$route.params.id
                if (!userId) {
                    this.error = 'No user ID provided'
                    return
                }

                // Fetch user data
                const response = await axios.get(`${this.baseUrl}/user`, {
                    params: { id: userId }
                })
                
                this.user = response.data

                // Fetch current logged in user if available
                if (this.$store.state.userId !== -1) {
                    const currentUserResponse = await axios.get(`${this.baseUrl}/user`, {
                        params: { id: this.$store.state.userId }
                    })
                    this.currentUser = currentUserResponse.data
                }

                // Check subscription status
                await this.checkSubscriptionStatus()
            } catch (error) {
                console.error('Error fetching user data:', error)
                this.error = 'Failed to load user data. Please try again later.'
            } finally {
                this.loading = false
            }
        },

        async fetchUserPosts() {
            try {
                this.loadingPosts = true
                const userId = this.$route.params.id
                
                const response = await axios.get(`${this.baseUrl}/posts/user`, {
                    params: { id: userId }
                })
                
                this.userPosts = response.data
            } catch (error) {
                console.error('Error fetching user posts:', error)
                this.postsError = 'Failed to load posts. Please try again later.'
            } finally {
                this.loadingPosts = false
            }
        },

        async fetchUsers() {
            try {
                const response = await axios.get(`${this.baseUrl}/users`)
                this.users = response.data
            } catch (error) {
                console.error('Error fetching users:', error)
            }
        },

        async checkSubscriptionStatus() {
            if (this.$store.state.userId === -1 || !this.user.idUser) return
            
            try {
                const response = await axios.get(`${this.baseUrl}/checkSubscription`, {
                    params: {
                        subscribedToID: this.user.idUser,
                        subscriberID: this.$store.state.userId
                    }
                })
                
                this.isSubscribed = response.data.subscribed
            } catch (error) {
                console.error('Error checking subscription status:', error)
            }
        },

        async subscribeUnsubscribe() {
            if (this.$store.state.userId === -1) {
                alert('You need to be logged in to subscribe')
                return
            }
            
            this.subscribing = true
            try {
                await axios.get(`${this.baseUrl}/subscribe`, {
                    params: {
                        subscribedToID: this.user.idUser,
                        subscriberID: this.$store.state.userId
                    }
                })
                
                // Toggle subscription status
                this.isSubscribed = !this.isSubscribed
            } catch (error) {
                console.error('Error changing subscription status:', error)
                alert('Failed to change subscription status. Please try again.')
            } finally {
                this.subscribing = false
            }
        },

        goToMessages() {
            this.$router.push(`/dm/${this.user.idUser}`)
        }
    },
    created() {
        this.fetchUserData()
        this.fetchUserPosts()
        this.fetchUsers()
        this.fetchNumberOfSubscribers()
        this.fetchNumberOfSubscribeTo()

        console.log("users data", this.user)
        
        // Set up periodic subscription checking
        this.subscriptionCheckInterval = setInterval(() => {
            this.checkSubscriptionStatus()
        }, 30000)
    },
    beforeUnmount() {
        // Clean up the interval when component is destroyed
        if (this.subscriptionCheckInterval) {
            clearInterval(this.subscriptionCheckInterval)
        }
    },
    watch: {
        // Watch for route changes to update user data when navigating between different user profiles
        '$route.params.id': {
            handler() {
                this.fetchUserData()
                this.fetchUserPosts()
            },
            immediate: true
        }
    }
}
</script>

<style scoped>
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

.message-btn {
    width: 100%;
    padding: 8px;
    border-radius: 20px;
    border: none;
    background: #1da1f2;
    color: white;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.2s ease;
}

.message-btn:hover {
    background: #1991da;
    transform: translateY(-1px);
}

.subscribe-container {
    margin: 15px 0;
}

.subscribe-button {
    padding: 8px 16px;
    border-radius: 20px;
    border: none;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    min-width: 120px;
}

.subscribe-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.subscribe-button:not(.subscribed) {
    background-color: #1da1f2;
    color: white;
}

.subscribe-button:not(.subscribed):hover {
    background-color: #1991da;
}

.subscribe-button.subscribed {
    background-color: #e8f5fe;
    color: #1da1f2;
    border: 1px solid #1da1f2;
}

.subscribe-button.subscribed:hover {
    background-color: #f8f9fa;
}

.subscribe-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
}

.content-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.posts-section {
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

.no-posts-message {
    text-align: center;
    padding: 30px;
    color: #536471;
    font-style: italic;
}

/* Styled user stats */
.user-stats {
    display: flex;
    justify-content: center;
    gap: 20px;
    margin: 25px 0;
}

.stat-card {
    padding: 12px;
    background: #f8f9fa;
    border-radius: 10px;
    width: 80px;
    text-align: center;
    box-shadow: 0 1px 3px rgba(0,0,0,0.05);
    transition: all 0.2s ease;
}

.stat-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 6px rgba(0,0,0,0.1);
}

.stat-number {
    font-size: 1.5rem;
    font-weight: 700;
    color: #1da1f2;
    margin-bottom: 4px;
}

.stat-label {
    font-size: 0.8rem;
    color: #536471;
    font-weight: 500;
}

/* Make sure responsive design still works */
@media (max-width: 768px) {
    .user-stats {
        gap: 10px;
    }
    
    .stat-card {
        width: 70px;
        padding: 10px;
    }
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