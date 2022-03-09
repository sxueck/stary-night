function post_sites_describe() {
    let uid = $("#uid").value;
    let mail = $("#mail").value;
    let url = $("#url").value;
    let remark = $("#remark").value;

    axios({
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        transformRequest: [function (data) {
            data = JSON.stringify(data)
            return data
        }],
        params: {},
        url: "/api/v1/site",
        data: {
            author: uid,
            contact: mail,
            description: remark,
            url: url,
        }
    })
        .then(function (res) {
            console.log(res);
            if (res.statusText === "OK") {
                // success
                $("#success_hint").hidden = false;
            } else {
                // failed
                $("#failed_hint").hidden = false;
                $("#error").innerText = res.statusText;
            }
        })
        .catch(function (e) {
            $("#failed_hint").hidden = false;
            $("#error").innerText = e;
        })
}

function eula() {
    let btn = $("#post_button")
    let s = $("#eula_status");

    if (!s.checked) {
        btn.disabled = true;
        btn.style.background = "gray";
        btn.innerText = "PROGRAM DISABLED"

        $(".remark").innerHTML = "请刷新网页后重试";
        s.hidden = true;
    }
}

function review() {
    $("#failed_hint").hidden = true;
}