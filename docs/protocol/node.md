## Econode Methods
Methods for econodes.

#### `fetchNode`
Fetches info about a node.

##### JSON Params
If you want to fetch another node by name and not the logged in client's node,
simply send the name. This is optional.

##### Returns
When successful, returns the following JSON object:  
```json5
{
	"name": "", // name of the node
	"ownerID": 0, // id of the owner
	"balance": 0 // balance of node
}
```

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

