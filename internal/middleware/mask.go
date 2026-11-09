package middleware

import (
	"fmt"
	"regexp"
)

// Tên field bị coi là nhạy cảm (không phân biệt hoa thường):
//   - chứa một trong các từ ở sensitiveContains, vd: new_password, access_token, X-Api-Key
//   - hoặc đúng bằng một từ ở sensitiveExact (so nguyên tên để "shipping" không bị che vì chứa "pin")
//
// Muốn che thêm field nào thì thêm vào 2 danh sách này.
const (
	sensitiveContains = `pass(?:word|wd|code)|secret|token|api[_-]?key|authorization|card[_-]?number|cvv`
	sensitiveExact    = `otp|pin|pwd`
	masked            = "***"
)

var (
	// JSON: "password": "abc"  |  "pin": 1234  |  "api_key": {...}  |  "tokens": [...]
	jsonSensitive = regexp.MustCompile(fmt.Sprintf(
		`(?i)("(?:[^"]*(?:%s)[^"]*|%s)"\s*:\s*)("(?:[^"\\]|\\.)*"?|\{[^}]*\}?|\[[^\]]*\]?|[^,}\]\s]+)`,
		sensitiveContains, sensitiveExact))

	// Form urlencoded / query string: password=abc&name=An  |  ...?token=abc (cả URL nằm trong JSON)
	formSensitive = regexp.MustCompile(fmt.Sprintf(
		`(?i)((?:^|[&?])(?:[\w.\[\]%%-]*(?:%s)[\w.\[\]%%-]*|%s)=)[^&\s"]*`,
		sensitiveContains, sensitiveExact))

	// Multipart form-data: Content-Disposition: form-data; name="password"\r\n\r\nabc
	multipartSensitive = regexp.MustCompile(fmt.Sprintf(
		`(?i)(name="(?:[^"]*(?:%s)[^"]*|%s)"[^\r\n]*\r\n(?:[^\r\n]+\r\n)*\r\n)[^\r\n]*`,
		sensitiveContains, sensitiveExact))

	// Chỉ tên field: "password", "new_password", "otp"...
	sensitiveName = regexp.MustCompile(fmt.Sprintf(
		`(?i)^(?:.*(?:%s).*|%s)$`,
		sensitiveContains, sensitiveExact))
)

// isSensitiveName: tên field có phải field nhạy cảm không.
func isSensitiveName(name string) bool {
	return sensitiveName.MatchString(name)
}

// maskSensitive thay giá trị của các field nhạy cảm bằng "***".
// Áp cả 3 dạng (JSON, form/query, multipart) lên mọi chuỗi, không dựa vào Content-Type
// vì client có thể gửi sai. Dùng regex thay vì parse nên vẫn che được body bị cắt ngắn.
func maskSensitive(s string) string {
	s = jsonSensitive.ReplaceAllString(s, `${1}"`+masked+`"`)
	s = formSensitive.ReplaceAllString(s, "${1}"+masked)
	s = multipartSensitive.ReplaceAllString(s, "${1}"+masked)
	return s
}
