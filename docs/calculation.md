# Calculations
Various calculations that an Econetwork implementation has to do: like CPS, item
prices, and the rest.

> ⚠️  Work In Progress!

# Item Pricing
Item prices should be calculated like:  
`(I * (1.20^N - 1)) / 0.20`  
Where `I` is an item's base cost and
`N` is the amount someone has + how much they want to get (which is usually 1)  

As an example, to buy a 2nd electron:  
I = 20, N = 2   
`(20 * (1.20^2 - 1)) / 0.20 = 44`, meaning a user will need $44 to get a 2nd electron.

