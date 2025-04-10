#!/bin/sh

input=./input.stat.json
output=./output.stat.asn1.der.dat

jsonStat(){
	echo generating input json...

	jq -c -n '{
		count: 42,
		minimum: 0.599,
		maximum: 3.776,
		average: 3.14,
		variance: 2.99792458,
	}' |
		dd \
			if=/dev/stdin \
			of="${input}" \
			status=none
}

test -f "${input}" || jsonStat

echo converting json to der bytes...
cat "${input}" |
	./jstat2packed |
	dd \
		if=/dev/stdin \
		of="${output}" \
		status=none

echo converting der bytes to json using asn1tools...
cat "${output}" |
	python3 der2json.py |
	jq

ls -l "${input}" "${output}"
