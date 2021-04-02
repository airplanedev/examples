# Query DB

This is a public pre-built Airplane Task for running a SQL query on a SQL DB.

## Configuration

The recommended method for using this task is to provide a DSN as an environment variable:

```sh
AIRPLANE_DSN="..." go run ./main.go --driver="postgres" --query "SELECT * FROM users"
```

Alternatively, this task can dynamically fetch a DSN from a remote source. If configured, this takes precedence over any environmental DSN configuration.

Currently, only files in S3 are supported. By default, it assumes the file itself just contains the DSN as a string. However, you can provide an optional [JSON pointer](https://tools.ietf.org/html/rfc6901) to select the DSN from a JSON document.

```sh
go run ./main.go \
  --query "SELECT * FROM users" \
  --source="s3" \
  --s3-url="s3://my-bucket/config.json" \
  --s3-json-path="/production/us-west-2/dsn"
```
