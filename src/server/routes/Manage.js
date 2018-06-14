import Express from "express";
import Bot from '../../discord/bot';

let router = Express.Router();

router.get(('/guilds'), (req, res) => {
    let guilds = Bot.getGuilds();
    let activeGuilds = Bot.getActiveGuilds();
    let guildsDTO = [];
    for (let guild of guilds) {
        guildsDTO.push(
            {
                id: guild.id,
                icon: guild.icon,
                name: guild.name,
                active: !!activeGuilds.find((el)=>el.id===guild.id)
            }
        );
    }
    res.send(guildsDTO);
});

export default router;