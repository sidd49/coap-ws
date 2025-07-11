package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Thing struct {
	DeviceId    string
	SessionID   string
	Description string
}

type RequestBody struct {
	DeviceID string `json:"deviceid"`
}

func ExtractDeviceID(reader io.ReadSeeker) (string, error) {
	// Reset the reader to the start before reading
	if reader == nil {
		return "", fmt.Errorf("there is no body for the req")
	}
	_, err := reader.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("error resetting reader: %v", err)
	}

	// Read the data from the io.ReadSeeker into a byte slice
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading data: %v", err)
	}

	// Print the raw JSON data for debugging
	fmt.Println("Raw JSON data:", buf.String())

	// Parse the JSON into the RequestBody struct
	var body RequestBody
	err = json.Unmarshal(buf.Bytes(), &body)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Return the device ID
	return body.DeviceID, nil
}
