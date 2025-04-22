# BrowserBase Go SDK (Generated)

This repository contains a generated Go client for the [BrowserBase API](https://docs.browserbase.com/), built on top of an OpenAPI specification and powered by [ogen](https://github.com/ogen-go/ogen).

## How is the code generated?

- The API surface and type schemas are described in the `openapi.yaml` (OpenAPI v3) file.
- Code is generated automatically by [ogen](https://github.com/ogen-go/ogen) using the `go generate` directive in `gen.go`:

    ```sh
    go generate ./...
    ```
    ...which runs:
    ```sh
    go run github.com/ogen-go/ogen/cmd/ogen@latest --target browserbase --clean openapi.yaml
    ```
  
- This generates the strongly-typed Go client and models for all defined endpoints.
- Example tests in `session_test.go` show how to authenticate and use the client.

## Learn more

For full API documentation and feature guides, see the official [BrowserBase documentation](https://docs.browserbase.com/).

---
This SDK is generated and not hand-edited; changes to the client code should be made via OpenAPI spec changes and regeneration.
