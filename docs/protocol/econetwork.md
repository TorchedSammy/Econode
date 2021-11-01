## Econetwork Moments
General methods which interact with the Econetwork itself.

#### `welcome`
The `welcome` method is an incoming response type which will be sent
whenever a client joins. It mainly has the MOTD which is used as a general greeting.

##### JSON Params
Sends the following object:  
```json5
{
	"motd": ""
}
```

#### `stats`
Gives stats about the Econetwork.

#### Returns
```json5
{
	"nodes": 0, // amount of nodes
	"accounts": 0, // number of created accounts
	"networkVersion": "impl v0.0.0" // version of implementation. for the base
									// implementation it will be "Econode vX.X.X"
}
```

