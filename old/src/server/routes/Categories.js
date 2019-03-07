import Express from 'express';
import Bot from '../../discord/bot';
import fs from 'fs';

const SONG_DIR = 'songs/';

let router = Express.Router();

router.post('/', (req, res) => {
    fs.mkdir(SONG_DIR + req.body.name, (err)=>{
        console.log(err);
    });
    res.sendStatus(200);
});

export default router;