package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"ui/services/model"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RegisterRequest struct {
	LoginID  string `json:"login_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  model.User `json:"user"`
	Token string     `json:"token"`
}

type FindAllUsersResponse struct {
	Users []model.User `json:"users"`
}

func (s *HTTPClientService) Register(reqData RegisterRequest) error {
	url := fmt.Sprintf("%s/auth/register", s.baseURL)

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 공통 함수 호출
	apiResp, err := getResponse[any](resp)
	if err != nil {
		return err
	}

	if apiResp.Error != nil {
		return setErrResponse(apiResp.Error)
	}

	return nil
}

func (s *HTTPClientService) Login(reqData LoginRequest) (*APIResponse[LoginResponse], error) {
	url := fmt.Sprintf("%s/auth/login", s.baseURL)

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	apiResp, err := getResponse[LoginResponse](resp)
	if err != nil {
		return nil, err
	}

	if apiResp.Error != nil {
		return nil, setErrResponse(apiResp.Error)
	}

	return apiResp, nil
}

func (s *HTTPClientService) FindAllUsers() (*APIResponse[FindAllUsersResponse], error) {
	url := fmt.Sprintf("%s/auth", s.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	apiResp, err := getResponse[FindAllUsersResponse](resp)
	if err != nil {
		return nil, err
	}

	if apiResp.Error != nil {
		return nil, setErrResponse(apiResp.Error)
	}

	return apiResp, nil
}

func (s *HTTPClientService) SendMessage(roomID uuid.UUID, senderID uuid.UUID, receiverID *uuid.UUID, content string) error {
	s.chatMu.Lock()
	conn := s.chatConn
	s.chatMu.Unlock()

	if conn == nil {
		return fmt.Errorf("웹소켓 연결이 끊겨있습니다")
	}

	msg := model.WSMessage{
		Type:       model.TypeTalkMessage,
		RoomID:     &roomID,
		ReceiverID: receiverID,
		SenderID:   senderID,
		Content:    content,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, payload)
}
