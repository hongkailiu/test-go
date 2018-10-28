$(document).ready(function(){
    $.get( "/whoami", function( data ) {
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
});