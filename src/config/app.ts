const config = {
  development: {
    port: 8000,
  },
  production: {
    port: 80,
  },
};

export default config[process.env.NODE_ENV || 'development'];
