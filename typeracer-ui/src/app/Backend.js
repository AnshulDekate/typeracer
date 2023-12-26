export const getRaces = async (setRaces) => {
    try {
        const response = await fetch('http://localhost:8081/races');
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
        const response = await fetch('http://localhost:8081/raceCompleted');
        setRaceCompleted("Race Completed!")
    } catch (error) {
        console.error('Error fetching data:', error);
    }
};