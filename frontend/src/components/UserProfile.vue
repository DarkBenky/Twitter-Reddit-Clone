<template>
    <div :class="['user-profile', { 'compact': compact, 'expanded': expanded }]">
        <div class="profile-content">
            <div class="avatar" @click="moveToDMs">
                {{ initials }}
            </div>
            
            <!-- Updated button with class -->
            <button @click="showUserProfile" class="profile-button">
                <span class="button-text">View Profile</span>
            </button>

            <div v-if="this.$store.state.userId != -1 && enableSubscribe && initials != '?'">
                <button v-if="this.$store.state.userId != -1 && enableSubscribe"
                    :class="['subscribe-button', { subscribed: isSubscribed }]" @click="SubscribeUnsubscribe">
                    {{ isSubscribed ? 'Unsubscribe' : 'Subscribe' }}
                </button>
            </div>

            <div class="user-info">
                <div v-if="!compact" class="username">
                    <p v-if="user">{{ user.displayName }}</p>
                    <div v-if="expanded">
                        <p>Username: {{ user.username }}</p>
                        <p>Email: {{ user.email }}</p>
                        <!-- <p>Display Name: {{ user.displayName }}</p> -->
                        <h1>Posts</h1>
                        <div>
                            <div v-for="post in usersPosts" :key="post.idPost" class="post">
                                <p>{{ post.content_text }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- Close Button for Full-Screen Mode -->
        <button v-if="expanded" class="close-button" @click="$emit('close')">Close</button>
    </div>
</template>

<script>
import axios from 'axios'
import api from '../services/api.js'

export default {
    name: 'UserProfile',

    props: {
        user: {
            type: Object,
            required: true
        },
        compact: {
            type: Boolean,
            default: false
        },
        expanded: {
            type: Boolean,
            default: false
        },
        enableSubscribe: {
            type: Boolean,
            default: true,
            required: false
        },
    },

    data() {
        return {
            url: "http://localhost:5533",
            usersPosts: [],
            SubscribedToUserOfPost: false,
            subscriptionCheckInterval: null, // Add this line
        }
    },

    computed: {
        initials() {
            if (!this.user || !this.user.displayName) return '?'
            return this.user.displayName
                .split(' ')
                .map(word => word[0])
                .join('')
                .toUpperCase()
                .slice(0, 2)
        },
        isSubscribed() {
            return this.SubscribedToUserOfPost;
        }
    },

    watch: {
        expanded: {
            immediate: true,
            handler(newVal) {
                if (newVal) {
                    this.getUserPosts()
                }
            }
        }
    },

    created() {
        this.CheckIfUserIsSubscribed();
        // Set up periodic checking every 30 seconds
        this.subscriptionCheckInterval = setInterval(() => {
            this.CheckIfUserIsSubscribed();
        }, 3000);
    },

    unmounted() {
        // Clean up the interval when component is destroyed
        if (this.subscriptionCheckInterval) {
            clearInterval(this.subscriptionCheckInterval);
        }
    },

    methods: {
        showUserProfile() {
            console.log('Showing user profile:', this.user)
            // move to path user/id
            this.$router.push(`/user/${this.user.idUser}`)
        },

        async SubscribeUnsubscribe() {
            if (this.$store.state.userId === -1) {
                console.error("User not logged in");
                return;
            }

            try {
                const response = await api.get(`${this.url}/subscribe`, {
                    params: {
                        subscribedToID: this.user.idUser,
                        subscriberID: this.$store.state.userId
                    }
                });

                // Show success message or handle response
                console.log(response.data.message);
                this.CheckIfUserIsSubscribed(); // Refresh subscription status
            } catch (error) {
                console.error("Error changing subscription status:", error);
            }
        },

        async CheckIfUserIsSubscribed() {
            if (this.$store.state.userId === -1 || !this.user) return;
            try {
                const response = await api.get(`${this.url}/checkSubscription`, {
                    params: {
                        subscribedToID: this.user.idUser,
                        subscriberID: this.$store.state.userId
                    }
                });
                console.log('Subscription check response:', response.data); // Add this log
                this.SubscribedToUserOfPost = response.data.subscribed;
            } catch (error) {
                console.error("Error checking subscription:", error);
                return false;
            }
        },

        moveToDMs() {
            console.log('Moving to DMs with user:', this.user)
            this.$router.push('dm/' + this.user.idUser)
        },

        async getUserPosts() {
            if (!this.expanded) return  // Prevent fetching if not expanded
            try {
                const response = await axios.get(`${this.url}/posts/user`, {
                    params: {
                        id: this.user.id
                    }
                })
                this.usersPosts = response.data
            } catch (error) {
                console.error('Error fetching user posts:', error)
                this.usersPosts = []
            }
        }
    }
}
</script>

<style scoped>
.user-profile {
    display: flex;
    align-items: center;
    padding: 0.5em;
    transition: all 0.3s ease;
}

.profile-content {
    display: flex;
    align-items: center;
    gap: 0.5em;
}

.avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: #2196f3;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 1em;
}

.compact .avatar {
    width: 30px;
    height: 30px;
    font-size: 0.8em;
}

.user-info {
    display: flex;
    flex-direction: column;
}

.display-name {
    font-weight: bold;
    color: #333;
}

.username {
    color: #666;
    font-size: 0.9em;
}

.compact .user-info {
    font-size: 0.9em;
}

.loading,
.error {
    font-size: 0.9em;
    color: #666;
}

.error {
    color: #d32f2f;
}

/* Full-screen mode */
.expanded {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.8);
    color: white;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.expanded .profile-content {
    flex-direction: column;
    align-items: center;
}

.expanded .avatar {
    width: 100px;
    height: 100px;
    font-size: 2em;
}

.expanded .display-name {
    font-size: 1.5em;
}

.expanded .username {
    font-size: 1.2em;
    color: #aaa;
}

/* Close button for full-screen mode */
.close-button {
    position: absolute;
    top: 2rem;
    right: 3rem;
    background: none;
    border: none;
    color: #fff;
    font-size: 1.5rem;
    cursor: pointer;
}

.subscribe-button {
    padding: 8px 16px;
    border-radius: 20px;
    border: none;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    min-width: 100px;
    margin: 0 8px;
}

.subscribe-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Styling for Subscribe button */
.subscribe-button:not(.subscribed) {
    background-color: #1da1f2;
    color: white;
}

.subscribe-button:not(.subscribed):hover {
    background-color: #1991da;
}

/* Styling for Unsubscribe button */
.subscribe-button.subscribed {
    background-color: #e8f5fe;
    color: #1da1f2;
    border: 1px solid #1da1f2;
}

.subscribe-button.subscribed:hover {
    background-color: #f8f9fa;
}

/* Disabled state */
.subscribe-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
}

.profile-button {
    padding: 8px 12px;
    border-radius: 20px;
    border: none;
    background-color: #f0f2f5;
    color: #1a1a1a;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 8px;
    font-size: 0.9rem;
}

.profile-button:hover {
    background-color: #e4e6eb;
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.profile-button:active {
    transform: translateY(0);
}

.button-text {
    margin-left: 4px;
}

/* For compact mode */
.compact .profile-button {
    padding: 6px 10px;
    font-size: 0.8rem;
}

/* For expanded mode (dark theme) */
.expanded .profile-button {
    background-color: #3a3b3c;
    color: #e4e6eb;
}

.expanded .profile-button:hover {
    background-color: #4a4b4c;
}
</style>
