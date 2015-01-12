# 42minutes Server API

42minutes should provide a common web interface for various Tv related things:

* Keep track of your favorite Tv Shows.
* Keep track of what you and your friends are watching.
* Stay up to date with what's out.
* Get recomendations based on your friends' suggestions etc.

A daemon is also being simultaneous developed to allow some more weird tings:

* Sync your local Tv shows with 42minutes so you know what you have locally and what's missing.
* Ability to download missing episodes; either one by one, by complete season or complete series.
* Ability to automatically download new episodes of your favorite Tv shows as they become available.

*The daemon needs to communicate with the 42minutes API and cannot be used stand alone.*

## Installtion & Setup

	go build && DB_USER=... DB_PASS=... TRAKT_API_KEY=... TRAKT_ACCESS_TOKEN=... ./42minutes-server-api 
