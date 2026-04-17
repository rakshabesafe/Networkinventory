# Cloud-Native Network Inventory System

This repository contains a modern, cloud-native network inventory system designed to expose TMF-compliant APIs and manage network states via a graph database.

## Functional Requirements

The system fulfills the following core functional requirements:

1. **TMF 639 Resource Inventory API:** Provides CRUD operations for network resources (both logical and physical).
2. **TMF 638 Service Inventory API:** Provides CRUD operations for network services.
3. **Adapter-based Discovery:** Includes a modular discovery adapter service that periodically scans network elements (simulated) to discover live resource states.
4. **Reconciliation & Audit:** Includes a reconciliation engine that periodically compares the discovered live state against the desired state stored in the inventory, producing audit reports.
5. **Graph Database Storage:** All inventory data (Services and Resources) is intended to be stored in a Neo4j Graph Database to easily model complex network topologies and relationships.

## Architecture

The project is built using a **Microservices Architecture** leveraging the following technologies:
* **Go (Golang):** The backend language. The repository is structured as a Go workspace (`go.work`) containing multiple independent modules.
* **Neo4j:** Used as the primary graph database. The Go microservices connect to it using the official Neo4j Go driver.
* **Docker & Docker Compose:** Used for cloud-native containerization and local orchestration.

### Directory Structure

* `pkg/db`: Shared Neo4j database connection and retry logic.
* `pkg/models`: Shared TMF data structures representing Services and Resources.
* `services/tmf638`: Microservice exposing the Service Inventory API.
* `services/tmf639`: Microservice exposing the Resource Inventory API.
* `services/adapter`: Microservice running background network discovery tasks.
* `services/reconciliation`: Microservice running background state comparison tasks.

## Configuration

The system is configured via Environment Variables. When running via Docker Compose, these are automatically set.

* `NEO4J_URI`: Connection URI for the database (default: `neo4j://localhost:7687`).
* `NEO4J_USER`: Database username (default: `neo4j`).
* `NEO4J_PASSWORD`: Database password (default: `password`).
* `PORT`: HTTP Port for the APIs (8081 for TMF639, 8082 for TMF638).

## Usage

### How to Start

You can start the entire stack, including the Neo4j database and all microservices, using Docker Compose:

```bash
docker-compose up --build
```

*Note: The microservices include robust retry logic. They will wait for the Neo4j database to initialize and become healthy before successfully connecting. If Neo4j is completely unavailable, the APIs will fall back to using an in-memory map to allow testing without a DB.*

### Endpoints

* **TMF 639 Resource Inventory:** `http://localhost:8081/tmf-api/resourceInventoryManagement/v4/resource`
* **TMF 638 Service Inventory:** `http://localhost:8082/tmf-api/serviceInventoryManagement/v4/service`

### How to Stop

To gracefully stop and remove the containers:

```bash
docker-compose down
```

To remove the persisted Neo4j data volume as well:

```bash
docker-compose down -v
```

### Testing

You can run the unit tests across the entire Go workspace using:

```bash
go test ./...
```
