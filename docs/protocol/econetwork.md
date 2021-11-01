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

