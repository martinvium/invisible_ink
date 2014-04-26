$(document).ready(function() {
  var ws = new WebSocket("ws://localhost:8080/save");
  var draw = false;

  $('.canvas').mousedown(function() {
    var id = Math.random().toString(36).substr(2, 9);
    ws.send('STARTDRAW: ' + id);
    draw = true;
  });

  $(document).mouseup(function() {
    draw = false;
    ws.send('ENDDRAW');
  });

  $(".canvas").mousemove(function(event) {
    if(draw) {
      var msg = event.pageX + "," + event.pageY;
      ws.send(msg);
    }
  });
});