## User Methods
Here documents the methods which relate to a user in the Econetwork.

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

#### `pm`
Sends a message to someone.

##### JSON Params
Expects a JSON object as follows:  
```json5
{
	"username": "", // username of who to send message
	"message": "" //
}
```

#### `fetchAccount`
Gets info about a user's account.

##### JSON Params
Expects a string of the username. If not provided, will send the session's account.

#### Returns
```json5
{
	"username": "",
	"id": 0, // id of account
	"nodeID": 0, // id of account's node
	"op": false // if the account has op status
}
```

