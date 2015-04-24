Neobombe Simulator
==================

This is a software simulator that attempts to replicate the behavior of Neobombe
which is a modern rebuilt of the famous bombes that were used by the Polish and
British to fight against German's Enigma cipher. Instead of using messages noted
from radios, this software simulator uses Twitter tweets that match certain
`track` keywords.

Build
-----
Before it can be built, your system must have the [Go compiler](https://golang.org/dl). After installing
it you, you can proceed to get the package by running:
`go get -v github.com/jackyb/neobombe-sim`.

Then go to the directory on `$GOPATH/src/github.com/jackyb/neobombe-sim` and
run `go build` to build the program.

Configuration
-------------
The Twitter credentials and track keywords are set in the twitter_settings.json
file. It has the following format:

	{
		"consumer_key": "yourconsumerkey",
		"consumer_secret": "yourconsumersecret",
		"access_token": "youraccesstoken",
		"access_token_secret": "youraccesstokensecret",
		"track": "#thanksobama"
	}

You can change the _track_ value to whatever stream you want to listen to on
Twitter. You can reference the Twitter documentation of
[track](https://dev.twitter.com/streaming/overview/request-parameters#track)
for further information on how to format _track_ to suit your needs.

Run
---
To run the program, type `./neobombe-sim` on Linux/OSX or `neobombe-sim` on
Windows in the directory. Double-clicking it on GUI environment should work too.

Then go to `localhost:8080` on your favorite browser (that supports WebGL).

Interfacing through OSC
-----------------------
The program sends information about the state of the machine through OSC
(Open Sound Control) messages to `localhost:8000` with addresses like
`/rotor/0` up to `/rotor/35`. The message format is `ff` with the first float
being the rotor's current angle and the second float is the rotor's speed.
