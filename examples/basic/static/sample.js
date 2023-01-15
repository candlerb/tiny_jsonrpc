function ping() {
    fetch('rpc', {
        'method': 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            'jsonrpc': '2.0',
            'method': 'ping',
        }),
    })
    .then((response) => response.json())
    .then(function (data) {
        alert(data.result);
    })
    .catch(function (err) {
        console.log('Something went wrong!', err);
    });
}

function add() {
    arg1 = parseFloat(document.getElementById('arg1').value);
    arg2 = parseFloat(document.getElementById('arg2').value);
    fetch('rpc', {
        'method': 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            'jsonrpc': '2.0',
            'method': 'add',
            'params': [arg1, arg2],
        }),
    })
    .then((response) => response.json())
    .then(function (data) {
        alert(data.result);
    })
    .catch(function (err) {
        console.log('Something went wrong!', err);
    });
}
