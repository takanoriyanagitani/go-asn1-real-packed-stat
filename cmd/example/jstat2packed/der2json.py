import asn1tools
import sys
import json
import struct

packed = asn1tools.compile_files("./pstat.asn")

encoded = sys.stdin.buffer.read()

decoded = packed.decode(
	"PackedStat",
	encoded,
)

count = decoded["count"]
packedStat = decoded["packedStat"]

s = struct.Struct(">dddd")
unpacked = s.unpack(packedStat)

stat = dict(
	count = count,
	minimum = unpacked[0],
	maximum = unpacked[1],
	average = unpacked[2],
	variance = unpacked[3],
)

json.dump(
	stat,
	fp = sys.stdout,
)
