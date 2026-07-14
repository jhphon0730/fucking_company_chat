import { useEffect, useState, useRef } from "react";
import { Send } from "lucide-react";

import useFriendStore from "../../stores/friendStore";
import useAuthStore from "../../stores/authStore";

import { model } from "../../../wailsjs/go/models";
import { EventsOn, EventsOff } from "../../../wailsjs/runtime"
import { GetChatMessages, SendMessage } from "../../../wailsjs/go/services/HTTPClientService";
import { TypeTalkMessage } from "../../types/msg";

interface ChatRoomPageProps {
  roomId: string;
  roomName: string;
  isGroupChat: boolean;

  onBack: () => void;
}

const ChatRoomPage = ({ roomId, roomName, isGroupChat, onBack }: ChatRoomPageProps) => {
  const [messages, setMessages] = useState<model.ChatMessage[]>([]);
  const [input, setInput] = useState("");
  const [targetFriend, setTargetFriend] = useState<model.User>()

  const scrollRef = useRef<HTMLDivElement>(null);

  const myUser = useAuthStore((state) => state.user);
  const getFriends = useFriendStore((state) => state.getFriends)


  const fetchMessages = async () => {
    try {
      const res = await GetChatMessages(roomId, "", 30);
      if (res?.data?.messages) {
        setMessages([...res.data.messages].reverse());
      }
    } catch (err) {
      console.error("메시지 로드 실패:", err);
    }
  };

  // 방 입장 시에 상대방 ID 추출 (그룹 아닐 시에)
  useEffect(() => {
    if (isGroupChat) { return }

    const findUser = getFriends().filter((item) => item.name == roomName)[0]
    setTargetFriend(findUser)
  }, [roomId])

  // 과거 메시지 로드
  useEffect(() => {
    void fetchMessages();
  }, [roomId]);

  // 자동 스크롤
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  // 메시지 전송
  const handleSend = async () => {
    if (!input.trim()) return;
    try {
      await SendMessage(roomId, myUser?.id, targetFriend?.id || null, input);
      setInput("");
    } catch (e) {
      console.error("전송 실패:", e);
    }
  };

  useEffect(() => {
    // TALK 타입 메시지만 수신
    const handleTalkMessage = (msg: model.WSMessage) => {
      // 현재 보고 있는 방과 같은 방 메시지인지 확인 후 반영
      if (msg.room_id?.toString() === roomId) {
        void fetchMessages()
      }
    };

    EventsOn(TypeTalkMessage, handleTalkMessage); // "TALK" 이벤트 구독

    return () => {
      EventsOff(TypeTalkMessage);
    };
  }, [roomId]);

  return (
    <div className="h-full flex flex-col bg-slate-50">
      <div className="flex items-center gap-2 border-b bg-white p-4">
        <button onClick={onBack} className="rounded-lg p-2 transition hover:bg-slate-100">
          ←
        </button>
        <h2 className="font-semibold text-slate-900">채팅방 {roomId.slice(0, 8)}...</h2>
      </div>

      <div ref={scrollRef} className="flex-1 overflow-auto p-4 space-y-4">
        {messages.map((msg) => {
          const isMyMessage = msg.sender_id === myUser?.id;

          return (
            <div
              key={msg.id.toString()}
              className={`flex w-full ${isMyMessage ? "justify-end" : "justify-start"}`}
            >
              <div
                className={`max-w-[70%] rounded-2xl p-3 text-sm shadow-sm ${
                  isMyMessage
                    ? "rounded-tr-none bg-slate-900 text-white"
                    : "rounded-tl-none border border-slate-200 bg-white text-slate-900"
                }`}
              >
                <p className="whitespace-pre-wrap">{msg.content}</p>
                <div
                  className={`mt-1 text-[10px] opacity-70 ${
                    isMyMessage ? "text-slate-300" : "text-slate-400"
                  }`}
                >
                  {new Date(msg.created_at).toLocaleTimeString([], {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="border-t bg-white p-2 lg:min-h-0 min-h-[120px]">
        <div className="flex gap-2">
          <input
            className="flex-1 min-w-0 rounded-lg border border-slate-200 p-2 outline-none"
            placeholder="메시지를 입력하세요..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
          />
          <button
            onClick={handleSend}
            className="rounded-lg bg-slate-900 p-2 text-white transition hover:bg-slate-800"
          >
            <Send size={18} />
          </button>
        </div>
      </div>
    </div>
  );
};

export default ChatRoomPage;