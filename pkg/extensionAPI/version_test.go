package extensionAPI

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidVersion(t *testing.T) {
	a := assert.New(t)
	err := checkCompatibleVersion("id", "invalid")
	a.EqualError(err, `extension "id" uses invalid API version number "invalid"`)
}

const currentMajorVersion = "0"

/* [utest -> dsn~extension-compatibility~1] */
func TestCompatibleNewerVersion(t *testing.T) {
	a := assert.New(t)
	err := checkCompatibleVersion("id", currentMajorVersion+".99.99")
	a.NoError(err)
}

/* [utest -> dsn~extension-compatibility~1] */
func TestCompatibleOlderVersion(t *testing.T) {
	a := assert.New(t)
	err := checkCompatibleVersion("id", currentMajorVersion+".0.0")
	a.NoError(err)
}

func TestCompatibleSameVersion(t *testing.T) {
	a := assert.New(t)
	err := checkCompatibleVersion("id", supportedApiVersion)
	a.NoError(err)
}

func TestIncompatibleOlderVersion(t *testing.T) {
	a := assert.New(t)
	err := checkCompatibleVersion("id", "99.0.0")
	a.EqualError(err, `extension "id" uses incompatible API version "99.0.0". Please update the extension to use supported version "0.1.16"`)
}
