Gobrid Client
=================
A CLI client for RealDebrid.

# Currently Supported
1. Adding torrents via magnet link
2. Downloading torrents

# Installation
`go install github.com/simantovyousoufov/rd-client/gobrid`

# Usage
Single link:
`gobrid magnet magnet:?xt=somemagnetlink....`

Multiple links from file where each magnet link is separated by a new line character:
`gobrid magnet -f filename.txt`