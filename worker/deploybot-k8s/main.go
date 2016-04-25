package main

// listen to a fifo, and create an ingress

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
	"text/template"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/dns/v1"

	"github.com/lestrrat/go-pdebug"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util/yaml"
)

func main() {
	os.Exit(_main())
}

type Incoming struct {
	Target  string            `json:"target"` // deployment, ingress, etc
	Mode    string            `json:"mode"`   // create, list, delete, etc
	Name    string            `json:"name"`
	Channel string            `json:"channel"`
	Args    map[string]string `json:"args"`
	slackgw string            // used internally
	token   string            // used internally
}

type Bot struct {
	dns       *dns.Service
	namespace string
	projectID string
	root      string
	zone      string
}

func _main() int {
	var authtokenf string
	var fifopath string
	var ns string
	var projectID string
	var slackgw string
	var zone string
	var root string
	flag.StringVar(&authtokenf, "authtokenfile", "", "File containing token used to authentication when posting")
	flag.StringVar(&fifopath, "fifopath", "", "path to fifo where jobs are pushed into")
	flag.StringVar(&projectID, "project_id", "", "GCE project ID")
	flag.StringVar(&zone, "zone", "", "DNS zone")
	flag.StringVar(&ns, "namespace", "default", "k8s application namespace")
	flag.StringVar(&slackgw, "slackgw", "http://slackgw:4979", "slack gateway url")
	flag.StringVar(&root, "root", "/tmpl", "location of templates")
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

	ctx := context.Background()
	httpcl, err := google.DefaultClient(ctx, dns.NdevClouddnsReadwriteScope)
	if err != nil {
		fmt.Printf("failed to create OAuth'ed HTTP client: %s", err)
		return 1
	}

	svc, err := dns.New(httpcl)
	if err != nil {
		fmt.Printf("failed to dns client: %s", err)
		return 1
	}

	bot := Bot{
		dns:       svc,
		namespace: ns,
		projectID: projectID,
		root:      root,
		zone:      zone,
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
		in.slackgw = slackgw
		in.token = authtoken

		switch in.Target {
		case "ingress":
			switch in.Mode {
			case "create":
				go bot.IngressCreate(in)
			case "delete":
				go bot.IngressDelete(in)
			case "list":
				go bot.IngressList(in)
			case "get":
				go bot.IngressGet(in)
			}
		case "deployment":
		}
	}

	return 0
}

func (in *Incoming) reply(s string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("Incoming.reply").BindError(&err)
		defer g.End()

		pdebug.Printf("Posting to '%s'", in.slackgw+"/post")
	}

	values := url.Values{
		"channel": []string{in.Channel},
		"message": []string{s},
	}
	buf := bytes.Buffer{}
	buf.WriteString(values.Encode())

	req, err := http.NewRequest("POST", in.slackgw+"/post", &buf)
	if err != nil {
		return err
	}

	if token := in.token; token != "" {
		req.Header.Set("X-Slackgw-Auth", token)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		if pdebug.Enabled {
			pdebug.Printf("Slack gateway returned with '%s'", res.Status)
		}
		return errors.New("bad response from slack gateway: " + res.Status)
	}

	return err
}

