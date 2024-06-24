## Core lib for microservices

We provide some helpers to simplify and unify microservice structure,
basic units and types

### Stream

We define a basic microservice message format.

Our brokers can enhance this message but must retain the core structure for pipeline processing.

Specific functions are provided for simplifying use of google's pubsub.

### Google Bigquery

Add time series queries for google big query
(this is temporarily PowerUsage based but should allow generalization eventually)



