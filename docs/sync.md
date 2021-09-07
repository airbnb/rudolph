# The Sync Process
Santa synching is documented here: https://santa.readthedocs.io/en/latest/introduction/syncing-overview/ 

The sync process requires a SyncBaseURL to be configured in the Configuration Profile. This sync server must be an API server that responds to the following 5 endpoints:

* POST /preflight/{MACHINE_ID}
* POST /eventupload/{MACHINE_ID}
* POST /logupload/{MACHINE_ID} - if enabled, as controlled via MDM profiles
* POST /ruledownload/{MACHINE_ID}
* POST /postflight/{MACHINE_ID}

Note: “MACHINE_ID” is configured in the Configuration Profile.

This flow of API calls allows the Santa sensor and the sync server to synchronize their rules and desired configurations. This process is complex, and is expanded upon below. 

## Synchronization Process
One key point about syncing: All API endpoints are HTTP POST by design: they are not intended to be idempotent. The Rudolph server maintains most of the state in each sync transaction and holds a significant amount of control over what happens during a sync. 

Additionally, API endpoints only write synchronization state back to the DynamoDB and API endpoints have no functionality to write/inject rules back into the DynamoDB table. 

### XSRF - CSRF
** This endpoint is not enabled in the deployed version of Rudolph, as it is not clear how this feature improves security **

If implemented, this API endpoint generates a token which is linked to the client sync session. This token will be sent in the header of each POST request to the sync server and will validate, on the server end, that this matches the current session XSRF token.

A simple API endpoint has been created that will always return status 200. 

Rudolph is a RestAPI for the Santa agent to synchronize rulesets from Rudolph. In this case, clients do not need to authenticate nor are cookies present during this synchronization session and the Santa agent manages states during the sync. XSRF/CSRF is not a vulnerable vector due to the way Santa agents communicate with Rudolph.

### Preflight
#### URL - HTTP POST /preflight/{machine_uuid}
This endpoint allows the Santa Sensor to upload its current configuration and sensor data, and allows the Rudolph server to send down desired configuration as well as instructions on whether to clean sync.
#### Reconfiguring the Sensor
    - The Rudolph server’s preflight response overrides many configuration parameters of the Santa sensor when it is received. It can change:
    - The sensor’s MODE (Monitor/lockdown)
    - The number of logs/events uploaded in each POST request
    - Regex-based DENY and ALLOW rules

More detailed documentation: https://santa.readthedocs.io/en/latest/deployment/configuration/#sync-server-provided-configuration

Client Requesting Clean Sync
The Santa sensor can “request” for a sync process to be “clean”. This shows up as a POST BODY parameter of request_clean_sync: true but has no further meaning unless the server respects this (see below).

NOTE: You can force this using: santactl sync --clean

Server Requesting Clean Sync
The Rudolph server can command the Santa sensor to destroy all of its local rules and re-download all rules from scratch. This is known as a “clean sync” and the server can mandate this by sending back clean_sync: true in the preflight response. Doing so will cause the client to clear its rules database immediately, prior to doing ruledownload.

NOTE: Even if the Client requests a clean sync, the server is not obligated to respect it.

#### Request - JSON

```json
{
    "os_build":"20D5029f",
    "santa_version":"2021.1",
    "hostname":"my-awesome-macbook-pro.attlocal.net",
    "transitive_rule_count":0,
    "os_version":"11.2",
    "certificate_rule_count":0,
    "client_mode":"MONITOR",
    "serial_num":"C02123456789",
    "binary_rule_count":0,
    "primary_user":"",
    "compiler_rule_count":0
}
```

### Response - JSON

More info found here: https://santa.readthedocs.io/en/latest/deployment/configuration/#sync-server-provided-configuration

```json
{
	"client_mode": "LOCKDOWN",
	"blocked_path_regex": "(\/tmp)|(\/Users\/.+\/trash)",
	"allowed_path_regex": "\/usr\/local\/bin",
	"batch_size": 37,
	"enable_bundles": true,
	"enable_transitive_rules": false,
	"upload_logs_url": "/aaa"
}
```


### Logupload - Not implemented

