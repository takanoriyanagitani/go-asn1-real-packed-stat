package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	ps "github.com/takanoriyanagitani/go-asn1-real-packed-stat"
	. "github.com/takanoriyanagitani/go-asn1-real-packed-stat/util"
)

func ReaderToBytesLimited(limit int64) func(io.Reader) IO[[]byte] {
	return Lift(func(rdr io.Reader) ([]byte, error) {
		limited := &io.LimitedReader{
			R: rdr,
			N: limit,
		}
		var buf bytes.Buffer
		_, e := io.Copy(&buf, limited)
		return buf.Bytes(), e
	})
}

var jsonStatBytesLimit int64 = 1048576

var reader2bytes func(io.Reader) IO[[]byte] = ReaderToBytesLimited(
	jsonStatBytesLimit,
)

var jsonBytesStdin IO[[]byte] = reader2bytes(os.Stdin)

var parsedStat IO[ps.SimpleStat] = Bind(
	jsonBytesStdin,
	Lift(ps.SimpleStatFromJson),
)

var packedDerBytes IO[[]byte] = Bind(
	parsedStat,
	Lift(ps.SimpleStatToPackedDerBytes),
)

func BytesToWriter(wtr io.Writer) func([]byte) IO[Void] {
	return Lift(func(dat []byte) (Void, error) {
		_, e := wtr.Write(dat)
		return Empty, e
	})
}

var bytes2stdout func([]byte) IO[Void] = BytesToWriter(os.Stdout)

var stdin2jstat2der2stdout IO[Void] = Bind(
	packedDerBytes,
	bytes2stdout,
)

func main() {
	_, e := stdin2jstat2der2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
