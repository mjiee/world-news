import { create } from "zustand";

interface NotificationState {
  message: string | null;
  showNotification: (message: string) => void;
  closeNotification: () => void;
}

export const useNotificationStore = create<NotificationState>((set) => ({
  message: null,
  showNotification: (message: string) => {
    set({ message });

    setTimeout(() => {
      set({ message: null });
    }, 3000); // Hide after 3 seconds
  },
  closeNotification: () => set({ message: null }),
}));
