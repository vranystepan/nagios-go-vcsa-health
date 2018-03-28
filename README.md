# nagios-go-vcsa-health

Simple nagios plugin to monitor vCenter Server Appliance (VCSA) via VAMI api. **Work in progress.**

# Sources
- [Official VMware API documentation](https://vdc-repo.vmware.com/vmwb-repository/dcr-public/1cd28284-3b72-4885-9e31-d1c6d9e26686/71ef7304-a6c9-43b3-a3cd-868b2c236c81/doc/index.html#PKG_com.vmware.cis)

### Open topics

1. address authentication issue com.vmware.vapi.std.errors.unauthenticated

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
2. add session deletion at the end of the program


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



