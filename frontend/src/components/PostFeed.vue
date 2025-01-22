<template>
    <div>
        <NavBar :user="getUserWithId($store.state.userId)"></NavBar>
        <div class="feed-container">
            <div class="sort-controls">
                <label>Sort by:</label>
                <select v-model="sortBy" @change="handleSort" class="sort-select">
                    <option value="newest">Newest First</option>
                    <option value="oldest">Oldest First</option>
                    <option value="mostLiked">Most Liked</option>
                    <option value="mostDisliked">Most Disliked</option>
                </select>
            </div>
            <div class="posts-section">
                <PostView v-for="post in sortedPosts" 
                         :key="post.idPost" 
                         :post="post" 
                         :user="getUserWithId(post.userID)"
                         :users="users" 
                         @post-deleted="removePost" />
            </div>
            <button 
                @click="loadMore" 
                class="load-more-btn"
                v-if="hasMorePosts"
            >
                {{ loading ? 'Loading...' : 'Load More' }}
            </button>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import NavBar from './NavBar.vue';
import PostView from './Post.vue';

export default {
    name: 'PostFeed',

    components: {
        NavBar,
        PostView
    },

    data() {
        return {
            posts: [],
            users: [],
            baseUrl: "http://localhost:5555",
            sortBy: 'newest',
            postsWithMetrics: [],
            offset: 0,
            loading: false,
            hasMorePosts: true,
            fetchInterval: null,
        };
    },

    created() {
        this.fetchPosts();
        this.fetchUsers();
    },

    mounted() {
        this.fetchPosts();
        // Set up interval to fetch posts every 2.5 seconds
        this.fetchInterval = setInterval(() => {
            this.fetchPosts();
        }, 2500);
    },

    beforeUnmount() {
        if (this.fetchInterval) {
            clearInterval(this.fetchInterval);
        }
    },

    computed: {
        sortedPosts() {
            return [...this.postsWithMetrics].sort((a, b) => {
                switch (this.sortBy) {
                    case 'newest':
                        return new Date(b.created_at) - new Date(a.created_at);
                    case 'oldest':
                        return new Date(a.created_at) - new Date(b.created_at);
                    case 'mostLiked':
                        return b.likes - a.likes;
                    case 'mostDisliked':
                        return b.dislikes - a.dislikes;
                    default:
                        return 0;
                }
            });
        }
    },

    methods: {
        removePost(postId) {
            // Remove the post from the local posts array
            this.posts = this.posts.filter(post => post.idPost !== postId);
        },
        getUserWithId(id) {
            if (this.users.length === 0) return null;
            return this.users.find(user => user.idUser === id);
        },

        async loadMore() {
            if (this.loading) return;
            
            this.loading = true;
            try {
                const response = await axios.get(`${this.baseUrl}/posts`, {
                    params: {
                        offset: this.offset
                    }
                });
                const newPosts = response.data;
                
                if (newPosts.length === 0) {
                    this.hasMorePosts = false;
                    return;
                }

                // Fetch metrics for new posts
                const postsWithMetrics = await Promise.all(
                    newPosts.map(async (post) => {
                        const metricsResponse = await axios.get(
                            `${this.baseUrl}/likesDislikes`,
                            { params: { idPost: post.idPost } }
                        );
                        return {
                            ...post,
                            likes: metricsResponse.data.likes || 0,
                            dislikes: metricsResponse.data.dislikes || 0
                        };
                    })
                );

                // Append new posts to existing posts
                this.postsWithMetrics = [...this.postsWithMetrics, ...postsWithMetrics];
                this.offset += 25; // Increment offset for next load
            } catch (error) {
                console.error('Error loading more posts:', error);
            } finally {
                this.loading = false;
            }
        },

        async fetchPosts() {
            // Only fetch if not currently loading
            if (this.loading) return;
            
            try {
                this.loading = true;
                this.offset = 0; // Reset offset when fetching initial posts
                const response = await axios.get(`${this.baseUrl}/posts`, {
                    params: {
                        offset: this.offset
                    }
                });
                const posts = response.data;
                
                // Fetch metrics for each post
                this.postsWithMetrics = await Promise.all(
                    posts.map(async (post) => {
                        const metricsResponse = await axios.get(
                            `${this.baseUrl}/likesDislikes`,
                            { params: { idPost: post.idPost } }
                        );
                        return {
                            ...post,
                            likes: metricsResponse.data.likes || 0,
                            dislikes: metricsResponse.data.dislikes || 0
                        };
                    })
                );
                this.offset = posts.length;
            } catch (error) {
                console.error('Error fetching posts:', error);
            } finally {
                this.loading = false;
            }
        },

        async fetchUsers() {
            try {
                const response = await axios.get(`${this.baseUrl}/users`);
                this.users = response.data;
                this.$store.state.users = response.data;
            } catch (error) {
                console.error('Error fetching users:', error);
            }
        },
    }
};
</script>

<style scoped>
.feed-container {
    max-width: 800px;
    margin: 2rem auto;
    padding: 1rem;
}

.sort-controls {
    background: white;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    gap: 1rem;
}

.sort-controls label {
    font-weight: 500;
    color: #333;
}

.sort-select {
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    background-color: white;
    font-size: 1rem;
    cursor: pointer;
    outline: none;
    transition: border-color 0.2s ease;
}

.sort-select:hover {
    border-color: #999;
}

.sort-select:focus {
    border-color: #0066cc;
    box-shadow: 0 0 0 2px rgba(0, 102, 204, 0.2);
}

.posts-section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.load-more-btn {
    width: 100%;
    padding: 1rem;
    margin-top: 1rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.2s ease;
}

.load-more-btn:hover {
    background-color: #45a049;
}

.load-more-btn:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
}
</style>
