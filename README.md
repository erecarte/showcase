# Numeral exercise
- Tried separating concerns into API, Domain/Service and Storage. The `payment_orders` package contains everything needed to process that resource with the idea that in the future the application could be extended with more resources
- Created an API client to talk to the API which encapsulates communication concerns such as authorization, marshalling, etc.
- Have hard-coded the credentials because of time. I would have done it better if I had more time
- The tests use this API client and created a couple of them:
  - Not authenticated user returns 401
  - Wrong JSON format returns 400
  - Success case
  - With more time, I would have done much more extensive testing
- Created a "bank" application which simulates responding to the XML request files.
## NOTES
- Amount is misspelled in the sample JSON and schema. Left them as they are
- Did not have prior experience with SQLITE...I hope I set it up right :)
- I think there is an error in the sample xml/xsd. Inside Dbtr, the child element is CdtrAcct where I guess it should be DbtrAcct? I left it as it is
- In order to run the tests, the following steps need to be followed:
  - In one tab, run `go run cmd/numeral/main.go`
  - In another tab, run `go run cmd/bank/main.go`
  - run `go test ./...`
