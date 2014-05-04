Invisible Ink
=============

Small demo that captures mousemove events, sends them through a Websocket to a small Go webservice that writes the 
coordinates per canvas to a Cassandra instance. It then requests the image for the History, which is rendered from the
data stored in Cassandra.

![Example](/assets/images/example.png?raw=true)

Installation
------------

    CREATE KEYSPACE mykeyspace WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
    
    CREATE TABLE coordinates (
      drawing_id text,
      timestamp timeuuid,
      x int, 
      y int,
      PRIMARY KEY (drawing_id, timestamp)
    );

Build and Run
-------------

    go build && ./invisible_ink

Visit http://127.0.0.1:8080/ in your browser.

Test Installation
-------------

Before you can run tests you must create a keyspace and tables:

    CREATE KEYSPACE testkeyspace WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

Now run the CREATE TABLE statement from installation.