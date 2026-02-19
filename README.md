# LeConfiguer

# API Gateway Service

## Overview
The API Gateway acts as the central entry point for all clients (other services or users) to interact with configuration data. It provides a JSON-based CRUD interface and routes requests to the Config Storage service. After successful operations, it triggers notifications via the Notifier service.

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


# Config Storage Service

## Overview
The Config Storage service is responsible for storing and retrieving configuration items. It supports CRUD operations and provides the current state of configurations. Each configuration can include environment, tags, and JSON content.

## Features
- Store configuration items with fields: id, name, environment, json_content, tags
- CRUD operations for API Gateway
- Filter configurations by environment or tags
- Maintain lightweight validation of JSON content
- Notify Versioning Service on changes (new version created)

## Requirements
- Database: PostgreSQL / SQLite / local JSON file for prototyping
- JSON API over HTTP for internal communication


# Versioning Service

## Overview
The Versioning Service tracks changes to configuration items and keeps a full history of all updates. It allows clients to retrieve version lists, inspect individual versions, and optionally rollback to previous configurations.

## Features
- Record every configuration change with metadata (version, timestamp, author)
- Retrieve list of versions per configuration
- View details of any configuration version
- Optional: rollback to a previous version
- Limit number of stored versions (e.g., last 10 changes)

## Requirements
- Receives update notifications from Config Storage or API Gateway


# Notifier Service

## Overview
The Notifier service sends JSON-based notifications to subscribed services whenever a configuration item is created or updated. This ensures that all dependent services are aware of configuration changes in real time.

## Features
- Receive events from API Gateway or Config Storage
- Format JSON notification:
  {
    "config_id": "1234",
    "action": "updated",
    "timestamp": "2026-02-19T14:00:00Z",
    "author": "service_xyz"
  }
- Send HTTP POST notifications to subscribed services
- Optional: retry failed notifications
- Log success/failure of notifications
