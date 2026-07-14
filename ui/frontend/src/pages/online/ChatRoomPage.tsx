import { useEffect, useState } from "react";
import useAuthStore from "../../stores/authStore";
import { GetChatMessages } from "../../../wailsjs/go/services/HTTPClientService";
import { model } from "../../../wailsjs/go/models";

interface ChatRoomPageProps {
  roomId: string;
  onBack: () => void;
}

const ChatRoomPage = ({ roomId, onBack }: ChatRoomPageProps) => {
  const [messages, setMessages] = useState<model.ChatMessage[]>([]);
  const myUser = useAuthStore((state) => state.user);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const res = await GetChatMessages(roomId, "", 30);
        if (res?.data?.messages) {
          // 메시지를 최신순으로 보여주기 위해 reverse() 사용
          setMessages([...res.data.messages].reverse());
        }
      } catch (err) {
        console.error("메시지 로드 실패:", err);
      }
    };
    void fetchMessages();
  }, [roomId]);

  return (
    <div className="flex h-full flex-col bg-slate-50">
      {/* 헤더 부분 */}
      <div className="p-4 border-b bg-white flex items-center gap-2">
        <button 
          onClick={onBack} 
          className="p-2 hover:bg-slate-100 rounded-lg transition"
        >
          ←
        </button>
        <h2 className="font-semibold text-slate-900">
          채팅방 {roomId.slice(0, 8)}...
        </h2>
      </div>

      {/* 메시지 리스트 영역 */}
      <div className="flex-1 overflow-auto p-4 space-y-4">
        {messages.map((msg) => {
          // 메시지의 sender_id와 내 ID를 비교하여 정렬 판단
          const isMyMessage = msg.sender_id === myUser?.id;

          return (
            <div 
              key={msg.id.toString()} 
              className={`flex w-full ${isMyMessage ? "justify-end" : "justify-start"}`}
            >
              <div 
                className={`p-3 rounded-2xl shadow-sm max-w-[70%] text-sm ${
                  isMyMessage 
                    ? "bg-slate-900 text-white rounded-tr-none" 
                    : "bg-white text-slate-900 rounded-tl-none border border-slate-200"
                }`}
              >
                <p className="whitespace-pre-wrap">{msg.content}</p>
                <div className={`text-[10px] mt-1 opacity-70 ${isMyMessage ? "text-slate-300" : "text-slate-400"}`}>
                  {new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default ChatRoomPage;