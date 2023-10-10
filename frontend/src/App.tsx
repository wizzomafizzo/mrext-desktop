import {useEffect, useState} from 'react';
import {WindowHide} from "../wailsjs/runtime";
import useWebSocket from "react-use-websocket";
import {GetHost, SetHost, SetIconOffline, SetIconOnline, WriteGame, WriteSystem} from "../wailsjs/go/main/App";
import axios from "axios";

function toWsUrl(mister: string) {
    return `ws://${mister}/api/ws`
}

function toApiUrl(mister: string) {
    return `http://${mister}/api`
}

interface Running {
    core: string;
    system: string;
    systemName: string;
    game: string;
    gameName: string;
}

async function getPlaying(mister: string): Promise<Running> {
    const resp = await axios.get(toApiUrl(mister) + "/games/playing");
    return resp.data;
}

function handleMessage(
    event: MessageEvent,
    mister: string,
    setCore: (core: string) => void,
    setGame: (game: string) => void) {
    const msg = event.data;
    const idx = msg.indexOf(":");
    const cmd = msg.substring(0, idx);
    // const data = msg.substring(idx + 1);

    console.log("cmd: " + cmd);

    if (cmd === "coreRunning" || cmd === "gameRunning") {
        getPlaying(mister).then((running) => {
            let core: string;
            if (running.systemName !== "") {
                core = running.systemName;
            } else {
                core = running.core;
            }
            setCore(core === "" ? "—" : core);
            WriteSystem(core).catch((err) => console.error(err));
            setGame(running.gameName === "" ? "—" : running.gameName);
            WriteGame(running.gameName).catch((err) => console.error(err));
        }).catch((err) => {
            console.error(err);
        });
    }
}

function App() {
    const [mister, setMister] = useState("");
    const [core, setCore] = useState("");
    const [game, setGame] = useState("");

    const ws = useWebSocket(toWsUrl(mister), {
        onMessage: (event) => handleMessage(event, mister, setCore, setGame),
        shouldReconnect: () => true,
        reconnectAttempts: Infinity,
        share: true,
    });

    useEffect(() => {
        if (ws.readyState === WebSocket.OPEN) {
            SetHost(mister).catch((err) => console.error(err));
            SetIconOnline().catch((err) => console.error(err));
        } else {
            SetIconOffline().catch((err) => console.error(err));
        }
    }, [ws.readyState]);

    useEffect(() => {
        GetHost().then((host) => {
            setMister(host);
        }).catch((err) => {
            console.error(err);
        });
    }, []);

    return (
        <div id="app">
            <div style={{padding: "10px"}}>
                <div className="form">

                    <div className="row">
                        <label className="label" htmlFor="mister">MiSTer address</label>
                        <input
                            className="value"
                            autoComplete="off"
                            name="mister"
                            type="text"
                            value={mister}
                            onChange={(e) => {
                                setMister(e.target.value);
                            }}
                            style={{
                                borderRadius: "5px",
                                backgroundColor: ws.readyState === WebSocket.OPEN ? "lightgreen" : "lightcoral"
                            }}
                        />
                    </div>

                    <div className="row">
                        <div className="label">
                            Active system
                        </div>
                        <div className="value">
                            {core === "" ? "—" : core}
                        </div>
                    </div>

                    <div className="row">
                        <div className="label">
                            Active game
                        </div>
                        <div className="value">
                            {game === "" ? "—" : game}
                        </div>
                    </div>
                    <div style={{flexGrow: 1, display: "flex", alignItems: "flex-end", justifyContent: "right"}}>
                        <button
                            onClick={() => {
                                WindowHide();
                            }}
                        >
                            Hide window
                        </button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App
