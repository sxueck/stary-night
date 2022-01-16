function post_sites_describe() {
    let url = $("#url").val()
    let author = $("#author").val()
    let contact = $("#contact").val()
    let description = $("#description").val()

    $("#remind").text("系统正在检查您的网站")

    $.ajax({
        url: "/api/v1/site",
        method: "POST",
        dataType: "text",
        contentType: "application/json; charset=utf-8",
        async: false,
        data: JSON.stringify({
            url: url,
            author: author,
            contact: contact,
            description: description
        }),
        success: function (res) {
            if (res.indexOf("[ERROR]") === -1) {
                $("#welcome-look").text("赞 比超棒还棒")
                $("#welcome-who").text(res + " 里边请")
                $("#remind").text("添加成功, 请等待下一个轮询器周期结束(~20min)即可")
                window.setTimeout(BackReturn, 2000);
            } else {
                $("#welcome-who").text("这是详细的错误消息")
                $("#welcome-prompt").text(res)
            }
        },
        error: function (e) {
            $("#welcome-who").text("嘿朋友 也许出现了点问题")
            $("#welcome-prompt").text("快来检查下自己的输入")
            console.log(e)
        }
    })
}

function BackReturn() {
    window.history.go(-1)
}