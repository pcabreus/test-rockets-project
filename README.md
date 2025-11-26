# test-rockets-project

## Description

This project implements a backend service for managing rocket data and launch scheduling. It receives rocket updates through a webhook and exposes a RESTful API to retrieve up-to-date rocket information. The service processes events asynchronously while maintaining strict ordering guarantees per event channel.

## Architecture Overview

### Rocket Event Webhook

The system provides a webhook endpoint to capture all rocket updates. Events are stored in a queue and processed asynchronously to update rocket statuses.

**Endpoint:** `POST http://localhost:8080/messages`

**Request Payload Example:**
```json
{
    "metadata": {
        "channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
        "messageNumber": 1,    
        "messageTime": "2022-02-02T19:39:05.86337+01:00",                                          
        "messageType": "RocketLaunched"                             
    },
    "message": {                                                    
        "type": "Falcon-9",
        "launchSpeed": 500,
        "mission": "ARTEMIS",
        "by": 3000,
        "reason": "PRESSURE_VESSEL_FAILURE",
        "newMission": "SHUTTLE_MIR"
    }
}
```

**Event Status:** Events have two states: `pending` and `processed`. Each event is processed only once, and certain event types are restricted to single processing.

### Event Processing

The event processor consumes the event queue asynchronously:

- **Processing Interval:** Every 1 second
- **Ordering Guarantee:** Events are processed sequentially per channel
- **Processing Logic:**
  1. Retrieve all `pending` events
  2. Iterate through the event list
  3. Process only if the event's message number matches the expected next number for that channel
  4. Apply domain logic based on the event type
  5. Persist updated rocket data
  6. Increment the expected message number for the channel

This ensures strict ordering and prevents out-of-sequence event processing.

### Rockets API

**Get All Rockets**

```
GET /rockets
```

Response (200 OK):
```json
[
    {
        "Channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
        "Type": "Falcon-9",
        "Speed": 500,
        "Mission": "ARTEMIS",
        "Status": "EXPLODED",
        "Reason": ""
    }
]
```

**Get Rocket by Channel**

```
GET /rockets/{channelId}
```

Response (200 OK):
```json
{
    "Channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
    "Type": "Falcon-9",
    "Speed": 500,
    "Mission": "ARTEMIS",
    "Status": "EXPLODED",
    "Reason": ""
}
```

## How to Run

### Prerequisites
- Go 1.16 or higher

### Option 1: Run with Go Directly

```bash
go run main.go
```

The service will start on `http://localhost:8080`

### Option 2: Build and Run Compiled Binary

```bash
# Build the executable
go build -o test-rockets-project

# Run the compiled binary
./test-rockets-project
```

### Configuration

The service listens on:
- **Host:** localhost
- **Port:** 8080

## Testing

### Manual Testing with cURL

**Send a rocket event:**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
      "messageNumber": 1,
      "messageTime": "2022-02-02T19:39:05.86337+01:00",
      "messageType": "RocketLaunched"
    },
    "message": {
      "type": "Falcon-9",
      "launchSpeed": 500,
      "mission": "ARTEMIS"
    }
  }'
```

**Retrieve all rockets:**
```bash
curl http://localhost:8080/rockets
```

**Retrieve specific rocket:**
```bash
curl http://localhost:8080/rockets/193270a9-c9cf-404a-8f83-838e71d9ae67
```

### Automated Testing

Unit and integration tests are not currently provided due to project timeline constraints. 

TODO tests for:
- Event webhook payload validation
- Event ordering per channel
- Rocket state transitions
- API response formatting