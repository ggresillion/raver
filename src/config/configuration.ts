import {Logger} from '@nestjs/common';

const logger = new Logger('ConfigModule');

const config = {
  development: {
    port: 8000,
    clientUrl: 'http://localhost:4200/#/',
    isDiscordIntegrationEnabled: true,
    database: {
      type: "postgres",
      host: "localhost",
      port: 5432,
      database: "dsb",
      username: "dsb",
      password: "dsb",
      autoLoadEntities: true
    }
  },
  production: {
    port: 8000,
    isDiscordIntegrationEnabled: true,
    database: {
      type: "postgres",
      host: process.env.DB_HOST || "localhost",
      port: process.env.DB_PORT || 5432,
      database: process.env.DB_DATABASE || "dsb",
      username: process.env.DB_USERNAME || "dsb",
      password: process.env.DB_PASSWORD,
      autoLoadEntities: true
    }
  },
};

const isProd = process.env.NODE_ENV === 'production';
if (isProd) {
  logger.log('Production configuration loaded');
} else {
  logger.log('Development configuration loaded');
}
export default () => config[process.env.NODE_ENV || 'development'];
