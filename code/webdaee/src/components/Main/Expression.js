import { FaCalculator, FaEquals, FaClock } from 'react-icons/fa'; 


function Expression(props) {
    return (
        <div className="expression">
            <div className="icon-label">
                <FaCalculator className="icon" />
                <span className="label">Expression:</span>
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
            <div className="icon-label">
                <FaClock className="icon" />
                <span className="label">Time:</span>
            </div>
            <div className="expression-content">
                <div className="expression-time">{props.time}</div>
            </div>
        </div>
    );
}

export default Expression