import {useEffect, useState} from 'react';
import {WindowHide} from "../wailsjs/runtime";
import useWebSocket from "react-use-websocket";
import {SetIconOffline, SetIconOnline} from "../wailsjs/go/main/App";

function toWsUrl(mister: string) {
  return `ws://${mister}/api/ws`
}

function formatGame(game: string) {
  let formatted = game;

  if (formatted === "") {
    return "";
  }

  let idx = formatted.indexOf("/");
  if (idx > -1) {
    formatted = formatted.substring(idx + 1);
  }

  idx = formatted.indexOf(".");
  if (idx > -1) {
    formatted = formatted.substring(0, idx);
  }

  return formatted;
}

function handleMessage(
  event: MessageEvent,
  setCore: (core: string) => void,
  setGame: (game: string) => void) {
  const msg = event.data;
  const idx = msg.indexOf(":");
  const cmd = msg.substring(0, idx);
  const data = msg.substring(idx + 1);

  switch (cmd) {
    case "coreRunning":
      setCore(data);
      break;
    case "gameRunning":
      setGame(formatGame(data));
      break;
  }
}

function App() {
  const [mister, setMister] = useState("mistuh:8182");
  const [core, setCore] = useState("");
  const [game, setGame] = useState("");

  const ws = useWebSocket(toWsUrl(mister), {
    onMessage: (event) => handleMessage(event, setCore, setGame),
    shouldReconnect: () => true,
    reconnectAttempts: Infinity,
    share: true,
  });

  useEffect(() => {
    if (ws.readyState === WebSocket.OPEN) {
      console.log("setting icon online");
      SetIconOnline().catch((err) => console.error(err));
    } else {
      console.log("setting icon offline");
      SetIconOffline().catch((err) => console.error(err));
    }
  }, [ws.readyState]);

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
              Active core
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
