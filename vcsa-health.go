package main

import (
  "os"
  "fmt"
  "flag"
  "encoding/json"
  "gopkg.in/resty.v1"
)

type VapiMessage struct {
  Value string `json:"value"`
}

var host string = ""
var hostPassword string = ""
var hostUsername string = ""

func main() {
  handleInput()

  // login to the VAPI
  c := resty.New()
  authResp, authErr := c.R().
    SetBasicAuth(hostUsername, hostPassword).
    Post("https://" + host + "/rest/com/vmware/cis/session")
  handleError(authErr)
  
  // parse auth token with encoding/json
  authData := VapiMessage{}
  authDataJsonErr := json.unmarshall(authResp.Body(), &authData)
  handleError(authDataJsonErr)
  authToken := authData.Value()

  // get health status
  healthResp, healthErr := c.R().
    SetHeader("vmware-api-session-id", authToken).
    Get("https://" + host + "/rest/appliance/health/applmgmt")
  handleError(healthErr)

  // parse health data with encoding/json
  
}

// custom functions

func handleError(err error) {
  if err != nil {
    exitUnknown(err.Error())
  }
}

func handleInput() {
  // specify commandline arguments
  hostPtr := flag.String("host", "", "IP or FQDN of VMware VCSA")
  usernamePtr := flag.String("username", "", "authorized user account name")
  passwordPtr := flag.String("password", "", "password in plain text")

  // parse command line arguments
  flag.Parse()

  // check command line arguments
  if *hostPtr == "" { exitUnknown("--host must be set") }
  if *usernamePtr == "" { exitUnknown("--username must be set") }
  if *passwordPtr == "" { exitUnknown("--password must be set") }
  
  // assign input params to variables  
  host, hostUsername, hostPassword = *hostPtr, *usernamePtr, *passwordPtr
}

func exitUnknown(msg string) {
  fmt.Printf("UNKNOWN: %s\n", msg)
  os.Exit(3)
}
