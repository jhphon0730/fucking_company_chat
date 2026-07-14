package services

import (
	"fmt"
	"net/http"
	"ui/services/model"
)

type GetUserRoomsResponse struct {
	Rooms []model.RoomListResponse `json:"rooms"`
}

type GetMessagesResponse struct {
	Messages   []model.ChatMessage `json:"messages"`
	NextCursor string              `json:"next_cursor"`
	Count      int                 `json:"count"`
}

func (s *HTTPClientService) GetChatMessages(roomID string, cursor string, limit int) (*APIResponse[GetMessagesResponse], error) {
	url := fmt.Sprintf("%s/room/%s/messages?limit=%d", s.baseURL, roomID, limit)
	if cursor != "" {
		url += "&cursor=" + cursor
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	apiResp, err := getResponse[GetMessagesResponse](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

func (s *HTTPClientService) GetUserRooms() (*APIResponse[GetUserRoomsResponse], error) {
	url := fmt.Sprintf("%s/room", s.baseURL)

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

	apiResp, err := getResponse[GetUserRoomsResponse](resp)
	if err != nil {
		return nil, err
	}

	if apiResp.Error != nil {
		return nil, setErrResponse(apiResp.Error)
	}

	return apiResp, nil
}
