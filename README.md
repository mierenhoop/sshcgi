# sshcgi (PROOF OF CONCEPT)
CGI like handling of connections in an SSH server.

It's not very useful or secure (yet) and only expirimental.

## Use cases
* Fetch data from a server by just ssh'ing into a server as an alternative to curl.
* Make fun multiplayer games to run in the terminal.
* Probably many more...

## Improvements
* I would prefer if this project is a wrapper (like fcgiwrap/spawn-fcgi/slowcgi for http) for ssh servers (like sshd) instead of what it currently is: a simple ssh server written in Go.
* Possibly make modifications to x/crypto/ssh to expose more internals.
