package main

import (
  "os"
  "fmt"
  "flag"
  "encoding/json"
  "gopkg.in/resty.v1"
  "crypto/tls"
)

type VapiMessage struct {
  Value string `json:"value"`
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

  // create and configure REST client
  c := resty.New()
  c.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })

  // login to the VAPI
  authResp, authErr := c.R().
    SetBasicAuth(hostUsername, hostPassword).
    Post("https://" + host + "/rest/com/vmware/cis/session")
  handleError("authentiation request", authErr)
  
  // parse auth token with encoding/json
  authData := VapiMessage{}
  authDataJsonErr := json.Unmarshal(authResp.Body(), &authData)
  handleError("authentication json parsing", authDataJsonErr)
  authToken := authData.Value

  // set variables for the avaluation
  overallStatus := "green"
  statusMessages := []string{}

  for _, vapiEndpointObj := range vapiEndpointList {
    // execute only one subcommand if specified
    if subcommand != "all" {
      if vapiEndpointObj.name != subcommand { continue }
    }
    
    // get health status
    healthResp, healthErr := c.R().
      SetHeader("vmware-api-session-id", authToken).
      Get("https://" + host + vapiEndpointObj.path)
    handleError("health request", healthErr)

    // parse health data with encoding/json
    healthData := VapiMessage{}
    healthDataJsonErr := json.Unmarshal(healthResp.Body(), &healthData)
    handleError("health json parsing", healthDataJsonErr)

    // append message
    statusMessages = append(statusMessages, vapiEndpointObj.name + " is " + healthData.Value)
    
    // green can be changed to any status
    if overallStatus == "green" { 
      if healthData.Value != "green" {
        overallStatus = healthData.Value
      }
    }

    // orange can be changed only to red
    if overallStatus == "orange" {
      if healthData.Value == "red" {
        overallStatus = healthData.Value
      }
    }

    // red can't be changed to statuses with lesser severity

  }
  
  // logout from the appliance
  _, deleteErr := c.R().
    SetHeader("vmware-api-session-id", authToken).
    Delete("https://" + host + "/rest/com/vmware/cis/session")
  handleError("logout", deleteErr)

  //evaluate overall health status
  switch overallStatus {
    case "green": exitFinal(statusMessages, "OK", 0)
    case "orange": exitFinal(statusMessages, "WARNING", 1)
    case "red": exitFinal(statusMessages, "CRITICAL", 2)
    default: exitUnknown("overall status is missing!")
  }  
}

// custom functions

func handleError(step string, err error) {
  if err != nil {
    exitUnknown(step + "; " + err.Error())
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
  if validateSubcommand(*subcommandPtr) == false { exitUnknown("incorrect subcommand name") }
  
  // assign input params to variables  
  host, hostUsername, hostPassword, subcommand = *hostPtr, *usernamePtr, *passwordPtr, *subcommandPtr
}

func validateSubcommand(s string) bool {
  status := false
  for _, vapiEndpointObj := range vapiEndpointList {
    if s == vapiEndpointObj.name { status = true }
  }
  if s == "all" { status = true }
  return status
}

func exitUnknown(msg string) {
  fmt.Printf("UNKNOWN: %s\n", msg)
  os.Exit(3)
}

func exitFinal(messages []string, status string, exitCode int) {

  // print nagios status
  fmt.Printf("%s: ", status)

  // go through messages
  for messageIndex, message := range messages {
    fmt.Printf("%s", message)
    if messageIndex < len(messages)-1 {
      fmt.Printf(", ")
    } else {
      fmt.Printf("\n")
    }
  }

  //exit the program
  os.Exit(exitCode)
}
