package httputil

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type APIResponse struct {
	Data    any       `json:"data"`
	Message string    `json:"message,omitempty"`
	Error   *APIError `json:"error"`
}

func OK(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{
		Data:    data,
		Message: "",
		Error:   nil,
	})
}

func OKMessage(c *gin.Context, status int, message string, data any) {
	c.JSON(status, APIResponse{
		Data:    data,
		Message: message,
		Error:   nil,
	})
}

func Fail(c *gin.Context, status int, message string, details any) {
	// runtime.Caller(1)은 이 Fail 함수를 호출한 바로 위(상위) 함수의 정보를 가져옴
	pc, file, line, ok := runtime.Caller(1)
	callerInfo := "unknown"

	if ok {
		// 함수명 가져오기 (예: user.(*userHandler).Login)
		fullFuncName := runtime.FuncForPC(pc).Name()
		// 패키지 전체 경로 제거하고 패키지명.함수명 형태로 추출
		funcParts := strings.Split(fullFuncName, "/")
		shortFuncName := funcParts[len(funcParts)-1]

		// 파일 경로 가공하기 (절대 경로 -> 프로젝트 상대 경로)
		// "internal/" 문자열이 시작되는 위치를 찾아 그 앞의 로컬 절대 경로를 통째로 잘라내기
		displayFile := file
		if idx := strings.Index(file, "internal/"); idx != -1 {
			displayFile = file[idx:]
		} else {
			// 만약 internal 폴더 외부라면 가독성을 위해 가장 마지막 폴더/파일 이름만 출력
			parts := strings.Split(file, "/")
			if len(parts) > 2 {
				displayFile = strings.Join(parts[len(parts)-2:], "/")
			}
		}

		// 결과 예시: user.(*userHandler).Login (internal/user/handler.go:82)
		callerInfo = fmt.Sprintf("%s (%s:%d)", shortFuncName, displayFile, line)
	}

	c.JSON(status, APIResponse{
		Data:    nil,
		Message: "",
		Error: &APIError{
			Code:    callerInfo,
			Message: message,
			Details: details,
		},
	})
}
