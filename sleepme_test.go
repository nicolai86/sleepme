package sleepme

import (
	"context"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

var (
	sleepMeAccessToken = os.Getenv("SLEEP_ME_ACCESS_TOKEN")
)

func skipIfAccessTokenMissing(t *testing.T) {
	if sleepMeAccessToken == "" {
		t.Skip("SLEEP_ME_ACCESS_TOKEN is not configured")
	}
}

func TestClient(t *testing.T) {
	skipIfAccessTokenMissing(t)

	c, err := New(sleepMeAccessToken)
	assert.NilError(t, err, "failed to create client")

	devices, err := c.ListDevices(context.Background())
	assert.NilError(t, err, "faild to list devices")

	assert.Assert(t, len(devices) > 0, "failed to return any devices")

	details, err := c.Get(context.Background(), devices[0].ID)
	assert.NilError(t, err)

	desiredTemperatureF := float64(details.Control.SetTemperatureF)
	err = c.Update(context.Background(), devices[0].ID, UpdateRequest{
		SetTemperatureF: &desiredTemperatureF,
	})
	assert.NilError(t, err)
}
