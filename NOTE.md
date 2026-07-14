# User
 [x] POST /api/v1/auth - users (사용자 모두 조회)
 [x] POST /api/v1/auth/register - register
 [x] POST /api/v1/auth/login - login

# Chat & Room (HTTP)
 [x] GET /api/v1/room - get rooms (사용자가 속한 방 조회)
 [x] POST /api/v1/room - create room (그룹 방 생성)
 [x] POST /api/v1/room/:room_id/participants - add participants (그룹 방에 사람 추가)
 [x] POST /api/v1/room/:room_id/read (방, 그룹 방 진입 시에 읽음 처리)

# Chat & Room (WebSocket)
 [x] WS TypeTalkMessage (1:1 채팅)
 [x] WS TypeReadMessage (방 진입 시에 읽음 처리)
 [ ] WS TypeTalkGroupMessage (그룹 채팅)
 [ ] WS TypeLeaveMessage (방 나가기 처리)

# WebSocket
 * /ws