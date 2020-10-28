package uuid

import (
	"crypto/rand"
	math_rand "math/rand"
	"time"
)

/*

UUID is a universally unique identifier

UUIDs are binary sortable.
The first 6 bytes constitutes a millisecond-precision timestamp in big-endian byte order.

Data layout:

	Byte 0-3  timestamp second, big endian
	Byte 4-5  timestamp millisecond, big endian
	Byte 6-15 random

Data layout example:

	00 31 04 39 02 c9 39 ce 14 6c 0b db a1 40 77 78
	~~~~~~~~~~~ ----- ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	 |            |           Random bytes
	 |            |
	 |           713 milliseconds
	 |
	3212345 seconds since 2020-09-13 12:26:40
	                    = 2020-10-20 16:45:45.713 UTC

Note that while it seems like we could use nanosecond for the timestamp to reduce the
random data needed, environments like JavaScript doesn't necessarily provide high-precision
clocks. Doing things this way means that we can generate and parse the embedded timestamp
in a wide variety of programming languages.

*/
type UUID [16]byte

// Min is the zero UUID (must not be modified)
var Min UUID

// Max is the largest possible UUID (must not be modified)
var Max = UUID{
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
}

// StringMaxLen is the maximum length of a string representation of a UUID,
// i.e. as returned by UUID.String()
const StringMaxLen = 22

// idEpochBase offsets the timestamp to provide a wider range.
// Effective range (0x0–0xFFFFFFFF): 2020-09-13 12:26:40 – 2156-10-20 18:54:55 (UTC)
const idEpochBase int64 = 1600000000

// Gen generates a universally unique UUID suitable to be used for sorted identity
func Gen() UUID {
	var id UUID

	t := time.Now()
	sec := uint32(t.Unix() - idEpochBase)
	ns := uint64(t.Nanosecond())
	ms := uint16(ns / uint64(time.Millisecond))

	// second part
	id[0] = byte(sec >> 24)
	id[1] = byte(sec >> 16)
	id[2] = byte(sec >> 8)
	id[3] = byte(sec)

	// millisecond part
	id[4] = byte(ms >> 8)
	id[5] = byte(ms)

	// Use middle bytes of nanosecond to reduce need for random bytes.
	// We pick the middle bytes so that we don't have to know the endianess of the host.
	// Note that Windows uses a low-res timer for time.Now (Oct 2020)
	// See https://go-review.googlesource.com/c/go/+/227499/ + github issue for discussion,
	// see https://go-review.googlesource.com/c/go/+/227499/1/src/testing/time_windows.go for patch.
	id[6] = byte(ns >> 24)
	id[7] = byte(ns >> 16)

	// rest are random bytes
	if _, err := rand.Read(id[8:16]); err != nil {
		// If crypto/rand fails, fall back to pseudo random number generator.
		// This is fine since the id is not used for anything critical and its uniqueness
		// is eventually verified (i.e. when inserting into a database.)
		math_rand.Read(id[8:16])
	}

	return id
}

// New creates a new UUID with specific Unix timestamp and random bytes.
//
// nsec is the nanosecond part of the timestamp and should be in the range [0, 999999999].
// It's valid to pass values outside this range for nsec. Only the millisecond part of nsec
// is actually used.
//
// To create an UUID with a time.Time object, do this:
//   New(t.Unix(), t.Nanosecond(), random)
//
// Up to 10 bytes is used from random.
// If len(random) < 10, the remaining "random" bytes of UUID are zero.
//
func New(sec int64, nsec int, random []byte) UUID {
	var id UUID

	s := uint32(sec - idEpochBase)
	ms := uint16(nsec / int(time.Millisecond))

	// second part
	id[0] = byte(s >> 24)
	id[1] = byte(s >> 16)
	id[2] = byte(s >> 8)
	id[3] = byte(s)

	// millisecond part
	id[4] = byte(ms >> 8)
	id[5] = byte(ms)

	copy(id[6:], random)

	return id
}

// FromBytes copies verbatim bytes into an UUID and returns that UUID.
// verbatim must be at least 16 bytes long or this will panic.
func FromBytes(verbatim []byte) UUID {
	var id UUID
	copy(id[:], verbatim[:16]) // intentionally 16 to cause panic for too-small arguments
	return id
}

