// TextBoxForm.js
'use client'
import React, { useState, useEffect } from 'react';
import {getRaces, postRaceCompleted} from './BackendREST';
import {Lobby} from './BackendWebSocket';

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
  console.log("Component rendered");
  const [countDown, setCountDown] = useState('');
  const [joinLobby, setJoinLobby] = useState(0);
  const [inputBoxDisabled, setInputBoxDisabled] = useState(true);
  const [speed, setSpeed] = useState(0); // check why we are getting -ve speed after first word
  const [startTime, setStartTime] = useState(new Date());
  const [endTime, setEndTime] = useState(new Date());
  const [inputValue, setInputValue] = useState('');
  const [givenText, setGivenText] = useState('');
  const [words, setWords] = useState([])
  const [raceCompleted, setRaceCompleted] = useState('');
  const [idx, setIdx] = useState(0);
  const [races, setRaces] = useState(0);

  useEffect(() => {
    console.log("after mount")
    var s = GiveText()
    setGivenText(GiveText())
    setWords(strToWords(s))
    getRaces(setRaces)
  }, []);

  useEffect(() => {
    getRaces(setRaces)
  }, [raceCompleted]);
  
  const handleStartRace = (event) => {
    event.preventDefault();
    // countDown(); // shouldn't be able to click on input till this ends
    setJoinLobby(0)
    setInputBoxDisabled(false);
    setIdx(0); 
    setStartTime(new Date())
    console.log("start time", startTime)
    setSpeed(0);
    setRaceCompleted("")
  }

  const handleInputChange = (event) => {
    console.log(idx)
    console.log(words[idx])
    console.log("current", event.target.value)
    console.log(event.target.value == words[idx])


    // if word matched clear 
    if (event.target.value==words[idx]){
      setInputValue("");

      setEndTime(new Date())
      console.log("end time", endTime)
      console.log((60*(idx+1)*1000), endTime-startTime)
      var s = Math.abs((60*(idx+1)*1000)/(endTime-startTime))
      setSpeed(s)

      if (idx==words.length-1){
        handleRaceCompleted(event)
      }
      setIdx(idx+1)
    }
    else{
      setInputValue(event.target.value);
    }
  };

  const handleRaceCompleted = (event) => {
    event.preventDefault();
    postRaceCompleted(setRaceCompleted);
    setInputBoxDisabled(true);
    // setIdx(0)
  };
  // space bar should be trigger submit and check if the word is correct or not
  
  const handleJoinLobby = (event) =>{
    setJoinLobby(joinLobby+1)
    setInputBoxDisabled(false)
  }

  return (
    <div>
      <div className='flex-container'>
        <div>Typing Showdown</div>
        <div>User: anshul
          {/* <div> {races} races </div> */}
        </div>
      </div>
      <div className="line-div"></div>
      
      <div className='filler-div'> </div> 
      <div className='filler-div'> </div> 

      <div className='main-container'>
        <div className='flex-container-begin'>
          <button name="practice" onClick={handleStartRace} className='practice'>
              Practice
          </button>  
          <div> or </div>
          <button name="join lobby" onClick={handleJoinLobby} className='join-lobby'>
              Join Lobby 
          </button>
        </div>

        <div className='filler-div'> </div> 
        <div>
          <Lobby idx={idx} words={words} joinLobby={joinLobby}/>
        </div>
        <div className='filler-div'> </div> 
        
        <div className='race-box'>
          <div className='actual-run'>
            <div className="car" style={{ position: 'relative', left: `${(idx*100)/words.length}%` }}>o^^^o</div>
            </div>
          <div className='invisible'>o^^^o</div>
        </div>

        {/*Given Text*/}
        <div className='given-text'>
          {givenText}
        </div>
        <div className='filler-div'> </div> 
        <form >
        <label>
            Type here :   
            <input
            type="text"
            id="textInput" 
            name="textInput" 
            value={inputValue}
            onChange={handleInputChange}
            disabled={inputBoxDisabled}
            autoCapitalize="none"
            style={{ marginLeft: '5px' }}
            />
        </label>
        </form>
        <div>
          Speed: {speed.toFixed(2)} WPM
        </div>
        <div className='filler-div'> </div> 
        <div>
            {raceCompleted}
        </div>
      </div>

    </div>
    
  );
};

export default TextBoxForm;
