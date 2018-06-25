import Express from "express";
import youtube from "youtube-dl";

const SONG_DIR = 'songs/';

let router = Express.Router();

router.get('/search', (req, res) => {
    let url = req.query.url;
    youtube.getInfo(url, (err, info) => {
        if (err) res.sendStatus(404);
        else res.send({
            id: info.id,
            title: info.title,
            url: info.url,
            thumbnail: info.thumbnail,
            description: info.description
        });
    })
});

router.get('/upload', (req, res) => {
    let url = req.query.url;
    let category = req.query.category;
    let name = req.query.name;
    if (!url || !category || !name) {
        res.sendStatus(400);
        return;
    }

    youtube.exec(url, ['-x', '--audio-format', 'mp3', '-o', name + '.%(ext)s'], {cwd: SONG_DIR + '/' + category}, (err, output) => {
        if (err) throw err;
        console.log(output.join('\n'));
        res.sendStatus(200);
    });
});

export default router;