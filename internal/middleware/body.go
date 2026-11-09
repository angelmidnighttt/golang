package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// Body dài hơn mức này chỉ log phần đầu, tránh file log phình to.
const maxBodyLogSize = 4 << 10 // 4KB

// describeBody trả về chuỗi để log cho một body (request hoặc response):
//   - multipart/form-data: "username=an&password=***&avatar=<file a.png, image/png, 160322 bytes>"
//   - dữ liệu nhị phân (ảnh, file...): "<binary image/png, 160322 bytes>"
//   - còn lại (JSON, form, text...): nội dung, cắt ở 4KB
//
// Mọi trường hợp đều che dữ liệu nhạy cảm.
//   - head: tối đa maxBodyLogSize+1 byte đầu của body
//   - size: tổng kích thước body (-1 nếu không biết)
//   - form: multipart form handler đã parse (nil nếu chưa) để lấy kích thước file chính xác
func describeBody(head []byte, contentType string, size int64, form *multipart.Form) string {
	if len(head) == 0 {
		return ""
	}
	if boundary, ok := multipartBoundary(contentType); ok {
		if s, ok := describeMultipart(head, boundary, form); ok {
			return s
		}
	}
	if looksBinary(head) {
		return binarySummary("", detectType(contentType, head), size)
	}
	return describeText(head)
}

// describeText: che dữ liệu nhạy cảm, dài hơn maxBodyLogSize thì cắt và gắn "...(truncated)".
// Che trước rồi mới gắn "(truncated)" để value bị cắt dở ở cuối vẫn bị che.
func describeText(head []byte) string {
	if len(head) > maxBodyLogSize {
		return maskSensitive(string(trimIncompleteRune(head[:maxBodyLogSize]))) + "...(truncated)"
	}
	return maskSensitive(string(head))
}

// describeMultipart liệt kê các field trong phần đầu body multipart.
// ok = false nếu không đọc được part nào (body hỏng / sai boundary).
func describeMultipart(head []byte, boundary string, form *multipart.Form) (string, bool) {
	reader := multipart.NewReader(bytes.NewReader(head), boundary)
	var fields []string
	truncated := false

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			truncated = true
			break
		}

		// readErr != nil nghĩa là part bị cắt ngang ở mốc 4KB.
		data, readErr := io.ReadAll(part)
		name := part.FormName()

		switch {
		case part.FileName() != "" || looksBinary(data):
			size := int64(len(data))
			if readErr != nil {
				size = uploadedFileSize(form, name, part.FileName())
			}
			ctype := detectType(part.Header.Get("Content-Type"), data)
			fields = append(fields, name+"="+binarySummary(part.FileName(), ctype, size))
		case isSensitiveName(name):
			fields = append(fields, name+"="+masked)
		default:
			fields = append(fields, name+"="+string(trimIncompleteRune(data)))
		}

		if readErr != nil {
			truncated = true
			break
		}
	}

	if len(fields) == 0 {
		return "", false
	}
	s := strings.Join(fields, "&")
	if truncated {
		s += "...(truncated)"
	}
	return s, true
}

// uploadedFileSize lấy kích thước file từ form handler đã parse, -1 nếu không có.
func uploadedFileSize(form *multipart.Form, field, filename string) int64 {
	if form == nil {
		return -1
	}
	for _, fh := range form.File[field] {
		if fh.Filename == filename {
			return fh.Size
		}
	}
	return -1
}

func multipartBoundary(contentType string) (string, bool) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") || params["boundary"] == "" {
		return "", false
	}
	return params["boundary"], true
}

// looksBinary: true nếu dữ liệu không phải text (không phải UTF-8 hợp lệ
// hoặc có ký tự điều khiển như \x00). Xem nội dung chứ không tin Content-Type.
func looksBinary(b []byte) bool {
	// Bỏ ký tự bị cắt dở ở cuối trước, tránh nhận nhầm text là binary.
	b = trimIncompleteRune(b)
	if !utf8.Valid(b) {
		return true
	}
	for _, c := range b {
		if (c < 0x20 && c != '\t' && c != '\n' && c != '\r') || c == 0x7f {
			return true
		}
	}
	return false
}

// trimIncompleteRune bỏ ký tự UTF-8 bị cắt dở ở cuối (do cắt ở mốc 4KB),
// vd chữ "ệ" (3 byte) chỉ còn 2 byte đầu. Dữ liệu hợp lệ thì giữ nguyên.
func trimIncompleteRune(b []byte) []byte {
	for i := 0; i < utf8.UTFMax-1 && len(b) > 0 && !utf8.Valid(b); i++ {
		b = b[:len(b)-1]
	}
	return b
}

// detectType: dùng Content-Type được khai báo, nếu trống hoặc chung chung
// (application/octet-stream) thì tự đoán từ nội dung, vd: "image/png".
func detectType(contentType string, data []byte) string {
	if contentType != "" && contentType != "application/octet-stream" {
		return contentType
	}
	return http.DetectContentType(data)
}

func binarySummary(filename, contentType string, size int64) string {
	sizeText := "unknown size"
	if size >= 0 {
		sizeText = fmt.Sprintf("%d bytes", size)
	}
	if filename != "" {
		return fmt.Sprintf("<file %s, %s, %s>", filename, contentType, sizeText)
	}
	return fmt.Sprintf("<binary %s, %s>", contentType, sizeText)
}

// peekBody đọc tối đa maxBodyLogSize+1 byte đầu của request body để log
// (mọi method, mọi Content-Type), rồi trả lại body đầy đủ cho handler phía sau.
func peekBody(r *http.Request) []byte {
	if r.Body == nil || r.Body == http.NoBody {
		return nil
	}

	// Lỗi đọc body (nếu có) handler sẽ tự gặp lại khi đọc, ở đây chỉ log phần đọc được.
	head, _ := io.ReadAll(io.LimitReader(r.Body, maxBodyLogSize+1))

	// Ghép phần vừa đọc với phần còn lại để handler vẫn đọc được đủ body.
	r.Body = struct {
		io.Reader
		io.Closer
	}{io.MultiReader(bytes.NewReader(head), r.Body), r.Body}

	return head
}

// responseBodyWriter ghi response ra client như bình thường,
// đồng thời chép lại tối đa maxBodyLogSize+1 byte đầu để log.
type responseBodyWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	if n := w.room(); n > 0 {
		w.body.Write(b[:min(n, len(b))])
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	if n := w.room(); n > 0 {
		w.body.WriteString(s[:min(n, len(s))])
	}
	return w.ResponseWriter.WriteString(s)
}

// room: còn chép thêm được bao nhiêu byte.
func (w *responseBodyWriter) room() int {
	return maxBodyLogSize + 1 - w.body.Len()
}
