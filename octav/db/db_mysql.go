// Note: add build tags if/when we support multiple databases

package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/lestrrat/go-pdebug"
	"github.com/shogo82148/go-sql-proxy"
)

type NullTime struct {
	mysql.NullTime
}

var errNoEnv = errors.New("no env")

func readEnvConfig(name, ename string, dst *string) error {
	if err := readEnvConfigName(name, ename, dst); err != nil {
		switch err {
		case errNoEnv:
		default:
			return err
		}
	}

	if err := readEnvConfigFile(name, ename+"_FILE", dst); err != nil {
		switch err {
		case errNoEnv:
		default:
			return err
		}
	}

	// Nothing to do
	return nil
}

func readEnvConfigName(name, ename string, dst *string) error {
	v := os.Getenv(ename)
	if v == "" {
		return errNoEnv
	}
	if pdebug.Enabled {
		pdebug.Printf("Using %s from file specified in environment variable %s", name, ename)
	}

	*dst = v
	return nil
}

func readEnvConfigFile(name, ename string, dst *string) error {
	f := os.Getenv(ename)
	if f == "" {
		return errNoEnv
	}

	if pdebug.Enabled {
		pdebug.Printf("Using %s from file specified in environment variable %s", name, ename)
	}

	v, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	*dst = strings.TrimSpace(string(v))
	return nil
}

func ConfigureDSN() (tlsname string, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("ConfigureDSN").BindError(&err)
		defer g.End()
	}

	c := mysql.Config{
		User:      "root",
		DBName:    "octav",
		Net:       "tcp",
		ParseTime: true,
	}

	if err = readEnvConfig("username", "OCTAV_MYSQL_USERNAME", &c.User); err != nil {
		return "", err
	}

	if err = readEnvConfig("password", "OCTAV_MYSQL_PASSWORD", &c.Passwd); err != nil {
		return "", err
	}

	if err = readEnvConfig("net", "OCTAV_MYSQL_NET", &c.Net); err != nil {
		return "", err
	}

	if err = readEnvConfig("address", "OCTAV_MYSQL_ADDRESS", &c.Addr); err != nil {
		return "", err
	}

	if err = readEnvConfig("dbname", "OCTAV_MYSQL_DBNAME", &c.DBName); err != nil {
		return "", err
	}
	if c.Addr != "" {
		if c.Net == "" {
			c.Net = "tcp"
		}
	}

	if c.Net == "tcp" {
		portSuffix, err := regexp.Compile(`:\d+$`)
		if err != nil {
			return "", err
		}
		if !portSuffix.MatchString(c.Addr) {
			c.Addr = c.Addr + ":3306"
		}
	}

	tlsname, err = trySetupTLS()
	switch err {
	case ErrNoTLSRequested:
		// no op. we weren't requested to do TLS
		if pdebug.Enabled {
			pdebug.Printf("TLS is not requested")
		}
	case nil:
		if pdebug.Enabled {
			pdebug.Printf("TLS enabled, going to add '%s' as TLSConfig", tlsname)
		}
		c.TLSConfig = tlsname
	default:
		// now *this* is an error
		if pdebug.Enabled {
			pdebug.Printf("Failed to setup TLS: %s", err)
		}
		return "", err
	}

	dsn := c.FormatDSN()
	return dsn, nil
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

func trySetupTLS() (string, error) {
	caCertFile := os.Getenv("OCTAV_MYSQL_CA_CERT_FILE")
	clientCertFile := os.Getenv("OCTAV_MYSQL_CLIENT_CERT_FILE")
	clientKeyFile := os.Getenv("OCTAV_MYSQL_CLIENT_KEY_FILE")

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
		return "", ErrNoTLSRequested
	}

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return "", err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return "", err
	}
	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return "", err
	}
	clientCert = append(clientCert, certs)

	tlsname := "custom-tls"
	mysql.RegisterTLSConfig(tlsname, &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	})

	return tlsname, nil
}

func onConnect(db *sql.DB) error {
	_, err := db.Exec(`SET SESSION sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'`)
	if err != nil {
		return err
	}
	return nil
}
