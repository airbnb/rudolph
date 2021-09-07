# System Risks and Security

## Unauthorized Clients
API Gateway is publicly reachable and implements no authorization by default.

Until a more robust authorization plane can be implemented externally to Santa, like enforcement of client-certificates, Rudolph API server is globally available to any resource on the internet.

Rudolph's API allows external resources to:
 - Upload OS data
 - Download sensor configuration
 - Upload logs
 - Download rules

The public API in no way allows users to modify rules or sensor configurations. The most concerning case is the
unauthorized downloading of rules in your system. Depending on your team's security posture, this may be
considered "within acceptable risks".


## Trusting the Wrong Server
A rare (and unconfirmed) vulnerability is when Santa agents are tricked into sync'ing with the wrong server. DNS poisoning
is one possible attack vector.

Santa agents support server trusts by specifying a ServerAuthRootsFile in the macOS configuration profile. This file contains a SSL Certificate chain which is used to verify the identity of the server.

However, there is a bug in Santa: When this file is malformed or if the certificate does not match the SyncBaseURL, Santa defaults back to the Apple MacOS keychain, which will blanket trust anything signed by Amazon Web Services.

Thus it is possible for an attacker to trick a Santa sensor into syncing with the wrong server, forcing it to download incorrect rules. During the Preflight step, the Santa sensor will also upload some basic information about the MacOS machine, such as OS version and Santa version.

We consider this risk to be relatively minimal, and as such have not prioritized fixing it.


## Some Mitigations
Placing your entire Rudolph environment behind a VPN will prevent unauthorized clients from reading your rules. This
DOES NOT protect you against DNS poisoning attacks.

To configur Rudolph to only trust VPN IPs, use `allowed_inbound_ips` to restrict access to specific CIDR blocks.
