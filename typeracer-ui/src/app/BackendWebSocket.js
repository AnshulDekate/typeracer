import React, {useState, useEffect} from "react";
import { FONT_MANIFEST } from "../../node_modules/next/dist/shared/lib/constants";


import {joinLobby} from "./BackendREST"

// global lobby for now
// Global lobby, races countdown starts when at least two people are there in lobby
// global rank 
// each player has there cnt - words traversed
// after completing you use mutex and obtain the rank
// last player to complete finishes the race

export const Lobby = () => {
    const socket = new WebSocket('ws://10.107.107.107:8081/ws')

    const [players, setPlayers] = useState(0)
    // const [nxtRank, setnxtRank] = useState(1)
    // const [rank, setRank] = useState([])
    // const [progress, setProgress] = useState([])

    useEffect(()=>{
        socket.addEventListener("message", (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.event === "players") {
                    console.log("Received players event:", data.data);
                    setPlayers(data.data);
                } else {
                    console.log("Unknown event:", data.event);
                }
            } catch (error) {
                console.error("Error parsing message:", error);
            }
        });

        return () => {
            socket.close()
        }
    });

    const handleJoinLobby = (event) => {
        event.preventDefault();
        const message = 'join';
        socket.send(message);
    }
    return (
        <div>
            Currently {players} players in lobby
            <div>
            <div className='filler-div'> </div>     
            <button
            onClick={handleJoinLobby}
            >
                Join Lobby
            </button>
            </div>
        </div>
    )
}