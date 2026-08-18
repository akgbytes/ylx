package email

import (
	"fmt"
	"html"
)

const (
	brandBlue    = "#0866FF"
	brandInk     = "#18181B"
	brandCanvas  = "#F4F7FB"
	brandBorder  = "#E3E8EF"
	brandTagline = "The Marketplace for Developers"
)

type EmailTemplate struct {
	Subject string
	Text    string
	HTML    string
}

type emailContent struct {
	preheader   string
	heading     string
	intro       string
	action      string
	securityTip string
}

// SignUpOTPTemplate creates the verification email sent during sign-up.
func SignUpOTPTemplate(otp string) EmailTemplate {
	escapedOTP := html.EscapeString(otp)

	return EmailTemplate{
		Subject: "Verify your email for YLX",
		Text: fmt.Sprintf(`Welcome to YLX — The Marketplace for Developers.

Your verification code is: %s

Enter this code to finish creating your account. Do not share it with anyone.

If you did not create a YLX account, you can safely ignore this email.`, otp),
		HTML: renderEmail(emailContent{
			preheader: "Use your verification code to finish creating your YLX account.",
			heading:   "Verify your email",
			intro:     "Enter this code to finish creating your developer marketplace account.",
			action: fmt.Sprintf(`
          <div style="margin:28px 0; padding:22px 16px; background:#F4F7FB; border:1px solid #E3E8EF; border-radius:12px; text-align:center;">
            <p style="margin:0; color:%s; font-family:'Courier New', Courier, monospace; font-size:34px; font-weight:700; letter-spacing:8px; line-height:1.2;">%s</p>
          </div>`, brandInk, escapedOTP),
			securityTip: "This code expires shortly. Never share it with anyone, including someone claiming to work for YLX.",
		}),
	}
}

func PasswordResetTemplate(link string) EmailTemplate {
	escapedLink := html.EscapeString(link)

	return EmailTemplate{
		Subject: "Reset your YLX password",
		Text: fmt.Sprintf(`We received a request to reset your YLX password.

Reset your password using this link: %s

If you did not request a password reset, you can safely ignore this email. Your password will not change.`, link),
		HTML: renderEmail(emailContent{
			preheader: "Reset your YLX password securely.",
			heading:   "Reset your password",
			intro:     "We received a request to reset your password. Use the secure link below to choose a new one.",
			action: fmt.Sprintf(`
          <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:28px auto;">
            <tr>
              <td align="center" bgcolor="%s" style="border-radius:10px;">
                <a href="%s" target="_blank" rel="noopener noreferrer" style="display:inline-block; padding:14px 26px; color:#FFFFFF; font-size:15px; font-weight:700; line-height:20px; text-decoration:none;">Reset password</a>
              </td>
            </tr>
          </table>
          <p style="margin:0; color:#667085; font-size:12px; line-height:18px; text-align:center; word-break:break-all;">If the button does not work, copy and paste this link:<br /><a href="%s" style="color:%s; text-decoration:underline;">%s</a></p>`, brandBlue, escapedLink, escapedLink, brandBlue, escapedLink),
			securityTip: "If you did not request this reset, no action is needed and your password will remain unchanged.",
		}),
	}
}

func renderEmail(content emailContent) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta name="color-scheme" content="light" />
  <title>%s</title>
</head>
<body style="margin:0; padding:0; background:%s; color:%s; font-family:Arial, Helvetica, sans-serif;">
  <div style="display:none; max-height:0; overflow:hidden; opacity:0;">%s</div>
  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" bgcolor="%s" style="width:100%%; background:%s; padding:32px 12px;">
    <tr>
      <td align="center">
        <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="width:100%%; max-width:600px; overflow:hidden; background:#FFFFFF; border:1px solid %s; border-radius:16px;">
          <tr>
            <td style="padding:26px 32px; background:%s;">
              <table role="presentation" cellspacing="0" cellpadding="0" border="0">
                <tr>
                  <td style="padding-right:12px; vertical-align:middle;">
                    <img src="cid:%s" width="48" height="48" alt="YLX panda" style="display:block; width:48px; height:48px; border:0; border-radius:50%%;" />
                  </td>
                  <td style="vertical-align:middle;">
                    <p style="margin:0; color:#FFFFFF; font-size:28px; font-weight:800; letter-spacing:-1px; line-height:32px;">YLX</p>
                    <p style="margin:3px 0 0; color:#EAF1FF; font-size:12px; font-weight:600; line-height:18px;">%s</p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:36px 32px 32px;">
              <h1 style="margin:0 0 14px; color:%s; font-size:26px; font-weight:750; letter-spacing:-0.5px; line-height:34px;">%s</h1>
              <p style="margin:0; color:#475467; font-size:15px; line-height:24px;">%s</p>
              %s
              <div style="margin-top:28px; padding:14px 16px; background:#FFF9EB; border-radius:10px;">
                <p style="margin:0; color:#694D12; font-size:12px; line-height:19px;"><strong>Security note:</strong> %s</p>
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 32px; background:#FAFBFC; border-top:1px solid %s;">
              <p style="margin:0; color:#667085; font-size:12px; line-height:18px;">YLX · %s</p>
              <p style="margin:4px 0 0; color:#98A2B3; font-size:11px; line-height:17px;">This is an automated account-security email.</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		html.EscapeString(content.heading),
		brandCanvas,
		brandInk,
		html.EscapeString(content.preheader),
		brandCanvas,
		brandCanvas,
		brandBorder,
		brandBlue,
		brandLogoContentID,
		brandTagline,
		brandInk,
		html.EscapeString(content.heading),
		html.EscapeString(content.intro),
		content.action,
		html.EscapeString(content.securityTip),
		brandBorder,
		brandTagline,
	)
}
