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
	Query         string `json:"query"`
	Driver        string `json:"driver"`
	Source        string `json:"source"`
	S3URL         string `json:"s3URL"`
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
		innerQuery := strings.TrimSuffix(strings.TrimSpace(params.Query), ";")
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
