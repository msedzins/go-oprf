# go-oprf

An implementation of the Oblivious Pseudorandom Function (OPRF) protocol using Ristretto255, based on [RFC 9497](https://datatracker.ietf.org/doc/html/rfc9497).

## Features
- OPRF key generation
- Input blinding and evaluation
- Output finalization
- Follows RFC 9497

## Usage Example

```go
package main

import (
	"fmt"
	"github.com/msedzins/go-oprf/oprf"
)

func main() {
	// Generate server keypair
	key, err := oprf.NewKeyPair()
	if err != nil {
		panic(err)
	}

	// Client blinds input
	input := []byte("my secret")
	r, blinded, err := oprf.Blind(input)
	if err != nil {
		panic(err)
	}

	// Server evaluates
	evaluated := oprf.BlindEvaluate(key.Private, blinded)

	// Client finalizes
	output := oprf.Finalize(r, evaluated)
	fmt.Printf("OPRF output: %x\n", output)
}
```

## API
- `NewKeyPair()`: Generates a random OPRF keypair
- `Blind(input []byte)`: Blinds the input and returns blinding scalar and blinded element
- `BlindEvaluate(sk, blinded)`: Server multiplies blinded element by secret key
- `Finalize(r, evaluated)`: Client unblinds and hashes the result

## Reference
- [RFC 9497: Oblivious Pseudorandom Functions (OPRFs) Using Prime-Order Groups](https://datatracker.ietf.org/doc/html/rfc9497)

## License
MIT
