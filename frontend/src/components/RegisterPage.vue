<template>
    <div class="register-container">
        <h2>Register</h2>
        <form @submit.prevent="handleRegister">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" v-model="username" required placeholder="Choose a username" />
            </div>
            <div class="form-group">
                <label for="displayName">Display Name</label>
                <input type="text" id="displayName" v-model="displayName" required
                    placeholder="Enter your display name" />
            </div>
            <div class="form-group">
                <label for="email">Email</label>
                <input type="email" id="email" v-model="email" required placeholder="Enter your email" />
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" v-model="password" required placeholder="Choose a password" />
            </div>
            <button type="submit" :disabled="loading">
                {{ loading ? "Registering..." : "Register" }}
            </button>
            <p v-if="error" class="error">{{ error }}</p>
            <p class="login-link">
                Already have an account?
                <router-link to="/login">Login here</router-link>
            </p>
        </form>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "RegisterPage",
    data() {
        return {
            username: "",
            displayName: "",
            email: "",
            password: "",
            loading: false,
            error: null,
            baseUrl: "http://localhost:5533",
        };
    },
    methods: {
        async handleRegister() {
            this.loading = true;
            this.error = null;

            // Create JSON payload instead of FormData
            const payload = {
                username: this.username,
                displayName: this.displayName,
                email: this.email,
                password: this.password
            };

            try {
                await axios.post(this.baseUrl + "/register", payload, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                this.$router.push("/login");
            } catch (err) {
                this.error = err.response?.data?.error || "Registration failed. Please try again.";
            } finally {
                this.loading = false;
            }
        }
    }
};
</script>

<style scoped>
.register-container {
    max-width: 400px;
    margin: 2rem auto;
    padding: 2rem;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
    text-align: center;
    margin-bottom: 1.5rem;
    color: #333;
}

.form-group {
    margin-bottom: 1rem;
}

label {
    display: block;
    margin-bottom: 0.5rem;
    color: #555;
}

input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
}

button {
    width: 100%;
    padding: 0.7rem;
    background-color: #4caf50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
}

button:disabled {
    background-color: #a5d6a7;
    cursor: not-allowed;
}

.error {
    margin-top: 1rem;
    color: #d32f2f;
    text-align: center;
}

.login-link {
    margin-top: 1rem;
    text-align: center;
}

.login-link a {
    color: #4caf50;
    text-decoration: none;
}

.login-link a:hover {
    text-decoration: underline;
}
</style>
