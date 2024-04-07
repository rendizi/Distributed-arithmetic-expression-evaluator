import React from "react";
import { FaLink, FaClock } from 'react-icons/fa'; 


function Agent(props){
    return (
        <div className="expression">
            <div className="icon-label">
                <FaLink className="icon" />
                <span className="label">Address:</span>
            </div>
            <div className="expression-content">
                <div className="expression-text">{props.address}</div>
            </div>
            <div className="icon-label">
                <FaClock className="icon" /> 
                <span className="label">Last ping:</span>
            </div>
            <div className="expression-content">
                <div className="expression-result">{props.last_ping}</div>
            </div>
            
        </div>
    );
}

function Agents(props) {
    const [agents, setAgents] = React.useState([]);

    React.useEffect(() => {
        const token = localStorage.getItem('token');

        if (token) {
            fetch('http://localhost:8080/agents', {
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                setAgents(data);
            })
            .catch(error => {
                console.error('Error fetching expressions:', error);
            });
        } else {
            console.error('Token not found in localStorage');
        }
    }, []); 


    const expressionComponents = agents.map((expr) => (
        <Agent key={expr.id} address={expr.address} last_ping={expr.last_ping} />
    ));
    
    return (
        <div className="expressions">
            {expressionComponents}
        </div>
    );
    
}

export default Agents