# GoMicroServices

## Question for Follow-Up
 - Email must be Unique in DB (contact purposes), responsibility is in usecase or Database ?
 - Update in database, if ID (primary Key) doesn't exist, responsibility to check presence of this ID is in usecase or Database ?
 - Use type from domain layer of other microservice ? (I.E: Account in Chat service)
 - Infra ? apart for database, what use ?
 - `Searching an object in a database based on a subquery string isn't a good idea.` What do you mean ?


## TODO
 - Small Dockerfile to build each MS
 - Small Docker Compose to link MS and use Docker container
 - Small Makefile to wrap some useful docker-compose cmd
 - Small CI/CD for building MS for now


## Architecture Thinking

### services:
- ????

### ORM:
- Gorm (EZ but not very customizable)
- Ent (medium, code generation)
- prisma (little hard, but swagger and it's rocks !) https://github.com/prisma/prisma-client-go

### Databases:
- Postgres (EZ)
- ElasticSearch (Medium)
- Cassandra (Hard)
- MongoDb (très peu use)

### Transports:
- Http (easy)
  Router/frameworks:
    - Mux (low level)
    - Gin (high level)
    - Gorilla (several problems encountered)
- Grpc (medium)
- Kafka (hard)
- socket (nul !)
- Rabbit MQ (never used)
- redis (bof)
- WebSocket (pk pas ?)

### Infra:
- dockerisation with docker-compose
- K8s ??
- Hosting on scaleway ? (ssh ez config)
- CI/CD with Github action ?

### Tests:
- Unit (tool intégré dans le binaire go)
- E2E ?

### Possible Bonus:
- Front (VueJs 3 ?)
- proxy, load balancing (nginx ?)
-

### Organisation:
- git avec Branch protection
- airtable ? Notion ? Rien ?