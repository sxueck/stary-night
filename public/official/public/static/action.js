async function postData(url = '', data = {}) {
    const resp = await fetch(url, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        headers: {
            'Content-Type': 'application/json',
        },
        redirect: 'follow',
        referrerPolicy: 'same-origin',
        body: JSON.stringify(data),
    });
    return resp.json();
}

function ObtainEmail() {
    let mail = $("#mail").text()
    if (mail === "") {
        console.log("mail is null");
        return 0
    }

    postData("https://")
}
