# Security policy

tezcar opens real cars and moves real money, so a security bug can mean a stolen car or
a wrong charge. Thank you for reporting one privately.

## Reporting

Report through GitHub: **Security → Report a vulnerability** on this repository. Do not
open a public issue, discussion or pull request for a vulnerability.

Please include the commit, the component, the steps to reproduce and what an attacker
could do.

## What counts

- Unlocking, starting or locking a car without an active rental of your own.
- Sending, forging or replaying commands to a vehicle, or reading its telemetry without
  authorisation.
- Renting, paying or finishing a rental on someone else's behalf, or changing the amount
  charged.
- Bypassing identity verification or the driver's-licence checks.
- Reading or changing other users' personal data, documents, selfies or trips.

## What happens next

We confirm receipt within 7 days and agree on a fix and a date. Disclosure is coordinated
and happens at the latest 90 days after the report, or earlier once a fix is deployed.
Reporters are credited in the advisory unless they ask otherwise.

## Supported versions

Until 1.0 only the latest commit on `main` receives security fixes.
