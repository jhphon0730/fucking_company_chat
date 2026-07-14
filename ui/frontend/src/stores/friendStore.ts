import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import { model } from "../../wailsjs/go/models";

type friendStatus = {
  friends: model.User[]
  getFriends: () => model.User[]
  setFriends: (friends: model.User[]) => void
  clearFriends: () => void
}

const useFriendStore = create<friendStatus>()(
  persist(
    (set, get) => ({
      friends: [],

      setFriends: (friends) =>
        set({
          friends,
        }),

      clearFriends: () =>
        set({
          friends: [],
        }),

      getFriends: () => get().friends,
    }),
    {
      name: "friend-store",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        friends: state.friends
      }),
    },
  ),
);

export default useFriendStore;
