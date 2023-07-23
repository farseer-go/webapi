package context

import (
	"bytes"
	"fmt"
	"github.com/farseer-go/fs/flog"
	"regexp"
	"strconv"
	"strings"
)

type RouteRegexpOptions struct {
	StrictSlash    bool
	UseEncodedPath bool
}

type RegexpType int

const (
	RegexpTypePath RegexpType = iota
	RegexpTypeHost
	RegexpTypePrefix
	RegexpTypeQuery
)

// 获取默认正则
func (receiver RegexpType) getDefaultPattern() string {
	if receiver == RegexpTypeQuery {
		return ".*"
	}

	if receiver == RegexpTypeHost {
		return "[^.]+"
	}

	// RegexpTypePath or RegexpTypePrefix
	return "[^/]+"
}

// routeRegexp stores a regexp to match a host or path and information to
// collect and validate route variables.
type routeRegexp struct {
	// The unmodified template.
	template string
	// The type of match
	regexpType RegexpType
	// Options for matching
	options RouteRegexpOptions
	// Expanded regexp.
	regexp *regexp.Regexp
	// Reverse template.
	reverse string
	// Variable names.
	varsN []string
	// Variable regexps (validators).
	varsR []*regexp.Regexp
	// Wildcard host-port (no strict port match in hostname)
	wildcardHostPort bool
}

// NewRouteRegexp 构建支持正则的路由
func NewRouteRegexp(tpl string, typ RegexpType, options RouteRegexpOptions) *routeRegexp {
	// 找出占位符位置
	placeholderIndexs, errBraces := getPlaceholderIndex(tpl)
	if errBraces != nil {
		flog.Panic(errBraces)
		return nil
	}

	// 复制原tpl
	template := tpl
	// 获取默认正则
	defaultPattern := typ.getDefaultPattern()
	// 如果不匹配，则只匹配 /
	if typ != RegexpTypePath {
		options.StrictSlash = false
	}

	// 末尾 /
	endSlash := false
	if options.StrictSlash && strings.HasSuffix(tpl, "/") {
		tpl = tpl[:len(tpl)-1]
		endSlash = true
	}
	// 匹配到的变量名称
	varsN := make([]string, len(placeholderIndexs)/2)
	// 匹配到的对应正则格式
	varsR := make([]*regexp.Regexp, len(placeholderIndexs)/2)
	// 整个Path的正则
	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')
	reverse := bytes.NewBufferString("")
	var end int
	var err error
	for i := 0; i < len(placeholderIndexs); i += 2 {
		// Set all values we are interested in.
		raw := tpl[end:placeholderIndexs[i]]
		end = placeholderIndexs[i+1]
		parts := strings.SplitN(tpl[placeholderIndexs[i]+1:end-1], ":", 2)
		name := parts[0]
		patt := defaultPattern
		if len(parts) == 2 {
			patt = parts[1]
		}
		// Name or pattern can't be empty.
		if name == "" || patt == "" {
			flog.Panicf("webapi: missing name or pattern in %q", tpl[placeholderIndexs[i]:end])
			return nil
		}
		// Build the regexp pattern.
		fmt.Fprintf(pattern, "%s(?P<%s>%s)", regexp.QuoteMeta(raw), varGroupName(i/2), patt)

		// Build the reverse template.
		fmt.Fprintf(reverse, "%s%%s", raw)

		// Append variable name and compiled pattern.
		varsN[i/2] = name
		varsR[i/2], err = regexp.Compile(fmt.Sprintf("^%s$", patt))
		if err != nil {
			flog.Panic(err)
			return nil
		}
	}
	// Add the remaining.
	raw := tpl[end:]
	pattern.WriteString(regexp.QuoteMeta(raw))
	if options.StrictSlash {
		pattern.WriteString("[/]?")
	}
	if typ == RegexpTypeQuery {
		// Add the default pattern if the query value is empty
		if queryVal := strings.SplitN(template, "=", 2)[1]; queryVal == "" {
			pattern.WriteString(defaultPattern)
		}
	}
	if typ != RegexpTypePrefix {
		pattern.WriteByte('$')
	}

	var wildcardHostPort bool
	if typ == RegexpTypeHost {
		if !strings.Contains(pattern.String(), ":") {
			wildcardHostPort = true
		}
	}
	reverse.WriteString(raw)
	if endSlash {
		reverse.WriteByte('/')
	}
	// Compile full regexp.
	reg, errCompile := regexp.Compile(pattern.String())
	if errCompile != nil {
		flog.Panic(errCompile)
		return nil
	}

	// Check for capturing groups which used to work in older versions
	if reg.NumSubexp() != len(placeholderIndexs)/2 {
		panic(fmt.Sprintf("route %s contains capture groups in its regexp. ", template) +
			"Only non-capturing groups are accepted: e.g. (?:pattern) instead of (pattern)")
	}

	// Done!
	return &routeRegexp{
		template:         template,
		regexpType:       typ,
		options:          options,
		regexp:           reg,
		reverse:          reverse.String(),
		varsN:            varsN,
		varsR:            varsR,
		wildcardHostPort: wildcardHostPort,
	}
}

// 从路由地址中找出{}的位置索引
func getPlaceholderIndex(tpl string) ([]int, error) {
	// 当前层级
	var level int
	// 占位符位置
	var idx int
	// 多组占位符位置
	var idxs []int
	for i := 0; i < len(tpl); i++ {
		switch tpl[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("webapi: unbalanced braces in %q", tpl)
			}
		}
	}
	// 占位符开头、结束不匹配
	if level != 0 {
		return nil, fmt.Errorf("webapi: unbalanced braces in %q", tpl)
	}
	return idxs, nil
}

// varGroupName builds a capturing group name for the indexed variable.
func varGroupName(idx int) string {
	return "v" + strconv.Itoa(idx)
}
