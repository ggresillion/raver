import Express from 'express';
import fs from 'fs';
import path from 'path';
import Bot from '../../discord/bot';
import multer from 'multer';

const SONG_DIR = 'songs/';

let router = Express.Router();
let storage = multer.diskStorage({
    destination: function (req, file, cb) {
        fs.mkdir(SONG_DIR + req.body.category, (err) => {
            if (err || err && err.code === 'EEXIST') {
                cb(null, SONG_DIR + req.body.category)
            }
            else {
                console.log(err);
            }
        });
    },
    filename: function (req, file, cb) {
        cb(null, file.originalname);
    }
});
let upload = multer({ storage: storage });

function extractSongName(song) {
    return song.substr(0, song.lastIndexOf('.'));
}

function getSongsFromDir(dir) {
    let songs = [];
    fs.readdirSync(dir).forEach((song) => {
        songs.push({ name: extractSongName(song), path: dir + '/' + song })
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
            categories.push({ name: category, songs: getSongsFromDir(SONG_DIR + category) });
        });
        res.json(categories);
    }
});

router.post('/', upload.single('songs'), (req, res) => {
    res.sendStatus(200);
})

export default router;