#### URL - HTTP POST /logupload/{machine_uuid}
** This endpoint is not enabled in the deployed version of Rudolph and must be enabled via configuration MDM profiles **

logupload - endpoint will upload all new changes to the Santa.log file and uploads the contents to the sync server for ingestion
Logs stored in /var/db/santa/santa.log
Work in progress to move this to JSON or protocol buffer
If KEXT is in use → /usr/bin/log show --info --debug --predicate 'senderImagePath == "/Library/Extensions/santa-driver.kext/Contents/MacOS/santa-driver"'

### Eventupload

#### URL - HTTP POST /eventupload/{machine_uuid}
This endpoint allows the Santa sensor to upload any DENY events to the Rudolph server. This endpoint is useful for gathering logs but is unimportant to the ruledownload part of sync.

This endpoints accepts events logs from the Santa agent which records all binary/application executions. Once the data is uploaded to Rudolph, it is sent via Kinesis Stream.

#### Request - JSON

```json
{
  "events": [
    {
      "parent_name":"","file_path":"\/usr\/local\/Cellar\/git\/2.30.0\/bin",
      "quarantine_timestamp":0,
      "logged_in_users":["derek_wang"],
      "signing_chain":[
        {
          "cn":"Developer ID Application: Hashicorp, Inc. (D38WU7D763)",
          "valid_until":1730930319,
          "org":"Hashicorp, Inc.",
          "valid_from":1573077519,
          "ou":"D38WU7D763",
          "sha256":"7576abcd4a9aa8a96f2860e6a0e38dabfd87491934201bf884835a5d911da500"
        },
        {
          "cn":"Developer ID Certification Authority",
          "valid_until":1801519935,
          "org":"Apple Inc.",
          "valid_from":1328134335,
          "ou":"Apple Certification Authority",
          "sha256":"7afc9d01a62f03a2de9637936d4afe68090d2de18d03f29c88cfb0b1ba63587f"
        },
        {
          "cn":"Apple Root CA",
          "valid_until":2054670036,
          "org":"Apple Inc.",
          "valid_from":1146001236,
          "ou":"Apple Certification Authority",
          "sha256":"b0b1730ecbc7ff4505142c49f1295e6eda6bcaed7e2c68c5be91b5a11001f024"
        }
      ],
      "ppid":69149,
      "executing_user":"derek_wang",
      "file_name":"git",
      "execution_time":1611254994.299253,
      "file_sha256":"d5af4e67bedacb605a415c992d31812a0ca5cfbacf306548d332531d27cd956f",
      "decision":"BLOCK_BINARY",
      "pid":70053,
      "current_sessions":["derek_wang@console","derek_wang@ttys000","derek_wang@ttys001"]
    }
  ]
}
```

#### RESPONSE - None

### Ruledownload

#### URL - HTTP POST /ruledownload/{machine_uuid}
This endpoint allows the Rudolph server to push down new rules to the Santa sensor. This process is complex and is explained within the [rule-pagination page](./rule-pagination.md)

#### Strategies Definitions:
- 1 = Downloads all GlobalRules
- 2 = Downloads all FeedRules which is determined via the last sync state time
- 3 = Downloads all MachineRules

#### Request - JSON
```json
{
	"cursor": {
		"strategy": 2,
		"batch_size": 7,
		"page": 2,
		"pk": "AAAA",
		"sk": "eeeeee"
	}
}
```

### Postflight

#### URL - HTTP POST /postflight/{machine_uuid}
This endpoint is called at the end of a successful sync process. It is the way for the Santa sensor to inform the Rudolph server that the sync process was successful and all rules were correctly committed.

Rudolph keeps state records of every synchronization event that occurs including information on time to process the entire synchronization process. 

#### Request - Blank

#### Response - Blank - Sends HTTP Status 200

# Configuration
Unlike many other sensors, Santa is not configured via a .conf or .yaml file on the disk; it is configured by a [MacOS Configuration Profile](https://developer.apple.com/business/documentation/Configuration-Profile-Reference.pdf).

An example configuration is provided in the docs.

The format of this profile is documented here: https://santa.readthedocs.io/en/latest/deployment/configuration/ 