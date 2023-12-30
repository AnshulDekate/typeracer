export const getRaces = async (setRaces) => {
    console.log("calling get races from ui")
    try {
        const response = await fetch('http://10.107.107.107:8081/races');
        var data = await response.text();
        data = Number(data)
        console.log("data ", data)
        setRaces(data)
        return data
    } catch (error) {
        console.error('Error fetching data:', error);
    }
};

export const postRaceCompleted = async (setRaceCompleted) => {
    try {
        const response = await fetch('http://10.107.107.107:8081/raceCompleted');
        setRaceCompleted("Race Completed!")
    } catch (error) {
        console.error('Error fetching data:', error);
    }
};

export const joinLobby = async (setPlayers) => {
    try {
        const response = await fetch('http://10.107.107.107:8081/joinLobby');
        var data = await response.text()
        console.log("players ", data)
        setPlayers(Number(data))
    } catch (error) {
        console.error('Error fetching data:', error);
    }
}
