import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import { model } from "../../wailsjs/go/models";

type AuthState = {
  user: model.User | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: model.User | null, token: string | null) => void;
  clearAuth: () => void;
  getToken: () => string | null;
  getUser: () => model.User | null;
};

const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      setAuth: (user, token) =>
        set({
          user,
          token,
          isAuthenticated: Boolean(user && token),
        }),

      clearAuth: () =>
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        }),

      getToken: () => get().token,
      getUser: () => get().user,
    }),
    {
      name: "auth-store",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
);

export default useAuthStore;
