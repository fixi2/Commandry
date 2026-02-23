package textblock

import (
	"errors"
	"strings"
)

var ErrMalformedMarkers = errors.New("hook block markers are malformed")

type Span struct {
	Start int
	End   int
}

func FindSingle(content, begin, end string) (Span, bool, error) {
	start := strings.Index(content, begin)
	if start < 0 {
		if strings.Contains(content, end) {
			return Span{}, false, ErrMalformedMarkers
		}
		return Span{}, false, nil
	}

	rest := content[start+len(begin):]
	relEnd := strings.Index(rest, end)
	if relEnd < 0 {
		return Span{}, false, ErrMalformedMarkers
	}

	if strings.Contains(rest[:relEnd], begin) || strings.Contains(rest[relEnd+len(end):], begin) {
		return Span{}, false, ErrMalformedMarkers
	}
	if strings.Contains(rest[relEnd+len(end):], end) {
		return Span{}, false, ErrMalformedMarkers
	}

	finish := start + len(begin) + relEnd + len(end)
	for finish < len(content) && (content[finish] == '\r' || content[finish] == '\n') {
		finish++
	}
	return Span{Start: start, End: finish}, true, nil
}

func Upsert(content, begin, end, block string) (string, bool, error) {
	span, exists, err := FindSingle(content, begin, end)
	if err != nil {
		return "", false, err
	}
	if exists {
		existing := content[span.Start:span.End]
		if normalizeComparable(existing) == normalizeComparable(block) {
			return content, false, nil
		}
		return content[:span.Start] + block + content[span.End:], true, nil
	}

	if content == "" {
		return block + "\n", true, nil
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + "\n" + block + "\n", true, nil
}

func Remove(content, begin, end string) (string, bool, error) {
	span, exists, err := FindSingle(content, begin, end)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return content, false, nil
	}

	left := content[:span.Start]
	right := content[span.End:]
	updated := strings.TrimRight(left, "\r\n")
	if right != "" {
		if updated != "" {
			updated += "\n"
		}
		updated += strings.TrimLeft(right, "\r\n")
	}
	return updated, true, nil
}

func normalizeLineEndings(v string) string {
	return strings.ReplaceAll(v, "\r\n", "\n")
}

func normalizeComparable(v string) string {
	return strings.TrimRight(normalizeLineEndings(v), "\n")
}
