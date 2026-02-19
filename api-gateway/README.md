# API Gateway

This is the API Gateway for the LeConfiguer project, responsible for handling requests and routing them to the appropriate services.

## Features
- Create, Read, Update, Delete configuration items
- Validate incoming JSON payloads
- Filter and sort configurations by name, environment, or tags
- Send notifications to Notifier after successful create/update
- Simple API key-based authentication

## Requirements
- JSON API over HTTP
- Connects to Config Storage via REST
- Connects to Notifier for sending updates

## Endpoints

### Create Configuration Item
- **POST** `/configurations`
- Request Body: JSON object representing the configuration item.

### Read Configuration Item
- **GET** `/configurations/{id}`

### Update Configuration Item
- **PUT** `/configurations/{id}`
- Request Body: JSON object representing the updated configuration item.

### Delete Configuration Item
- **DELETE** `/configurations/{id}`

### Filter and Sort Configurations
- **GET** `/configurations?name={name}&environment={environment}&tags={tags}`

## Authentication
- API key required in the header for all requests.

## Notifications
- Sends notifications to Notifier service after successful create/update operations.

## Connecting to Services
- Connects to Config Storage via REST API.
- Connects to Notifier for sending updates.