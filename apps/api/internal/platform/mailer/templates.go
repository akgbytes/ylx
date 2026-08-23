package mailer

import (
	"fmt"
	"html"
)

const (
	BrandBlue    = "#0866FF"
	BrandInk     = "#18181B"
	BrandCanvas  = "#F4F7FB"
	BrandBorder  = "#E3E8EF"
	brandTagline = "The Marketplace for Developers"
)

func SignupEmailTemplate(
	otp string,
	recipientName string,
) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="color-scheme" content="light" />
    <title>Your verification code</title>
  </head>

  <body
    style="
      margin: 0;
      padding: 0;
      background: %s;
      color: %s;
      font-family: Arial, Helvetica, sans-serif;
    "
  >
    <!-- Preheader -->
    <div style="display: none; max-height: 0; overflow: hidden; opacity: 0">
      Your YLX verification code is %s.
    </div>

    <table
      role="presentation"
      width="100%%"
      cellspacing="0"
      cellpadding="0"
      border="0"
      bgcolor="%s"
      style="width: 100%%; background: %s; padding: 32px 12px"
    >
      <tr>
        <td align="center">
          <!-- Main card -->
          <table
            role="presentation"
            width="100%%"
            cellspacing="0"
            cellpadding="0"
            border="0"
            style="
              width: 100%%;
              max-width: 600px;
              overflow: hidden;
              background: #ffffff;
              border: 1px solid %s;
              border-radius: 16px;
            "
          >
            <!-- Header -->
            <tr>
              <td style="padding: 26px 32px; background: %s">
                <table
                  role="presentation"
                  cellspacing="0"
                  cellpadding="0"
                  border="0"
                >
                  <tr>
                    <!-- Logo -->
                    <td style="padding-right: 12px; vertical-align: middle">
                      <img
                        src="cid:%s"
                        width="48"
                        height="48"
                        alt="YLX panda"
                        style="
                          display: block;
                          width: 48px;
                          height: 48px;
                          border: 0;
                          border-radius: 50%%;
                        "
                      />
                    </td>

                    <!-- Brand -->
                    <td style="vertical-align: middle">
                      <p
                        style="
                          margin: 0;
                          color: #ffffff;
                          font-size: 28px;
                          font-weight: 800;
                          letter-spacing: -1px;
                          line-height: 32px;
                        "
                      >
                        YLX
                      </p>

                      <p
                        style="
                          margin: 3px 0 0;
                          color: #eaf1ff;
                          font-size: 12px;
                          font-weight: 600;
                          line-height: 18px;
                        "
                      >
                        %s
                      </p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Content -->
            <tr>
              <td style="padding: 36px 32px 32px">
                <!-- Heading -->
                <h1
                  style="
                    margin: 0 0 14px;
                    color: %s;
                    font-size: 26px;
                    font-weight: 700;
                    letter-spacing: -0.5px;
                    line-height: 34px;
                  "
                >
                  Your verification code
                </h1>

                <!-- Recipient -->
                <p
                  style="
                    margin: 0 0 18px;
                    font-size: 15px;
                    line-height: 24px;
                  "
                >
                  Hi, %s
                </p>

                <!-- Intro -->
                <p
                  style="
                    margin: 0 0 28px;
                    font-size: 15px;
                    line-height: 24px;
                  "
                >
                  Thanks for signing up for YLX! Please use the 6-digit
                  verification code below to verify your email address and
                  complete your sign up for YLX.
                </p>

                <!-- OTP -->
                <p
                  style="
                    margin: 24px 0 28px;
                    color: %s;
                    font-size: 32px;
                    font-weight: 700;
                    letter-spacing: 8px;
                    line-height: 40px;
                    text-align: center;
                  "
                >
                  %s
                </p>

                <!-- Expiry -->
                <p
                  style="
                    margin: 0 0 18px;
                    font-size: 15px;
                    line-height: 24px;
                  "
                >
                  This code will expire in
                  <strong>10 minutes.</strong>
                </p>

                <!-- Security message -->
                <p
                  style="
                    margin: 0 0 28px;
                    font-size: 15px;
                    line-height: 24px;
                  "
                >
                  If you didn't request this code, you can safely ignore this
                  email.
                </p>

                <!-- Sign off -->
                <p
                  style="
                    margin: 0;
                    font-size: 15px;
                    line-height: 24px;
                  "
                >
                  Thanks,<br />
                  The YLX Team
                </p>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td
                style="
                  padding: 20px 32px;
                  background: #fafbfc;
                  border-top: 1px solid %s;
                  text-align: center;
                "
              >
                <p
                  style="
                    margin: 0;
                    color: #a2aaba;
                    font-size: 12px;
                    line-height: 18px;
                  "
                >
                  YLX · %s
                </p>

                <p
                  style="
                    margin: 4px 0 0;
                    color: #a2aaba;
                    font-size: 11px;
                    line-height: 17px;
                  "
                >
                  This is an automated account-security email.
                </p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`,
		BrandCanvas,                      // body background
		BrandInk,                         // body color
		html.EscapeString(otp),           // preheader OTP
		BrandCanvas,                      // table bgcolor
		BrandCanvas,                      // table background
		BrandBorder,                      // card border
		BrandBlue,                        // header background
		LogoContentID,                    // logo CID
		brandTagline,                     // tagline
		BrandInk,                         // heading color
		html.EscapeString(recipientName), // recipient
		BrandInk,                         // OTP color
		html.EscapeString(otp),           // OTP
		BrandBorder,                      // footer border
		brandTagline,                     // footer tagline
	)
}
