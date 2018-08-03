import express from 'express';
import bodyParser from 'body-parser';
import Bot from '../discord/bot';
import Songs from './routes/Songs';

class Server {
    constructor() {
        console.log('Server init ...');
        this.init();
        this.bindRoutes();
        this.start();
    }

    init() {
        this.app = express();
        this.app.use((req, res, next) => {
            res.header("Access-Control-Allow-Origin", "*");
            res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
            res.setHeader('Content-Type', 'application/json');
            next();
        });
        this.app.use(bodyParser.urlencoded({ extended: true }));
        this.app.use(bodyParser.json());
        this.app.use(express.static('../../web-client/lib'));
    }

    bindRoutes() {
        this.app.use('/api/songs', Songs);
    }

    start() {
        this.app.listen(4000, () => {
            console.log('Server listening ...');
        })
    }
}

new Server();