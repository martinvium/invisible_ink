Invisible Ink
=============

Small demo that captures mousemove events, sends them through a Websocket to a small Go webservice, which in turn
renders the coordinates as dots in a PNG.

![Example](/assets/images/example.png?raw=true)

Installation
------------

    create table coordinates (coordinate_id int primary key, drawing_id text, x int, y int); 

Build and Run
-------------

    go build && ./invisible_ink

Visit http://127.0.0.1:8080/ in your browser.