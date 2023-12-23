// TextBoxForm.js
'use client'
import React, { useState, useEffect } from 'react';

const strToWords = (s) => {
  var words = []
  var word = ""
  for (let i=0; i < s.length; i++){
    if (s[i]==' ') {
      if (word.length) {
        word+=s[i]
        words.push(word)
      }
      word = ""
    }
    else word += s[i]
  }
  words.push(word)
  console.log(words)
  return words
}

const GiveText = () => {
  var s = "Type this text correctly you monkey"
  return (
    s
  )
}


const TextBoxForm = () => {
  const [countDown, setCountDown] = useState('');
  const [speed, setSpeed] = useState(0); // check why we are getting -ve speed after first word
  const [startTime, setStartTime] = useState(new Date());
  const [endTime, setEndTime] = useState(new Date());
  const [inputValue, setInputValue] = useState('');
  const [givenText, setGivenText] = useState('');
  const [words, setWords] = useState([])
  const [raceComplete, setRaceComplete] = useState('');
  const [idx, setIdx] = useState(0);
  const [races, setRaces] = useState(0);

  useEffect(() => {
    var s = GiveText()
    setGivenText(GiveText())
    setWords(strToWords(s))
  }, []);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://localhost:8081/races');

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const data = await response.text();
        console.log("data ", data)
        setRaces(data)
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [raceComplete]);


  const handleInputChange = (event) => {
    console.log(idx)
    console.log(words[idx])
    console.log("current", event.target.value)
    console.log(event.target.value == words[idx])


    // if word matched clear 
    if (event.target.value==words[idx]){
      setInputValue("");
      if (idx==words.length-1){
        handleSubmit(event)
      }
      else{
        setIdx(idx+1)
      }
      setEndTime(new Date())

      console.log("end time", endTime)
      console.log((60*(idx+1)*1000), endTime-startTime)
      
      var s = Math.abs((60*(idx+1)*1000)/(endTime-startTime))
      setSpeed(s)
    }
    else{
      setInputValue(event.target.value);
    }
  };

  const handleStartRace = (event) => {
    event.preventDefault();
    setIdx(0);
    // countDown(); // shouldn't be able to click on input till this ends
    setStartTime(new Date())
    console.log("start time", startTime)
    setSpeed(0);
    setRaceComplete("")
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    setRaceComplete("Race Completed!")

    try {
      const response = fetch('http://localhost:8081/raceCompleted');
      console.log(response)
    } catch (error) {
      console.error('Error fetching data:', error);
    }

    setIdx(0)
  };
  // space bar should be trigger submit and check if the word is correct or not


  
  return (
    <div>
      <button 
        name="start race"
        onClick={handleStartRace}>
          Start Race
      </button>  
      <div className='filler-div'> </div> 
      <div>
        User: anshul
      </div>
      <div>
        {races} races
      </div>
      <div className='filler-div'> </div> 
      <div>
        {countDown}
      </div>
      {/*Given Text*/}
      <div className='given-text'>
        {givenText}
      </div>
      <div className='filler-div'> </div> 
      <form >
      <label>
          Type here:   
          <input
          type="text"
          id="textInput" 
          name="textInput" 
          value={inputValue}
          onChange={handleInputChange}
          />
      </label>
      </form>
      <div>
        Speed: {speed.toFixed(2)} WPM
      </div>
      <div className='filler-div'> </div> 
      <div>
          {raceComplete}
      </div>
    </div>
    
  );
};

export default TextBoxForm;
