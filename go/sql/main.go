package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	airplane "github.com/airplanedev/go-sdk"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func main() {
	airplane.Run(run)
}

type Parameters struct {
	// The SQL query to execute.
	SQL string `json:"sql"`
	// The type of SQL database to run the query on. Currently this can only be "postgres".
	Driver string `json:"driver"`
	// Where to get the DB DSN from. Defaults to "env", which pulls it from an `AIRPLANE_DSN`
	// environment variable. Can also be set to "s3" to pull the DSN from an s3 bucket.
	Source string `json:"source"`
	// If source="s3", this should be the full URL (s3://...) of an S3 document.
	S3URL string `json:"s3URL"`
	// If source="s3", this can be used to fetch a specific element from the above S3 document
	// where the contents are parsed as JSON. If not set, defaults to treating the document's
	// contents as the DSN.
	S3JSONPointer string `json:"s3JSONPointer"`
}

func run(ctx context.Context) error {
	var params Parameters
	if err := airplane.Parameters(&params); err != nil {
		return err
	}
	// TODO: would be great if Parameters() could also read env vars for you,
	// f.e. AIRPLANE_ + struct tag or field name.
	if driver := os.Getenv("AIRPLANE_DRIVER"); driver != "" {
		params.Driver = driver
	}
	if source := os.Getenv("AIRPLANE_SOURCE"); source != "" {
		params.Source = source
	}
	if s3URL := os.Getenv("AIRPLANE_S3_URL"); s3URL != "" {
		params.S3URL = s3URL
	}
	if s3JSONPointer := os.Getenv("AIRPLANE_S3_JSON_POINTER"); s3JSONPointer != "" {
		params.S3JSONPointer = s3JSONPointer
	}

	dsn, err := getDSN(params)
	if err != nil {
		return err
	}

	switch params.Driver {
	case "postgres":
		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			return errors.Wrap(err, "connecting to db")
		}
		defer db.Close()

		// Since we are wrapping the user-provided query, we need to remove the semicolon, if present.
		innerQuery := strings.TrimSuffix(strings.TrimSpace(params.SQL), ";")
		rows, err := db.Queryx(fmt.Sprintf(`SELECT row_to_json(t) FROM (%s) t`, innerQuery))
		if err != nil {
			return errors.Wrap(err, "running query on db")
		}
		defer rows.Close()

		for rows.Next() {
			var buf []byte
			if err := rows.Scan(&buf); err != nil {
				return errors.Wrap(err, "scanning row")
			}

			airplane.MustNamedOutput("rows", string(buf))
		}
	default:
		return errors.Errorf(`invalid driver: expected one of ["postgres"] got %s`, params.Driver)
	}

	return nil
}
