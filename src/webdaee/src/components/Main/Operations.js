import React from "react";
import { FaCalculator, FaEquals } from 'react-icons/fa'; 


function Operation(props){
    return (
        <div className="expression">
            <div className="icon-label">
                <FaCalculator className="icon" />
                <span className="label">Operation:</span>
            </div>
            <div className="expression-content">
                <div className="expression-text">{props.expression}</div>
            </div>
            <div className="icon-label">
                <FaEquals className="icon" /> 
                <span className="label">Result:</span>
            </div>
            <div className="expression-content">
                <div className="expression-result">{props.result}</div>
            </div>

        </div>
    );
}

function Operations(props) {
    const [operations, setOperations] = React.useState([]);

    React.useEffect(() => {
        const token = localStorage.getItem('token');

        if (token) {
            fetch('http://localhost:8080/operations', {
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                setOperations(data);
            })
            .catch(error => {
                console.error('Error fetching expressions:', error);
            });
        } else {
            console.error('Token not found in localStorage');
        }
    }, []); 


    const expressionComponents = operations.map((expr) => (
        <Operation key={expr.id} expression={expr.operation} result={expr.result} />
    ));
    
    return (
        <div className="expressions">
            {expressionComponents}
        </div>
    );
    
}
export default Operations