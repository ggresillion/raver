import {OnGatewayConnection, SubscribeMessage, WebSocketGateway, WebSocketServer} from '@nestjs/websockets';
import {Client, Server, Socket} from 'socket.io';
import {ClientEvents} from './dto/client-events.enum';
import {TrackInfos} from './dto/track-infos';
import {ServerEvents} from './dto/server-events.enum';
import {YoutubeService} from './youtube.service';
import {forwardRef, Inject, Logger} from '@nestjs/common';
import {PlayerStatus} from './model/player-status';
import {PlayerState} from './model/player-state';

@WebSocketGateway({namespace: 'player'})
export class YoutubeGateway implements OnGatewayConnection {

  private readonly logger = new Logger(YoutubeGateway.name);

  @WebSocketServer()
  private server: Server;
  private addToPlaylistListeners = [];

  constructor(
    @Inject(forwardRef(() => YoutubeService))
    private readonly youtubeService: YoutubeService,
  ) {
  }

  public handleConnection(client: Socket) {
    client.emit(
      ServerEvents.SYNC,
      {state: this.youtubeService.getState()});
  }

  public sendStatusUpdate(status: PlayerStatus) {
    this.server.emit(ServerEvents.STATUS_UPDATED, {status});
  }

  public sendStateUpdate(state: PlayerState) {
    this.logger.log('Sending state update');
    this.server.emit(ServerEvents.STATE_UPDATED, {state});
  }

  public onAddToPlaylist(cb: (track: TrackInfos) => void) {
    this.addToPlaylistListeners.push(cb);
  }

  public sendProgressUpdate(progressSeconds: number): void {
    this.server.emit(ServerEvents.PROGRESS_UPDATED, {progressSeconds});
  }

  // public sendAddToPlaylist(track: TrackInfos) {
  //   this.server.emit(ServerEvents.YT_ADD_TO_PLAYLIST, {track});
  // }

  @SubscribeMessage(ClientEvents.ADD_TO_PLAYLIST)
  private addToPlaylistAction(client: Client, data: { track: TrackInfos }) {
    this.logger.log(`Received event : ${ClientEvents.ADD_TO_PLAYLIST} (${data.track.title})`);
    if (!!data.track) {
      this.addToPlaylistListeners.forEach(cb => cb(data.track));
    }
  }

  @SubscribeMessage(ClientEvents.GET_PLAYLIST)
  private getPlaylistAction(): TrackInfos[] {
    this.logger.log(`Received event : ${ClientEvents.GET_PLAYLIST}`);
    return this.youtubeService.getPlaylist();
  }

  @SubscribeMessage(ClientEvents.PLAY)
  private playAction() {
    this.logger.log(`Received event : ${ClientEvents.PLAY}`);
    this.youtubeService.play();
  }

  @SubscribeMessage(ClientEvents.PAUSE)
  private plauseAction() {
    this.logger.log(`Received event : ${ClientEvents.PAUSE}`);
    this.youtubeService.pause();
  }

  @SubscribeMessage(ClientEvents.NEXT)
  private nextAction() {
    this.logger.log(`Received event : ${ClientEvents.NEXT}`);
    this.youtubeService.next();
  }
}
