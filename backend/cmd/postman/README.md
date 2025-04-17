# Postman Collection Sync

This tool syncs the OpenAPI specification with a Postman collection using the Postman API.

## Prerequisites

1. A Postman API key
2. A Postman collection ID
3. Go 1.21 or later

## Setup

1. Get your Postman API key:

   - Go to [Postman API Keys](https://web.postman.co/settings/me/api-keys)
   - Create a new API key

2. Get your collection ID:

   - Open your collection in Postman
   - Click on the three dots (...) next to the collection name
   - Select "Share"
   - Copy the collection ID from the URL

3. Set environment variables:
   ```bash
   export POSTMAN_API_KEY="your-api-key"
   export POSTMAN_COLLECTION_ID="your-collection-id"
   ```

## Usage

1. Build the tool:

   ```bash
   go build
   ```

2. Run the sync:
   ```bash
   ./postman
   ```

The tool will:

1. Read the OpenAPI specification from `../../http/openapi.yaml`
2. Convert it to a Postman collection format
3. Update your Postman collection via the Postman API

## Features

- Automatically converts OpenAPI paths to Postman requests
- Preserves request methods, descriptions, and parameters
- Adds authentication headers where required
- Handles request bodies for POST/PUT requests

## Error Handling

The tool will exit with an error code if:

- Required environment variables are missing
- The OpenAPI spec cannot be read or parsed
- The Postman API request fails

Error messages will provide details about what went wrong.
