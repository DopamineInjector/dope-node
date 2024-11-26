# dope-node
Blockchain node

## API for external use (love women, get money #thuglife)
### POST /api/account
Ayo we be creating dat account my fella. We be posting that keypair.
Body: 
```ts
class CreateAccountRequest {
    publicKey: string,
    privateKey: string,
}
```

### PUT /api/account/info
Because honestly who cares about REST convention #swampofpox.
Body:
```ts
class GetAccountRequest {
    publicKey: string
}
```
Response:
```ts
class GetAccountResponse {
    publicKey: string,
    balance: number,
}
```

### POST /api/transfer
We be sending that money to fellas, we be trappin forreal #TakeRisksAndProsper
```ts
class TransferRequest {
    payload: {
        sender: string,
        recipient: string,
        amount: number
    },
    signature: string // Signed stringified payload
}
```

Returns 403 if signature is bad

### POST /api/smartContract
We be on that SC grind, a lil scam if ya know what I be sayin. We out here, we in the streets.
```ts
class SmartContractRequest {
    payload: {
        sender: string,
        contract: string, // Smart contract address on da blockchain
        entrypoint: string, // SC entrypoint
        args: string, // Stringified args to the function
    }
    signature?: string // Signed stringified payload
    view?: boolean // if false or undefined, then it is treated as normal transaction
}
```
Response:
```ts
class SmartContractResponse {
    output: string
}
```
