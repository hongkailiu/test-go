$(document).ready(function () {
    $.get("/whoami", function (data) {
        var username = data.username
        $("#username").text(username);
        if (username == "") {
            $("#login").show()
            $("#logout").hide()
        } else {
            $("#login").hide()
            $("#logout").show()
        }
    });


    $.ajax({
        type: 'get',
        url: '/token',
        statusCode: {
            200: function (data) {
                $("#token").text(data.token);
                $("#header").text("Authorization: " + "Bearer " + data.token);
            }
        }
    });

});