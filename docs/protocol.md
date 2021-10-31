# Econode Protocol
Here is documented the protocol for an Econode client and server.

> ⚠️  Work In Progress!

## Client Response
`!%` = required
```js
{
	"sessionID": "d4a44b2ef2da23", // randomly generated session id string
	"method": "", // !% - our route basically, what to do with the data
	"data": null // can be any json type
}
```

## Server Response
```js
{
	"code": "", // %! - string of a result code, ie `success`, `failure`, `forbidden`
	"method": "", // %! - method the client sent which caused this response
	"data": null // any other data the server wants to send, any json type
}
```
## Session IDs
Session IDs are the way clients will authorize with the Econode server. They are used
for most methods for authentication and to determine what client to act on.
When a client logins in or registers, they will be sent the session ID as the single string
argument to the `data` field.

## Response Codes
An Econode server will send a response code for every method and request. They are as
follows:
- `success` - Request has completed successfully.
- `accepted` - Request has completed successfully.
  Like the `success` response, but doesn't send any `data`.
- `fail` - Request finished but failed.
- `error` - An error occurred during the processing of a method. Usually an internal
  server error.
- `forbidden` - Client is not allowed to use method.
- `malformed` - Request is missing fields.

## Methods
Methods are, as explained in the comments of responses, basically like HTTP routes.  
They determine what happens with the data.

Each of the JSON Params goes into the `data` field. Since it accepts any JSON type,
some of these methods may require just a string or integer instead of a full JSON object.

Optional fields will have comments that start with `%?`.

Server responses with the `error` code usually only have a single error message in the
`data` field.

#### `register`
The `register` method is fairly self explainable, as are most (or all) other methods.  
This is used to register a *user* to the Econetwork.

##### JSON Params
Expects a JSON object as follows:  
```json5
{
	"username": "", // user's username
	"password": "" // and their password
}
```

##### Returns
When successful, it returns a response with the `success` code, and the session ID
in the `data` field as a string.

#### `login`
The `login` method logs a user into the Econetwork.

##### JSON Params
Expects a JSON object as follows:  
```json5
{
	"username": "", // user's username
	"password": "" // and their password
}
```

##### Returns
When successful, it returns a response with the `success` code, and the session ID
in the `data` field as a string.

#### `buyItem`
Buys an [item](store.md) for the logged in user's node.

##### JSON Params
Expects a JSON object as follows:  
```json5
{
	"itemname": "", // name of the item
	"amount": 1 // amount to buy
}
```

##### Returns
When successful, it gives a response with the `success` code.

