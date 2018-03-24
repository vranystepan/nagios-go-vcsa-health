package main

import (
  "os"
  "fmt"
  "flag"
//  "encoding/json"
  "gopkg.in/resty.v1"
)

type Message struct {
  Value string `json:"value"`
}

var host string = ""
var hostPassword string = ""
var hostUsername string = ""

func main() {
  handleInput()

  // login to the VAPI
  client1 := resty.New()
  client1.SetBasicAuth(hostUsername, hostPassword)
  authResp, authErr := client1.R().Post("https://" + host + "")
  handleError(authErr)
  fmt.Println(authResp)
  
  // parse auth token with encoding/json

  // assemble second web client

  // do the checks

}

// custom functions

func handleError(err error) {
  // a generic function to handle response errors
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
