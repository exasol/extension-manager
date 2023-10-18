# Embedding Extension Manager in Other Go Programs

This page describes how to integrate Extension Manager functionality into another Go program.

## Embedding the REST API

You can embed the Extension Manager's REST API in other programs that use the [Nightapes/go-rest](https://github.com/Nightapes/go-rest) library:

```go
// Create an instance of `openapi.API`
api := openapi.NewOpenAPI()
// Create a new configuration
config := extensionController.ExtensionManagerConfig{
    ExtensionRegistryURL: "https://example.com/registry.json", 
    BucketFSBasePath: "/buckets/bfsdefault/default/",
    ExtensionSchema: "EXA_EXTENSIONS",
}
// Add endpoints
err := restAPI.AddPublicEndpoints(api, config)
if err != nil {
    return err
}
// Start the server
```

## Embedding the Extension Controller

If you want to directly use the extension controller in your application you can use the following code as an example:

```go
// Create a new configuration
config := extensionController.ExtensionManagerConfig{
    ExtensionRegistryURL: "https://example.com/registry.json", 
    BucketFSBasePath: "/buckets/bfsdefault/default/",
    ExtensionSchema: "EXA_EXTENSIONS",
}
// Create controller and handle configuration validation error
ctrl, err := extensionController.CreateWithValidatedConfig(config)
if err != nil {
    return err
}

// Create database connection (required as an argument for all controller methods)
var db *sql.DB = createDBConnection()

// Call controller method and process result. Use a custom context if available.
extensions, err := ctrl.GetAllExtensions(context.Background(), db)
// ...
```
