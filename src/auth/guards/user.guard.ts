import {ExecutionContext, Injectable, UnauthorizedException} from '@nestjs/common';
import {AuthService} from '../auth.service';
import {Request} from 'express';

@Injectable()
export class UserGuard {

  public async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest<Request>();
    const authorizationHeader = request.headers.authorization as string;
    if (!authorizationHeader) {
      throw new UnauthorizedException('Missing Authorization header');
    }
    if (!authorizationHeader.startsWith('Basic ')) {
      throw new UnauthorizedException('Wrong Authorization header');
    }
    const token = authorizationHeader.slice(6);
    const user = await AuthService.getDiscordUser(token);
    if (!user) {
      throw new UnauthorizedException('Failed to contact Discord API');
    }
    request.params.user = user;
    return !!user;
  }
}
