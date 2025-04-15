import axios from 'axios'
import config from '../config'

// Create axios instance with dynamic baseURL
const api = axios.create({
  baseURL: config.apiUrl,
  withCredentials: true,
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