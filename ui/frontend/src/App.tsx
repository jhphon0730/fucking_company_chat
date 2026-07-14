// import {Greet} from "../wailsjs/go/main/App";

import { useEffect } from "react";
import { HashRouter, Route, Routes } from "react-router-dom";

import OfflineLayout from "./components/layouts/OfflineLayout";
import LoginPage from "./pages/offline/LoginPage";
import RegisterPage from "./pages/offline/RegisterPage";

import FriendsPage from "./pages/online/FriendsPage";
import OnlineLayout from "./components/layouts/OnlineLayout";
import useAuthStore from "./stores/authStore";
import { ConnectWebSocket, DisconnectWebSocket } from "../wailsjs/go/services/HTTPClientService";

function App() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const token = useAuthStore((state) => state.token);

  useEffect(() => {
    if (!isAuthenticated || !token) {
      void DisconnectWebSocket().catch(() => undefined);
      return;
    }

    void ConnectWebSocket().catch(() => undefined);
  }, [isAuthenticated, token]);

  return (
    <HashRouter>
      <Routes>
        <Route element={<OfflineLayout />}>
          <Route path="/" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />
        </Route>

        <Route element={<OnlineLayout />}>
          <Route path="online/friends" element={<FriendsPage />} />
          <Route path="online/chats" element={<>추가 예정</>} />
        </Route>
      </Routes>
    </HashRouter>
  );
}

export default App;
