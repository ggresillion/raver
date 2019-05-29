import { WebSocketGateway, WebSocketServer, SubscribeMessage, OnGatewayConnection } from '@nestjs/websockets';
import { Server, Client, Socket } from 'socket.io';
import { PlayerInfos } from './model/player-infos';
import { ClientEvents } from './dto/client-events.enum';
import { TrackInfos } from './dto/track-infos';
import { ServerEvents } from './dto/server-events.enum';
import { YoutubeService } from './youtube.service';
import { Inject, forwardRef } from '@nestjs/common';
import { PlayerStatus } from './model/player-status';

@WebSocketGateway({ namespace: 'player' })
export class YoutubeGateway implements OnGatewayConnection {

  @WebSocketServer()
  private server: Server;
  private addToPlaylistListeners = [];

  constructor(
    @Inject(forwardRef(() => YoutubeService))
    private readonly youtubeService: YoutubeService,
  ) { }

  public handleConnection(client: Socket) {
    client.emit(
      ServerEvents.SYNC,
      {
        status: this.youtubeService.getStatus(),
        playlist: this.youtubeService.getPlaylist(),
      });
  }

  public sendStatusUpdate(status: PlayerStatus) {
    this.server.emit(ServerEvents.STATUS_UPDATED, {status});
  }

  public onAddToPlaylist(cb: (track: TrackInfos) => void) {
    this.addToPlaylistListeners.push(cb);
  }

  public sendAddToPlaylist(track: TrackInfos) {
    this.server.emit(ServerEvents.ADD_TO_PLAYLIST, { track });
  }

  @SubscribeMessage(ClientEvents.ADD_TO_PLAYLIST)
  private addToPlaylistAction(client: Client, data: { track: TrackInfos }) {
    if (!!data.track) {
      this.addToPlaylistListeners.forEach(cb => cb(data.track));
    }
  }

  @SubscribeMessage(ClientEvents.GET_PLAYLIST)
  private getPlaylistAction(): TrackInfos[] {
    return this.youtubeService.getPlaylist();
  }

  @SubscribeMessage(ClientEvents.PLAY)
  private playAction() {
    this.youtubeService.play();
  }

  @SubscribeMessage(ClientEvents.PAUSE)
  private plauseAction() {
    this.youtubeService.pause();
  }
}
