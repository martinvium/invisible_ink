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
      // var link = '<li><a href="' + url + '">' + id + '</a></li>';
      var img = '<div class="image"><img src="' + url + '"/></div>';
      $(".log").prepend(img);
    }
  });

  $(".canvas").mousemove(function(event) {
    if(draw) {
      var msg = event.pageX + "," + event.pageY;
      ws.send(msg);
    }
  });
});