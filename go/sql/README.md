# Run SQL

This Airplane task allows you to quickly execute SQL queries against a database.

To get started, [install the Airplane CLI](https://docs.airplane.dev/cli/airplane-in-60s), then run:

```
airplane tasks deploy -f github.com/airplanedev/examples/go/sql/airplane.yml
```

Currently, this task only supports queries on PostgreSQL DBs. PRs are welcome to extend this support to MySQL or other kinds of SQL databases!
