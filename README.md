Invisible Ink
=============

Small demo that captures mousemove events, sends them through a Websocket to a small Go webservice that writes the 
coordinates per canvas to a Cassandra instance. It then requests the image for the History, which is rendered from the
data stored in Cassandra.

![Example](/assets/images/example.png?raw=true)

Installation
------------

    create table coordinates (id uuid primary key, drawing_id text, x int, y int); 
    create index on coordinates (drawing_id);

Build and Run
-------------

    go build && ./invisible_ink

Visit http://127.0.0.1:8080/ in your browser.