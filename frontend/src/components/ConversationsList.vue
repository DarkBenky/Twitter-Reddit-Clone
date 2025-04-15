<template>
    <div class="conversations-list">
        <NavBar :user="this.$store.state.currentUser"></NavBar>
        
        <!-- New conversation alert -->
        <div v-if="showNewConversationAlert" class="new-conversation-alert">
            <span>You have new messages!</span>
            <button @click="dismissNewConversationAlert" class="dismiss-button">×</button>
        </div>
        
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
        <div v-if="conversations && conversations.length > 0" class="conversations-container">
            <router-link 
                v-for="conversation in conversations" 
                :key="conversation.id"
                :to="`/dm/${conversation.senderID !== this.$store.state.userId ? conversation.senderID : conversation.receiverID}`"
                class="conversation-item"
                @click="dismissNewConversationAlert">
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
// import axios from 'axios';
import NavBar from './NavBar.vue';
import api from '../services/api.js';

export default {
    components: {
        NavBar
    },
    name: 'ConversationsList',
    data() {
        return {
            conversations: [],
            previousConversations: [], // Store previous state for comparison
            searchQuery: '',
            searchResults: [],
            fetchInterval: null,
            showNewConversationAlert: false
        };
    },
    methods: {
        async getConversations() {
            try {
                const response = await api.get('/conversations', {
                    params: {
                        userID: this.$store.state.userId
                    }
                });
                
                // Store current conversations for later comparison
                this.previousConversations = [...this.conversations];
                
                // Get new conversations
                const newConversations = response.data;
                
                // Check for new messages by comparing with previous state
                if (this.previousConversations.length > 0) {
                    const hasNewMessages = this.detectNewMessages(this.previousConversations, newConversations);
                    
                    if (hasNewMessages) {
                        this.showNewConversationAlert = true;
                        // Save notification state to cookies
                        this.$cookies.set('new_conversation_notification', true, { expires: '30d' });
                    }
                }
                
                // Store the serialized conversations in cookies for persistence
                if (newConversations.length > 0) {
                    this.$cookies.set('last_conversations', JSON.stringify(newConversations), { expires: '30d' });
                }
                
                // Update the conversations list
                this.conversations = newConversations;
                console.log('Conversations updated:', this.conversations);
            } catch (error) {
                console.error('Error fetching conversations:', error);
            }
        },
        
        // Compare previous and current conversations to detect new messages
        detectNewMessages(previous, current) {
            // Case 1: More conversations now than before
            if (current.length > previous.length) {
                return true;
            }
            
            // Case 2: Same number of conversations, but check for new messages
            for (const newConv of current) {
                // Find matching conversation in previous state
                const prevConv = previous.find(p => 
                    (p.senderID === newConv.senderID && p.receiverID === newConv.receiverID) || 
                    (p.senderID === newConv.receiverID && p.receiverID === newConv.senderID)
                );
                
                // If no previous conversation found, or lastMessage changed, it's new
                if (!prevConv || prevConv.lastMessage !== newConv.lastMessage) {
                    return true;
                }
            }
            
            return false;
        },
        
        startPeriodicFetching(interval = 10000) {
            // Clear any existing interval first
            this.stopPeriodicFetching();
            
            // Set up new interval
            this.fetchInterval = setInterval(() => {
                this.getConversations();
            }, interval);
            
            console.log(`Started periodic fetching every ${interval/1000} seconds`);
        },
        
        stopPeriodicFetching() {
            if (this.fetchInterval) {
                clearInterval(this.fetchInterval);
                this.fetchInterval = null;
                console.log('Stopped periodic fetching');
            }
        },
        
        dismissNewConversationAlert() {
            try {
                console.log('Dismissing notification');
                this.showNewConversationAlert = false;
                // Remove the notification from cookies
                this.$cookies.remove('new_conversation_notification');
                console.log('Notification dismissed');
            } catch (error) {
                console.error('Error dismissing notification:', error);
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
            try {
                console.log('Starting conversation with user:', user);
                
                // Dismiss any notification
                this.dismissNewConversationAlert();
                
                // Clear search
                this.searchQuery = '';
                this.searchResults = [];
                
                // Navigate to conversation
                this.$router.push(`/dm/${user.idUser}`);
                console.log('Navigation successful');
            } catch (error) {
                console.error('Error in startConversation:', error);
            }
        },
        
        loadStateFromCookies() {
            try {
                console.log('Loading state from cookies');
                
                // Load notification state from cookies
                const hasNewNotification = this.$cookies.get('new_conversation_notification');
                console.log('New notification cookie value:', hasNewNotification);
                if (hasNewNotification) {
                    this.showNewConversationAlert = true;
                }
                
                // Load previous conversations
                const savedConversations = this.$cookies.get('last_conversations');
                if (savedConversations) {
                    try {
                        this.previousConversations = JSON.parse(savedConversations);
                        console.log('Loaded previous conversations from cookies:', this.previousConversations);
                    } catch (e) {
                        console.error('Error parsing saved conversations:', e);
                    }
                }
                
                console.log('State loaded from cookies successfully');
            } catch (error) {
                console.error('Error loading state from cookies:', error);
            }
        }
    },
    watch: {
        searchQuery() {
            this.searchUsers();
        }
    },
    mounted() {
        try {
            console.log('Component mounted');
            
            // Load state from cookies first
            this.loadStateFromCookies();
            
            // Then get fresh data
            this.getConversations();
            
            // Start fetching conversations every 5 seconds
            this.startPeriodicFetching(5000);
            
            console.log('Component mounted successfully');
        } catch (error) {
            console.error('Error in mounted hook:', error);
        }
    },
    beforeUnmount() {
        // Important: clean up the interval when component is destroyed
        this.stopPeriodicFetching();
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

.new-conversation-alert {
    background-color: #4CAF50;
    color: white;
    padding: 12px 15px;
    margin-bottom: 15px;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    animation: fadeIn 0.3s;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
}

.dismiss-button {
    background: none;
    border: none;
    color: white;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0 5px;
}

.dismiss-button:hover {
    opacity: 0.8;
}
</style>