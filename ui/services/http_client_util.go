package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getResponse[T any](resp *http.Response) (*APIResponse[T], error) {
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("서버 내부 오류 발생 (Status: %d)", resp.StatusCode)
	}

	var apiResp APIResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패 (Status: %d): %w", resp.StatusCode, err)
	}

	return &apiResp, nil
}

func setErrResponse(apiErr *APIError) error {
	if apiErr == nil {
		return fmt.Errorf("알 수 없는 서버 오류")
	}

	errBytes, jsonErr := json.Marshal(apiErr)
	if jsonErr != nil {
		return fmt.Errorf("알 수 없는 에러: %w", jsonErr)
	}

	// JSON 문자열을 error 타입으로 변환하여 리턴
	return fmt.Errorf("%s", string(errBytes))
}
