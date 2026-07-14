import { Outlet } from "react-router-dom";

export default function OfflineLayout() {
  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(14,165,233,0.16),_transparent_38%),linear-gradient(135deg,_#f8fbff_0%,_#eef4ff_45%,_#f7f9ff_100%)] px-2 py-6 sm:px-4">
      <div className="mx-auto flex min-h-screen max-w-7xl items-center justify-center">
        <Outlet />
      </div>
    </div>
  );
}