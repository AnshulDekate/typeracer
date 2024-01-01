import React, {useState, useEffect} from "react";

export const Lobby = ({idx, words, joinLobby, practice}) => {
    const [socket, setSocket] = useState(null);
    const [playerID, setPlayerID] = useState()
    const [sessionID, setSessionID] = useState()
    const [players, setPlayers] = useState(0)
    // const [nxtRank, setnxtRank] = useState(1)
    // const [rank, setRank] = useState([])
    const [progress, setProgress] = useState([])
    const [rank, setRank] = useState([])
    
    useEffect(()=>{
        const newSocket = new WebSocket('ws://10.107.107.107:8081/ws')
        newSocket.addEventListener("message", (event) => {
            try {
                const resp = JSON.parse(event.data);
                if (resp.event == "players") {
                    console.log("Received players event:", resp.data);
                    console.log(resp.data.N)
                    console.log(resp.data.progress)
                    setPlayers(resp.data.N);
                    setProgress(resp.data.progress);
                    setRank(resp.data.rank);
                } else if (resp.event == "joined") {
                    console.log("Joined the lobby with player id: ", resp.data);
                    setPlayerID(resp.data.playerID);
                    setSessionID(resp.data.sessionID)
                } else {
                    console.log("Unknown event:", data.event);
                }
            } catch (error) {
                console.error("Error parsing message:", error);
            }
        });

        newSocket.addEventListener('open', () => {
            setSocket(newSocket)
            console.log('WebSocket connection established.', newSocket.readyState, WebSocket.OPEN);
        });
        
        newSocket.addEventListener('error', (error) => {
            console.error('WebSocket error:', error);
        });

        newSocket.addEventListener('close', (event) => {
            console.log('WebSocket connection closed:', event.code, event.reason);
            setSocket(null);
        });

        return () => {
            if (newSocket.readyState === WebSocket.OPEN) {
                console.log('WebSocket connection closed')
                newSocket.close();
            }
        }
    }, []);

    useEffect(() => {
        console.log("useeffect sending progress")
        const data = {
            "sessionID": sessionID,
            "playerID": playerID,
            "percentage": parseInt((idx*100)/words.length),
        };
          
        if (socket && socket.readyState == WebSocket.OPEN) {
            console.log("IN")
            console.log("sending progress", socket.readyState)
            console.log(JSON.stringify(data))
            socket.send(JSON.stringify(data));
        }
    }, [idx])

    useEffect(() => {
        const message = 'join';
        if (socket && socket.readyState === WebSocket.OPEN) {
         socket.send(message);
        }
    }, [joinLobby])

    return (
        <div>
            {practice ? null:(
                <div>
                    <div className='filler-div'> </div>     
                    {/* Currently {players} players in lobby */}
                    <div>
                        Your player ID is {playerID} and session ID is {sessionID}
                    </div>
                    <div>
                        Stats:
                        <ul>
                            {/* Render the players as an array of React elements */}
                            {Object.keys(progress).map((key) => (
                                <li key={key}>{`Player ${key}: Rank ${rank[key]!== undefined ? rank[key] : ''}`}</li>
                                ))}
                        </ul>
                    </div>
                    <div>
                    <div className='filler-div'> </div>     
                    </div>
                </div>
            )}
        </div>
    )
}