import { OnGatewayConnection, SubscribeMessage, WebSocketGateway, WebSocketServer } from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { ServerEvents } from './dto/server-events.enum';
import { BotStateDTO } from './dto/bot-state.dto';
import { forwardRef, Inject, Logger } from '@nestjs/common';
import { YoutubeService } from '../youtube/youtube.service';
import { BotService } from './bot.service';
import { ClientEvents } from './dto/client-events.enum';

@WebSocketGateway({ namespace: 'bot' })
export class BotGateway implements OnGatewayConnection {

  private readonly logger = new Logger(BotGateway.name);

  @WebSocketServer()
  private server: Server;

  constructor(
    @Inject(forwardRef(() => BotService))
    private readonly botService: BotService
  ) {
  }

  public sendStateUpdate(guildId: string, state: BotStateDTO) {
    this.server.to(guildId).emit(ServerEvents.STATE_UPDATE, { state });
  }

  handleConnection(client: any, ...args: any[]): any {
    // client.emit(this.botService.getState());
  }

  @SubscribeMessage(ClientEvents.JOIN_ROOM)
  private joinRoom(socket: Socket, data: { guildId: string }) {
    this.logger.log(`Received event : ${ClientEvents.JOIN_ROOM} (${data.guildId})`);
    // socket.leaveAll();
    socket.join(data.guildId);
    socket.emit(
      ServerEvents.STATE_UPDATE,
      { state: this.botService.getState(data.guildId) });
  }

  @SubscribeMessage(ClientEvents.CHANGE_VOLUME)
  private addToPlaylistAction(socket: Socket, data: { volume: number }) {
    const guildId = this.getGuildId(socket);
    this.logger.debug(`Received event : ${ClientEvents.CHANGE_VOLUME} (${guildId} - ${data.volume})`);
    this.botService.setVolume(guildId, data.volume);
  }

  private getGuildId(socket: Socket): string {
    return Object.keys(socket.rooms)[0];
  }
}
