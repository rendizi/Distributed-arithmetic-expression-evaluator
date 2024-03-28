import React from "react";
import Expression from "./Expression"

function Expressions() {
    const [expressions, setExpressions] = React.useState([]);

    React.useEffect(() => {
        const token = localStorage.getItem('token');

        if (token) {
            fetch('http://localhost:8080/expression', {
                headers: {
                    Authorization: token
                }
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                setExpressions(data);
            })
            .catch(error => {
                console.error('Error fetching expressions:', error);
            });
        } else {
            console.error('Token not found in localStorage');
        }
    }, []); 

    console.log(expressions)

    const expressionComponents = expressions.map((expr) => (
        <Expression key={expr.id} expression={expr.expression} result={expr.result} time={expr.time} />
    ));
    
    return (
        <div className="expressions">
            {expressionComponents}
        </div>
    );
    
}
export default Expressions