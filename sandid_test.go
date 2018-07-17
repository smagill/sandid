package sandid

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	sIDs := make([]SandID, 0, 1000000)
	for i := 0; i < cap(sIDs); i++ {
		sID := New()
		assert.NotZero(t, sID)
		sIDs = append(sIDs, sID)
	}
	assert.Len(t, sIDs, cap(sIDs))
	sort.Slice(sIDs, func(i, j int) bool {
		return Compare(sIDs[i], sIDs[j]) < 0
	})
	for i := 1; i < len(sIDs); i++ {
		assert.NotEqual(t, sIDs[i-1], sIDs[i])
	}
}

func TestParse(t *testing.T) {
	sID, err := Parse("00000000000000000000000000000000")
	assert.Zero(t, sID)
	assert.NoError(t, err)
	sID, err = Parse("ffffffffffffffffffffffffffffffff")
	assert.NotZero(t, sID)
	assert.NoError(t, err)
	sID, err = Parse("000000000000000g0000000000000000")
	assert.Zero(t, sID)
	assert.Error(t, err)
}

func TestMustParse(t *testing.T) {
	assert.NotPanics(t, func() {
		MustParse("00000000000000000000000000000000")
	})
	assert.Panics(t, func() {
		MustParse("000000000000000g0000000000000000")
	})
}

func TestSandIDIsZero(t *testing.T) {
	assert.True(t, SandID{}.IsZero())
	assert.False(t, New().IsZero())
}

func TestSandIDString(t *testing.T) {
	assert.Equal(t, "00000000000000000000000000000000", SandID{}.String())
}

func TestSandIDScan(t *testing.T) {
	sID := SandID{}
	assert.NoError(t, sID.Scan("ffffffffffffffffffffffffffffffff"))
	assert.Equal(t, "ffffffffffffffffffffffffffffffff", sID.String())
	assert.Error(t, sID.Scan([]byte{
		255, 255, 255, 255,
		255, 255, 255, 255,
		255, 255, 255, 255,
		255, 255, 255, 255,
		255,
	}))
	assert.Error(t, sID.Scan(0))
}

func TestSandIDValue(t *testing.T) {
	sID := New()
	v, err := sID.Value()
	assert.NoError(t, err)
	b, ok := v.([]byte)
	assert.True(t, ok)
	assert.Equal(t, sID[:], b)
}

func TestSandIDMarshalText(t *testing.T) {
	sID := New()
	b, err := sID.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, sID.String(), string(b))
}

func TestSandIDUnmarshalText(t *testing.T) {
	sID := SandID{}
	assert.NoError(
		t,
		sID.UnmarshalText([]byte("ffffffffffffffffffffffffffffffff")),
	)
	for _, b := range sID {
		assert.Equal(t, byte(255), b)
	}
}

func TestSandIDMarshalBinary(t *testing.T) {
	sID := New()
	b, err := sID.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, sID[:], b)
}

func TestSandIDUnmarshalBinary(t *testing.T) {
	sID := SandID{}
	assert.NoError(t, sID.UnmarshalBinary([]byte{
		255, 255, 255, 255,
		255, 255, 255, 255,
		255, 255, 255, 255,
		255, 255, 255, 255,
	}))
	assert.Equal(t, "ffffffffffffffffffffffffffffffff", sID.String())
}

func TestSandIDMarshalJSON(t *testing.T) {
	sID := New()
	b, err := sID.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\""+sID.String()+"\"", string(b))
}

func TestSandIDUnmarshalJSON(t *testing.T) {
	sID := SandID{}
	assert.Error(t, sID.UnmarshalJSON(nil))
	assert.NoError(t, sID.UnmarshalJSON([]byte(
		"\"ffffffffffffffffffffffffffffffff\"",
	)))
	assert.Equal(t, "ffffffffffffffffffffffffffffffff", sID.String())
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(SandID{}, SandID{}))
	assert.False(t, Equal(New(), New()))
}

func TestCompare(t *testing.T) {
	assert.Equal(t, -1, Compare(
		MustParse("00000000000000000000000000000001"),
		MustParse("00000000000000000000000000000002"),
	))
	assert.Equal(t, 0, Compare(SandID{}, SandID{}))
	assert.Equal(t, 1, Compare(
		MustParse("00000000000000000000000000000002"),
		MustParse("00000000000000000000000000000001"),
	))
}

func TestNullSandIDScan(t *testing.T) {
	nsID := NullSandID{}
	assert.NoError(t, nsID.Scan(nil))
	assert.True(t, nsID.SandID.IsZero())
	assert.False(t, nsID.Valid)
	assert.NoError(t, nsID.Scan([]byte{
		255, 255, 255, 255,
		255, 255, 255, 255,
		255, 255, 255, 255,
		255, 255, 255, 255,
	}))
	assert.False(t, nsID.SandID.IsZero())
	assert.True(t, nsID.Valid)
	assert.Equal(
		t,
		"ffffffffffffffffffffffffffffffff",
		nsID.SandID.String(),
	)
}

func TestNullSandIDValue(t *testing.T) {
	nsID := NullSandID{}
	v, err := nsID.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
	nsID.SandID = MustParse("ffffffffffffffffffffffffffffffff")
	nsID.Valid = true
	v, err = nsID.Value()
	assert.NoError(t, err)
	assert.NotNil(t, v)
	b, ok := v.([]byte)
	assert.True(t, ok)
	assert.Equal(t, nsID.SandID[:], b)
}