// FromString decodes a string representation of an UUID (i.e. from String())
func FromString(encoded string) UUID {
	var id UUID
	id.DecodeString([]byte(encoded))
	return id
}

// String returns a string representation of the UUID.
// The returned string is sortable with the same order as the "raw" UUID bytes and is URL safe.
func (id UUID) String() string {
	var buf [StringMaxLen]byte
	n := id.EncodeString(buf[:])
	return string(buf[n:])
}

// Bytes returns the IDs natural 16 byte long value.
// The returned slice's bytes must not be modified.
func (id UUID) Bytes() []byte {
	return id[:]
}

// Time returns the time portion of the UUID
func (id UUID) Time() time.Time {
	sec, ms := id.Timestamp()
	return time.Unix(int64(sec)+idEpochBase, int64(ms)*int64(time.Millisecond))
}

// Timestamp returns the timestamp portion of the UUID
func (id UUID) Timestamp() (sec uint32, millisec uint16) {
	sec = uint32(id[0])<<24 | uint32(id[1])<<16 | uint32(id[2])<<8 | uint32(id[3])
	millisec = uint16(id[4])<<8 | uint16(id[5])
	return
}

/*
EncodeString and DecodeString have been adapted from the ksuid project,
licensed as follows:

MIT License

Copyright (c) 2017 Segment.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

const base62Characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// EncodeString writes the receiver to dst which must be at least StringMaxLen (22) bytes.
// Returns the start offset (this function starts writing at the end of dst.)
func (id UUID) EncodeString(dst []byte) int {
	const srcBase = 0x100000000
	const dstBase = 62

	parts := [4]uint32{
		uint32(id[0])<<24 | uint32(id[1])<<16 | uint32(id[2])<<8 | uint32(id[3]),
		uint32(id[4])<<24 | uint32(id[5])<<16 | uint32(id[6])<<8 | uint32(id[7]),
		uint32(id[8])<<24 | uint32(id[9])<<16 | uint32(id[10])<<8 | uint32(id[11]),
		uint32(id[12])<<24 | uint32(id[13])<<16 | uint32(id[14])<<8 | uint32(id[15]),
	}

	n := len(dst)
	bp := parts[:]
	bq := [4]uint32{}
	dst[0] = '0'

	for len(bp) != 0 {
		quotient := bq[:0]
		remainder := uint64(0)

		for _, c := range bp {
			value := uint64(c) + uint64(remainder)*srcBase
			digit := value / dstBase
			remainder = value % dstBase

			if len(quotient) != 0 || digit != 0 {
				quotient = append(quotient, uint32(digit))
			}
		}

		// Writes at the end of the destination buffer because we computed the
		// lowest bits first.
		n--
		dst[n] = base62Characters[remainder]
		bp = quotient
	}

	return n
}

// DecodeString sets the receiving UUID to the decoded value of src, which is expected to be a
// string previously encoded using EncodeString (base62 0-9A-Za-z)
func (id *UUID) DecodeString(src []byte) {
	const srcBase = 62
	const dstBase = 0x100000000

	parts := [StringMaxLen]byte{}

	partsIndex := 21
	for i := len(src); i > 0; {
		// offsets into base62Characters
		const offsetUppercase = 10
		const offsetLowercase = 36

		i--
		b := src[i]
		switch {
		case b >= '0' && b <= '9':
			b -= '0'
		case b >= 'A' && b <= 'Z':
			b = offsetUppercase + (b - 'A')
		default:
			b = offsetLowercase + (b - 'a')
		}
		parts[partsIndex] = b
		partsIndex--
	}

	n := len(id)
	bp := parts[:]
	bq := make([]byte, 0, len(src))

	for len(bp) > 0 {
		quotient := bq[:0]
		remainder := uint64(0)

		for _, c := range bp {
			value := uint64(c) + uint64(remainder)*srcBase
			digit := value / dstBase
			remainder = value % dstBase

			if len(quotient) != 0 || digit != 0 {
				quotient = append(quotient, byte(digit))
			}
		}

		id[n-4] = byte(remainder >> 24)
		id[n-3] = byte(remainder >> 16)
		id[n-2] = byte(remainder >> 8)
		id[n-1] = byte(remainder)
		n -= 4
		bp = quotient
	}

	var zero [16]byte
	copy(id[:n], zero[:])
}
