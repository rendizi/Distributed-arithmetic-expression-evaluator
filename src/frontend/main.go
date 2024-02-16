package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, website)
	})
	http.ListenAndServe(":8079", nil)
}

var website = `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Request Sender</title>
    <style>
		
.container {
    max-width: 600px;
    margin: 0 auto;
    text-align: center;
    padding: 20px;
}

button {
    padding: 10px 20px;
    margin: 10px;
    font-size: 16px;
    cursor: pointer;
}

#response-box {
    margin-top: 20px;
    border: 1px solid #ccc;
    padding: 10px;
    min-height: 100px;
}

</style>
</head>
<body>
<div class="container">
    <h1>Request Sender</h1>
    <button onclick="getExpression()">Expressions</button>
    <button onclick="sendRequest('http://127.0.0.1:8080/operations')">Operations</button>
    <button onclick="sendRequestMachines('http://127.0.0.1:8080/machines')">Machines</button>
    <div id="response-box"></div>

    <div class="container">
        <h1>Request Sender</h1>
        <div>
            <label for="uid">UID:</label>
            <input type="text" id="uid" value="1">
        </div>
        <div>
            <label for="task">Task:</label>
            <input type="text" id="task" value="2 + 2 * 2">
        </div>
        <div>
            <label for="settings">Settings (comma separated: *,/,+,-):</label>
            <input type="text" id="settings" value="1,1,1,1">
        </div>
        <button onclick="sendExpression()">Send Request</button>
    </div>
</div>
<script>
    function sendRequest(url) {
        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                document.getElementById("response-box").innerText = xhr.responseText;
            }
        };
        xhr.send();
    }

    function sendExpression() {
        var uid = document.getElementById("uid").value;
        var task = document.getElementById("task").value;
        var settings = document.getElementById("settings").value;

        // Format the task parameter
        task = encodeURIComponent(task);

        var url = "http://127.0.0.1:8080/expression?uid=" + uid + "&task=" + task + "&settings=" + settings;

        var xhr = new XMLHttpRequest();
        xhr.open("POST", url, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                document.getElementById("response-box").innerText = xhr.responseText;
            }
        };
        xhr.send();
    }
	function sendRequestMachines(){
		var xhr = new XMLHttpRequest();
        xhr.open("GET", "http://127.0.0.1:8080/machines", true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                document.getElementById("response-box").innerText = xhr.responseText;
            }
        };
        xhr.send();
	}
		
    function getExpression() {
        var uid = document.getElementById("uid").value;
        var url = 'http://127.0.0.1:8080/expression?id=' + uid;

        var xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                document.getElementById("response-box").innerText = xhr.responseText;
            }
        };
        xhr.send();
    }
</script>

</body>
</html>

`
