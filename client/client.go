package client

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"time"
	"log"
	"strings"
)

// Original Python API at:
// https://github.com/ttroy50/pyephember

const apiBaseUrl = "https://eu-https.topband-cloud.com/ember-back/"

type EmberClient struct {
	Credentials Credentials
	httpClient *http.Client

	Homes []Home
	Zones []Zone 
}

func (e *EmberClient) Login(username string, password string) (error) {
	loginUrl := apiBaseUrl + "appLogin/login";

	requestBody, _ := json.Marshal(map[string]string{"userName": username, "password": password})

	resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Errorf("%s", err)
		panic("Couldn't log in")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	loginResult := LoginResult{}

	json.Unmarshal(body, &loginResult)

	e.Credentials = loginResult.Credentials

	e.httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	return nil
}

type Credentials struct {
	RefreshToken string `json:"refresh_token"`
	Token string `json:"token"`
}


type LoginResult struct {
	Credentials Credentials `json:"data"`
}


type Home struct {
	// 			"deviceType":1,
	// 			"gatewayid":"1234",
	// 			"invitecode":"C70DE",
	// 			"name":"Home",
	// 			"uid":null,
	// 			"zoneCount":2

	GatewayId string `json:"gatewayid"`
	Name string `json:"name"`
	ZoneCount string `json:"zoneCount"`
}

type ListHomesResponse struct {
	Homes []Home `json:"data"`
}


func (e *EmberClient) ListHomes() ([]Home, error) {

	if len(e.Homes) > 0 {
		return e.Homes, nil
	}

	listHomesUrl := apiBaseUrl + "homes/list";

	req, err := http.NewRequest("GET", listHomesUrl, nil)

	if err != nil {
		log.Fatal("Error fetching homes. ", err)
	}

	e.setRequestHeaders(req)
 
	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
 
	// fmt.Printf("%s\n", body)

	listHomesResponse := ListHomesResponse{}

	json.Unmarshal(body, &listHomesResponse)

	e.Homes = listHomesResponse.Homes

	return e.Homes, nil
}

func (e *EmberClient) setRequestHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", e.Credentials.Token)
}

type GetZonesResponse struct {
	Zones []Zone `json:"data"`
}

type Zone struct {
	Name string `json:"name"`
	ZoneId int `json:"zoneid"`
	CurrentTemperature float64 `json:"currenttemperature"`
	TargetTemperature float64 `json:"targettemperature"`
	IsHotWater bool `json:"ishotwater"`
	IsBoostActive bool `json:"isboostactive"`
	IsAdvanceActive bool `json:"isadvanceactive"`
	Status int `json:"status"` // 1 seems to be off, 2 is on
	Prefix string `json:"prefix"`
}

// """
// Taken verbatim from // https://github.com/ttroy50/pyephember
// Check if the zone is on.
// This is a bit of a hack as the new API doesn't have a currently
// active variable
// """
func (z *Zone) IsOn() (bool){

	if z.Prefix != "" {
		if strings.Contains(z.Prefix, " off ") {
			return false
		}

		if strings.Contains(z.Prefix, "active ") {
			return true
		}

		if strings.Contains(z.Prefix, "ON mode") {
			return true
		}
	}

	return z.IsBoostActive || z.IsAdvanceActive
}


func (e *EmberClient) GetZones(gatewayId string) ([]Zone, error) {
	getHomeUrl := apiBaseUrl + "zones/polling";

	requestBody, _ := json.Marshal(map[string]string{
		"gateWayId": gatewayId,
	})

	req, err := http.NewRequest("POST", getHomeUrl, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal("Error fetching home. ", err)
	}

	e.setRequestHeaders(req)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
 
	fmt.Printf("%s\n", body)

	getZonesResponse := GetZonesResponse{}

	json.Unmarshal(body, &getZonesResponse)

	e.Zones = getZonesResponse.Zones

	return e.Zones, nil
}

func (e *EmberClient) BoostZone(zoneId int, hours int, temperature int) (error) {
	boostZoneUrl := apiBaseUrl + "zones/boost";

	requestBody, _ := json.Marshal(map[string]int{
		"zoneid": zoneId,
	    "hours": hours,
	    "temperature": temperature,
	})

	req, err := http.NewRequest("POST", boostZoneUrl, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal("Error fetching home. ", err)
	}

	e.setRequestHeaders(req)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
 
	fmt.Printf("%s\n", body)

	return nil
}

func (e *EmberClient) DeactivateBoostForZone(zoneId int) (error) {
	deactivateBoostZoneUrl := apiBaseUrl + "zones/cancelBoost";

	requestBody, _ := json.Marshal(map[string]int{
		"zoneid": zoneId,
	})

	req, err := http.NewRequest("POST", deactivateBoostZoneUrl, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal("Error fetching home. ", err)
	}

	e.setRequestHeaders(req)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
 
	fmt.Printf("%s\n", body)

	return nil
}

func (e * EmberClient) SetTargetTemperatureForZone(zoneId int, temperature int) {
	setTempUrl := apiBaseUrl + "zones/setTargetTemperature";

	requestBody, _ := json.Marshal(map[string]int{
		"zoneid": zoneId,
		"temperature": temperature,
	})

	req, err := http.NewRequest("POST", setTempUrl, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal("Error fetching home. ", err)
	}

	e.setRequestHeaders(req)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
 
	fmt.Printf("%s\n", body)
}

func (e *EmberClient) ZoneByName(name string) (Zone) {
	for _, zone := range e.Zones {
		if zone.Name == name {
			return zone
		}
	}

	return Zone{}
}

