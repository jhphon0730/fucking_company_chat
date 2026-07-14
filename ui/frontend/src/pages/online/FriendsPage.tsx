import { useEffect, useState } from "react";
import { MessageCircleMore, Search, Users } from "lucide-react";

import { model } from "../../../wailsjs/go/models";
import { FindAllUsers } from "../../../wailsjs/go/services/HTTPClientService";

import useAuthStore from "../../stores/authStore";
import { ShowAPIError } from "../../utils/alert";

const FriendsPage = () => {
  const user = useAuthStore((state) => state.user);
  const [friends, setFriends] = useState<model.User[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetchFriends()
  }, []);

  const fetchFriends = async () => {
    setIsLoading(true);
    try {
      const res = await FindAllUsers();
      const list = Array.isArray(res?.data?.users) ? res.data.users : [];
      setFriends(list);
    } catch (error) {
      ShowAPIError(error, "친구 목록 조회 실패")
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex h-full flex-col bg-white">
      <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
        <div>
          <p className="text-sm font-semibold text-slate-500">친구 목록</p>
          <h1 className="text-xl font-semibold text-slate-900">모든 사용자</h1>
        </div>
        <div className="flex items-center gap-2 rounded-2xl bg-slate-100 px-3 py-2 text-sm text-slate-600">
          <Users className="h-4 w-4" />
          {friends.length}명
        </div>
      </div>

      <div className="px-4 py-3">
        <div className="flex items-center gap-2 rounded-2xl border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-500">
          <Search className="h-4 w-4" />
          <input
            type="text"
            placeholder="사용자 검색"
            className="w-full bg-transparent text-sm text-slate-900 outline-none placeholder:text-slate-400"
          />
        </div>
      </div>

      <div className="flex-1 space-y-2 overflow-auto px-3 pb-4">
        {isLoading ? (
          <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4 text-sm text-slate-500">
            사용자 목록을 불러오는 중입니다...
          </div>
        ) : friends.length === 0 ? (
          <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4 text-sm text-slate-500">
            표시할 사용자가 없습니다.
          </div>
        ) : (
          friends.map((friend) => {
            const isMe = user?.login_id === friend.login_id;
            const initials = friend.name?.slice(0, 2) ?? "유";

            return (
              <button
                key={friend.id.toString()}
                type="button"
                className="flex w-full items-center gap-3 rounded-2xl border border-slate-200 bg-white p-3 text-left shadow-sm transition hover:shadow-md"
              >
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900 font-semibold text-white">
                  {initials}
                </div>

                <div className="min-w-0 flex-1">
                  <div className="flex items-center justify-between gap-2">
                    <p className="truncate font-semibold text-slate-900">{friend.name}</p>
                    <span className="text-xs text-slate-500">
                      {isMe ? "나" : "대화 가능"}
                    </span>
                  </div>
                  <p className="truncate text-sm text-slate-500">@{friend.login_id}</p>
                </div>

                <div className="flex h-9 w-9 items-center justify-center rounded-2xl bg-slate-100 text-slate-600">
                  <MessageCircleMore className="h-4 w-4" />
                </div>
              </button>
            );
          })
        )}
      </div>
    </div>
  );
};

export default FriendsPage;