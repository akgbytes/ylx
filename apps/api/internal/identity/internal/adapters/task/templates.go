package task

import (
	"fmt"

	"github.com/akgbytes/ylx/internal/platform/mailer"
)

func signupOTPTemplate(otp, recipient string) mailer.Template {
	return mailer.Template{
		Subject: "Verify your email for YLX",
		Text: fmt.Sprintf(`Welcome to YLX — The Marketplace for Developers.

Your verification code is: %s

Enter this code to finish creating your account. Do not share it with anyone.

If you did not create a YLX account, you can safely ignore this email.`, otp),
		HTML: mailer.SignupEmailTemplate(
			otp,
			recipient,
		),
	}
}
