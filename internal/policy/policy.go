package policy

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	RedactedValue     = "[REDACTED]"
	DeniedPlaceholder = "[REDACTED BY POLICY]"
)

type Result struct {
	Command string
	Denied  bool
}

type redactor struct {
	re   *regexp.Regexp
	repl string
}

type Policy struct {
	denylist []*regexp.Regexp
	redact   []redactor
}

var credentialInImageRef = regexp.MustCompile(`^[^/\s:@]+:[^/\s@]+@`)

func NewDefault() *Policy {
	return &Policy{
		denylist: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\bcat\s+~\/\.ssh\/`),
			regexp.MustCompile(`(?i)\bid_rsa\b`),
			regexp.MustCompile(`(?i)\.(pem|key)(\s|$)`),
			regexp.MustCompile(`(?i)\bkubectl\s+get\s+secret\b.*\s-o\s+(yaml|json)\b`),
			regexp.MustCompile(`(?i)\bgcloud\s+auth\s+print-access-token\b`),
		},
		redact: []redactor{
			{
				re:   regexp.MustCompile(`(?i)(authorization\s*:\s*bearer\s+)([^\s"']+)`),
				repl: `${1}` + RedactedValue,
			},
			{
				re:   regexp.MustCompile(`(?i)(--(?:token|password|passwd|api[_-]?key|apikey|secret|private[_-]?key)=)([^\s]+)`),
				repl: `${1}` + RedactedValue,
			},
			{
				re:   regexp.MustCompile(`(?i)(--(?:token|password|passwd|api[_-]?key|apikey|secret|private[_-]?key)\s+)([^\s]+)`),
				repl: `${1}` + RedactedValue,
			},
			{
				re:   regexp.MustCompile(`(?i)(-p\s+)([^\s]+)`),
				repl: `${1}` + RedactedValue,
			},
			{
				re:   regexp.MustCompile(`(?i)(\b(?:token|secret|password|passwd|api[_-]?key|apikey|private[_-]?key)\b\s*[:=]\s*)([^\s]+)`),
				repl: `${1}` + RedactedValue,
			},
			{
				re:   regexp.MustCompile(`(?i)(\b[A-Za-z_][A-Za-z0-9_]*=)"[^"]*"`),
				repl: `${1}"` + RedactedValue + `"`,
			},
			{
				re:   regexp.MustCompile(`(?i)(\b[A-Za-z_][A-Za-z0-9_]*=)'[^']*'`),
				repl: `${1}'` + RedactedValue + `'`,
			},
			{
				re:   regexp.MustCompile(`(?i)(\b[A-Za-z_][A-Za-z0-9_]*=)([^\s"']+)`),
				repl: `${1}` + RedactedValue,
			},
		},
	}
}

func (p *Policy) Apply(rawCommand string, args []string) Result {
	if p.isDenied(rawCommand, args) {
		return Result{
			Command: DeniedPlaceholder,
			Denied:  true,
		}
	}

	sanitized, preserved := preserveKubectlSetImageAssignments(rawCommand, args)
	for _, rule := range p.redact {
		sanitized = rule.re.ReplaceAllString(sanitized, rule.repl)
	}
	for placeholder, original := range preserved {
		sanitized = strings.ReplaceAll(sanitized, placeholder, original)
	}

	return Result{
		Command: sanitized,
		Denied:  false,
	}
}

func preserveKubectlSetImageAssignments(rawCommand string, args []string) (string, map[string]string) {
	if !isKubectlSetImage(args) {
		return rawCommand, nil
	}

	sanitized := rawCommand
	preserved := make(map[string]string)
	index := 0

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") || !strings.Contains(arg, "=") {
			continue
		}
		if !isSafeImageAssignment(arg) {
			continue
		}

		placeholder := "__INFRATRACK_IMG_ASSIGN_" + strconv.Itoa(index) + "__"
		index++
		sanitized = strings.ReplaceAll(sanitized, arg, placeholder)
		preserved[placeholder] = arg
	}

	return sanitized, preserved
}

func isKubectlSetImage(args []string) bool {
	if len(args) < 3 {
		return false
	}

	binary := strings.ToLower(filepath.Base(args[0]))
	if binary != "kubectl" && binary != "kubectl.exe" {
		return false
	}

	return strings.EqualFold(args[1], "set") && strings.EqualFold(args[2], "image")
}

func isSafeImageAssignment(arg string) bool {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return false
	}

	key := strings.ToLower(parts[0])
	value := strings.ToLower(parts[1])
	if key == "" || value == "" {
		return false
	}

	if containsSensitiveKeyword(key) || containsSensitiveKeyword(value) {
		return false
	}

	return !credentialInImageRef.MatchString(parts[1])
}

func containsSensitiveKeyword(v string) bool {
	keywords := []string{
		"token", "secret", "password", "passwd", "authorization", "bearer", "api_key", "apikey", "private_key",
	}
	for _, keyword := range keywords {
		if strings.Contains(v, keyword) {
			return true
		}
	}
	return false
}

func (p *Policy) isDenied(rawCommand string, args []string) bool {
	if len(args) > 0 {
		binary := strings.ToLower(filepath.Base(args[0]))
		if binary == "env" || binary == "printenv" {
			return true
		}
	}

	for _, rule := range p.denylist {
		if rule.MatchString(rawCommand) {
			return true
		}
	}

	return false
}
