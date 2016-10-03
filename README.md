# zazkia - tcp proxy for simulating transport errors

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