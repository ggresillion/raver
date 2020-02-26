import {CanActivate, ExecutionContext, Injectable, UnauthorizedException} from '@nestjs/common';
import {Request} from 'express';
import {DiscordService} from '../../discord/discord.service';
import {ConfigService} from '@nestjs/config';

@Injectable()
export class UserGuard implements CanActivate {

  constructor(
    private readonly discordService: DiscordService,
    private readonly configService: ConfigService,
  ) {
  }

  public async canActivate(context: ExecutionContext) {
    if (!this.configService.get('isDiscordIntegrationEnabled')) {
      return true;
    }
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
