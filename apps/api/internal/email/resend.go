package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"

	_ "embed"
)

const brandLogoContentID = "ylx-panda"

//go:embed assets/ylx-panda.png
var brandLogo []byte

type Sender interface {
	SendSignUpOTP(ctx context.Context, recipient, otp string) error
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

func (s *ResendSender) SendSignUpOTP(ctx context.Context, recipient, otp string) error {
	template := SignUpOTPTemplate(otp)

	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{recipient},
		Subject: template.Subject,
		Text:    template.Text,
		Html:    template.HTML,
		Attachments: []*resend.Attachment{
			{
				Content:     brandLogo,
				Filename:    "ylx-panda.png",
				ContentType: "image/png",
				ContentId:   brandLogoContentID,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send signup OTP email: %w", err)
	}

	return nil
}
