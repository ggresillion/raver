import Express from 'express';
import fs from 'fs';
import Bot from '../../discord/bot';

let router = Express.Router();

router.get('/', (req, res) => {
    fs.readdir('./src/songs', (err, items) => {
        let songs = [];
        if (!items) {
            res.sendStatus(500);
        }
        items.forEach((item) => {
            songs.push({name: item.substr(0, item.lastIndexOf('.'))});
        });
        res.json(songs);
    });
});

router.get('/:sound', (req, res) => {
    try {
        Bot.play(req.params.sound);
        res.sendStatus(200);
    } catch (err) {
        res.status(404);
        res.send(err.message);
    }
});

export default router;