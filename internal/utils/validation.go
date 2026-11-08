package utils

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// Validator gom lỗi theo field, ví dụ: {"id": "id must be a positive integer"}
type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Check thêm lỗi cho field nếu ok == false.
// Mỗi field chỉ giữ lỗi đầu tiên.
func (v *Validator) Check(ok bool, field, message string) {
	if ok {
		return
	}
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// Valid trả về true nếu không có lỗi nào.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Required: chuỗi không rỗng (bỏ qua khoảng trắng hai đầu).
func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinLength: độ dài tối thiểu, đếm theo ký tự (rune) để chữ có dấu tính đúng 1 ký tự.
func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxLength: độ dài tối đa, đếm theo ký tự (rune).
func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// ParsePositiveInt: chuỗi phải là số nguyên > 0. Trả về số đã parse và ok.
func ParsePositiveInt(value string) (int, bool) {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}
