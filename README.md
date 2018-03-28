# nagios-go-vcsa-health

Work in progress. 

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



