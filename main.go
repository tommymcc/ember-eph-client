package main

import (
	"fmt"
	"github.com/tommymcc/ember-eph-client/client"
	"os"
)

func main() {
	var userName string = os.Getenv("EMBER_USERNAME")
	var password string = os.Getenv("EMBER_PASSWORD")

	if userName == "" || password == "" {
		panic("EMBER_USERNAME and EMBER_PASSWORD environment variables must be set")
	}

	client := client.EmberClient{}
	client.Login(userName, password)

	homes, _ := client.ListHomes()

	fmt.Println("Found the following homes:")

	for _, home := range homes {
		fmt.Printf("Home ID: %s\n", home.GatewayId)

		zones, _ := client.GetZones(home.GatewayId)

		for _, zone := range zones {
			fmt.Printf(" -- Zone: '%s' - %f*C  - On: %v\n", zone.Name, zone.CurrentTemperature, zone.IsOn())
		}
	}

	// heating := client.ZoneByName("Heating")
	// hotWater := client.ZoneByName("Hot Water")

	// client.BoostZone(heating.ZoneId, 1, 27)
	// client.BoostZone(hotWater.ZoneId, 1, 60)

	// client.DeactivateBoostForZone(heating.ZoneId)
}
