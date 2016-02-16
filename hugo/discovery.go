package hugo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HubAddress struct {
	Id                string `json:"id`
	InternalIpAddress string `json:"internalipaddress`
}

// Find all hubs on the local network.
// NOTE: We use a simple service provided by Philips:
// https://www.meethue.com/api/nupnp to discover hubs. As such, this requires
// an internet connection.
// TODO(rileya): Make this work without an internet connection.
func FindHubAddresses() ([]HubAddress, error) {
	var nupnpUrl = "https://www.meethue.com/api/nupnp"

	// Fire off an http Get() at the Hue discovery service.
	// This hands us back JSON results like this:
	// [{"id":"001788fffe0999be","internalipaddress":"192.168.86.150"}]
	response, err := http.Get(nupnpUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bodyContents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	// Dig out results from the JSON and return 'em.
	var addresses []HubAddress
	json.Unmarshal(bodyContents, &addresses)
	return addresses, nil
}
