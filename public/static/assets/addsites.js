function post_sites_describe() {
    let url = $("#url").val()
    let author = $("#author").val()
    let contact = $("#contact").val()
    let description = $("#description").val()

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
        success: function (res){
            if (res.indexOf("[ERROR]") === -1) {
                $("#welcome-look").text("赞 比超棒还棒")
                $("#welcome-who").text(res + " 里边请")
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