# Docker coscheduler
Implementation of a computing task co-scheduler system (with other strategies) on a computing cluster, the main idea of which is to manage idle computer resources and reassign tasks between cluster nodes for greater efficiency. 

Using the blackbox approach for tasks, which is implemented using Docker.

# Description of the system
The co-scheduler consists of 2 main components:
* Worker - the software implementation of processes that run on nodes and are responsible for completing tasks
* The scheduler itself - runs on the main node and implements cost-management strategies
  
# Necessary tools to run the system
* Docker Daemon must be installed and running on all nodes of the system.;
* The perf console utility must be installed on the nodes.

# How to run
1. To start a project, you need to create `scheduler.yaml` and `worker.yaml` files in the root of the repository, an example can be taken from `worker.example.yaml` or `scheduler.example.yaml`.
2. In the deployments folder, you need to create a `.env` file (an example can be taken from `env.example`).
3. Next, you need to generate the necessary proto modules in the project.:
> make generate
4. It is necessary to up the database with the command:
> make run-database
5. Import the corresponding environment variable in the terminal:
> export POSTGRES_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"  # db connection string example  
6. Do migrations:
> make migrate
7. Run the project:
* Worker:
> make run-worker
* Scheduler:
> make run-main

# Migrate creation
1. To perform database migration, you need to install the goose library.:
> go install github.com/pressly/goose/v3/cmd/goose@latest

2. Go to the migration folder of the corresponding microservice and run the command:
> goose create \*имя миграции\* sql

3. Import the corresponding environment variable in the terminal:
> POSTGRES_URL  # db connection string  
4. Make command:
> make migrate
