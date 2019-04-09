package ebpf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinuxKernelVersionCode(t *testing.T) {
	// Some sanity checks
	assert.Equal(t, linuxKernelVersionCode(2, 6, 9), uint32(132617))
	assert.Equal(t, linuxKernelVersionCode(3, 2, 12), uint32(197132))
	assert.Equal(t, linuxKernelVersionCode(4, 4, 0), uint32(263168))

	assert.Equal(t, stringToKernelCode("2.6.9"), uint32(132617))
	assert.Equal(t, stringToKernelCode("3.2.12"), uint32(197132))
	assert.Equal(t, stringToKernelCode("4.4.0"), uint32(263168))
}

func TestUbuntu44119NotSupported(t *testing.T) {
	for i := uint32(119); i < 127; i++ {
		ok, err := verifyOSVersion(linuxKernelVersionCode(4, 4, i), nil)
		assert.False(t, ok)
		assert.Error(t, err)
	}
}

func TestExcludedKernelVersion(t *testing.T) {
	exclusionList := []string{"5.5.1", "6.3.2"}
	ok, err := verifyOSVersion(linuxKernelVersionCode(4, 4, 121), exclusionList)
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = verifyOSVersion(linuxKernelVersionCode(5, 5, 1), exclusionList)
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = verifyOSVersion(linuxKernelVersionCode(6, 3, 2), exclusionList)
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = verifyOSVersion(linuxKernelVersionCode(6, 3, 1), exclusionList)
	assert.True(t, ok)
	assert.Nil(t, err)

	ok, err = verifyOSVersion(linuxKernelVersionCode(5, 5, 2), exclusionList)
	assert.True(t, ok)
	assert.Nil(t, err)
}
