import { OnGatewayConnection, SubscribeMessage, WebSocketGateway, WebSocketServer } from '@nestjs/websockets';
import { Client, Server, Socket } from 'socket.io';
import { ClientEvents } from './dto/client-events.enum';
import { TrackInfos } from './dto/track-infos';
import { ServerEvents } from './dto/server-events.enum';
import { YoutubeService } from './youtube.service';
import { forwardRef, Inject, Logger } from '@nestjs/common';
import { PlayerStatus } from './model/player-status';
import { PlayerState } from './model/player-state';
import { Guild } from 'discord.js';
import { Video } from 'ytsr';

@WebSocketGateway({ namespace: 'player' })
export class YoutubeGateway implements OnGatewayConnection {

  private readonly logger = new Logger(YoutubeGateway.name);

  @WebSocketServer()
  private server: Server;
  private addToPlaylistListeners: ((guildId: string, track: TrackInfos) => void)[] = [];

  constructor(
    @Inject(forwardRef(() => YoutubeService))
    private readonly youtubeService: YoutubeService,
  ) {
  }

  public handleConnection(client: Socket) {
  }

  public sendStatusUpdate(guildId: string, status: PlayerStatus) {
    this.server.to(guildId).emit(ServerEvents.STATUS_UPDATED, { status });
  }

  public sendStateUpdate(guildId: string, state: PlayerState) {
    this.logger.log('Sending state update');
    this.server.to(guildId).emit(ServerEvents.STATE_UPDATED, { state });
  }

  public onAddToPlaylist(cb: (guildId: string, track: Video) => void) {
    this.addToPlaylistListeners.push(cb);
  }

  public sendProgressUpdate(guildId: string, progressSeconds: number): void {
    this.server.to(guildId).emit(ServerEvents.PROGRESS_UPDATED, { progressSeconds });
  }

  // public sendAddToPlaylist(track: TrackInfos) {
  //   this.server.emit(ServerEvents.YT_ADD_TO_PLAYLIST, {track});
  // }

  @SubscribeMessage(ClientEvents.ADD_TO_PLAYLIST)
  private addToPlaylistAction(socket: Socket, data: { track: TrackInfos }) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.ADD_TO_PLAYLIST} (${guildId} - ${data.track.title})`);
    if (!!data.track) {
      this.addToPlaylistListeners.forEach(cb => cb(guildId, data.track));
    }
  }

  @SubscribeMessage(ClientEvents.GET_PLAYLIST)
  private getPlaylistAction(socket: Socket): TrackInfos[] {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.GET_PLAYLIST} (${guildId})`);
    return this.youtubeService.getPlaylist(guildId);
  }

  @SubscribeMessage(ClientEvents.PLAY)
  private playAction(socket: Socket) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.PLAY} (${guildId})`);
    this.youtubeService.play(guildId);
  }

  @SubscribeMessage(ClientEvents.PAUSE)
  private plauseAction(socket: Socket) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.PAUSE} (${guildId})`);
    this.youtubeService.pause(guildId);
  }

  @SubscribeMessage(ClientEvents.NEXT)
  private nextAction(socket: Socket) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.NEXT} (${guildId})`);
    this.youtubeService.next(guildId);
  }

  @SubscribeMessage(ClientEvents.JOIN_ROOM)
  private joinRoom(socket: Socket, data: { guildId: string }) {
    this.logger.log(`Received event : ${ClientEvents.JOIN_ROOM} (${data.guildId})`);
    socket.leaveAll();
    socket.join(data.guildId);
    socket.emit(
      ServerEvents.SYNC,
      { state: this.youtubeService.getState(data.guildId) });
  }

  @SubscribeMessage(ClientEvents.MOVE_UPWARDS)
  private moveUpwards(socket: Socket, data: { index: number }) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.MOVE_UPWARDS} (${guildId})`);
    this.youtubeService.moveUpwards(guildId, data.index);
  }

  @SubscribeMessage(ClientEvents.MOVE_DOWNWARDS)
  private moveDownwards(socket: Socket, data: { index: number }) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.MOVE_DOWNWARDS} (${guildId})`);
    this.youtubeService.moveDownwards(guildId, data.index);
  }

  @SubscribeMessage(ClientEvents.REMOVE_FROM_PLAYLIST)
  private removeFromPlaylist(socket: Socket, data: { index: number }) {
    const guildId = this.getGuildId(socket);
    this.logger.log(`Received event : ${ClientEvents.REMOVE_FROM_PLAYLIST} (${guildId})`);
    this.youtubeService.removeFromPlaylist(guildId, data.index);
  }

  private getGuildId(socket: Socket): string {
    return Object.keys(socket.rooms)[0];
  }
}
