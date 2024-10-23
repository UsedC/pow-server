# pow-server

I've chosen Hashcash as a POW algorithm because of it's simplicity and because it's sufficient for the task. It's also possible to create a mechanism to dynamically modify the difficulty depending on the server load.

There's a Makefile to quickly run tests/builds.

The integration between the server and the client can be tested by running docker-compose. There you can specify the number of requests from the client and the amount of nonces .