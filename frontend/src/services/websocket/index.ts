import { v4 as uuidv4 } from 'uuid';
import { AppDispatch } from '../../store';
import { WebSocketTypes } from './types';

export interface Message<T> {
  id?: string;
  type: string;
  payload?: T;
  room?: string;
  token?: string;
}

class WebSocketClient {

  private socket: WebSocket;

  private dispatch: AppDispatch;

  private messageSubscribers: (<T>(m: Message<T>) => void)[] = [];

  constructor(url: string, dispatch: AppDispatch) {
    this.dispatch = dispatch;
    this.socket = new WebSocket(url);
    this.socket.onmessage = (m) => {
      const message = JSON.parse(m.data);
      this.dispatch({ type: WebSocketTypes.RECEIVED, payload: message });
      this.dispatch({ type: message.type, payload: message.payload });
      this.messageSubscribers.forEach(cb => cb(message));
    };
  }

  async join(room: string): Promise<void> {
    return this.sendWithResponse(WebSocketTypes.JOIN, { room });
  }

  send(type: string, payload: any, room?: string, id?: string) {
    const message = { type, payload, room, id, token: localStorage.getItem('accessToken') };
    this.dispatch({ type: WebSocketTypes.SEND, payload: message });
    this.socket.send(JSON.stringify(message));
  }

  async sendWithResponse(type: string, payload: any, room?: string): Promise<any> {
    return new Promise<any>((resolve, reject) => {
      const id = uuidv4();
      this.send(type, payload, room, id);
      setTimeout(() => reject(), 5000);
      this.onMessage(m => {
        if (m.id === id) {
          if (m.type === WebSocketTypes.OK) {
            resolve(m.payload);
          } else {
            reject(m.payload);
          }
        }
      });
    });
  }

  onMessage(cb: (message: Message<any>) => void) {
    this.messageSubscribers.push(cb);
  }
}

export let webSocketClient: WebSocketClient;

export function connect(url: string, dispatch: AppDispatch) {
  webSocketClient = new WebSocketClient(url, dispatch);
}
