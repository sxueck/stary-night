const allSitesCookieName = "AllResult"

function lightning() {
    let v = getCookie(allSitesCookieName)
    if (v === "") {
        axios.get('/api/v1/list', {
            responseType: 'json',
        }).then(function (res) {
            v = JSON.stringify(res.data)
            setCookie(allSitesCookieName, v, 60 * 3) // 3min

        }).catch(function (err) {
            console.log(err);
        });
    }

    let resCookie = getCookie(allSitesCookieName)
    console.log(resCookie)
    let res = JSON.parse(resCookie)
    let ran = randomNumber(res.length)
    window.location.href = res[ran].url;
}

function randomNumber(max) {
    return Math.floor(Math.random() * max)
}

function setCookie(cName, cValue, exSecond) {
    let d = new Date();
    d.setTime(d.getTime() + (exSecond));
    let expires = "max-age=" + d.toUTCString();
    document.cookie = cName + "=" + cValue + "; " + expires;
}

function getCookie(cName) {
    let name = cName + "=";
    let ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i].trim();
        if (c.indexOf(name) === 0) return c.substring(name.length, c.length);
    }
    return "";
}
