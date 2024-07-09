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

### Device

The framework for processing device data. 
**Device** provides a pipeline with various fields for storing and processing the data.

Typically this is taken from a store/cache and attached to the interpreted device payload.
Pipeline elements can then fork, process, store based around these 
e.g 
* a microservice that stores all device information for a location
* a microservice that sums a field over all messages in a system

etc
