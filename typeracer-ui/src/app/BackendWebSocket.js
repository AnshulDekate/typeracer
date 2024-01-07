import React, {useState, useEffect} from "react";
import { RaceBox } from "./RaceBox";

export const Lobby = ({idx, words, joinLobby, practice, setInputBoxDisabled, setStartTime}) => {
    const [socket, setSocket] = useState(null);
    const [playerID, setPlayerID] = useState()
    const [sessionID, setSessionID] = useState()
    const [players, setPlayers] = useState(0)
    const [timer, setTimer] = useState(20)
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
                    setTimer(20-resp.data.timer)
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
          
        if (sessionID>0 && socket && socket.readyState == WebSocket.OPEN) {
            console.log("sending progress", socket.readyState)
            console.log(JSON.stringify(data))
            socket.send(JSON.stringify(data));
        }
    }, [idx])

    useEffect(() => {
        console.log("sending join")
        const message = 'join';
        if (joinLobby==1 && socket && socket.readyState === WebSocket.OPEN) {
         socket.send(message);
        }
    }, [joinLobby])

    useEffect(()=>{
        if (timer==0) {
            setInputBoxDisabled(false)
            setStartTime(new Date())
        }
    }, [timer])

    return (
        <div>
            {practice ? null:(
                <div>
                    <div> 
                    {(progress!=null && progress[playerID]== 100)? `You got rank ${rank[playerID]}`: (timer == 0 ? "Start!": `Timer: ${timer}`) }
                    </div>
                    <div className='filler-div'> </div>     
                    {/* Currently {players} players in lobby */}
                    <div>
                        {/* <div>
                            Stats: Your player ID is {playerID} and session ID is {sessionID}
                        </div>
                        <ul>
                            {Object.keys(progress).map((key) => (
                                <li key={key}>{`Player ${key}: progress ${progress[key]!== undefined ? progress[key] : ''}%, rank ${rank[key]!== undefined ? rank[key] : '-'}`}</li>
                                ))}
                        </ul> */}
                        
                        <RaceBox pos={progress[playerID]} /> 

                        {Object.keys(progress).map((key) => (
                            key==playerID? null : (<div key={key} className='dimmer'>
                                 <RaceBox pos={progress[key]} /> 
                            </div>)
                        ))}   
                        <div> {players<=3? "waiting for other players to join ... \n" : null} </div>
                        <div> {timer<=12? "adding two bots ..." : null} </div>
                        
                    </div>
                    <div className='filler-div'> </div>         
                </div>
            )}
        </div>
    )
}