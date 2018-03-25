package main

import (
  "os"
  "fmt"
  "flag"
  "encoding/json"
  "gopkg.in/resty.v1"
)

type vapiMessage struct {
  value string `json:"value"`
}

type vapiEndpoint struct {
  name string
  path string
}

var host string = ""
var hostPassword string = ""
var hostUsername string = ""
var subcommand string = ""

// static VAPI resource mapping
var vapiEndpointList = []vapiEndpoint{
  vapiEndpoint{
    name: "mgmt",
    path: "/rest/appliance/health/applmgmt",
  },
  vapiEndpoint{
    name: "database", 
    path: "/rest/appliance/health/database-storage",
  },
  vapiEndpoint{
    name: "load", 
    path: "/rest/appliance/health/load",
  },
  vapiEndpoint{
    name: "storage", 
    path: "/rest/appliance/health/storage",
  },
  vapiEndpoint{
    name: "swap", 
    path: "/rest/appliance/health/swap",
  },
  vapiEndpoint{
    name: "system", 
    path: "/rest/appliance/health/system",
  },
}

func main() {
  handleInput()

  // login to the VAPI
  c := resty.New()
  authResp, authErr := c.R().
    SetBasicAuth(hostUsername, hostPassword).
    Post("https://" + host + "/rest/com/vmware/cis/session")
  handleError(authErr)
  
  // parse auth token with encoding/json
  authData := vapiMessage{}
  authDataJsonErr := json.Unmarshal(authResp.Body(), &authData)
  handleError(authDataJsonErr)
  authToken := authData.value

  // set variables for the avaluation
  overallStatus := ""
  statusMessages := []string{}

  for vapiEndpointIndex, vapiEndpointObj := range vapiEndpointList {
    if subcommand != "all" {
      if vapiEndpointObj.name != subcommand { continue }
    }
    
    // get health status
    healthResp, healthErr := c.R().
      SetHeader("vmware-api-session-id", authToken).
      Get("https://" + host + vapiEndpointObj.path)
    handleError(healthErr)

    // parse health data with encoding/json
    healthData := vapiMessage{}
    healthDataJsonErr := json.Unmarshal(healthResp.Body(), &healthData)
    handleError(healthDataJsonErr)

    // append message
    statusMessages = append(statusMessages, vapiEndpointObj.name + " is " + healthData.value)
    
    

  }

  // evaluate messages    
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
  subcommandPtr := flag.String("subcommand", "all", "subcommand you want to execute <all|mgmt|database|load|storage|swap|system>")

  // parse command line arguments
  flag.Parse()

  // check command line arguments
  if *hostPtr == "" { exitUnknown("--host must be set") }
  if *usernamePtr == "" { exitUnknown("--username must be set") }
  if *passwordPtr == "" { exitUnknown("--password must be set") }
  if *subcommandPtr == "" { exitUnknown("--subcommand can't be empty")  }
  
  // assign input params to variables  
  host, hostUsername, hostPassword, subcommand = *hostPtr, *usernamePtr, *passwordPtr, *subcommandPtr
}

func exitUnknown(msg string) {
  fmt.Printf("UNKNOWN: %s\n", msg)
  os.Exit(3)
}
