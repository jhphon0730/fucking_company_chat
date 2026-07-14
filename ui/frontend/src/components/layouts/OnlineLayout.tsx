import { Link, Outlet, useNavigate } from "react-router-dom";
import { LogOut, Menu, MessageSquareMore, Plus, Users } from "lucide-react";
import { useState } from "react";
import useAuthStore from "../../stores/authStore";

export default function OnlineLayout() {
  const navigate = useNavigate();
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const user = useAuthStore((state) => state.user);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleLogout = () => {
    clearAuth();
    navigate("/");
  };

  const menuItems = [
    { to: "/online/friends", label: "친구 목록", icon: Users },
    { to: "/online/chats", label: "채팅방", icon: MessageSquareMore },
    { to: "/online/friends", label: "방 생성", icon: Plus },
  ];

  return (
    <div className="bg-slate-100 text-slate-800 h-screen max-h-screen overflow-hidden p-3">
      <div className="mx-auto flex max-w-7xl overflow-hidden rounded-[24px] border border-slate-200 bg-white shadow-[0_8px_30px_rgba(15,23,42,0.08)] h-full">
        {/* 데스크톱 사이드바 */}
        <aside className="hidden w-80 flex-col border-r border-slate-200 bg-slate-50 lg:flex">
          <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900 text-white">
                <MessageSquareMore className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm font-semibold text-slate-900">G-kk-ch</p>
                <p className="text-xs text-slate-500">사내 채팅</p>
              </div>
            </div>
          </div>

          <div className="flex-1 space-y-2 px-3 py-3">
            {menuItems.map((item) => {
              const Icon = item.icon;
              return (
                <Link
                  key={item.label}
                  to={item.to}
                  className="flex items-center gap-3 rounded-2xl px-3 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-200"
                >
                  <div className="flex h-9 w-9 items-center justify-center rounded-2xl bg-white shadow-sm">
                    <Icon className="h-4 w-4" />
                  </div>
                  {item.label}
                </Link>
              );
            })}
          </div>

          <div className="border-t border-slate-200 px-4 py-4">
            <div className="flex items-center justify-between rounded-2xl bg-white px-3 py-3 shadow-sm">
              <div>
                <p className="text-sm font-semibold text-slate-900">{user?.name ?? "사용자"}</p>
                <p className="text-xs text-slate-500">@{user?.login_id ?? "unknown"}</p>
              </div>
              <button
                type="button"
                onClick={handleLogout}
                className="rounded-xl p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
                title="로그아웃"
              >
                <LogOut className="h-4 w-4" />
              </button>
            </div>
          </div>
        </aside>

        {/* 모바일 오버레이 */}
        {isMenuOpen && (
          <div
            className="fixed inset-0 z-40 bg-black/50 lg:hidden"
            onClick={() => setIsMenuOpen(false)}
          />
        )}

        {/* 모바일 슬라이드 메뉴 */}
        <aside
          className={`fixed left-0 top-0 z-50 h-screen w-72 transform border-r border-slate-200 bg-slate-50 shadow-2xl transition-transform duration-300 ease-in-out lg:hidden ${
            isMenuOpen ? "translate-x-0" : "-translate-x-full"
          }`}
        >
          <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900 text-white">
                <MessageSquareMore className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm font-semibold text-slate-900">G-kk-ch</p>
                <p className="text-xs text-slate-500">사내 채팅</p>
              </div>
            </div>
            <button
              type="button"
              onClick={() => setIsMenuOpen(false)}
              className="rounded-xl p-2 text-slate-500 transition hover:bg-slate-100"
            >
              <Menu className="h-5 w-5" />
            </button>
          </div>

          <div className="flex-1 space-y-2 px-3 py-3 overflow-y-auto">
            {menuItems.map((item) => {
              const Icon = item.icon;
              return (
                <Link
                  key={item.label}
                  to={item.to}
                  className="flex items-center gap-3 rounded-2xl px-3 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-200"
                  onClick={() => setIsMenuOpen(false)}
                >
                  <div className="flex h-9 w-9 items-center justify-center rounded-2xl bg-white shadow-sm">
                    <Icon className="h-4 w-4" />
                  </div>
                  {item.label}
                </Link>
              );
            })}
          </div>

          <div className="border-t border-slate-200 px-4 py-4">
            <div className="flex items-center justify-between rounded-2xl bg-white px-3 py-3 shadow-sm">
              <div>
                <p className="text-sm font-semibold text-slate-900">{user?.name ?? "사용자"}</p>
                <p className="text-xs text-slate-500">@{user?.login_id ?? "unknown"}</p>
              </div>
              <button
                type="button"
                onClick={handleLogout}
                className="rounded-xl p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
                title="로그아웃"
              >
                <LogOut className="h-4 w-4" />
              </button>
            </div>
          </div>
        </aside>

        <section className="flex-1 max-h-full bg-white">
          <div className="flex items-center justify-between border-b border-slate-200 px-4 py-3 lg:hidden">
            <button
              type="button"
              onClick={() => setIsMenuOpen((prev) => !prev)}
              className="rounded-xl p-2 text-slate-700 transition hover:bg-slate-100"
            >
              <Menu className="h-5 w-5" />
            </button>
            <div className="text-sm font-semibold text-slate-900">G-kk-ch</div>
          </div>

          <Outlet />
        </section>
      </div>
    </div>
  );
}