package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/lestrrat/go-pdebug"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {
	os.Exit(_main())
}

type Incoming struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Cert    string `json:"tls.crt"`
	Key     string `json:"tls.key"`
}

func _main() int {
	var authtokenf string
	var fifopath string
	var ns string
	var slackgw string
	flag.StringVar(&authtokenf, "authtokenfile", "", "File containing token used to authentication when posting")
	flag.StringVar(&fifopath, "fifopath", "", "path to fifo where jobs are pushed into")
	flag.StringVar(&ns, "namespace", "default", "k8s application namespace")
	flag.StringVar(&slackgw, "slackgw", "http://slackgw:4979", "slack gateway url")
	flag.Parse()

	var authtoken string
	if authtokenf != "" {
		buf, err := ioutil.ReadFile(authtokenf)
		if err != nil {
			fmt.Printf("Failed to open file '%s': %s", authtokenf, err)
			return 1
		}
		authtoken = string(buf)
	}

	if _, err := os.Stat(fifopath); err != nil { // doesn't exist
		if err := syscall.Mknod(fifopath, syscall.S_IFIFO|0666, 0); err != nil {
			// Failed to create... timing problem?
			if _, err := os.Stat(fifopath); err != nil {
				// Hmm, weird. bail
				fmt.Printf("failed to create fifo %s", fifopath)
				return 1
			}
		}
	}

	if pdebug.Enabled {
		pdebug.Printf("Reading from FIFO at %s", fifopath)
	}

	f, err := os.Open(fifopath)
	if err != nil {
		fmt.Printf("failed to open '" + fifopath + "': " + err.Error())
		return 1
	}

	r := bufio.NewReader(f)
	t := time.Tick(time.Second)
	dec := json.NewDecoder(r)
	for {
		select {
		case <-t:
		}

		// There should be at least 50 bytes
		if _, err := r.Peek(50); err != nil {
			continue
		}

		var in Incoming
		if err := dec.Decode(&in); err != nil {
			if pdebug.Enabled {
				pdebug.Printf("Failed to decode JSON: %s", err)
			}
			continue
		}
		go upload(slackgw, authtoken, ns, in)
	}

	return 0
}

func upload(slackgw, authtoken, ns string, in Incoming) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("upload %s", in.Name).BindError(&err)
		defer g.End()
	}

	s := api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name: in.Name,
			Labels: map[string]string{
				"group": "ssl-cert",
			},
		},
		Data: map[string][]byte{
			"tls.crt": []byte(in.Cert),
			"tls.key": []byte(in.Key),
		},
		Type: "Opaque",
	}

	cl, err := unversioned.NewInCluster()
	if err != nil {
		fmt.Printf("Failed upload secret %s: failed to create k8s client: %s", in.Name, err)
		// send notice to slack
		return err
	}

	svc := cl.Secrets(ns)
	if _, err := svc.Create(&s); err != nil {
		fmt.Printf("Failed upload secret %s: failed to create k8s secret: %s", in.Name, err)
		return err
	}

	values := url.Values{
		"channel": []string{in.Channel},
		"message": []string{":tada: Successfully created secret " + in.Name},
	}
	buf := bytes.Buffer{}
	buf.WriteString(values.Encode())

	req, err := http.NewRequest("POST", slackgw+"/post", &buf)
	if err != nil {
		return err
	}
	if authtoken != "" {
		req.Header.Set("X-Slackgw-Auth", authtoken)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}