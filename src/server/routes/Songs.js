import Express from 'express';
import fs from 'fs';
import path from 'path';
import Bot from '../../discord/bot';

let router = Express.Router();

const SONG_DIR = 'songs/';

function extractSongName(song) {
    return song.substr(0, song.lastIndexOf('.'));
}

function getSongsFromDir(dir) {
    let songs = [];
    fs.readdirSync(dir).forEach((song) => {
        songs.push({name: extractSongName(song), path: dir + '/' + song})
    });
    return songs;
}

router.get('/', (req, res) => {
    if (req.query.category && req.query.song) {
        try {
            Bot.play(req.query.category, req.query.song);
            res.sendStatus(200);
        } catch (err) {
            res.status(404);
            res.send(err.message);
        }
    } else {
        let categories = [];
        fs.readdirSync(SONG_DIR).map((category) => {
            categories.push({name: category, songs: getSongsFromDir(SONG_DIR + category)});
        });
        res.json(categories);
    }
});

router.get('/', (req, res) => {
    console.log(req.params)

});

export default router;