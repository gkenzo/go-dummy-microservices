# go-dummy-microservices
Getting to know what's up with go's microservices


"project" folder contains docker-compose instructions to execute microservices;
to run all microservices run: ``cd project && docker compose up``

### TODO:

- set up a proper read.me file
  - describe Make and Makefiles for managing services
  - describe services responsibilities
  - generate architectural docs to better illustrate how services communicate

- run seed sql when running pg containers
- implement tests of each service
- generate endpoints docs

- study and validate implementation of load balancing for internal communication