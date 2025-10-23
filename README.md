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
- `BlindConstantTime(input []byte)`: Blinds the input using a constant-time algorithm (POC)
- `Finalize(r, evaluated)`: Client unblinds and hashes the result

## Reference
- [RFC 9497: Oblivious Pseudorandom Functions (OPRFs) Using Prime-Order Groups](https://datatracker.ietf.org/doc/html/rfc9497)

## Constant-Time Blinding (POC)
This implementation includes a proof-of-concept constant-time blinding function `BlindConstantTime(input []byte)`. Function is designed to mitigate timing side-channel attacks by ensuring that the execution time does not depend on the input values.

It is done by padding the input to a fixed length (the maximum expected input length) before processing it. This way, the function always processes the same amount of data, regardless of the actual input length.

### Benchmarking Methodology

The measurements are done using Go's built-in benchmarking framework. To run the benchmark, use the following command:

```bash
go test -bench=<BenchmarkName> -count 10 > bench_results.txt
```

Results are analyzed with the `benchstat` tool from the `golang.org/x/perf/cmd/benchstat` package.

What we want to do is to compare the performance of the constant-time blinding function against the standard blinding function. We do that by feeding them with random inputs of different lengths, running the benchmamrks and then analyzing the results:

```bash
# First we run benchmarks for the standard blinding function
# We need two files:
# First file contains measurements for the function with ascending input lengths
# Second file contains measurements for the function with descending input lengths
# benchstat is comparing the entries in both files based on the name of the benchmark so
# the names must match. 
go test -bench=AscBlind -count 10 > asc_blind.txt
go test -bench=DescBlind -count 10 > desc_blind.txt
sed 's/BenchmarkDescBlind/BenchmarkAscBlind/g' desc_blind.txt > desc_blind_fixed.txt
benchstat asc_blind.txt desc_blind_fixed.txt

# Secondly, we run benchmarks for the constant-time blinding function
go test -bench=AscBlindConstantTime -count 10 > asc_blind_constant_time.txt
go test -bench=DescBlindConstantTime -count 10 > desc_blind_constant_time.txt
sed 's/BenchmarkDescBlindConstantTime/BenchmarkAscBlindConstantTime/g' desc_blind_constant_time.txt > desc_blind_constant_time_fixed.txt
benchstat asc_blind_constant_time.txt desc_blind_constant_time_fixed.txt
```

#### The results

For the standard blinding function, we expect to see a significant difference in execution time between ascending and descending input lengths, indicating that the function's performance is input-dependent. 

The results confirm this expectation (see last column of the measurements below):

```
cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
                    │ asc_blind.txt  │              desc_blind_fixed.txt        │
                    │     sec/op     │     sec/op      vs base                  │
Blind/len_1-8           68.04µ ± 11%   3488.51µ ±  1%  +5027.25% (p=0.000 n=10)
Blind/len_8-8           67.99µ ±  1%   1759.69µ ± 10%  +2488.11% (p=0.000 n=10)
Blind/len_16-8          68.04µ ±  0%    918.19µ ±  9%  +1249.46% (p=0.000 n=10)
Blind/len_32-8          68.40µ ± 13%    494.33µ ±  2%   +622.71% (p=0.000 n=10)
Blind/len_64-8          68.79µ ± 12%    280.92µ ±  1%   +308.35% (p=0.000 n=10)
Blind/len_128-8         68.76µ ±  2%    173.85µ ±  1%   +152.84% (p=0.000 n=10)
Blind/len_256-8         68.88µ ± 12%    121.96µ ±  2%    +77.08% (p=0.000 n=10)
Blind/len_512-8         69.09µ ± 10%     94.86µ ±  1%    +37.30% (p=0.000 n=10)
Blind/len_1024-8        69.95µ ±  1%     81.81µ ± 13%    +16.96% (p=0.000 n=10)
Blind/len_2048-8        71.54µ ±  2%     75.21µ ±  3%     +5.13% (p=0.001 n=10)
Blind/len_4096-8        74.66µ ± 11%     72.04µ ± 12%          ~ (p=0.052 n=10)
Blind/len_8192-8        81.61µ ± 11%     69.79µ ±  1%    -14.48% (p=0.000 n=10)
Blind/len_16384-8       98.06µ ± 11%     69.15µ ± 10%    -29.48% (p=0.000 n=10)
Blind/len_32768-8      121.90µ ± 11%     68.50µ ±  1%    -43.80% (p=0.000 n=10)
Blind/len_65536-8      174.75µ ± 11%     68.61µ ± 10%    -60.74% (p=0.000 n=10)
Blind/len_131072-8     280.22µ ±  9%     68.26µ ±  1%    -75.64% (p=0.000 n=10)
Blind/len_262144-8     493.22µ ± 12%     68.19µ ±  1%    -86.17% (p=0.000 n=10)
Blind/len_524288-8     926.41µ ± 13%     68.82µ ± 12%    -92.57% (p=0.000 n=10)
Blind/len_1048576-8   1764.67µ ± 11%     68.20µ ±  2%    -96.14% (p=0.000 n=10)
Blind/len_2097152-8   3476.79µ ± 12%     68.03µ ±  1%    -98.04% (p=0.000 n=10)
geomean                 147.7µ           147.4µ           -0.15%
```

The constant-time blinding function shows no significant difference in execution time between ascending and descending input lengths, confirming its input-independent performance:

```
cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
                                   │ asc_blindconstanttime.txt │     desc_blindconstanttime_fixed.txt│
                                   │          sec/op           │    sec/op     vs base               │
AscBlindConstantTime/len_1-8                      70.35µ ± 19%   68.10µ ±  2%  -3.20% (p=0.005 n=10)
AscBlindConstantTime/len_8-8                      68.39µ ±  4%   67.75µ ±  3%       ~ (p=0.105 n=10)
AscBlindConstantTime/len_16-8                     69.18µ ± 12%   67.96µ ± 12%       ~ (p=0.315 n=10)
AscBlindConstantTime/len_32-8                     68.33µ ± 11%   67.57µ ±  0%  -1.10% (p=0.000 n=10)
AscBlindConstantTime/len_64-8                     68.40µ ± 11%   67.89µ ±  1%       ~ (p=0.052 n=10)
AscBlindConstantTime/len_128-8                    68.53µ ±  2%   67.65µ ±  3%  -1.29% (p=0.029 n=10)
AscBlindConstantTime/len_256-8                    70.34µ ± 11%   67.65µ ±  1%  -3.82% (p=0.001 n=10)
AscBlindConstantTime/len_512-8                    69.07µ ±  7%   67.58µ ±  0%  -2.17% (p=0.000 n=10)
AscBlindConstantTime/len_1024-8                   67.95µ ±  1%   67.69µ ±  1%       ~ (p=0.218 n=10)
AscBlindConstantTime/len_2048-8                   69.10µ ±  9%   67.68µ ±  1%  -2.05% (p=0.000 n=10)
AscBlindConstantTime/len_4096-8                   68.77µ ±  9%   67.83µ ±  1%  -1.36% (p=0.003 n=10)
AscBlindConstantTime/len_8192-8                   68.35µ ± 10%   67.76µ ±  1%  -0.86% (p=0.000 n=10)
AscBlindConstantTime/len_16384-8                  68.09µ ±  1%   67.79µ ±  1%       ~ (p=0.063 n=10)
AscBlindConstantTime/len_32768-8                  68.45µ ±  9%   67.50µ ±  0%  -1.38% (p=0.002 n=10)
AscBlindConstantTime/len_65536-8                  68.39µ ± 12%   67.99µ ±  1%  -0.60% (p=0.043 n=10)
AscBlindConstantTime/len_131072-8                 68.79µ ±  9%   67.79µ ±  2%  -1.45% (p=0.043 n=10)
AscBlindConstantTime/len_262144-8                 69.34µ ±  1%   67.37µ ±  3%  -2.85% (p=0.009 n=10)
AscBlindConstantTime/len_524288-8                 68.64µ ±  2%   67.31µ ±  1%  -1.94% (p=0.001 n=10)
AscBlindConstantTime/len_1048576-8                68.40µ ±  4%   67.33µ ±  8%  -1.55% (p=0.015 n=10)
AscBlindConstantTime/len_2097152-8                68.38µ ±  0%   67.32µ ±  1%  -1.56% (p=0.000 n=10)
geomean                                           68.76µ         67.67µ        -1.58%
```

## License
MIT
