<template>
    <div class="conversations-list">
        <NavBar :user="this.$store.state.currentUser"></NavBar>
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
            baseUrl: "http://localhost:5050",
            conversations: []
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
</style>