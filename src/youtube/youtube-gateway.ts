import {WebSocketGateway, WebSocketServer} from '@nestjs/websockets';
import {Server} from 'socket.io';
import {PlayerStatus} from './model/player-status';

@WebSocketGateway({namespace: 'player'})
export class YoutubeGateway {

  @WebSocketServer()
  private server: Server;

  public sendStatusUpdate(status: PlayerStatus) {
    this.server.emit('status', status);
  }
}
