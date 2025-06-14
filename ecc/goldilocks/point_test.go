package goldilocks_test

import (
	"crypto/rand"
	"encoding"
	"testing"

	"github.com/quantumcoinproject/circl/ecc/goldilocks"
	"github.com/quantumcoinproject/circl/internal/test"
)

func randomPoint() *goldilocks.Point {
	var k goldilocks.Scalar
	_, _ = rand.Read(k[:])
	return goldilocks.Curve{}.ScalarBaseMult(&k)
}

func TestPoint(t *testing.T) {
	c := goldilocks.Curve{}
	t.Run("IsOnCurve(ok)", func(t *testing.T) {
		goodGen := c.Generator()
		test.CheckOk(c.IsOnCurve(goodGen), "valid point should pass", t)
	})

	t.Run("IsOnCurve(zero)", func(t *testing.T) {
		var allZeros goldilocks.Point
		test.CheckOk(!c.IsOnCurve(&allZeros), "invalid point should be detected", t)
	})
}

func TestPointAdd(t *testing.T) {
	const testTimes = 1 << 10
	var e goldilocks.Curve
	for i := 0; i < testTimes; i++ {
		P := randomPoint()
		// 16P = 2^4P
		got := e.Double(e.Double(e.Double(e.Double(P))))
		// 16P = P+P...+P
		Q := e.Identity()
		for j := 0; j < 16; j++ {
			Q = e.Add(Q, P)
		}
		want := Q
		if !e.IsOnCurve(got) || !e.IsOnCurve(want) || !got.IsEqual(want) {
			test.ReportError(t, got, want, P)
		}
	}
}

func TestPointNeg(t *testing.T) {
	const testTimes = 1 << 10
	var e goldilocks.Curve
	for i := 0; i < testTimes; i++ {
		P := randomPoint()
		Q := *P
		Q.Neg()
		R := e.Add(P, &Q)
		got := R.IsIdentity()
		want := true
		if got != want {
			test.ReportError(t, got, want, P)
		}
	}
}

func TestPointAffine(t *testing.T) {
	const testTimes = 1 << 10
	for i := 0; i < testTimes; i++ {
		got := randomPoint()
		x, y := got.ToAffine()
		want, err := goldilocks.FromAffine(&x, &y)
		if !got.IsEqual(want) || err != nil {
			test.ReportError(t, got, want)
		}
	}
}

func TestPointMarshal(t *testing.T) {
	const testTimes = 1 << 10
	var want error
	for i := 0; i < testTimes; i++ {
		var P interface{} = randomPoint()
		mar, _ := P.(encoding.BinaryMarshaler)
		data, got := mar.MarshalBinary()
		if got != want {
			test.ReportError(t, got, want, P)
		}
		unmar, _ := P.(encoding.BinaryUnmarshaler)
		got = unmar.UnmarshalBinary(data)
		if got != want {
			test.ReportError(t, got, want, P)
		}
	}
}

func BenchmarkPoint(b *testing.B) {
	P := randomPoint()
	Q := randomPoint()
	b.Run("ToAffine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.ToAffine()
		}
	})
	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.Add(Q)
		}
	})
	b.Run("Double", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.Double()
		}
	})
}
