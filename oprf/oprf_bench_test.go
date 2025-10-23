package oprf

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/gtank/ristretto255"
)

func BenchmarkAscBlind(b *testing.B) {
	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
	benchmars := make([]struct {
		name  string
		input []byte
	}, len(lengths))

	for i, l := range lengths {
		input := make([]byte, l)
		_, _ = rand.Read(input)
		benchmars[i].name = fmt.Sprintf("len_%d", l)
		benchmars[i].input = input
	}

	b.ResetTimer()

	for _, bm := range benchmars {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := Blind(bm.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDescBlind(b *testing.B) {
	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
	benchmars := make([]struct {
		name  string
		input []byte
	}, len(lengths))

	for i, _ := range lengths {
		input := make([]byte, lengths[len(lengths)-1-i])
		_, _ = rand.Read(input)
		benchmars[i].name = fmt.Sprintf("len_%d", lengths[i]) // Note: lengths[i] to force benchstat to compare short with long inputs
		benchmars[i].input = input
	}

	b.ResetTimer()

	for _, bm := range benchmars {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := Blind(bm.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkAscBlindConstantTime(b *testing.B) {
	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
	benchmars := make([]struct {
		name  string
		input []byte
	}, len(lengths))

	for i, l := range lengths {
		input := make([]byte, l)
		_, _ = rand.Read(input)
		benchmars[i].name = fmt.Sprintf("len_%d", l)
		benchmars[i].input = input
	}

	b.ResetTimer()

	for _, bm := range benchmars {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := BlindConstantTime(bm.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDescBlindConstantTime(b *testing.B) {
	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
	benchmars := make([]struct {
		name  string
		input []byte
	}, len(lengths))

	for i, _ := range lengths {
		input := make([]byte, lengths[len(lengths)-1-i])
		_, _ = rand.Read(input)
		benchmars[i].name = fmt.Sprintf("len_%d", lengths[i]) // Note: lengths[i] to force benchstat to compare short with long inputs
		benchmars[i].input = input
	}

	b.ResetTimer()

	for _, bm := range benchmars {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := BlindConstantTime(bm.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDescEvaluate(b *testing.B) {
	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
	benchmars := make([]struct {
		name  string
		input *ristretto255.Element
		sk    *ristretto255.Scalar
	}, len(lengths))

	for i, _ := range lengths {
		input := make([]byte, lengths[len(lengths)-1-i])
		_, _ = rand.Read(input)
		_, blinded, _ := Blind(input)

		keys, _ := NewKeyPair()

		benchmars[i].name = fmt.Sprintf("len_%d", lengths[i])
		benchmars[i].input = blinded
		benchmars[i].sk = keys.Private
	}

	b.ResetTimer()

	for _, bm := range benchmars {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				BlindEvaluate(bm.sk, bm.input)
			}
		})
	}
}

func BenchmarkAscEvaluate(b *testing.B) {
	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
	benchmars := make([]struct {
		name  string
		input *ristretto255.Element
		sk    *ristretto255.Scalar
	}, len(lengths))

	for i, l := range lengths {
		input := make([]byte, l)
		_, _ = rand.Read(input)
		_, blinded, _ := Blind(input)

		keys, _ := NewKeyPair()

		benchmars[i].name = fmt.Sprintf("len_%d", l)
		benchmars[i].input = blinded
		benchmars[i].sk = keys.Private
	}

	b.ResetTimer()

	for _, bm := range benchmars {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				BlindEvaluate(bm.sk, bm.input)
			}
		})
	}
}

// func BenchmarkHashToGroup(b *testing.B) {
// 	lengths := []int{1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152}
// 	benchmars := make([]struct {
// 		name  string
// 		input []byte
// 	}, len(lengths))

// 	for i, l := range lengths {
// 		input := make([]byte, l)
// 		_, _ = rand.Read(input)
// 		benchmars[i].name = fmt.Sprintf("len_%d", l)
// 		benchmars[i].input = input
// 	}

// 	b.ResetTimer()

// 	for _, bm := range benchmars {
// 		b.Run(bm.name, func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				_, err := hashToGroup(bm.input)
// 				if err != nil {
// 					b.Fatal(err)
// 				}
// 			}
// 		})
// 	}
// }
