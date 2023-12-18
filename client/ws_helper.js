
function connectToWS(wsURI, stockPrice) {
    return new WebSocket(wsURI+ `?price=${stockPrice}`);
}