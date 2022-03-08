import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Loader } from "../../../components/Loader";
import { HttpClient } from "../../../services/http";
import { Guild } from "../../../services/model/guild";
import { RootState } from "../../../store";
import './Guilds.scss';
import { joinGuild } from "../../../services/guilds";
import { setGuild } from "../../../slices/guilds";

export function Guilds() {

    const [guilds, setGuilds] = useState<Guild[]>();

    const dispatch = useDispatch();
    const selectedGuildId = useSelector((state: RootState) => state.guilds.guildId);

    useEffect(() => {
        const getGuilds = async () => {
            const http = new HttpClient();
            const res = await http.get<Guild[]>('/guilds');

            const guildId = localStorage.getItem('selectedGuildId');
            if (!!guildId && res.some(g => g.id === guildId)) {
                selectGuild(guildId);
            } else {
                dispatch(setGuild({ guildId: null }));
            }

            setGuilds(res);
        };
        getGuilds();
    }, []);

    async function selectGuild(guildId: string) {
        try {
            await joinGuild(guildId);
            dispatch(setGuild({ guildId }));
            localStorage.setItem('selectedGuildId', guildId);

        } catch (e) {
            console.log(e);
        }
    }

    if (!guilds) {
        return (
            <div className="guilds">
                <Loader />
            </div>
        )
    }

    return (
        <div className="guilds">
            {guilds.map(g => {
                const style = {
                    backgroundColor: stringToRGB(g.name)
                };
                return (
                    <div className={g.id === selectedGuildId ? 'guild selected' : 'guild'} key={g.id} style={style} title={g.name} onClick={() => selectGuild(g.id)}>
                        {
                            !g.icon ?
                                <span>{g.name.charAt(0).toUpperCase()}</span> :
                                <img className="guild-icon" src={`https://cdn.discordapp.com/icons/${g.id}/${g.icon}`} />
                        }
                    </div>
                )
            })}
        </div>
    )
}

function stringToRGB(str: string): string {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    const c = (hash & 0x00FFFFFF)
        .toString(16)
        .toUpperCase();
    const color = '00000'.substring(0, 6 - c.length) + c;
    return '#' + color;
}