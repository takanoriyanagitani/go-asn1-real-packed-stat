package pstat

import (
	"encoding/asn1"
	"encoding/binary"
	"encoding/json"
	"math"
)

type RawDouble [8]byte

type PackedDouble256 [32]byte

func RawToPacked(raw [4]RawDouble) PackedDouble256 {
	var ret [32]byte
	for i := range 4 {
		var start int = i * 8
		var end int = start + 8
		var dst []byte = ret[start:end]
		var src [8]byte = raw[i]
		copy(dst, src[:])
	}
	return ret
}

func (p PackedDouble256) ToRawDoubles() [4]RawDouble {
	var ret [4]RawDouble
	for i := range 4 {
		var start int = i * 8
		var end int = start + 8
		var s []byte = p[start:end]
		copy(ret[i][:], s)
	}
	return ret
}

func (r RawDouble) ToDouble() float64 {
	var s []byte = r[:]
	var u uint64 = binary.BigEndian.Uint64(s)
	return math.Float64frombits(u)
}

type Double float64

func (d Double) ToBeBytes() [8]byte {
	var u uint64 = math.Float64bits(float64(d))
	var ret [8]byte
	var s []byte = ret[:]
	binary.BigEndian.PutUint64(s, u)
	return ret
}

type SimpleStatJson []byte

func (j SimpleStatJson) ToStat() (SimpleStat, error) {
	var ret SimpleStat
	e := json.Unmarshal(j, &ret)
	return ret, e
}

func SimpleStatFromJson(jbytes []byte) (SimpleStat, error) {
	return SimpleStatJson(jbytes).ToStat()
}

// Simple stat which can't be directly converted to der(using encoding/asn1).
type SimpleStat struct {
	Count    int64   `json:"count"`
	Minimum  float64 `json:"minimum"`
	Maximum  float64 `json:"maximum"`
	Average  float64 `json:"average"`
	Variance float64 `json:"variance"`
}

func (s SimpleStat) ToCount() int64 { return s.Count }

func (s SimpleStat) RawMaximum() RawDouble {
	return Double(s.Maximum).ToBeBytes()
}

func (s SimpleStat) RawMinimum() RawDouble {
	return Double(s.Minimum).ToBeBytes()
}

func (s SimpleStat) RawAverage() RawDouble {
	return Double(s.Average).ToBeBytes()
}

func (s SimpleStat) RawVariance() RawDouble {
	return Double(s.Variance).ToBeBytes()
}

func (s SimpleStat) ToPacked() PackedDouble256 {
	var src [4]RawDouble
	src[0] = s.RawMinimum()
	src[1] = s.RawMaximum()
	src[2] = s.RawAverage()
	src[3] = s.RawVariance()
	return RawToPacked(src)
}

func (s SimpleStat) WithCount(count int64) SimpleStat {
	s.Count = count
	return s
}

// Converts to a der serializable [PackedStat].
func (s SimpleStat) ToPackedStat() PackedStat {
	var ret PackedStat

	ret.Count = s.Count

	var buf [32]byte
	ret.PackedRealStat = buf[:]

	var packed PackedDouble256 = s.ToPacked()
	copy(ret.PackedRealStat, packed[:])
	return ret
}

func (s SimpleStat) ToPackedDerBytes() ([]byte, error) {
	var packed PackedStat = s.ToPackedStat()
	return packed.toDerBytes()
}

func SimpleStatToPackedDerBytes(s SimpleStat) ([]byte, error) {
	return s.ToPackedDerBytes()
}

// Creates [SimpleStat] with empty Count(0).
func PackedToCount(packed PackedDouble256) SimpleStat {
	var ret SimpleStat
	var doubles [4]RawDouble = packed.ToRawDoubles()

	ret.Minimum = doubles[0].ToDouble()
	ret.Maximum = doubles[1].ToDouble()
	ret.Average = doubles[2].ToDouble()
	ret.Variance = doubles[3].ToDouble()

	return ret
}

// Der serializable simple stat.
type PackedStat struct {
	Count          int64
	PackedRealStat []byte
}

func (s PackedStat) ToPackedRealStat() PackedDouble256 {
	var ret PackedDouble256
	copy(ret[:], s.PackedRealStat)
	return ret
}

func (s PackedStat) toDerBytes() ([]byte, error) {
	return asn1.Marshal(s)
}

func (s PackedStat) ToStat() SimpleStat {
	var packed PackedDouble256 = s.ToPackedRealStat()
	return PackedToCount(packed).
		WithCount(s.Count)
}

type PackedStatDer []byte

func (d PackedStatDer) ToStat() (PackedStat, error) {
	var ret PackedStat
	var s []byte = d
	_, e := asn1.Unmarshal(s, &ret)
	return ret, e
}
