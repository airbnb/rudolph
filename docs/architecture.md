# Architecture & Design Philosophies

Rudolph is designed for Amazon Web Services. It is designed with the following principles in mind:

* Ease of deployment
* Low maintenance
* Cost Conscious

To achieve these, it uses almost exclusively serverless AWS products. Of them the most prominent are:

* Amazon API Gateway
* Amazon Lambda
* Amazon DynamoDB
* Amazon S3


## API Gateway
API Gateway houses the public hostname, endpoints, and TLS. It also implements resource validations.

All API Gateway interactions are proxied to Lambda.


## Lambda
Lambda serves as the "server" component of Rudolph. Web requests invoke Lambda functions which interact with DynamoDB and
return response structures. These structures are translated into HTTP responses that are returned to the Santa clients.


## DynamoDB
All rules, machine configurations, and uploaded sensor data are housed in DynamoDB.

To modify sensor configurations or to edit rules, you would use the Santa golang code to make edits to this DynamoDB.
You **_would not_** go through Rudolph's public API, as that API does not have any endpoints that would implement this
use case.

## S3
Amazon S3 is used to store the compiled golang service binaries that are executed by the Lambda functions.


