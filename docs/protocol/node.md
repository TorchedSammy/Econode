## Econode Methods
Methods for econodes.

#### `buyItem`
Buys an [item](../store.md) for the logged in user's node.

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

