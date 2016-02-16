package main

import (
	"fmt"
	"github.com/rileya/hugo"
	"io/ioutil"
	"os"
	"strconv"
)

// Super-simple CLI for controlling lights.
func main() {
	// Try to find hubs on the local network.
	hubAddresses, err := hugo.FindHubAddresses()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(hubAddresses) == 0 {
		fmt.Println("No Hue hubs found.")
		return
	}

	var selectedAddress hugo.HubAddress

	// Prompt the user to pick onen hub, if multiple are present.
	if len(hubAddresses) > 1 {
		for {
			fmt.Println("\nFound multiple hubs, pick one:")
			for i, hub := range hubAddresses {
				fmt.Printf("\t(%d) %s - %s\n", i, hub.InternalIpAddress, hub.Id)
			}

			var selection string
			fmt.Scanln(&selection)
			i, err := strconv.Atoi(selection)

			if err != nil {
				fmt.Println(err)
				continue
			} else if i < 0 || i >= len(hubAddresses) {
				fmt.Printf("%d is out of bounds, expected range: [0,%d].\n", i, len(hubAddresses)-1)
				continue
			}
			selectedAddress = hubAddresses[i]
			break
		}
	} else {
		selectedAddress = hubAddresses[0]
	}

	// Construct a Hub struct.
	fmt.Println("Using hub with id:  " + selectedAddress.Id + " ip: " + selectedAddress.InternalIpAddress)
	hub := hugo.CreateHubWithAddress(selectedAddress)

	// The Hue API requires a username to be generated, with a press of
	// a physical button on the Hue hub. After one is generated we'll store
	// it at ~/.hugo_username
	configPath := os.Getenv("HOME") + "/.hugo_username"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println(err)
			return
		}
		deviceType := "hugo_cli#" + hostname
		fmt.Println("Authenticating new user...")
		username, err := hub.AuthenticateNewUser(deviceType, func() {
			fmt.Println("Press the hub button within 30 seconds, then press ENTER")
			fmt.Scanln()
		})

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("New user authenticated: " + username)
		fmt.Println("Saving username to: " + configPath)
		err = ioutil.WriteFile(configPath, []byte(username), 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Reading username from: " + configPath)
		configContents, err := ioutil.ReadFile(configPath)
		fmt.Println("Username: " + string(configContents))
		if err != nil {
			fmt.Println(err)
			return
		}
		err = hub.AuthenticateExistingUser(string(configContents))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// If we got this far without returning, then we've got a valid hub
	// address + username.
	fmt.Println("Authentication successful.")

	// Toggle all reachable lights.
	// TODO(rileya): Actually make this an interactive interface-y thing.
	for name, light := range hub.Lights {
		newState := light.State
		if light.State.Reachable {
			newState.On = !newState.On
			hub.SetLightState(name, newState)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
