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

const k8sNamespace = "default" // TODO make configurable

type watchctx struct {
	known map[string]struct{}
}

func StartWatch(done chan struct{}) (err error) {
	defer close(done)
	if pdebug.Enabled {
		g := pdebug.Marker("k8s.Watch").BindError(&err)
		defer g.End()
	}

	k8sClient, err := k8sc.NewInCluster()
	if err != nil {
		return err
	}

	label, err := labels.Parse("group=octav")
	if err != nil {
		return err
	}

	// First, get the list of all known pods, so to not to flood
	// the notification.
	podlist, err := k8sClient.Pods(k8sNamespace).List(api.ListOptions{})
	if err != nil {
		return err
	}
	ctx := watchctx{}
	ctx.known = map[string]struct{}{}
	for _, pod := range podlist.Items {
		ctx.known[pod.Name] = struct{}{}
	}

	w, err := k8sClient.Pods(k8sNamespace).Watch(api.ListOptions{
		LabelSelector: label,
	})
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
				if err := notifyPod(&ctx, e, e.Object.(*api.Pod)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func notifyPod(ctx *watchctx, e watch.Event, pod *api.Pod) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("k8s.notifyPod").BindError(&err)
		defer g.End()
	}

	switch e.Type {
	case watch.Added:
		if _, ok := ctx.known[pod.Name]; ok {
			if pdebug.Enabled {
				pdebug.Printf("ignoring pod %s", pod.Name)
			}
			// If this is a known pod, then we probably
			// don't want to report its creation
			delete(ctx.known, pod.Name)
			return nil
		}

		msgbuf := bytes.Buffer{}
		msgbuf.WriteString("Pod started: name=")
		msgbuf.WriteString(pod.Name)
		params := slack.NewPostMessageParameters()
		params.Username = "GKE Status"
		params.Markdown = true
		params.Attachments = append(params.Attachments, slack.Attachment{
			Color:    "good",
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
			Color:    "warning",
			Fallback: msgbuf.String(),
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Pod Name",
					Value: pod.Name,
					Short: true,
				},
			},
			Title:    "Pod Deleted",
			ThumbURL: "https://pbs.twimg.com/media/Bt_pEfqCAAAiVyz.png",
		})

		_, _, err := slackClient.PostMessage("#gcp-status", "*GKE Status Changed*", params)
		if err != nil {
			return err
		}
	default:
		if pdebug.Enabled {
			pdebug.Printf("Unknown event")
			pdebug.Printf("%#v", e)
			pdebug.Printf("%#v", pod)
		}
	}

	return nil
}
