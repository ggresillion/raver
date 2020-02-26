import {Logger} from '@nestjs/common';

const logger = new Logger('ConfigModule');

const config = {
  development: {
    port: 8000,
    clientUrl: 'http://localhost:4200/#/',
    isDiscordIntegrationEnabled: false,
  },
  production: {
    port: 80,
    isDiscordIntegrationEnabled: false,
  },
};

const isProd = process.env.NODE_ENV === 'production';
if (isProd) {
  logger.log('Production configuration loaded');
} else {
  logger.log('Development configuration loaded');
}
export default () => config[process.env.NODE_ENV || 'development'];
