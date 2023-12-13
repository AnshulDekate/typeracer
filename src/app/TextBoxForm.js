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
  const [inputValue, setInputValue] = useState('');
  const [givenText, setGivenText] = useState('');
  const [words, setWords] = useState([])
  const [accepted, setResult] = useState([]);
  const [idx, setIdx] = useState(0);

  useEffect(() => {
    var s = GiveText()
    setGivenText(GiveText())
    setWords(strToWords(s))
  }, []);


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
    }
    else{
      setInputValue(event.target.value);
    }
  };


  const handleSubmit = (event) => {
    event.preventDefault();
    console.log('Input value:', inputValue);
    setResult("Accepted")
    setIdx(0)
  };
  // space bar should be trigger submit and check if the word is correct or not
  return (
    <div>
      {/*Given Text*/}
        <div>
          {givenText}
        </div>
        <form >
        <label>
            Enter Text:
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
            {accepted}
        </div>
    </div>
    
  );
};

export default TextBoxForm;
