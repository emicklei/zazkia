function closeAllWithError() {
    $.ajax({
        type: "POST",
        url: "/links/closeAllWithError",
        success: function() { window.location.reload(); }
    });
}

function doLinkAction(id, action) {
    $.ajax({
        type: "POST",
        url: "/links/" + id + "/" + action,
        success: function() { window.location.reload(); }
    });
}

function doRouteAction(label, action) {
    $.ajax({
        type: "POST",
        url: "/routes/" + label + "/" + action,
        success: function() { window.location.reload(); }
    });
}