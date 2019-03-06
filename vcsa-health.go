package main

import (
  "os"
  "fmt"
  "encoding/json"
  "gopkg.in/resty.v1"
  "crypto/tls"
  "github.com/jessevdk/go-flags"
)

type VapiMessage struct {
  Value string `json:"value"`
}

type vapiEndpoint struct {
  name string
  path string
}

// opts specification
type ProgamOptions struct {
  Host string `short:"H" long:"host" description:"IP address or FQDN of vCSA you want to connect to" required:"yes"`
  Username string `short:"u" long:"username" description:"username of auhorized user account" required:"yes"`
  Password string `short:"p" long:"password" description:"password of authorized user account" required:"yes"`
  Subcommand string `short:"s" long:"subcommand" description:"name of check you want to execute against vCSA" required:"no" default:"all" choice:"all" choice:"mgmt" choice:"database" choice:"load" choice:"storage" choice:"swap" choice:"system"`
  Verbose bool `short:"v" long:"verbose" description:"verbose output for debug" required:"no"`
}

// static VAPI resource mapping
var vapiEndpointList = []vapiEndpoint{
  vapiEndpoint{ name: "mgmt", path: "/rest/appliance/health/applmgmt", },
  vapiEndpoint{ name: "database", path: "/rest/appliance/health/database-storage", },
  vapiEndpoint{ name: "load", path: "/rest/appliance/health/load", },
  vapiEndpoint{ name: "storage", path: "/rest/appliance/health/storage", },
  vapiEndpoint{ name: "swap", path: "/rest/appliance/health/swap", },
  vapiEndpoint{ name: "system", path: "/rest/appliance/health/system", },
}

func main() {
  //handle commandline params  
  host, hostUsername, hostPassword, subcommand := handleInput(os.Args[1:])

  // create and configure REST client
  c := resty.New()
  c.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })

  // login to the VAPI
  authResp, authErr := c.R().
    SetBasicAuth(hostUsername, hostPassword).
    Post("https://" + host + "/rest/com/vmware/cis/session")
  handleError("authentiation request", authErr)
  handleHttpStatus(authResp.StatusCode(), authResp.Body())
  
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
    handleHttpStatus(healthResp.StatusCode(), healthResp.Body())

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

    // orange, yellow and grey can be changed only to red
    if overallStatus == "orange" || overallStatus == "yellow" || overallStatus == "grey" {
      if healthData.Value == "red" {
        overallStatus = healthData.Value
      }
    }

  }
  
  // logout from the appliance
  deleteResp, deleteErr := c.R().
    SetHeader("vmware-api-session-id", authToken).
    Delete("https://" + host + "/rest/com/vmware/cis/session")
  handleError("logout", deleteErr)
  handleHttpStatus(deleteResp.StatusCode(), deleteResp.Body())

  //evaluate overall health status
  switch overallStatus {
    case "green": exitFinal(statusMessages, "OK", 1)
    case "yellow": exitFinal(statusMessages, "WARNING", 1)
    case "orange": exitFinal(statusMessages, "WARNING", 1)
    case "grey": exitFinal(statusMessages, "WARNING", 1)
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

func handleHttpStatus(statusCode int, statusBody []byte) {
  if statusCode != 200 {
    statusBodyString := string(statusBody[:])
    exitCritical(statusBodyString)
  }
}

func handleInput(args []string) (string, string, string, string) {
  opts := ProgamOptions{}

  // we don't want to print error message automatically
  parser := flags.NewParser(&opts, (flags.HelpFlag | flags.PassDoubleDash))
  _, err := parser.ParseArgs(args);
  handleError("input params", err)

  return opts.Host, opts.Username, opts.Password, opts.Subcommand
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

func exitCritical(msg string) {
  fmt.Printf("CRITICAL: %s\n", msg)
  os.Exit(2)
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
