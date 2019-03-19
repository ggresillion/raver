import {BadRequestException, Controller, Get, Query, Req, Res} from '@nestjs/common';
import {Request, Response} from 'express';
import {AuthService} from './auth.service';
import fetch from 'node-fetch';

@Controller('auth')
export class AuthController {

  private readonly CLIENT_ID = process.env.CLIENT_ID;
  private readonly CLIENT_SECRET = process.env.CLIENT_SECRET;
  private readonly redirectEndpoint = 'api/auth/callback';
  private readonly scopes = ['identify'];

  constructor(private readonly authService: AuthService) {
  }

  @Get('login')
  public async login(@Res() res: Response, @Req() req: Request): Promise<any> {

    const redirect = this.getRedirectURI(req) + '/login';
    res.redirect(`https://discordapp.com/oauth2/authorize?client_id=${this.CLIENT_ID}&scope=${this.scopes.join()}` +
      `&response_type=code&redirect_uri=${redirect}`);
  }

  @Get('callback')
  public async callback(@Req() req: Request, @Query('code')code: string) {
    if (!code) {
      throw new BadRequestException();
    }
    const redirect = this.getRedirectURI(req);
    const creds = Buffer.from(`${this.CLIENT_ID}:${this.CLIENT_SECRET}`).toString('base64');
    const response = await fetch(`https://discordapp.com/api/oauth2/token?grant_type=authorization_code&code=${code}` +
      `&redirect_uri=${redirect}&scope=${this.scopes.join()}`,
      {
        method: 'POST',
        headers: {
          Authorization: `Basic ${creds}`,
        },
      });
    const body = await response.json();
    const token = body.access_token;
    const refreshToken = body.refresh_token;
    const jwtToken = this.authService.generateToken(token, refreshToken);
  }

  private getRedirectURI(req: Request) {
    return req.protocol + '://' + req.get('Host');
  }
}
