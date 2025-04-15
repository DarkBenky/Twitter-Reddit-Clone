<template>
  <div class="login-container">
    <h2>Login</h2>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">Username or Email</label>
        <input type="text" id="username" v-model="username" required placeholder="Enter your username or email" />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" required placeholder="Enter your password" />
      </div>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
    <div class="login-container">
      <!-- ...existing template code... -->
      <p class="register-link">
        Don't have an account? <router-link to="/register">Register here</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import config from '../config.js'

export default {
  name: 'LoginPage',
  data() {
    return {
      username: '',
      password: '',
      loading: false,
      error: null,
      baseUrl: config.apiUrl,
    }
  },
  methods: {
    async handleLogin() {
      this.loading = true;
      this.error = null;
      try {
        const response = await axios.post(this.baseUrl + '/login', {
          username: this.username,
          password: this.password,
        });

        // Store token in cookies first
        window.$cookies.set('auth_token', response.data.token);

        // Verify it was set
        const storedToken = window.$cookies.get('auth_token');
        console.log('Token stored in cookies:', storedToken);

        // Store user data in Vuex
        this.$store.commit('setUserId', response.data.user.idUser);
        this.$store.commit('setCurrentUser', response.data.user);
        this.$store.commit('setToken', response.data.token);

        console.log('User data stored in Vuex:', response.data);

        this.$router.push('/');
      } catch (err) {
        console.error('Login error:', err);
        this.error = err.response?.data?.error || 'Login failed. Please try again.';
      } finally {
        this.loading = false;
      }
    },
  },
}
</script>

<style scoped>
.register-link {
  margin-top: 1rem;
  text-align: center;
}

.register-link a {
  color: #4CAF50;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}

.login-container {
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
  background-color: #4CAF50;
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
</style>