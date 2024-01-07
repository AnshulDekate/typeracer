'use client'
import React, { useState, useEffect } from 'react';
import {getRaces, postRaceCompleted} from './BackendREST';
import {Lobby} from './BackendWebSocket';
import {RaceBox} from './RaceBox';

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
  var s = "Feel the thunder of engines and the anticipation at the starting line. Paint a vivid picture of the moments before a race begins."
  return (
    s
  )
}


const MainBox = () => {
  // don't allow copy paste
  console.log("Component rendered");
  const [countDown, setCountDown] = useState('');
  const [practice, setPractice] = useState(1);
  const [joinLobby, setJoinLobby] = useState(0);
  const [inputBoxDisabled, setInputBoxDisabled] = useState(true);
  const [speed, setSpeed] = useState(0); // check why we are getting -ve speed after first word
  const [startTime, setStartTime] = useState(new Date());
  const [endTime, setEndTime] = useState(new Date());
  const [inputValue, setInputValue] = useState('');
  const [givenText, setGivenText] = useState(GiveText()); // ux is faster
  const [words, setWords] = useState([])
  const [raceCompleted, setRaceCompleted] = useState('');
  const [idx, setIdx] = useState(0);
  const [races, setRaces] = useState(0);

  useEffect(() => {
    console.log("after mount")
    setWords(strToWords(givenText))
  }, []);

  useEffect(() => {
    // getRaces(setRaces)
  }, [raceCompleted]);

  const start = () => {
    setIdx(0); setSpeed(0);
    console.log("start time", startTime)
    setRaceCompleted("")
  }

  const handlePractice = (event) => {
    event.preventDefault();
    // countDown(); // shouldn't be able to click on input till this ends
    setPractice(1);
    start()
    setStartTime(new Date()) // change it to first input
    setInputBoxDisabled(false);
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
    // postRaceCompleted(setRaceCompleted);
    setInputBoxDisabled(true);
    setJoinLobby(0);
  };
  // space bar should be trigger submit and check if the word is correct or not
  
  const handleJoinLobby = (event) =>{
    event.preventDefault();
    setPractice(0)
    setJoinLobby(joinLobby+1)
    start()
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      // Add your custom logic here
    }
  };

  return (
    <div>
      <div className='flex-container'>
        <div>Typing Showdown</div>
        <div>User: Guest
          {/* <div> {races} races </div> */}
        </div>
      </div>
      <div className="line-div"></div>
      
      <div className='filler-div'> </div> 
      <div className='filler-div'> </div> 

      <div className='main-container'>
        <div className='flex-container-begin'>
          <button name="practice" onClick={handlePractice} className='practice' disabled={joinLobby}>
              Practice
          </button>  
          <div> or </div>
          <button name="join lobby" onClick={handleJoinLobby} className='join-lobby' disabled={joinLobby}>
              Join Lobby 
          </button>
        </div>

        <div className='filler-div'> </div> 
        <div>
          <Lobby idx={idx} words={words} joinLobby={joinLobby} practice={practice} setInputBoxDisabled={setInputBoxDisabled} setStartTime={setStartTime}/>
        </div>
        <div className='filler-div'> </div> 
        
        {practice ? (<RaceBox pos={(idx*100)/words.length}/>): null}

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
            onKeyDown={handleKeyDown}
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

export default MainBox;
