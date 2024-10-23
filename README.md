### Observability POC

This repo contains code for a crud app written in golang and postgres. The metrics for the app is collected using prometheus and visualised using grafana. 

Metric collected:
- Query execution time
- Errors related to DB
- Active HTTP requests
- Total HTTP requests
- HTTP request duration

Query execution time:
<img width="1066" alt="Screenshot 2024-10-23 at 3 39 03 PM" src="https://github.com/user-attachments/assets/89c5c3bd-e507-4236-8c9b-ef5bb6b9756a">

HTTP response time:
<img width="1061" alt="Screenshot 2024-10-23 at 3 40 57 PM" src="https://github.com/user-attachments/assets/17b9fa2c-875f-4613-82e6-42d0cd499725">


To use this app, docker must be installed on the computer.

Steps to start the application:

Clone the repo
```bash
  git clone https://github.com/Dev79844/observability-poc.git
```

Start the application using docker compose

```bash
  docker compose up -d
```
