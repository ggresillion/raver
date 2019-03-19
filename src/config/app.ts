const config = {
  development: {
    port: 8000,
    clientUrl: 'http://localhost:4200/',
  },
  production: {
    port: 80,
    clientUrl: '',
  },
};

export default config[process.env.NODE_ENV || 'development'];
