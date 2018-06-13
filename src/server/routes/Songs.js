import Express from 'express';
import fs from 'fs';
import Bot from '../../discord/bot';
import multer from 'multer';
import * as Finder from "fs-finder";

const SONG_DIR = 'songs/';

let router = Express.Router();
let storage = multer.diskStorage({
    destination: function (req, file, cb) {
        console.log(file)
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
let upload = multer({storage: storage});

function extractSongName(song) {
    return song.substr(0, song.lastIndexOf('.'));
}

function extractExtension(song) {
    return song.substr(song.lastIndexOf('.'), song.length);
}

function getSongsFromDir(dir) {
    let songs = [];
    fs.readdirSync(dir).forEach((song) => {
        songs.push({name: extractSongName(song), path: dir + '/' + song})
    });
    return songs;
}

router.get('/', (req, res) => {
    let categories = [];
    fs.readdirSync(SONG_DIR).map((category) => {
        categories.push({name: category, songs: getSongsFromDir(SONG_DIR + category)});
    });
    res.json(categories);
});

router.post('/', upload.single('songs'), (req, res) => {
    res.sendStatus(200);
});

router.put('/', (req, res) => {
    if (req.query.category && req.query.song) {
        let songs = Finder.from(SONG_DIR).findFiles(req.query.song + '.*');
        try {
            let oldPath = songs[0];
            let newPath = SONG_DIR + req.query.category + '/' + req.query.song + extractExtension(songs[0]);
            console.log('Moving ' + oldPath + ' to ' + newPath);
            fs.renameSync(oldPath, newPath);
            res.sendStatus(200);
        } catch (err) {
            console.log(err);
            res.json(err);
        }
    } else {
        res.sendStatus(404);
    }
});

router.get('/play', (req, res) => {
    if (req.query.category && req.query.song) {
        try {
            Bot.play(req.query.category, req.query.song);
            res.sendStatus(200);
        } catch (err) {
            res.status(404);
            res.send(err.message);
        }
    } else {
        res.sendStatus(404);
    }
});

export default router;