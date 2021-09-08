# Deploying Santa Agents
For reference, it's recommended to read over how to [configure Santa](https://santa.dev/deployment/configuration.html).

Provided here are some accelerated instructions.

## Full Disk Access and Sysext Approval
For unattended installs of Santa you'll want to approve Santa's system extension and full disk access. For convenience
we've provided some sample configuration profiles here to get you started:

* [configs/santa-sysext.mobileconfig](configs/santa-sysext.mobileconfig)
* [configs/santa-tcc.mobileconfig](configs/santa-tcc.mobileconfig)

Make sure to doctor them up a bit prior to a real production deploy.


## Configuration Profile
The Santa agent is configured via MacOS Profiles. Take a look at the example `.mobileconfig` file we've provided
[configs/santa-configuration.mobileconfig](configs/santa-configuration.mobileconfig). Below we'll go over some of
the **_more important_** values to pay attention to.

### `SyncBaseURL`
The most important parameter to consider is the `SyncBaseURL`. The value here must match the `sync_base_url`
that is output by your `make deploy` command. This will instruct the Santa sensor to treat your Rudolph sync server
as its authority.

### `MachineIDPlist` and `MachineIDKey`
Specifies the `.plist` file and relevant XML key located on a MacOS machine's filesystem that specifies the `MachineID`
of the Santa agent. This `MachineID` is used to uniquely identify the MacOS machine and all rules, configurations, and
logs are tied to this `MachineID`.

### `ClientMode`
We **_highly recommend_** the integer value `1`. This will initialize all Santa sensors in `MONITOR` mode. The mode can be
later remotely changed by Rudolph once everything is set up, but an initial default setting of `MONITOR` will reduce the chances for problems.


## Plist File
Deploy a `.plist` file to the `MachineIDPlist` location, using your MDM or otherwise. We've included an example file,
[configs/com.google.santa.machine-mapping.plist](configs/com.google.santa.machine-mapping.plist).

We **_highly recommend_** using **UPPER-CASE-HEXADECIMAL-UUID-WITH-SPACES** (e.g. `AAAABBBB-CCCC-DDDD-EEEE-123456780000`)
for the `MachineID`.


## Installation
Grab the installation `.dmg` over from https://github.com/google/santa/releases and install it.


## Verify it Works
After Santa is installed, you can perform your first sync via:

```
sudo santactl sync --debug
```

If you receive some messages like:

> Received 0 rules
> Received 0 rules
> Sync completed successfully

Then you should be good to go.


## Other Miscellaneous
Santa dumps tons of logs into `/var/db/santa.log`. These logs are extremely useful but this file can fill up very quickly.

We recommend a newsyslog configuration: `/private/etc/newsyslog.d/com.google.santa.newsyslog.conf`
```
# logfilename             [owner:group]       mode   count size(KiB) when   flags [/pid_file] # [sig_num]
/var/db/santa/santa.log   root:wheel          644    10    5000     *      NZ
```