func (b *Bot) IngressDelete(in Incoming) (err error) {
	// Name is NOT a domain: it's the actual kubernetes ingress' name
	if pdebug.Enabled {
		g := pdebug.Marker("b.IngressDelete %s", in.Name).BindError(&err)
		defer g.End()
	}

	cl, err := unversioned.NewInCluster()
	if err != nil {
		in.reply(":exclamation: failed to create a client: " + err.Error())
		return err
	}

	// Find the previous Ingress instance that was running for this resource
	ingress, err := cl.Ingress(b.namespace).Get(in.Name)
	if err != nil {
		in.reply(":exclamation: failed to lookup ingress '" + in.Name + "'")
		return err
	}

	hostname := ingress.Labels["hostname"]
	if hostname != "" {
		// Does this ingress have an IP address?
		switch len(ingress.Status.LoadBalancer.Ingress) {
		case 0:
			in.reply(":exclamation: no ip address associated with this ingress...")
		case 1:
			addr := ingress.Status.LoadBalancer.Ingress[0].IP
			in.reply(":white_check_mark: ingress has an IP address '" + addr + "'")

			rrslist, err := b.dns.ResourceRecordSets.List(b.projectID, b.zone).Name(hostname + ".").Type("A").Do()
			if err != nil {
				in.reply(":exclamation: failed to list changes for '" + in.Name + ".': " + err.Error())
				return err
			}

			// Theres should be just one

			oldips := make([]string, 0, len(rrslist.Rrsets[0].Rrdatas))
			newips := make([]string, 0, len(rrslist.Rrsets[0].Rrdatas))
			for _, rrd := range rrslist.Rrsets[0].Rrdatas {
				oldips = append(oldips, rrd)
				if rrd != addr {
					newips = append(newips, rrd)
				}
			}

			in.reply(":white_check_mark: removing '" + addr + "' from DNS entries for '" + hostname + "'")
			ch := dns.Change{
				Deletions: []*dns.ResourceRecordSet{
					&dns.ResourceRecordSet{
						Kind:    "dns#resourceRecordSet",
						Name:    hostname + ".", // need to terminate with a "."
						Rrdatas: oldips,
						Ttl:     60,
						Type:    "A",
					},
				},
				Additions: []*dns.ResourceRecordSet{
					&dns.ResourceRecordSet{
						Kind:    "dns#resourceRecordSet",
						Name:    hostname + ".", // need to terminate with a "."
						Rrdatas: newips,
						Ttl:     60,
						Type:    "A",
					},
				},
			}

			if _, err := b.dns.Changes.Create(b.projectID, b.zone, &ch).Do(); err != nil {
				in.reply(":exclamation: failed to delete DNS entries for '" + hostname + "' (" + addr + ")")
				return err
			}
		default:
			in.reply(":exclamation: ingress '" + in.Name + "' has more than one IP address. Don't know what to do:")
			for _, ingaddr := range ingress.Status.LoadBalancer.Ingress {
				in.reply(":point_right: " + ingaddr.IP)
			}
			in.reply(":exclamation: make sure to fix DNS records")
		}
	}

	in.reply(":white_check_mark: deleting ingress '" + in.Name + "'")
	if err := cl.Ingress(b.namespace).Delete(in.Name, nil); err != nil {
		in.reply(":exclamation: failed to delete ingress '" + in.Name + "'")
		return err
	}
	in.reply(":tada: successfully deleted ingress and related DNS entries")
	return nil
}

func (b *Bot) IngressGet(in Incoming) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("b.IngressGet %s", in.Name).BindError(&err)
		defer g.End()
	}

	cl, err := unversioned.NewInCluster()
	if err != nil {
		in.reply(":exclamation: failed to create a client: " + err.Error())
		return err
	}

	label, err := labels.Parse("hostname=" + in.Name)
	if err != nil {
		in.reply(":exclamation: failed to parse label 'hostname=" + in.Name + "': " + err.Error())
		return err
	}

	list, err := cl.Ingress(b.namespace).List(api.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		in.reply(":exclamation: failed to list ingress")
		return err
	}

	if len(list.Items) <= 0 {
		in.reply(":exclamation: failed to find matching ingress")
		return
	}

	buf := bytes.Buffer{}
	buf.WriteString(":white_check_mark: Existing ingress objects:\n")
	for _, oldingress := range list.Items {
		jsonbuf, err := json.MarshalIndent(oldingress, "", "  ")
		if err != nil {
			buf.WriteString("Failed to encode ingress object: " + err.Error())
		} else {
			buf.WriteString("```")
			buf.Write(jsonbuf)
			buf.WriteString("```")
		}
		buf.WriteString("\n---\n")
	}
	in.reply(buf.String())
	return nil
}

func (b *Bot) IngressList(in Incoming) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("b.IngressList").BindError(&err)
		defer g.End()
	}

	cl, err := unversioned.NewInCluster()
	if err != nil {
		in.reply(":exclamation: failed to create a client: " + err.Error())
		return err
	}

	list, err := cl.Ingress(b.namespace).List(api.ListOptions{})
	if err != nil {
		in.reply(":exclamation: failed to list ingress")
		return err
	}

	if len(list.Items) <= 0 {
		in.reply(":exclamation: No ingress found")
		return
	}
	buf := bytes.Buffer{}
	buf.WriteString(":white_check_mark: Existing ingress objects:\n")
	for _, oldingress := range list.Items {
		buf.WriteString("    ")
		buf.WriteString(oldingress.ObjectMeta.Name)
		buf.WriteByte('\n')
	}
	in.reply(buf.String())
	return nil
}

