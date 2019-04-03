import {ExecutionContext, Injectable, UnauthorizedException} from '@nestjs/common';
import {Request} from 'express';
import {DiscordService} from '../../discord/discord.service';

@Injectable()
export class UserGuard {

  constructor(private readonly discordService: DiscordService) {
  }

  public async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest<Request>();
    const authorizationHeader = request.headers.authorization as string;
    if (!authorizationHeader) {
      throw new UnauthorizedException('Missing Authorization header');
    }
    if (!authorizationHeader.startsWith('Bearer ')) {
      throw new UnauthorizedException('Wrong Authorization header');
    }
    const token = authorizationHeader.slice(7);
    this.discordService.setToken(token);
    const user = await this.discordService.getUser();
    request.params.user = user;
    return !!user;
  }
}
