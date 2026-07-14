import { useEffect, useState } from "react";
import { Plus, Search } from "lucide-react";

import ChatRoomPage from "./ChatRoomPage";

import { GetUserRooms } from "../../../wailsjs/go/services/HTTPClientService";
import { ShowAPIError } from "../../utils/alert";

interface RoomListItem {
  id: string;
  name: string;
  is_group: boolean;
  last_message: string;
  last_message_at?: string;
  unread_count: number;
}

const ChatsPage = () => {
  const [rooms, setRooms] = useState<RoomListItem[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null);

  useEffect(() => {
    void fetchRooms();
  }, []);

  const fetchRooms = async () => {
    setIsLoading(true);

    try {
      const res = await GetUserRooms();
      const list = Array.isArray(res?.data?.rooms) ? res.data.rooms : [];
      setRooms(
        list.map((room: any) => ({
          id: room.id,
          name: room.name,
          is_group: room.is_group ?? false,
          last_message: room.last_message ?? "",
          last_message_at: room.last_message_at ?? undefined,
          unread_count: room.unread_count ?? 0,
        })),
      );
    } catch (error) {
      ShowAPIError(error, "채팅방 목록 조회 실패");
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  if (selectedRoomId) {
    return <ChatRoomPage roomId={selectedRoomId} onBack={() => setSelectedRoomId(null)} />;
  }

  return (
    <div className="flex h-full flex-col bg-white">
      <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
        <div>
          <p className="text-sm font-semibold text-slate-500">채팅</p>
          <h1 className="text-xl font-semibold text-slate-900">내 채팅방</h1>
        </div>
        <button
          type="button"
          className="flex items-center gap-2 rounded-2xl bg-slate-100 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-200"
          onClick={() => void fetchRooms()}
        >
          새로고침
        </button>
      </div>

      <div className="px-4 py-3">
        <div className="flex items-center gap-2 rounded-2xl border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-500">
          <Search className="h-4 w-4" />
          <input
            type="text"
            placeholder="채팅방 검색"
            className="w-full bg-transparent text-sm text-slate-900 outline-none placeholder:text-slate-400"
            readOnly
          />
        </div>
      </div>

      <div className="flex-1 space-y-2 overflow-auto px-3 pb-4">
        {isLoading ? (
          <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4 text-sm text-slate-500">
            채팅방을 불러오는 중입니다...
          </div>
        ) : rooms.length === 0 ? (
          <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4 text-sm text-slate-500">
            참여 중인 채팅방이 없습니다.
          </div>
        ) : (
          rooms.map((room) => {
            const title = room.name || (room.is_group ? "그룹 채팅" : "1:1 채팅");
            const subtitle = room.last_message || "아직 메시지가 없습니다.";
            const receivedAt = room.last_message_at
              ? new Date(room.last_message_at).toLocaleString()
              : "";

            return (
              <button
                key={room.id}
                onClick={() => setSelectedRoomId(room.id)}
                type="button"
                className="flex w-full flex-col gap-3 rounded-2xl border border-slate-200 bg-white p-4 text-left shadow-sm transition hover:shadow-md"
              >
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="text-sm font-semibold text-slate-900">{title}</p>
                    <p className="mt-1 text-sm text-slate-500 truncate">{subtitle}</p>
                  </div>
                  <span className="rounded-2xl bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-600">
                    {room.unread_count}개
                  </span>
                </div>
                <div className="flex items-center justify-between text-xs text-slate-400">
                  <span>{room.is_group ? "그룹" : "1:1"}</span>
                  <span>{receivedAt}</span>
                </div>
              </button>
            );
          })
        )}
      </div>
    </div>
  );
};

export default ChatsPage