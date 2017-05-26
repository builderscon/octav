package assets_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/builderscon/octav/octav/assets"
	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			asset, err := assets.Asset(path)
			if !assert.NoError(t, err, "Asset(%s) should succeed", path) {
				return
			}

			expected, err := ioutil.ReadFile(path)
			if !assert.NoError(t, err, "reading from %s should succeed", path) {
				return
			}

			if !assert.Equal(t, expected, asset, "Asset() should match file content") {
				return
			}
		})
		return nil
	})
}
