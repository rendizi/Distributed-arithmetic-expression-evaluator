import React from "react";


function Default(){
    const [formData, setFormData] = React.useState({
        expression: "",
        settings: {
            plus: 1,
            minus: 1,
            mult: 1,
            div: 1,
        }
    })

    function handleSubmit(event){
        event.preventDefault();
        fetch('http://localhost:8080/expression', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': localStorage.getItem("token"),
            },
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (response.ok) {
                return response.json().then(data => {
                    alert(data.id, data.message);
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
            setFormData(prevFormData => ({
                ...prevFormData,
                [name]: value
            }));
    }
    function handleSettings(event){
        const { name, value } = event.target;
        setFormData(prevFormData => ({
            ...prevFormData,
            settings: {
                ...prevFormData.settings,
                [name]: parseInt(value)
            }
        }));
    }
    

    return (
        <div class="default">
            <form class="form" onSubmit={handleSubmit}>
                <label htmlFor="exor">Expression:</label>
                <input 
                    type="text" 
                    placeholder="Ex: 2 + 2 * 2"
                    className="form--input"
                    name="expression"
                    onChange={handleChange}
                    value={formData.expression} />
                <div class="form--settings">
                    <label for="plus">Plus:</label>
                    <input 
                        type="number" 
                        className="form--input"
                        name="plus"
                        min="0"
                        step="1"
                        onChange={handleSettings}
                        value={formData.plus} />

                    <label for="minus">Minus:</label>
                    <input 
                        type="number" 
                        className="form--input"
                        name="minus"
                        min="0"
                        step="1"
                        onChange={handleSettings}
                        value={formData.minus} />

                    <label for="mult">Mult:</label>
                    <input 
                        type="number" 
                        className="form--input"
                        name="mult"
                        min="0"
                        step="1"
                        onChange={handleSettings}
                        value={formData.mult} />

                    <label for="div">Div:</label>
                    <input 
                        type="number" 
                        className="form--input"
                        name="div"
                        min="0"
                        step="1"
                        onChange={handleSettings}
                        value={formData.div} />
                </div>

                <button className="default--submit">
                    Send
                </button>
            </form>
</div>

    )
}

export default Default