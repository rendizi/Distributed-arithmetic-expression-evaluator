import React from "react";


export default function Header(props) {
    const [prompt, setPrompt] = React.useState("");

    function updateResults() {
        fetch(`http://localhost:8080/expression?id=${prompt}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': localStorage.getItem("token"),
            },
        })
        .then(response => {
            if (response.ok) {
                return response.json().then(data => {
                    props.setSearch((prev)=>({...prev,expression: data.expression, result: data.result, time: data.time}))
                });
            } else {
                return response.json().then(error => {
                    alert(`Error: ${error.message}`);
                });
            }
        })
        .catch(error => {
            console.error("Network Error:", error);
            alert("Network Error. Please try again later.");
        });
    }

    function handleChange(event) {
        const { name, value } = event.target;
        if (name === 'prompt') {
            setPrompt(value);
        }
    }
    
    function handleKeyDown(event) {
        if (event.key === 'Enter') {
            updateResults();
        }
    }

    function logout(){
        localStorage.removeItem("login")
        localStorage.removeItem("password")
        sessionStorage.removeItem("token")
        window.location.reload()
    }
    
    return (
        <nav>
        <a href="#" onClick={()=>props.setWindow("home")}>Home</a>
        <a href="#" onClick={()=>props.setWindow("expressions")}>Expressions</a>
        <a href="#" onClick={()=>props.setWindow("operations")}>Operations</a>
        <a href="#" onClick={()=>props.setWindow("agents")}>Agents</a>
        <input 
            type="text" 
            placeholder="Search..." 
            value={prompt} 
            name="prompt" 
            onChange={handleChange} 
            onKeyDown={handleKeyDown} 
        />
        <button onClick={logout}>Log Out</button>
    </nav>
    

    );
}