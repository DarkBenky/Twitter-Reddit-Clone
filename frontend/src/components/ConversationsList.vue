<template>
    <div class="conversations-list">
        <NavBar :user="this.$store.state.currentUser"></NavBar>
        
        <!-- Add search section -->
        <div class="search-section">
            <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="Search users..."
                class="search-input"
            />
            <div v-if="searchResults.length > 0" class="search-results">
                <div 
                    v-for="user in searchResults" 
                    :key="user.idUser" 
                    @click="startConversation(user)"
                    class="search-result-item"
                >
                    <span class="username">{{ user.username }}</span>
                </div>
            </div>
        </div>

        <!-- Existing conversations list -->
        <div v-if="conversations" class="conversations-container">
            <router-link 
                v-for="conversation in conversations" 
                :key="conversation.id"
                :to="`/dm/${conversation.senderID}`"
                class="conversation-item">
                <div class="conversation-content">
                    <span v-if="conversation.senderID != this.$store.state.userId" class="username">{{ conversation.senderName }}</span>
                    <span v-else class="username">{{ conversation.receiverName }}</span>
                    <span class="last-message">{{ conversation.lastMessage }}</span>
                </div>
            </router-link>
        </div>
        <div v-else class="no-conversations">
            No conversations yet
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import NavBar from './NavBar.vue';

export default {
    components: {
        NavBar
    },
    name: 'ConversationsList',
    data() {
        return {
            baseUrl: "http://localhost:5533",
            conversations: [],
            searchQuery: '',
            searchResults: []
        };
    },
    methods: {
        async getConversations() {
            try {
                const response = await axios.get(`${this.baseUrl}/conversations`, {
                    params: {
                        userID: this.$store.state.userId
                    }
                });
                this.conversations = response.data;
                console.log('Conversations:', this.conversations);
            } catch (error) {
                console.error('Error fetching conversations:', error);
            }
        },
        getUserWithId(id) {
            return this.$store.state.users.find(user => user.idUser === id);
        },
        searchUsers() {
            if (!this.searchQuery) {
                this.searchResults = [];
                return;
            }
            const query = this.searchQuery.toLowerCase();
            this.searchResults = this.$store.state.users.filter(user => 
                user.username.toLowerCase().includes(query) && 
                user.idUser !== this.$store.state.userId
            );
        },
        startConversation(user) {
            this.$router.push(`/dm/${user.idUser}`);
            this.searchQuery = '';
            this.searchResults = [];
        }
    },
    watch: {
        searchQuery() {
            this.searchUsers();
        }
    },
    mounted() {
        this.getConversations();
    }
};
</script>

<style scoped>
.conversations-list {
    margin: 0 auto;
    padding: 20px;
}

.conversations-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.conversation-item {
    padding: 15px;
    border-radius: 8px;
    background-color: white;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    text-decoration: none;
    color: inherit;
    transition: transform 0.2s ease;
}

.conversation-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.conversation-content {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.username {
    font-weight: bold;
    color: #333;
}

.last-message {
    color: #666;
    font-size: 0.9em;
}

.no-conversations {
    text-align: center;
    padding: 20px;
    color: #666;
    font-style: italic;
}

.search-section {
    margin-bottom: 20px;
    position: relative;
}

.search-input {
    width: 100%;
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
}

.search-results {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    margin-top: 5px;
    max-height: 200px;
    overflow-y: auto;
    z-index: 1000;
}

.search-result-item {
    padding: 10px;
    cursor: pointer;
    transition: background-color 0.2s;
}

.search-result-item:hover {
    background-color: #f5f5f5;
}
</style>