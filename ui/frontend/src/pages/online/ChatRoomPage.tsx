import { useEffect, useState, useRef } from "react";
import { Send } from "lucide-react";
import useAuthStore from "../../stores/authStore";
import { GetChatMessages, SendMessage } from "../../../wailsjs/go/services/HTTPClientService";
import { model } from "../../../wailsjs/go/models";

interface ChatRoomPageProps {
  roomId: string;
  onBack: () => void;
}

const ChatRoomPage = ({ roomId, onBack }: ChatRoomPageProps) => {
  const [messages, setMessages] = useState<model.ChatMessage[]>([]);
  const [input, setInput] = useState("");
  const myUser = useAuthStore((state) => state.user);
  const scrollRef = useRef<HTMLDivElement>(null);

  // 과거 메시지 로드
  useEffect(() => {
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
    void fetchMessages();
  }, [roomId]);

  // 메시지 추가 시 하단 스크롤 자동 이동
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  // 메시지 전송
  const handleSend = async () => {
    if (!input.trim()) return;
    try {
      await SendMessage(roomId, myUser?.id, null, input);
      setInput("");
    } catch (e) {
      console.error("전송 실패:", e);
    }
  };

  return (
    <div className="h-full flex flex-col flex-1 bg-slate-50">
      {/* 헤더 */}
      <div className="flex items-center gap-2 border-b bg-white p-4">
        <button onClick={onBack} className="rounded-lg p-2 transition hover:bg-slate-100">
          ←
        </button>
        <h2 className="font-semibold text-slate-900">채팅방 {roomId.slice(0, 8)}...</h2>
      </div>

      {/* 메시지 영역: flex-1과 overflow-auto로 고정된 스크롤 생성 */}
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

      {/* 입력창 (하단 고정) */}
      <div className="border-t bg-white p-4">
        <div className="flex gap-2">
          <input
            className="flex-1 rounded-lg border border-slate-200 p-2 outline-none"
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