import React from "react";
import '../styles/Auth.css';

function Auth() {
    const [formData, setFormData] = React.useState({
        login: "",
        password: ""
    });

    function handleChange(event) {
        const { name, value } = event.target;
        setFormData(prevFormData => ({
            ...prevFormData,
            [name]: value
        }));
    }

    function login() {
        const { login, password } = formData;
        const requestBody = {
            login: login,
            password: password
        };

        fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody)
        })
        .then(response => {
            if (response.ok) {
                return response.json().then(data => {
                    localStorage.setItem("token", data.token);
                    localStorage.setItem("login",login);
                    localStorage.setItem("password",password)
                    window.location.reload();
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

    function register() {
        const { login, password } = formData;
        const requestBody = {
            login: login,
            password: password
        };

        fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody)
        })
        .then(response => {
            if (response.ok) {
                return response.json().then(data => {
                    alert(data.message);
                    localStorage.setItem("login",login);
                    localStorage.setItem("password",password)
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

    function handleSubmit(event) {
        event.preventDefault();
        register(); // Register the user first
    }

    function handleLogin(event) {
        event.preventDefault();
        login(); // Then login
    }

    return (
        <div className="form-container">
            <form className="form" onSubmit={handleSubmit}>
                <input 
                    type="text" 
                    placeholder="Login"
                    className="form--input"
                    name="login"
                    onChange={handleChange}
                    value={formData.login}
                />
                <input 
                    type="password" 
                    placeholder="Password"
                    className="form--input"
                    name="password"
                    onChange={handleChange}
                    value={formData.password}
                />
                <button 
                    className="form--submit"
                    type="submit"
                >
                    Sign up
                </button>
                <a 
                href="#"
                    onClick={handleLogin}
                >
                    Sign in
                </a>
            </form>
        </div>
    );
}

export default Auth;
