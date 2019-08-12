package utils

import (
	"strings"
)

//upperCaseWords 全大写的字符串
var upperCaseWords = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
	"CN":    true,
	"JP":    true,
	"NO":    true,
}

//ConvertCamel2Line 大驼峰式转下划线式
func ConvertCamel2Line(camel string) string {
	bytes := []byte(camel)
	if len(bytes) > 0 {
		newBytes := make([]byte, 0)

		var ws, we int
		ws = -1
		for i := 0; i < len(bytes); i++ {
			if ws == -1 {
				ws = i
				we = i
				continue
			}

			we = i
			if bytes[i] >= 'A' && bytes[i] <= 'Z' {
				if we-ws >= 1 &&
					bytes[we-1] >= 'A' && bytes[we-1] <= 'Z' {
					//连续大写跳过
				} else {
					if ws > 0 {
						newBytes = append(newBytes, '_')
					}

					newBytes = append(newBytes, bytes[ws]+32)

					for j := ws + 1; j < we; j++ {
						newBytes = append(newBytes, bytes[j])
					}

					ws = we
				}
			} else {
				if we-ws >= 2 &&
					bytes[we-1] >= 'A' && bytes[we-1] <= 'Z' &&
					bytes[we-2] >= 'A' && bytes[we-2] <= 'Z' {

					word := string(bytes[ws : we-1])
					if ws > 0 {
						newBytes = append(newBytes, '_')
					}

					if upperCaseWords[word] {
						for j := ws; j < we-1; j++ {
							newBytes = append(newBytes, bytes[j]+32)
						}
					} else {
						newBytes = append(newBytes, bytes[ws]+32)
						for j := ws + 1; j < we-1; j++ {
							newBytes = append(newBytes, '_')
							newBytes = append(newBytes, bytes[j]+32)
						}
					}

					ws = we - 1
				}
			}

			if i == len(bytes)-1 {
				if ws > 0 {
					newBytes = append(newBytes, '_')
				}

				if we-ws >= 1 && bytes[we-1] >= 'A' && bytes[we-1] <= 'Z' {
					word := string(bytes[ws:])
					if upperCaseWords[word] {
						for j := ws; j <= we; j++ {
							newBytes = append(newBytes, bytes[j]+32)
						}
					} else {
						newBytes = append(newBytes, bytes[ws]+32)
						for j := ws + 1; j <= we; j++ {
							newBytes = append(newBytes, '_')
							newBytes = append(newBytes, bytes[j]+32)
						}
					}
				} else {
					newBytes = append(newBytes, bytes[ws]+32)

					for j := ws + 1; j <= we; j++ {
						newBytes = append(newBytes, bytes[j])
					}
				}
			}
		}

		if len(newBytes) > 0 {
			return string(newBytes)
		}
	}

	return camel
}

//ConvertLine2Camel 下划线式转大驼峰式
func ConvertLine2Camel(line string) string {
	bytes := []byte(line)
	if len(bytes) > 0 {
		newbytes := make([]byte, 0)

		var ws, we int
		ws = -1

		for i := 0; i < len(bytes); i++ {
			if ws == -1 {
				ws = i
				we = i
				continue
			}

			we = i
			if i == len(bytes)-1 || (bytes[i] == '_' && i > 0) {
				if i == len(bytes)-1 {
					we++
				}

				word := string(bytes[ws:we])

				if u := strings.ToUpper(word); upperCaseWords[u] {
					b := []byte(u)
					for j := 0; j < len(b); j++ {
						newbytes = append(newbytes, b[j])
					}
				} else {
					newbytes = append(newbytes, bytes[ws]-32)
					for j := ws + 1; j < we; j++ {
						newbytes = append(newbytes, bytes[j])
					}
				}

				ws = i + 1
				i++
			}
		}

		if len(newbytes) > 0 {
			return string(newbytes)
		}
	}

	return line
}
