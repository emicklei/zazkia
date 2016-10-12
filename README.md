# zazkia - tcp proxy for simulating transport errors

Zazkia is meant to intercept byte sequences send and received between a client and a service.
By specifying routes, you can tell zazkia on what ports it should listen and upon connection to what service it needs to connect.

The tranport part of the route can be used to setup the initial behavior of new connection pairs (called links).


### routes.json example
	[
	    {
	        "label": "oracle",
	        "target-domain": "some-host-name",
	        "target-port": 1521,
	        "listen-port": 49997,
	        "transport": {
	            "delay-target-response": 1000,
	            "sending-to-client": true,
	            "receivingFromClient": true,
	            "sendingToService": true,
	            "receivingFromService": true
	        }
	    }
	]
	
	
### How to build, test and run it

	sh run.sh	