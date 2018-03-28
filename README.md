# nagios-go-vcsa-health

Simple nagios plugin to monitor vCenter Server Appliance (VCSA) via VAMI api. 
The motivation for writing this plugin was pretty straightforward - there's no sufficient Nagios plugin for the VCSA monitoring. 

**Work in progress.**

### Sources
- [Official VMware API documentation](https://vdc-repo.vmware.com/vmwb-repository/dcr-public/1cd28284-3b72-4885-9e31-d1c6d9e26686/71ef7304-a6c9-43b3-a3cd-868b2c236c81/doc/index.html#PKG_com.vmware.cis)

### Sample usage #1

```bash
./vcsa-health --host=vcenter.fqdn --username=user_name --password=top_secret_pass
OK: mgmt is green, database is green, load is green, storage is green, swap is green, system is green
```

### Sample usage #2

```bash
./vcsa-health --host=vcenter.fqdn --username=user_name --password=top_secret_pass --subcommand=database
OK: database is green
```

### Sample error
```bash
./vcsa-health --host=vcenter.fqdn --username=user_name --password=wrong_pass --subcommand=database
CRITICAL: {"type":"com.vmware.vapi.std.errors.unauthenticated","value":{"messages":[{"args":[],"default_message":"Authentication required.","id":"com.vmware.vapi.endpoint.method.authentication.required"}]}}
```

### Do you want to contribute?
That's exactly why I put my code here. Do you have some idea how to improve this plugin? I guess you know the drill. Fork it, create a new branch with your edits and submit a new PR (And please, do not forget to add description of your changes).

### Open topics

- [x] address authentication issue com.vmware.vapi.std.errors.unauthenticated

```json
{  
   "type":"com.vmware.vapi.std.errors.unauthenticated",
   "value":{  
      "messages":[  
         {  
            "args":[  

            ],
            "default_message":"This method requires authentication.",
            "id":"vapi.method.authentication.required"
         }
      ]
   }
```
- [x] add session deletion at the end of the program
- [ ] Implement some basic automated testing. But how??
