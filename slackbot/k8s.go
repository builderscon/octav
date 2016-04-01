package slackbot

import (
	"bytes"

	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
	"k8s.io/kubernetes/pkg/api"
	k8sc "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func Watch() (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("k8s.Watch").BindError(&err)
		defer g.End()
	}

	k8sConfig, err := k8sc.InClusterConfig()
	if err != nil {
		return err
	}
	k8sClient, err := k8sc.New(k8sConfig)
	if err != nil {
		return err
	}

	label, err := labels.Parse("group=octav")
	if err != nil {
		return err
	}

	w, err := k8sClient.Pods("default").Watch(label, nil, "")
	if err != nil {
		return err
	}

	c := w.ResultChan()
	for {
		select {
		case e, ok := <-c:
			if !ok {
				break
			}
			switch e.Object.(type) {
			case *api.Pod:
				if err := notifyPod(e, e.Object.(*api.Pod)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func notifyPod(e watch.Event, pod *api.Pod) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("k8s.notifyPod").BindError(&err)
		defer g.End()
	}

	switch e.Type {
	case watch.Added:
		msgbuf := bytes.Buffer{}
		msgbuf.WriteString("Pod started: name=")
		msgbuf.WriteString(pod.Name)
		params := slack.NewPostMessageParameters()
		params.Username = "GKE Status"
		params.Markdown = true
		params.Attachments = append(params.Attachments, slack.Attachment{
			Color: "good",
			Fallback: msgbuf.String(),
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Pod Name",
					Value: pod.Name,
					Short: true,
				},
			},
			Title:    "New Pod Added",
			ThumbURL: "https://pbs.twimg.com/media/Bt_pEfqCAAAiVyz.png",
		})

		_, _, err := slackClient.PostMessage("#gcp-status", "*GKE Status Changed*", params)
		if err != nil {
			return err
		}
	case watch.Deleted:
		msgbuf := bytes.Buffer{}
		msgbuf.WriteString("Pod deleted: name=")
		msgbuf.WriteString(pod.Name)
		params := slack.NewPostMessageParameters()
		params.Username = "GKE Status"
		params.Markdown = true
		params.Attachments = append(params.Attachments, slack.Attachment{
			Color: "warning",
			Fallback: msgbuf.String(),
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Pod Name",
					Value: pod.Name,
					Short: true,
				},
			},
			Title: "Pod Deleted",
      ThumbURL: "https://pbs.twimg.com/media/Bt_pEfqCAAAiVyz.png",
		})

		_, _, err := slackClient.PostMessage("#gcp-status", "*GKE Status Changed*", params)
		if err != nil {
			return err
		}
	default:
		pdebug.Printf("Unknown event")
		pdebug.Printf("%#v", e)
		pdebug.Printf("%#v", pod)
	}

	return nil
}
