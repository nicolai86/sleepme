// Package sleepme implements the sleep.me API, as described on their developer portal
// https://sleep.me/account/developers
package sleepme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ProductionAPIEndpoint = "https://api.developer.sleep.me/v1"
)

type Client struct {
	APIEndpoint string
	token       string
	*http.Client
}

// New creates a new client and validates the provided token
func New(token string, opts ...func(*Client) error) (*Client, error) {
	c := &Client{
		token:       token,
		APIEndpoint: ProductionAPIEndpoint,
		Client:      &http.Client{},
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Device represents a Dock Pro unit
type Device struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Attachments []string `json:"attachments"`
}

// ListDevices lists all Dock Pro units available with the active user
func (c *Client) ListDevices(ctx context.Context) ([]Device, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices", c.APIEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req = req.WithContext(ctx)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var res []Device
	return res, json.NewDecoder(resp.Body).Decode(&res)
}

// DeviceDetails contains all the details available via the API
type DeviceDetails struct {
	About struct {
		FirmwareVersion string `json:"firmware_version"`
		IpAddress       string `json:"ip_address"`
		LanAddress      string `json:"lan_address"`
		MacAddress      string `json:"mac_address"`
		Model           string `json:"model"`
		SerialNumber    string `json:"serial_number"`
	} `json:"about"`

	Control struct {
		BrightnessLevel        int    `json:"brightness_level"`
		DisplayTemperatureUnit string `json:"display_temperature_unit"`
		SetTemperatureC        int    `json:"set_temperature_c"`
		SetTemperatureF        int    `json:"set_temperature_f"`
		ThermalControlStatus   string `json:"thermal_control_status"`
		TimeZone               string `json:"time_zone"`
	} `json:"control"`

	Status struct {
		IsConnected       bool    `json:"is_connected"`
		IsWaterLow        bool    `json:"is_water_low"`
		WaterLevel        int     `json:"water_level"`
		WaterTemperatureF int     `json:"water_temperature_f"`
		WaterTemperatureC float64 `json:"water_temperature_c"`
	} `json:"status"`
}

// Get fetches details for a specific Dock Pro unit
func (c *Client) Get(ctx context.Context, deviceID string) (*DeviceDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices/%s", c.APIEndpoint, deviceID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req = req.WithContext(ctx)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var res DeviceDetails
	return &res, json.NewDecoder(resp.Body).Decode(&res)
}

// ThermalControlStatus configures if the unit is active or not
type ThermalControlStatus string

var (
	// ThermalControlStatusActive sets the unit to active
	ThermalControlStatusActive ThermalControlStatus = "active"
	// ThermalControlStatusStandby sets the unit to standby
	ThermalControlStatusStandby ThermalControlStatus = "standby"
)

// DisplayTemperatureUnit configures what unit is used in the Dock Pro display
type DisplayTemperatureUnit string

var (
	//DisplayTemperatureUnitC displays temperature in Celsius
	DisplayTemperatureUnitC DisplayTemperatureUnit = "c"
	//DisplayTemperatureUnitF displays temperature in Fahrenheit
	DisplayTemperatureUnitF DisplayTemperatureUnit = "f"
)

// UpdateRequest contains all the fields that can be changed via the API on a Dock Pro
type UpdateRequest struct {
	ThermalControlStatus   *ThermalControlStatus   `json:"thermal_control_status,omitempty"`
	SetTemperatureF        *float64                `json:"set_temperature_f,omitempty"`
	SetTemperatureC        *float64                `json:"set_temperature_c,omitempty"`
	DisplayTemperatureUnit *DisplayTemperatureUnit `json:"display_temperature_unit,omitempty"`
	TimeZone               *string                 `json:"time_zone,omitempty"`
}

// Update reconfigures a Dock Pro
func (c *Client) Update(ctx context.Context, deviceID string, r UpdateRequest) error {
	bs := bytes.Buffer{}
	if err := json.NewEncoder(&bs).Encode(r); err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/devices/%s", c.APIEndpoint, deviceID), &bs)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req = req.WithContext(ctx)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	return nil
}
