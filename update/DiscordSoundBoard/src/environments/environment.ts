// This file can be replaced during build by using the `fileReplacements` array.
// `ng build --prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  api: 'http://localhost:8000/api',
  websocket: 'http://localhost:8000/',
  isDiscordIntegrationEnabled: true,
  discord: {
    api: 'https://discordapp.com/api/oauth2/authorize',
    clientId: '497474636133564417',
    permissions: '3147776',
    scope: 'bot'
  }
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.
