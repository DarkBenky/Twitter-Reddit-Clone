const apiUrl = process.env.VUE_APP_API_URL || 'http://localhost:5533'; // Fallback
const env = process.env.NODE_ENV || 'development';

export default {
  apiUrl: apiUrl, // Use the value from the environment variable

  // You can also export a method to check current environment
  isDev: () => env === 'development',
  isProd: () => env === 'production',

  // Get current environment name
  env: env
};