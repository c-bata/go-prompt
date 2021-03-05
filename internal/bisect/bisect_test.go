package bisect_test

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"

	crypto_rand "crypto/rand"
	math_rand "math/rand"

	"github.com/c-bata/go-prompt/internal/bisect"
)

func Example() {
	in := []int{1, 2, 3, 3, 3, 6, 7}
	fmt.Println("Insertion position for 0 in the slice is", bisect.Right(in, 0))
	fmt.Println("Insertion position for 4 in the slice is", bisect.Right(in, 4))

	// Output:
	// Insertion position for 0 in the slice is 0
	// Insertion position for 4 in the slice is 5
}

func TestBisectRight(t *testing.T) {
	// Thanks!! https://play.golang.org/p/y9NRj_XVIW
	in := []int{1, 2, 3, 3, 3, 6, 7}

	r := bisect.Right(in, 0)
	if r != 0 {
		t.Errorf("number 0 should inserted at 0 position, but got %d", r)
	}

	r = bisect.Right(in, 4)
	if r != 5 {
		t.Errorf("number 4 should inserted at 5 position, but got %d", r)
	}
}

func BenchmarkRight(b *testing.B) {
	// https://stackoverflow.com/a/54491783/6309
	var bb [8]byte
	_, err := crypto_rand.Read(bb[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	// Avoid G404: Use of weak random number generator (math/rand instead of crypto/rand)
	math_rand.Seed(int64(binary.LittleEndian.Uint64(bb[:])))

	for _, l := range []int{10, 1e2, 1e3, 1e4} {
		x := rand.Perm(l)
		insertion := rand.Int()

		b.Run(fmt.Sprintf("arrayLength=%d", l), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				bisect.Right(x, insertion)
			}
		})
	}
}
