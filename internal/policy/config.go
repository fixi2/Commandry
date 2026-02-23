package policy

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Denylist          []string
	RedactionKeywords []string
	EnforceDenylist   bool
}

func ParseConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read policy config: %w", err)
	}
	return ParseConfig(string(data))
}

func ParseConfig(content string) (Config, error) {
	content = strings.TrimPrefix(content, "\ufeff")
	cfg := Config{
		Denylist:          append([]string(nil), defaultDenylistPatterns...),
		RedactionKeywords: append([]string(nil), defaultRedactionKeywords...),
		EnforceDenylist:   false,
	}

	var (
		inPolicy      bool
		currentList   string
		denylistSeen  bool
		keywordsSeen  bool
		parsedDeny    []string
		parsedKeyword []string
	)

	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	for idx, raw := range lines {
		raw = strings.TrimPrefix(raw, "\ufeff")
		line := strings.TrimRight(raw, " \t")
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}

		if !strings.HasPrefix(line, " ") {
			currentList = ""
			inPolicy = strings.HasPrefix(trim, "policy:")
			continue
		}
		if !inPolicy {
			continue
		}

		if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
			currentList = ""
			key, value, hasValue := splitKeyValue(strings.TrimSpace(line))
			switch key {
			case "denylist":
				denylistSeen = true
				parsedDeny = parsedDeny[:0]
				currentList = "denylist"
			case "redaction_keywords":
				keywordsSeen = true
				parsedKeyword = parsedKeyword[:0]
				currentList = "redaction_keywords"
			case "enforce_denylist":
				if !hasValue {
					return Config{}, fmt.Errorf("parse policy config line %d: enforce_denylist requires a boolean value", idx+1)
				}
				value = strings.ToLower(strings.TrimSpace(value))
				switch value {
				case "true":
					cfg.EnforceDenylist = true
				case "false":
					cfg.EnforceDenylist = false
				default:
					return Config{}, fmt.Errorf("parse policy config line %d: enforce_denylist must be true or false", idx+1)
				}
			}
			continue
		}

		if strings.HasPrefix(line, "    - ") && currentList != "" {
			item := strings.TrimSpace(strings.TrimPrefix(line, "    - "))
			item = trimMatchingQuotes(item)
			if item == "" {
				continue
			}
			switch currentList {
			case "denylist":
				parsedDeny = append(parsedDeny, item)
			case "redaction_keywords":
				parsedKeyword = append(parsedKeyword, item)
			}
		}
	}

	if denylistSeen && len(parsedDeny) > 0 {
		cfg.Denylist = parsedDeny
	}
	if keywordsSeen && len(parsedKeyword) > 0 {
		cfg.RedactionKeywords = parsedKeyword
	}

	return cfg, nil
}

func splitKeyValue(line string) (key string, value string, hasValue bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return strings.TrimSpace(line), "", false
	}
	key = strings.TrimSpace(parts[0])
	value = strings.TrimSpace(parts[1])
	if value == "" {
		return key, "", false
	}
	return key, trimMatchingQuotes(value), true
}

func trimMatchingQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
