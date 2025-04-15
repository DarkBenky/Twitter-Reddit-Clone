import axios from 'axios'

// Create axios instance
const api = axios.create({
  baseURL: 'http://localhost:5533'
})

// Add request interceptor to include token with every request
api.interceptors.request.use(
  config => {
    const token = window.$cookies.get('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

export default api