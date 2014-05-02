$(document).ready(function() {
  var ws = new WebSocket("ws://localhost:8080/save");
  var draw = false;
  var id;

  $('.canvas').mousedown(function() {
    id = Math.random().toString(36).substr(2, 9);
    ws.send('STARTDRAW: ' + id);
    draw = true;
  });

  $(document).mouseup(function() {
    if(draw) {
      draw = false;
      ws.send('ENDDRAW');
      var url = '/drawing/' + id;
      var img = '<div class="image"><img src="' + url + '"/></div>';
      $(".log").prepend(img);
    }
  });

  $(".canvas").mousemove(function(e) {
    if(draw) {
      var offset = $(this).offset();
      var relX = e.pageX - offset.left;
      var relY = e.pageY - offset.top;
      var msg = relX + "," + relY;
      ws.send(msg);
    }
  });
});