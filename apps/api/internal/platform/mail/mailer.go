package mail

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"

	_ "embed"
)

type Template struct {
	Subject string
	Text    string
	HTML    string
}

// Used to reference logo directly in html using CID
const LogoContentID = "ylx-panda"

//go:embed assets/logo.png
var brandLogo []byte

type Sender interface {
	Send(ctx context.Context, recipient string, template Template) error
}

type ResendSender struct {
	client *resend.Client
	from   string
}

func NewResendSender(apiKey, from string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (s *ResendSender) Send(ctx context.Context, recipient string, template Template) error {
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{recipient},
		Subject: template.Subject,
		Text:    template.Text,
		Html:    template.HTML,
		Attachments: []*resend.Attachment{{
			Content:     brandLogo,
			Filename:    "logo.png",
			ContentType: "image/png",
			ContentId:   LogoContentID,
		}},
	})
	if err != nil {
		return fmt.Errorf("send email to recipient: %w", err)
	}

	return nil
}
