package email_test

import (
	"strings"
	"testing"

	"github.com/akgbytes/ylx/internal/email"
)

func TestSignUpOTPTemplate(t *testing.T) {
	t.Parallel()

	template := email.SignUpOTPTemplate("482913")

	assertContains(t, template.Subject, "YLX")
	assertContains(t, template.Text, "482913")
	assertContains(t, template.HTML, "482913")
	assertContains(t, template.HTML, "The Marketplace for Developers")
	assertContains(t, template.HTML, "#0866FF")
	assertContains(t, template.HTML, "cid:ylx-panda")
	assertNotContains(t, template.HTML, "SecureAuth")
}

func TestSignUpOTPTemplateEscapesHTML(t *testing.T) {
	t.Parallel()

	template := email.SignUpOTPTemplate(`<script>alert("oops")</script>`)

	assertNotContains(t, template.HTML, "<script>")
	assertContains(t, template.HTML, "&lt;script&gt;")
}

func TestPasswordResetTemplate(t *testing.T) {
	t.Parallel()

	link := `https://ylx.example/reset?token=a&next="/account"`
	template := email.PasswordResetTemplate(link)

	assertContains(t, template.Subject, "YLX")
	assertContains(t, template.Text, link)
	assertContains(t, template.HTML, "Reset password")
	assertContains(t, template.HTML, "token=a&amp;next=&#34;/account&#34;")
	assertNotContains(t, template.HTML, `token=a&next="/account"`)
}

func assertContains(t *testing.T, value, want string) {
	t.Helper()

	if !strings.Contains(value, want) {
		t.Fatalf("expected %q to contain %q", value, want)
	}
}

func assertNotContains(t *testing.T, value, unwanted string) {
	t.Helper()

	if strings.Contains(value, unwanted) {
		t.Fatalf("expected %q not to contain %q", value, unwanted)
	}
}
