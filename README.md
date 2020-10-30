# uuid

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/rsms/go-uuid.svg)][godoc]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsms/go-uuid)][godoc]
[![Go Report Card](https://goreportcard.com/badge/github.com/rsms/go-uuid)](https://goreportcard.com/report/github.com/rsms/go-uuid)

[godoc]: https://pkg.go.dev/github.com/rsms/go-uuid

- Binary sortable universally unique identifier
- 16 bytes long, 6 bytes millisecond timestamp + 10 random bytes
- Efficient base-62 string encoding that is URL safe


## go doc

[View as HTML on go.dev →][godoc]

```
CONSTANTS

const StringMaxLen = 22
    StringMaxLen is the maximum length of a string representation of a UUID,
    i.e. as returned by UUID.String()


VARIABLES

var Max = UUID{
    255, 255, 255, 255, 255, 255, 255, 255,
    255, 255, 255, 255, 255, 255, 255, 255,
}
    Max is the largest possible UUID (must not be modified)


TYPES

type UUID [16]byte
    UUID is a universally unique identifier

    UUIDs are binary sortable. The first 6 bytes constitutes a
    millisecond-precision timestamp in big-endian byte order.

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

    Note that while it seems like we could use nanosecond for the timestamp to
    reduce the random data needed, environments like JavaScript doesn't
    necessarily provide high-precision clocks. Doing things this way means that
    we can generate and parse the embedded timestamp in a wide variety of
    programming languages.

var Min UUID
    Min is the zero UUID (must not be modified)

func FromBytes(verbatim []byte) UUID
    FromBytes copies verbatim bytes into an UUID and returns that UUID. verbatim
    must be at least 16 bytes long or this will panic.

func FromString(encoded string) UUID
    FromString decodes a string representation of an UUID (i.e. from String())

func Gen() (UUID, error)
    Gen generates a universally unique UUID suitable to be used for sorted
    identity. An error is returned only in the case that the host system's
    random source fails.

func MustGen() UUID
    MustGen calls Gen and panics if Gen fails

func New(sec int64, nsec int, random []byte) UUID
    New creates a new UUID with specific Unix timestamp and random bytes.

    nsec is the nanosecond part of the timestamp and should be in the range [0,
    999999999]. It's valid to pass values outside this range for nsec. Only the
    millisecond part of nsec is actually used.

    To create an UUID with a time.Time object, do this:

        New(t.Unix(), t.Nanosecond(), random)

    Up to 10 bytes is used from random. If len(random) < 10, the remaining
    "random" bytes of UUID are zero.

func (id UUID) Bytes() []byte
    Bytes returns the IDs natural 16 byte long value. The returned slice's bytes
    must not be modified.

func (id *UUID) DecodeString(src []byte)
    DecodeString sets the receiving UUID to the decoded value of src, which is
    expected to be a string previously encoded using EncodeString (base62
    0-9A-Za-z)

func (id UUID) EncodeString(dst []byte) int
    EncodeString writes the receiver to dst which must be at least StringMaxLen
    (22) bytes. Returns the start offset (this function starts writing at the
    end of dst.)

func (id UUID) String() string
    String returns a string representation of the UUID. The returned string is
    sortable with the same order as the "raw" UUID bytes and is URL safe.

func (id UUID) Time() time.Time
    Time returns the time portion of the UUID

func (id UUID) Timestamp() (sec uint32, millisec uint16)
    Timestamp returns the timestamp portion of the UUID

```
