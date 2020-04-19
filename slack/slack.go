package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hb9tf/wiresx2influx/wiresx"
)

const (
	httpPOST        = "POST"
	httpContentType = "Content-Type"
	httpJSON        = "application/json"

	slackColorGood = "good"

	// timePostFormat is the date/time format presented in the Slack post.
	timePostFormat = "2006-01-02 15:04:05"
)

type Message struct {
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color    string `json:"color,omitempty"`
	Fallback string `json:"fallback"`

	CallbackID string `json:"callback_id,omitempty"`
	ID         int    `json:"id,omitempty"`

	AuthorID      string `json:"author_id,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	AuthorLink    string `json:"author_link,omitempty"`
	AuthorIcon    string `json:"author_icon,omitempty"`

	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Text      string `json:"text"`

	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`

	//Fields     []AttachmentField  `json:"fields,omitempty"`
	//Actions    []AttachmentAction `json:"actions,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`

	Footer     string `json:"footer,omitempty"`
	FooterIcon string `json:"footer_icon,omitempty"`

	Ts json.Number `json:"ts,omitempty"`
}

// Slacker is a super simple Slack bot which allows to post messages using a webhook.
type Slacker struct {
	Webhook string
	Client  *http.Client
}

// Post sends the provided message to the webhook, posting it in the channel.
func (s *Slacker) Post(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(httpPOST, s.Webhook, bytes.NewBuffer(data))
	req.Header.Set(httpContentType, httpJSON)
	_, err = s.Client.Do(req)
	return err
}

func getSlackMsg(log *wiresx.Log) *Message {
	loc := "unknown"
	if log.Loc != nil {
		loc = log.Loc.String()
		if log.Loc.Latitude != 0 && log.Loc.Longitude != 0 {
			loc = fmt.Sprintf("<https://www.google.com/maps/place/%f+%f|%s>", log.Loc.Latitude, log.Loc.Longitude, loc)
		}
	}
	return &Message{
		Attachments: []Attachment{
			{
				Pretext: fmt.Sprintf("%s (%s, %s): Location %s", log.Callsign, log.Dev.InferDevice(), log.Source, loc),
				Ts:      json.Number(strconv.FormatInt(log.Timestamp.Unix(), 10)),
			},
		},
	}
}

func Feed(ctx context.Context, logChan chan *wiresx.Log, client *Slacker, tags map[string]string, dry bool) {
	for l := range logChan {
		msg := getSlackMsg(l)
		if dry {
			log.Printf("DRY: Sending message to Slack: %+v", msg)
			continue
		}
		if err := client.Post(msg); err != nil {
			log.Printf("Error posting to slack: %s", err)
		}
	}
}
