# System Risks and Security

## Public Rudolph instance
Until a more robust authorization plane can be implemented externally to Santa like enforcement of client-certificates, Rudolph API server is globally available to any resource on the internet.

Downloading of rules is not considered sensitive data neither is knowing what is allowed/denied from a binary/application perspective. One risk of having a public API endpoint is the availability being compromised. 
Trusting the Wrong Server
How do we prevent DNS poisoning or MITM attacks between the Santa sensor and FIRST registration/authentication with the Rudolph server?

Currently, Santa supports server trusts by specifying a ServerAuthRootsFile in the macOS configuration profile. This file contains a SSL Certificate chain which is used to verify the identity of the server. 


However, there is a bug in Santa: When this file is malformed or if the certificate does not match the SyncBaseURL, Santa defaults back to the Apple MacOS keychain, which will blanket trust anything signed by Amazon Web Services. Thus it is possible for an attacker to trick a Santa sensor into syncing with the wrong server, forcing it to download incorrect rules. During the Preflight step, the Santa sensor will also upload some basic information about the MacOS machine, such as OS version and Santa version. We consider this risk to be relatively minimal, and as such have not prioritized fixing it.

### Some risks via a public API gateway:
- Download and examine (NOT MODIFY) all of our rules
- DDOS or DOS by impersonating other existing machines and corrupting the internal sync cursors
- Flood Rudolph DynamoDB with invalid or nonexistent machine IDs that do not belong to us

## Trusting Unknown Clients
How do we prevent unauthorized users from CURL’ing the ruledownload endpoint to download all of our rules?

At the moment, inbound requests must match a machine UUID. No further validation is in place to check if this matches a system inventory. Again, the risk here that anyone could download rulesets are considered a minimal risk to accept. 

Santa provides a ClientAuthCertificateFile, ClientAuthCertificatePassword, ClientAuthCertificateCN, and ClientAuthCertificateIssuerCN to do sync authentication, but this is currently not implemented yet as of 2021-04-05. Because attackers have no way to change rules and have extremely limited data to exfiltrate, we consider the risk to be minimal.