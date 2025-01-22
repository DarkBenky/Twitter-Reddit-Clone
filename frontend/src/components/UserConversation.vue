<template>
    <div class="user-conversation">
        <NavBar :user="currentUser"></NavBar>
        <p v-if="this.$store.state.userId == -1">You need to be logged in</p>
        <div class="messages-container" ref="messagesContainer">
            <div v-for="message in messages" 
                 :key="message.idMessage" 
                 :class="['message', message.senderID === $store.state.userId ? 'sent' : 'received']">
                <div class="message-content">
                    <UserProfile :user="getUserWithId(message.senderID)" />
                    <span v-html="formatMessageWithLinks(message.content)"></span>
                </div>
                <div class="message-timestamp">{{ message.created_at }}</div>
            </div>
        </div>
        <div class="container">
            <p v-if="!messages">No messages yet</p>
            <button class="scroll-btn" @click="scrollToBottom" v-if="messages">Scroll to bottom</button>
        </div>
        <div class="message-input-container">
            <input 
                type="text" 
                v-model="newMessage" 
                @keyup.enter="sendMessage"
                placeholder="Type a message..."
                class="message-input"
            />
            <button @click="sendMessage" class="send-button">Send</button>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import NavBar from './NavBar.vue';
import UserProfile from './UserProfile.vue';

export default {
    name: 'UserConversation',

    components: {
        NavBar,
        UserProfile
    },

    data() {
        return {
            baseUrl: "http://localhost:5050",
            messages: [],
            pollingInterval: null,
            intervalTime: 1000, // 1 second
            newMessage: '',
            currentUser: null
        };
    },

    methods: {
        getUserWithId(id) {
            console.log(this.$store.state.users);
            return this.$store.state.users.find(user => user.idUser === id);
        },

        async getMessages() {
            try {
                const response = await axios.get(`${this.baseUrl}/messages`, {
                    params: {
                        senderID: String(this.$store.state.userId),
                        receiverID: String(this.$route.params.id)
                    }
                });
                this.messages = response.data;
            } catch (error) {
                console.error('Error fetching messages:', error);
            }
        },

        async sendMessage() {
            if (!this.newMessage.trim()) return;
            
            try {
                await axios.post(`${this.baseUrl}/sendMessage`, {
                    senderID: String(this.$store.state.userId),
                    receiverID: String(this.$route.params.id),
                    content: this.newMessage
                });
                this.newMessage = '';
                await this.getMessages();
                this.scrollToBottom();
            } catch (error) {
                console.error('Error sending message:', error);
            }
        },

        scrollToBottom() {
            if (this.$refs.messagesContainer) {
                this.$refs.messagesContainer.scrollTop = this.$refs.messagesContainer.scrollHeight;
            }
        },

        startPolling() {
            this.pollingInterval = setInterval(this.getMessages, this.intervalTime);
        },

        formatMessageWithLinks(content) {
            const urlRegex = /(https?:\/\/[^\s]+)/g;
            return content.replace(urlRegex, url => 
                `<a href="${url}" 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    class="message-link">${url}</a>`
            );
        }
    },

    mounted() {
        this.getMessages();
        this.startPolling();
        this.currentUser = this.$store.state.currentUser;
    },

    beforeUnmount() {
        if (this.pollingInterval) {
            clearInterval(this.pollingInterval);
        }
    }
}
</script>

<style scoped>

.container {
    display: flex;
    justify-content: center;
}

.scroll-btn {
    padding: 12px 24px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 24px;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.3s ease;
    margin: 0px 0px 10px 0px;
}

.user-conversation {
    margin: 0 auto;
    padding: 20px;
    height: calc(100vh - 60px);
    display: flex;
    flex-direction: column;
    background-color: #f8f9fa;
}

.messages-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 20px;
    height: calc(100vh - 180px);
    overflow-y: auto;
    margin: 20px 0;
    border-radius: 12px;
    background-color: white;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.message {
    max-width: 70%;
    padding: 12px 16px;
    border-radius: 20px;
    position: relative;
    animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}

.sent {
    align-self: flex-end;
    background-color: #007bff;
    color: white;
    border-bottom-right-radius: 4px;
    margin-left: auto;
}

.received {
    align-self: flex-start;
    background-color: #e9ecef;
    color: #212529;
    border-bottom-left-radius: 4px;
    margin-right: auto;
}

.message-content {
    margin-bottom: 4px;
    word-wrap: break-word;
}

.message-timestamp {
    font-size: 0.75rem;
    opacity: 0.7;
}

.message-input-container {
    display: flex;
    gap: 12px;
    padding: 16px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.1);
}

.message-input {
    flex-grow: 1;
    padding: 12px 16px;
    border: 2px solid #dee2e6;
    border-radius: 24px;
    font-size: 16px;
    outline: none;
    transition: all 0.3s ease;
}

.message-input:focus {
    border-color: #007bff;
    box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.send-button {
    padding: 12px 24px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 24px;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.3s ease;
}

.send-button:hover {
    background-color: #0056b3;
    transform: translateY(-1px);
}

.send-button:active {
    transform: translateY(1px);
}

@media (max-width: 768px) {
    .user-conversation {
        padding: 10px;
    }

    .message {
        max-width: 85%;
    }

    .message-input-container {
        padding: 10px;
    }
}

:deep(.message-link) {
    color: #0066cc;
    text-decoration: underline;
    cursor: pointer;
}

:deep(.message-link:hover) {
    color: #004499;
}

.sent :deep(.message-link) {
    color: #ffffff;
    opacity: 0.9;
}

.sent :deep(.message-link:hover) {
    opacity: 1;
}
</style>