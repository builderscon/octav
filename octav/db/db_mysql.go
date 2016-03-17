// Note: add build tags if/when we support multiple databases

package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/lestrrat/go-pdebug"
	"github.com/shogo82148/go-sql-proxy"
)

type NullTime struct {
	mysql.NullTime
}

func readEnvConfig(name, ename string, dst *string) error {
	f := os.Getenv(ename)
	if f == "" {
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("Using %s from file specified in environment variable %s", name, ename)
	}

	v, err := ioutil.ReadFile(f)
	if err != nil {
		if pdebug.Enabled {
			pdebug.Printf("Failed to read file %s: %s", v, err)
		}
		return err
	}
	*dst = string(v)
	return nil
}

func DSNConfig() (*mysql.Config, error) {
	if pdebug.Enabled {
		g := pdebug.Marker("DSNConfig")
		defer g.End()
	}

	c := mysql.Config{
		User:      "root",
		DBName:    "octav",
		ParseTime: true,
	}

	if err := readEnvConfig("username", "OCTAV_MYSQL_USERNAME", &c.User); err != nil {
		return nil, err
	}

	if err := readEnvConfig("password", "OCTAV_MYSQL_PASSWORD", &c.User); err != nil {
		return nil, err
	}

	if err := readEnvConfig("address", "OCTAV_MYSQL_ADDRESS", &c.User); err != nil {
		return nil, err
	}

	if err := readEnvConfig("dbname", "OCTAV_MYSQL_DBNAME", &c.User); err != nil {
		return nil, err
	}

	return &c, nil
}

func driverName() string {
	driverName := "mysql"
	if Trace {
		driverName = "mysql-trace"
		var out io.Writer
		if pdebug.Enabled {
			// Send the output to the same place as pdebug
			out = pdebug.DefaultCtx.Writer
		} else {
			out = os.Stderr
		}
		sql.Register(driverName, proxy.NewTraceProxy(&mysql.MySQLDriver{}, log.New(out, "", 0)))
	}
	return driverName
}

func trySetupTLS() error {
	caCertFile := os.Getenv("OCTAV_MYSQL_CA_CERT")
	clientCertFile := os.Getenv("OCTAV_MYSQL_CLIENT_CERT")
	clientKeyFile := os.Getenv("OCTAV_MYSQL_CLIENT_KEY")

	if pdebug.Enabled {
		pdebug.Printf("Setting up TLS...")
		pdebug.Printf("   -> CA Cert: %s", caCertFile)
		pdebug.Printf("   -> Client Cert: %s", clientCertFile)
		pdebug.Printf("   -> Client Key: %s", clientKeyFile)
	}

	if caCertFile == "" || clientCertFile == "" || clientKeyFile == "" {
		if pdebug.Enabled {
			pdebug.Printf("Some or all fields are not specified. Abort setting up TLS...")
		}
		return ErrNoTLSRequested
	}

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return err
	}
	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return err
	}
	clientCert = append(clientCert, certs)
	mysql.RegisterTLSConfig("custom-tls", &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	})

	return nil
}

func onConnect(db *sql.DB) error {
	_, err := db.Exec(`SET SESSION sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'`)
	if err != nil {
		return err
	}
	return nil
}
