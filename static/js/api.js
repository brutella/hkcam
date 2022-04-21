function APIRequest(method, url) {
    "use strict";
    var xhttp = new XMLHttpRequest();
    xhttp.open(method, url);
    return xhttp;
}