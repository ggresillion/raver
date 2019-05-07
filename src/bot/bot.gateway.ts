import {WebSocketGateway, WebSocketServer} from '@nestjs/websockets';
import {Server} from 'socket.io';
import {BotStatus} from './bot.service';

@WebSocketGateway({namespace: 'bot'})
export class BotGateway {

  @WebSocketServer()
  private server: Server;

  public sendStatusUpdate(botStatus: BotStatus) {
    this.server.emit('status', botStatus);
  }
}