func (b *Bot) IngressCreate(in Incoming) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("b.IngressCreate %s", in.Name).BindError(&err)
		defer g.End()
	}

	// Look for templates
	tmplf := filepath.Join(b.root, "ingress", in.Name+".yaml.tmpl")
	in.reply(":white_check_mark: Looking for template file...")
	if _, err := os.Stat(tmplf); err != nil {
		in.reply(":exclamation: template file '" + tmplf + "' not found")
		return nil
	}

	t, err := template.New("t").ParseFiles(tmplf)
	if err != nil {
		in.reply(":exclamation: failed to parse template '" + tmplf + "': " + err.Error())
		return err
	}

	outf, err := ioutil.TempFile("", "octav-ingress-k8s-")
	if err != nil {
		in.reply(":exclamation: ingress template for '" + in.Name + "' not found")
		return err
	}

	if in.Args == nil {
		in.Args = map[string]string{}
	}
	in.Args["timestamp"] = time.Now().Format("20060102-150405")
	if err := t.ExecuteTemplate(outf, in.Name+".yaml.tmpl", in.Args); err != nil {
		in.reply(":exclamation: failed to execute template: " + err.Error())
		return err
	}
	outf.Sync()
	outf.Seek(0, 0)

	var ingress extensions.Ingress
	if err := yaml.NewYAMLToJSONDecoder(outf).Decode(&ingress); err != nil {
		in.reply(":exclamation: failed to decode ingress template for '" + in.Name + "': " + err.Error())
		return err
	}

	cl, err := unversioned.NewInCluster()
	if err != nil {
		in.reply(":exclamation: failed to create a client: " + err.Error())
		return err
	}

	label, err := labels.Parse("name=" + in.Name)
	if err != nil {
		in.reply(":exclamation: failed to parse label 'name=" + in.Name + "': " + err.Error())
		return err
	}

	in.reply(":white_check_mark: Looking for existing ingress resources...")
	// Find the previous Ingress instance that was running for
	// this resource, so we can tell the user to shut it down later
	list, err := cl.Ingress(b.namespace).List(api.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		in.reply(":exclamation: failed to list ingress")
		return err
	}

	if len(list.Items) > 0 {
		buf := bytes.Buffer{}
		buf.WriteString(":white_check_mark: Existing ingress objects:\n")
		for _, oldingress := range list.Items {
			buf.WriteString("    ")
			buf.WriteString(oldingress.ObjectMeta.Name)
			buf.WriteByte('\n')
		}
		in.reply(buf.String())
	}

	in.reply(":white_check_mark: Creating new ingress...")
	newingress, err := cl.Ingress(b.namespace).Create(&ingress)
	if err != nil {
		in.reply(":exclamation: failed to create a new ingress")
		return err
	}

	newname := newingress.ObjectMeta.Name
	in.reply(":white_check_mark: Request to create sent, waiting to get an IP address...")
	// Wait until we have an ip address
	gotIP := make(chan error)
	go func() {
		tick := time.Tick(5 * time.Second)
		timeout := time.After(5 * time.Minute)

		for {
			select {
			case <-tick:
				// fallthrough to end of select
			case <-timeout:
				// uh-oh
				gotIP <- errors.New("timed out while waiting for new ingress")
				return
			}

			newingress, err = cl.Ingress(b.namespace).Get(newname)
			if err != nil {
				// ignoring...
				continue
			}

			if len(newingress.Status.LoadBalancer.Ingress) > 0 {
				close(gotIP)
				return
			}
		}
	}()

	err = <-gotIP
	if err != nil {
		in.reply(":exclamation: " + err.Error())
		return err
	}

	in.reply(":white_check_mark: Updating DNS records...")
	// dns client doesn't exist (not now, at least), so we access
	// the api directly. note that latest google.golang.org/cloud is
	// incompatible with k8s, but google.golang.org/api is OK

	// TODO: handle paging
	rrslist, err := b.dns.ResourceRecordSets.List(b.projectID, b.zone).Name(in.Name + ".").Type("A").Do()
	if err != nil {
		in.reply(":exclamation: failed to list changes for '" + in.Name + ".': " + err.Error())
		return err
	}

	// Should be one
	oldips := make([]string, len(rrslist.Rrsets[0].Rrdatas))
	newips := make([]string, len(rrslist.Rrsets[0].Rrdatas)+1)
	for i, rrd := range rrslist.Rrsets[0].Rrdatas {
		oldips[i] = rrd
		newips[i] = rrd
	}
	newips[len(newips)-1] = newingress.Status.LoadBalancer.Ingress[0].IP

	ch := dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Kind:    "dns#resourceRecordSet",
				Name:    in.Name + ".", // need to terminate with a "."
				Rrdatas: oldips,
				Ttl:     60,
				Type:    "A",
			},
		},
		Additions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Kind:    "dns#resourceRecordSet",
				Name:    in.Name + ".", // need to terminate with a "."
				Rrdatas: newips,
				Ttl:     60,
				Type:    "A",
			},
		},
	}

	if _, err := b.dns.Changes.Create(b.projectID, b.zone, &ch).Do(); err != nil {
		in.reply(":exclamation: failed to change DNS entries for '" + in.Name + "'")
		return err
	}

	// It is just way risky to let an automat-piloted program to
	// disable access to previous versions of the ingress, so
	// we let the user know, and let them do it manually

	in.reply(":tada: deploying ingress complete: you must shutdown the previous ingresses manually")
	for _, oldingress := range list.Items {
		in.reply(":white_check_mark: <botname> ingress delete '" + oldingress.ObjectMeta.Name + "'")
	}

	return nil
}