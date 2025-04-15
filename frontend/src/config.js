const config = {
  development: {
    apiUrl: 'http://localhost:5533',
  },
  production: {
    apiUrl: 'http://138.68.76.63:5533',
  }
};

// Determine current environment
// You can use process.env.NODE_ENV which Vue CLI sets automatically,
// or set a custom ENV variable during build
const env = process.env.NODE_ENV || 'development';

export default {
  // Export the config for the current environment
  ...config[env],
  
  // You can also export a method to check current environment
  isDev: () => env === 'development',
  isProd: () => env === 'production',
  
  // Get current environment name
  env: env
};