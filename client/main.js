const wsURI = "ws://192.168.1.101:3550/stock";
const wsList = [];
const stockPriceDisplay = document.getElementById("stock-price");
const table = document.getElementById("stock-table");

for(const stock of stocks) {
    wsList.push(connectToWS(wsURI, stock.price));
}

// handle Error on Web Socket
for(const ws of wsList){
    ws.onerror = function () {
        console.log("Web Socket Connection Error!");
    }
}

// Map Stock data to the Html Table
for(const stock of stocks) {
    let rowCount = table.rows.length;
    let row = table.insertRow(rowCount);
    let stockName = row.insertCell(0);
    stockName.innerHTML = stock.name;
    let marketCap = row.insertCell(1);
    marketCap.innerHTML = stock.market_cap;
    let stockPrice = row.insertCell(2);
    stockPrice.innerHTML = stock.price;
    stockPrice.className = marketCap.className = "whitespace-nowrap px-3 py-4 text-sm text-gray-300";
}

// receive data
for(let i = 0; i < wsList.length; i++){
    wsList[i].onmessage = function (event) {
        const previousData = table.rows[i + 2].cells[2].innerHTML;
        table.rows[i + 2].cells[2].innerHTML = event.data + (previousData > event.data ? "🔽" : "🔼");
        if(previousData > event.data){
            table.rows[i + 2].cells[2].style.color = "red";
        } else {
            table.rows[i + 2].cells[2].style.color = "green";
        }
    }
}

// close Web Socket Connection
window.addEventListener("unload", function (){
    wsList.map((ws)=> ws.close());
});