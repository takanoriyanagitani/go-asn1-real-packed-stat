package pstat_test

import (
	"testing"

	ps "github.com/takanoriyanagitani/go-asn1-real-packed-stat"
)

func TestPacked(t *testing.T) {
	t.Parallel()

	t.Run("RawDouble", func(t *testing.T) {
		t.Parallel()

		t.Run("positive zero", func(t *testing.T) {
			t.Parallel()

			var rd ps.RawDouble
			var f float64 = rd.ToDouble()
			if 0.0 != f {
				t.Fatalf("unexpected value: %v\n", f)
			}
		})

		t.Run("positive half", func(t *testing.T) {
			t.Parallel()

			var rd ps.RawDouble
			rd[0] = 0x3f
			rd[1] = 0xe0
			var f float64 = rd.ToDouble()
			if 0.5 != f {
				t.Fatalf("unexpected value: %v\n", f)
			}
		})

		t.Run("negative half", func(t *testing.T) {
			t.Parallel()

			var rd ps.RawDouble
			rd[0] = 0xbf
			rd[1] = 0xe0
			var f float64 = rd.ToDouble()
			if -0.5 != f {
				t.Fatalf("unexpected value: %v\n", f)
			}
		})
	})

	t.Run("Double", func(t *testing.T) {
		t.Parallel()

		t.Run("positive zero", func(t *testing.T) {
			t.Parallel()

			var d ps.Double
			var rd ps.RawDouble = d.ToBeBytes()
			if 0.0 != rd.ToDouble() {
				t.Fatalf("unexpected value: %v\n", rd.ToDouble())
			}
		})

		t.Run("positive half", func(t *testing.T) {
			t.Parallel()

			var d ps.Double = 0.5
			var rd ps.RawDouble = d.ToBeBytes()
			if 0.5 != rd.ToDouble() {
				t.Fatalf("unexpected value: %v\n", rd.ToDouble())
			}
		})

		t.Run("negative half", func(t *testing.T) {
			t.Parallel()

			var d ps.Double = -0.5
			var rd ps.RawDouble = d.ToBeBytes()
			if -0.5 != rd.ToDouble() {
				t.Fatalf("unexpected value: %v\n", rd.ToDouble())
			}
		})
	})

	t.Run("PackedDouble256", func(t *testing.T) {
		t.Parallel()

		t.Run("zero", func(t *testing.T) {
			t.Parallel()

			var packed ps.PackedDouble256
			var rd4 [4]ps.RawDouble = packed.ToRawDoubles()
			for i := range 4 {
				var rd ps.RawDouble = rd4[i]
				if 0.0 != rd.ToDouble() {
					t.Fatalf("unexpected value: %v\n", rd.ToDouble())
				}
			}
		})
	})

	t.Run("SimpleStat", func(t *testing.T) {
		t.Parallel()

		t.Run("zero", func(t *testing.T) {
			t.Parallel()

			var s ps.SimpleStat

			var zero ps.PackedDouble256
			var packed ps.PackedDouble256 = s.ToPacked()

			if zero != packed {
				t.Fatal("zero value expected")
			}
		})

		t.Run("stat", func(t *testing.T) {
			t.Parallel()

			stat := ps.SimpleStat{
				Count:    634,
				Minimum:  0.599,
				Maximum:  3.776,
				Average:  3.141,
				Variance: 0.333,
			}
			var packed ps.PackedDouble256 = stat.ToPacked()

			var rd4 [4]ps.RawDouble = packed.ToRawDoubles()

			expected := [4]float64{
				0.599,
				3.776,
				3.141,
				0.333,
			}

			for i := range 4 {
				var rd ps.RawDouble = rd4[i]
				var ex float64 = expected[i]
				if ex != rd.ToDouble() {
					t.Errorf("expected %v != got %v", ex, rd.ToDouble())
				}
			}
		})

		t.Run("ToPackedDerBytes", func(t *testing.T) {
			t.Parallel()

			t.Run("zero", func(t *testing.T) {
				t.Parallel()

				var zero ps.SimpleStat
				der, e := zero.ToPackedDerBytes()
				if nil != e {
					t.Fatalf("unexpected error: %v\n", e)
				}

				if 0 == len(der) {
					t.Fatal("zero bytes got")
				}

				stat, e := ps.PackedStatDer(der).ToStat()
				if nil != e {
					t.Fatalf("unexpected err: %v\n", e)
				}

				if 0 != stat.Count {
					t.Fatalf("unexpected count: %v", stat.Count)
				}

				var ss ps.SimpleStat = stat.ToStat()
				var zs ps.SimpleStat
				if ss != zs {
					t.Fatal("zero stat expected")
				}
			})

			t.Run("stat", func(t *testing.T) {
				t.Parallel()

				stat := ps.SimpleStat{
					Count:    299792458,
					Minimum:  2.99792458,
					Maximum:  3.776,
					Average:  3.14,
					Variance: 0.634,
				}

				der, e := stat.ToPackedDerBytes()
				if nil != e {
					t.Fatalf("unexpected error: %v", e)
				}

				parsed, e := ps.PackedStatDer(der).ToStat()
				if nil != e {
					t.Fatalf("unexpected error: %v", e)
				}

				var ss ps.SimpleStat = parsed.ToStat()
				if ss != stat {
					t.Fatalf("unexpected stat: %v", ss)
				}
			})
		})
	})

	t.Run("PackedStatDer", func(t *testing.T) {
		t.Parallel()

		t.Run("ToStat", func(t *testing.T) {
			t.Parallel()

			var der []byte
			der = append(der, 0x30)
			der = append(der, 0x25)

			der = append(der, 0x02)
			der = append(der, 0x01)
			der = append(der, 0x00)

			der = append(der, 0x04)
			der = append(der, 0x20)
			var packed [32]byte
			der = append(der, packed[:]...)

			stat, e := ps.PackedStatDer(der).ToStat()
			if nil != e {
				t.Errorf("unexpected error: %v", e)
				t.Fatalf("bytes: %v", der)
			}

			if 0 != stat.Count {
				t.Fatalf("unexpected count: %v\n", stat.Count)
			}

			var zero ps.PackedDouble256
			if zero != stat.ToPackedRealStat() {
				t.Fatalf("zero value expected: %v", stat.ToPackedRealStat())
			}
		})
	})
}
