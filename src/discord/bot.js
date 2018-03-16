import Discord from 'discord.js';
import fs from 'fs';

const token = 'NDIzNTgwOTkzMzU0NzI3NDM0.DYsawQ.x9r0hgNg7H7J5v7y_pEJOPoXF68';

class Bot {

    constructor() {
        this.client = new Discord.Client();
        this.connections = [];
        this.configureClient();
    }

    configureClient() {
        this.client.on('ready', () => {
            console.log('Client Ready !');
        });

        this.client.on('message', message => {
            let command = message.content;
            if (command === '!join') {
                if (message.member.voiceChannel) {
                    message.member.voiceChannel.join()
                        .then(connection => {
                            this.connections.push(connection);
                        })
                        .catch(console.log);
                } else {
                    message.reply('You need to join a voice channel first!');
                }
            }
            else if (command.startsWith('!play')) {
                this.play(command.substring(6));
            }
            else if (command === '!stop' || command === '!quit' || command === '!leave') {
                if (message.member.voiceChannel) {
                    message.member.voiceChannel.leave();
                }
            }
        });

        this.client.login(token)
    }

    play(song) {
        if (this.connections.length === 0) {
            throw new Error('No connection !');
        }
        const songsDir = './src/songs/';
        fs.readdir(songsDir, (err, items) => {
            const filepath = songsDir + items.find((el) => {
                return el.toString().startsWith(song);
            });
            fs.stat(filepath, (err, stats) => {
                if (stats) {
                    this.connections.forEach((con) => {
                        let dispatcher = con.playFile(filepath);
                        dispatcher.on('end', () => {
                            console.log('Song ' + song + ' played !');
                        });
                    });
                }
                else {
                    throw new Error('No such song !');
                }

            });
        });
    }
}

export default new Bot();