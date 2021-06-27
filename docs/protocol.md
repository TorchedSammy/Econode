# Econode Protocol
This document is WIP, it'll be better soon

## Client Response
`!%` = required
```json5
{
	"sessionID": "12weejr", // !% - randomly generated session id string
	"method": "", // !% - our route basically, what to do with the data
	"data": null // can be any json type
```

## Server Response
```json5
{
	"result": "", // %! - string of a result code, ie `success`, `failure`, `forbidden`
	"method": "", // %! - method the client sent which caused this response
	"data": null // any other data the server wants to send, any json type
}
```

