# Setup Commands

## 1. Prepare Dependencies
Ensure Go modules are tidy before building.
```bash
go mod tidy
```

## 2. Build Images Locally
Since we are using local images in `docker-compose.yml`, build them first.
```bash
docker build -t orderservice:local -f docker/orderservice.Dockerfile .
docker build -t frontend:local -f docker/sws.Dockerfile .
```

## 3. Initialize Swarm
Initialize the swarm on the manager node.
```bash
docker swarm init
```

## 4. Deploy Stack
Deploy the stack to the swarm.
```bash
docker stack deploy -c docker-compose.yml ordersystem
```

## 5. Verify Deployment
Check the services and nodes.

```bash
docker service ls
docker node ls
```

### Current Node Status
```
ID                            HOSTNAME         STATUS    AVAILABILITY   MANAGER STATUS   ENGINE VERSION
0opoowzqpmfnqcz3h0o5jhgsw *   docker-desktop   Ready     Active         Leader           28.5.1
```
