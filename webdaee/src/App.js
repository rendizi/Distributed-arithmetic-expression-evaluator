import './styles/App.css';
import Auth from "./components/Auth"
import Header from './components/Header';
import Main from "./components/Main"
import React from 'react';

function App() {
  const [window,setWindow] = React.useState("home")
  const [search, setSearch] = React.useState({
    expression: "",
    result: "",
    time: ""
  })
  const [isInitialRender, setIsInitialRender] = React.useState(true);

  

  React.useEffect(() => {
    if (!isInitialRender) {
      console.log(search);
      setWindow("expression");
    } else {
      setIsInitialRender(false);
    }
  }, [search, isInitialRender]);

  if (localStorage.getItem("login") === null || localStorage.getItem("password") === null){
    console.log("null")
    return (
      <Auth />
    )
  }

 
  

  return (
    <div className="App">
      <Header setWindow = {setWindow} setSearch = {setSearch} />
      <Main window={window} expression={search.expression} result = {search.result} time = {search.time} />
    </div>
  );
}

export default App;
