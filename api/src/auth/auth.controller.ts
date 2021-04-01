import { BadRequestException, Controller, Get, Query, Req, Res, UnauthorizedException } from '@nestjs/common';
import { Request, Response } from 'express';
import { AuthService } from './auth.service';
import fetch from 'node-fetch';
import { ConfigService } from '@nestjs/config';
import { URLSearchParams } from 'url';

@Controller('auth')
export class AuthController {

  private readonly CLIENT_ID = process.env.CLIENT_ID;
  private readonly CLIENT_SECRET = process.env.CLIENT_SECRET;
  private readonly scopes = ['identify', 'guilds'];

  constructor(
    private readonly configService: ConfigService,
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
    const response = await fetch('https://discordapp.com/api/oauth2/token',
      {
        method: 'POST',
        headers: {
          Authorization: `Basic ${creds}`,
        },
        body: params
      });
    const body = await response.json();
    if (!response.ok) {
      throw new UnauthorizedException(body.error_description);
    }
    const accessToken = body.access_token;
    const refreshToken = body.refresh_token;
    res.redirect(`${this.configService.get('clientUrl') ? this.configService.get('clientUrl') : `${host}/#/`}` +
      `login?accessToken=${accessToken}&refreshToken=${refreshToken}`);
  }
}
