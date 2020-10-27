package uuid

import (
	"testing"
	"time"

	"github.com/rsms/go-testutil"
)

func TestUUID(t *testing.T) {
	assert := testutil.NewAssert(t)

	assert.Eq("Min encoding", Min.String(), "0")
	assert.Eq("Max encoding", Max.String(), "7n42DGM5Tflk9n8mt7Fhc7")

	smallId := UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	id1 := Gen()
	id2 := Gen()

	// t.Logf("id1: % x  %q", id1[:], id1)
	// t.Logf("id2: %x  %q", id2[:], id2)

	assert.Eq("Min decode(encode())", Min, FromString(Min.String()))
	assert.Eq("Max decode(encode())", Max, FromString(Max.String()))
	assert.Eq("smallId decode(encode())", smallId, FromString(smallId.String()))
	assert.Eq("id1 decode(encode())", id1, FromString(id1.String()))
	assert.Eq("id1 decode(encode())", id1, FromString(id1.String()))
	assert.Eq("id2 decode(encode())", id2, FromString(id2.String()))

	// DecodeString: Check that DecodeString rewrites all bytes, not just some
	var id3 UUID
	assert.Eq("id3", Min, id3)
	id3.DecodeString([]byte("7n42DGM5Tflk9n8mt7Fhc7"))
	assert.Eq("id3 DecodeString should work", id3, Max)
	id3.DecodeString([]byte("A"))
	assert.Eq("id3 should be zeroed", id3, UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xA})

	// Check that Time() and Timestamp() returns the correct time as encoded
	tm := time.Unix(1600000000+3212345, 123*1000000) // 2020-10-20 16:45:45.123 UTC
	idt := New(tm.Unix(), tm.Nanosecond(), []byte{})
	// t.Logf("tm  %s", tm.UTC())
	// t.Logf("idt %x", idt[:])
	tssec, tsms := idt.Timestamp()
	assert.Eq("Timestamp() tssec", tssec, uint32(3212345))
	assert.Eq("Timestamp() tsms", tsms, uint16(123))
	// must check with rounded time since UUID timestamp has only millisecond precision
	assert.Eq("Time()", idt.Time().Unix(), tm.Unix())
	assert.Eq("Time()",
		idt.Time().Nanosecond()/int(time.Millisecond),
		tm.Nanosecond()/int(time.Millisecond))

	// // Generate UUIDs for documentation or demo
	// for i := 0; i < 5; i++ {
	// 	id := Gen()
	// 	// 2020-10-20 16:45:45.xxx UTC
	// 	id = New(1603212345+int64(i), time.Now().Nanosecond()+i*1000000, id[6:])
	// 	t.Logf("% x  %q  %s", id[:], id, id.Time().UTC())
	// }
}
