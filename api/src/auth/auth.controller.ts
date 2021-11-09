import { BadRequestException, Controller, Get, Query, Req, Res, UnauthorizedException } from '@nestjs/common';
import { Request, Response } from 'express';
import { ConfigService } from '@nestjs/config';
import { URLSearchParams } from 'url';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';

@Controller('auth')
export class AuthController {

  private readonly CLIENT_ID = process.env.CLIENT_ID;
  private readonly CLIENT_SECRET = process.env.CLIENT_SECRET;
  private readonly scopes = ['identify', 'guilds'];

  constructor(
    private readonly configService: ConfigService,
    private readonly httpService: HttpService
  ) {
  }

  @Get('login')
  public async login(@Res() res: Response, @Req() req: Request): Promise<any> {

    const redirect = req.protocol + '://' + req.get('Host') + '/api/auth/token';
    res.redirect(`https://discordapp.com/oauth2/authorize?client_id=${this.CLIENT_ID}&scope=${this.scopes.join('%20')}` +
      `&response_type=code&redirect_uri=${redirect}`);
  }

  @Get('token')
  public async token(@Req() req: Request, @Res() res: Response, @Query('code') code: string) {
    if (!code) {
      throw new BadRequestException();
    }
    const host = `${req.protocol}://${req.get('Host')}/`;
    const redirect = host + 'api/auth/token';
    const creds = Buffer.from(`${this.CLIENT_ID}:${this.CLIENT_SECRET}`).toString('base64');
    const params = new URLSearchParams();
    params.append('grant_type', 'authorization_code');
    params.append('code', code);
    params.append('redirect_uri', redirect);
    params.append('scope', this.scopes.join('%20'));

    try {
      const response = await firstValueFrom(this.httpService.post('https://discordapp.com/api/oauth2/token', params,
        {
          headers: {
            Authorization: `Basic ${creds}`,
            'Content-Type': 'application/x-www-form-urlencoded'
          }
        }));
      const accessToken = response.data['access_token'];
      const refreshToken = response.data['refresh_token'];
      res.redirect(`${this.configService.get('clientUrl') ? this.configService.get('clientUrl') : `${host}/#/`}` +
        `login?accessToken=${accessToken}&refreshToken=${refreshToken}`);
    } catch (e) {
      throw new UnauthorizedException(e.response.data.error_description);
    }

  }
}
