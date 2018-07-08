var websocket;
var output;

function WebSocketSupport(message, action)
{
    if (browserSupportsWebSockets() === false) {
        document.getElementById("ws_support").innerHTML = "<h2>Sorry! Your web browser does not supports web sockets</h2>";

        var element = document.getElementById("wrapper");
        element.parentNode.removeChild(element);

        return;
    }

    output = document.getElementById("output");

    websocket = new WebSocket('ws:' +  window.location.hostname + ':' + window.location.port  + '/websocket');

    websocket.onopen = function(e) {
        writeToScreen("You have have successfully connected to the server");
        doSend(message + ',' + action, output);
    };


    websocket.onmessage = function(e) {
        onMessage(e)
    };

    websocket.onerror = function(e) {
        onError(e)
    };
}

function onMessage(e) {
    writeToScreen('<span style="color: blue;"> ' + e.data + '</span>');
}

function onError(e) {
    writeToScreen('<span style="color: red;">ERROR:</span> ' + e.data);
}

function doSend(message, output) {
    output.value = "";
    output.focus();
    websocket.send(message);
    writeToScreen(message);
}

function writeToScreen(message) {
    var pre = document.createElement("p");
    pre.style.wordWrap = "break-word";
    pre.innerHTML = message;
    //output.appendChild(pre);
    output = pre;
}

function userInputSupplied() {
    var chatname = document.getElementById('chatname').value;
    var msg = document.getElementById('msg').value;
    if (chatname === '') {
        return 'Please enter your username';
    } else if (msg === '') {
        return 'Please the message to send';
    } else {
        return '';
    }
}

function browserSupportsWebSockets() {
    if ("WebSocket" in window)
    {
        return true;
    }
    else
    {
        return false;
    }
}
