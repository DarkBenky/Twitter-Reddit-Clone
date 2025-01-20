<template>
    <div class="user-conversation">
        <h1>Messages</h1>
        <p v-if="this.$store.state.userId == -1">You need to be logged in</p>
        <div class="messages-container">
            <div v-for="message in messages" 
                 :key="message.idMessage" 
                 :class="['message', message.senderID === $store.state.userId ? 'sent' : 'received']">
                <div class="message-content">{{ message.content }}</div>
                <div class="message-timestamp">{{ message.created_at }}</div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'UserConversation',
    data() {
        return {
            baseUrl: "http://localhost:5050",
            messages: [],
            pollingInterval: null,
            intervalTime: 1000 // 1 second
        };
    },
    methods: {
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
        startPolling() {
            this.pollingInterval = setInterval(this.getMessages, this.intervalTime);
        }
    },
    created() {
        this.getMessages();
        this.startPolling();
    },
    beforeUnmount() {
        if (this.pollingInterval) {
            clearInterval(this.pollingInterval);
        }
    },
    formatDate(dateString) {
        return new Date(dateString).toLocaleString();
    }
}
</script>

<style scoped>
.user-conversation {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
}

.messages-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.message {
    padding: 10px;
    border-radius: 8px;
    max-width: 70%;
    margin: 5px 0;
}

.sent {
    align-self: flex-end;
    background-color: #007bff;
    color: white;
}

.received {
    align-self: flex-start;
    background-color: #e9ecef;
    color: black;
}

.message-content {
    margin-bottom: 5px;
}

.message-timestamp {
    font-size: 0.8em;
    opacity: 0.8;
}
</style>