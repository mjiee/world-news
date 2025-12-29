import { PodcastAudio } from "@/services";
import { create } from "zustand";

interface AudioItem {
  id: number;
  audio: PodcastAudio;
  progress: number;
}

interface AudioPlayStore {
  playing: boolean;
  playlist: AudioItem[];

  addAudio: (id: number, audio: PodcastAudio) => void;
  removeAudio: (id: number) => void;
  inPlayList: (id: number) => boolean;

  play: () => void;
  pause: (progress: number) => void;
  playNext: () => void;

  getCurrentAudio: () => AudioItem | null;
  setPlaying: (isPlaying: boolean) => void;
}

export const useAudioPlayStore = create<AudioPlayStore>((set, get) => ({
  playing: false,
  playlist: [],

  addAudio: (id, audio) => {
    const state = get();

    if (state.inPlayList(id)) {
      return;
    }

    const newItem: AudioItem = { id, audio, progress: 0 };
    const isFirstAudio = state.playlist.length === 0;

    set({
      playlist: [...state.playlist, newItem],
      playing: isFirstAudio ? true : state.playing,
    });
  },

  removeAudio: (id) => {
    const state = get();
    const filteredPlaylist = state.playlist.filter((item) => item.id !== id);

    if (filteredPlaylist.length === state.playlist.length) {
      return;
    }

    const isRemovingCurrent = state.playlist[0]?.id === id;
    const hasRemainingAudios = filteredPlaylist.length > 0;

    set({
      playlist: filteredPlaylist,
      playing: hasRemainingAudios && (isRemovingCurrent || state.playing),
    });
  },

  inPlayList: (id) => {
    return get().playlist.some((item) => item.id === id);
  },

  play: () => {
    if (get().playlist.length > 0) {
      set({ playing: true });
    }
  },

  pause: (progress) => {
    const state = get();

    if (state.playlist.length === 0) {
      set({ playing: false });
      return;
    }

    const [currentAudio, ...rest] = state.playlist;
    const updatedCurrent = { ...currentAudio, progress };

    set({
      playing: false,
      playlist: [updatedCurrent, ...rest],
    });
  },

  playNext: () => {
    const state = get();
    const [, ...remainingPlaylist] = state.playlist;

    set({
      playlist: remainingPlaylist,
      playing: remainingPlaylist.length > 0,
    });
  },

  getCurrentAudio: () => {
    return get().playlist[0] ?? null;
  },

  setPlaying: (isPlaying) => {
    set({ playing: isPlaying });
  },
}));

export function buildAudioSrc(audio: PodcastAudio): string {
  if (audio.url) {
    return audio.url;
  }

  const mimeTypes: Record<string, string> = {
    mp3: "audio/mpeg",
    wav: "audio/wav",
  };

  const mimeType = mimeTypes[audio.type];

  if (mimeType && audio.data) {
    return `data:${mimeType};base64,${audio.data}`;
  }

  return "";
}
