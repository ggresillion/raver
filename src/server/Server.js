import express from 'express';
import bodyParser from 'body-parser';
import Bot from '../discord/bot';
import Songs from './routes/Songs';
import Categories from "./routes/Categories";

class Server {
    constructor() {
        console.log('Server init ...');
        this.init();
        this.bindRoutes();
        this.start();
    }

    init() {
        this.app = express();
        this.app.use(express.static('web-client/build'));
        this.app.use((req, res, next) => {
            res.header("Access-Control-Allow-Origin", "*");
            res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
            res.header('Access-Control-Allow-Methods', 'PUT, POST, GET, DELETE, OPTIONS');
            next();
        });
        this.app.use(bodyParser.urlencoded({ extended: true }));
        this.app.use(bodyParser.json());
    }

    bindRoutes() {
        this.app.use('/api/songs', Songs);
        this.app.use('/api/categories', Categories);
    }

    start() {
        this.app.listen(4000, () => {
            console.log('Server listening');
        })
    }
}

new Server();