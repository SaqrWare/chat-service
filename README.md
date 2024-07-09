# Chat Service

This README provides instructions on setting up and running the Chat Service with the requirement sent in the assessment
developed in Go.

## Prerequisites

Before you begin, ensure you have installed Go (version 1.15 or later) on your system. You can download it
from [https://golang.org/dl/](https://golang.org/dl/).

## Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/SaqrWare/chat-service.git
cd chat-service
```

2. Install the dependencies:

```bash
go mod tidy
```

## Configuration

The required environment variables to run the app

    ```bash
    export PORT=8080
    CASSANDRA_HOST=127.0.0.1
    CASSANDRA_PORT=9042
    CASSANDRA_KEYSPACE=chat
    REDIS_ADDR=localhost:6379
    REDIS_PASSWORD=
    REDIS_DB=0
    ```

## Running the App

1. Build the app:

```bash
go build -o bin/chat-service main.go
```

2. Run the app:

```bash
./bin/chat-service
```

## Running with docker

1. Build the docker image:

```bash
docker build -t chat-service .
```

2. Run the docker container:

```bash
docker run -p 8080:8080 chat-service
```

## API Endpoints

### Messages

Messages APIs require authentication token in the headers, Login API returns the authentication token

#### Send Message

- **URL:** `/api/v1/message/send`
- **Method:** POST
- **Headers:**
    - Authorization: $AUTH_TOKEN // auth token of the sender
- **Body:**

```json
{
  "receiver": "9c82f29a-3cc7-11ef-98fc-30f6ef0622cd",
  // sample uuid
  "message": "laskndwelkdfewjkfn"
}
```

#### Get Messages

- **URL:** `/api/v1/message`
- **Method:** GET
- **Headers:**
    - Authorization: $AUTH_TOKEN // auth token of the receiver
- **Query Params:**
    - page: 1
    - limit: 10
    - lastID: 9c82f29a-3cc7-11ef-98fc-30f6ef0622cd // optional

### Users

#### Register

- **URL:** `/api/v1/user/register`
- **Method:** POST
- **Body:**

```json
{
  "username": "sampleUsername",
  "firstName": "Sample",
  "lastName": "User",
  "email": "sampleuser@example.com",
  "password": "securePassword123"
}
```

#### Login

Login api uses email or username as the identifier, and return the authentication token

- **URL:** `/api/v1/user/login`
- **Method:** POST
- **Body:**

```json
{
  "identifier": "sampleuser@example.com",
  "password": "securePassword123"
}
```

## Database Schema

### Cassandra

```sql
CREATE KEYSPACE chat WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}
                      AND durable_writes = true;
CREATE TABLE chat."user"
(
    id         uuid PRIMARY KEY,
    created_at timestamp,
    email      text,
    first_name text,
    last_name  text,
    "password" text,
    username   text
);

CREATE TABLE chat.message
(
    id          uuid PRIMARY KEY,
    content     text,
    created_at  timestamp,
    delivered   boolean,
    receiver_id uuid,
    sender_id   uuid
);

// indexes
CREATE INDEX IF NOT EXISTS ON chat.message (sender_id);
CREATE INDEX IF NOT EXISTS ON chat.message (receiver_id);
CREATE INDEX IF NOT EXISTS ON chat.message (created_at);
CREATE INDEX IF NOT EXISTS ON chat."user" (username);
CREATE INDEX IF NOT EXISTS ON chat."user" (email);
```

You can create the database schema with tools like cassandra-web

```bash
docker run  --name cassandra-web \
-e CASSANDRA_HOSTS=127.0.0.1 \
-e CASSANDRA_PORT=9042 \
-e CASSANDRA_USERNAME=user \
-e CASSANDRA_PASSOWRD=pass \
-p 3000:3000 \
--net=host \
delermando/docker-cassandra-web:v0.4.0
```

## Caching

The app uses Redis for caching the messages and user data after authentication

### Caching mechanism for messages

in the API `GET /api/v1/message` the app uses the query params and user id to generate a key for the cache, the key
is in the format `user_id:page:limit:lastID` and the value is the messages in the format of `[]byte` marshalled to json,
if the key is found in the cache the app will return the value from the cache, if not the app will query the database

in the API `POST /api/v1/message/send` the app will delete the cache for the sender and receiver to update the cache

The cache is set to expire after 1 minute