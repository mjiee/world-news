import { ActionIcon, Box, Group, Stack, Text, Tooltip } from "@mantine/core";
import { IconPlayerPause, IconPlayerPlay, IconPlayerSkipForward } from "@tabler/icons-react";
import { useEffect, useRef, useState } from "react";

import styles from "@/assets/styles/appLayout.module.css";
import { buildAudioSrc, useAudioPlayStore } from "@/stores";

const formatTime = (seconds: number): string => {
  if (!isFinite(seconds) || seconds < 0) {
    return "0:00";
  }
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, "0")}`;
};

export const useAudioPlayer = () => {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [duration, setDuration] = useState(0);
  const [currentTime, setCurrentTime] = useState(0);
  const [audioSrc, setAudioSrc] = useState<string>("");
  const endedRef = useRef(false);
  const loadedRef = useRef(false);
  const srcRef = useRef<string>("");
  const canPlayHandlerRef = useRef<(() => void) | null>(null);

  const playing = useAudioPlayStore((s) => s.playing);
  const playlist = useAudioPlayStore((s) => s.playlist);
  const currentAudio = useAudioPlayStore((s) => s.playlist[0] ?? null);
  const play = useAudioPlayStore((s) => s.play);
  const pause = useAudioPlayStore((s) => s.pause);
  const playNext = useAudioPlayStore((s) => s.playNext);
  const setPlaying = useAudioPlayStore((s) => s.setPlaying);

  useEffect(() => {
    const url = currentAudio?.audio?.url;
    const format = currentAudio?.audio?.format ?? "";

    const audio = audioRef.current;
    if (audio && canPlayHandlerRef.current) {
      audio.removeEventListener("canplay", canPlayHandlerRef.current);
      canPlayHandlerRef.current = null;
    }

    if (!url) {
      if (srcRef.current.startsWith("blob:")) {
        URL.revokeObjectURL(srcRef.current);
      }
      srcRef.current = "";
      setAudioSrc("");
      return;
    }

    let cancelled = false;

    (async () => {
      const src = await buildAudioSrc(format, url);
      if (cancelled) return;

      if (srcRef.current.startsWith("blob:")) {
        URL.revokeObjectURL(srcRef.current);
      }
      srcRef.current = src;
      setAudioSrc(src);
    })();

    return () => {
      cancelled = true;
    };
  }, [currentAudio?.id]);

  useEffect(() => {
    if (!audioSrc) return;
    const audio = audioRef.current;
    if (!audio) return;

    if (canPlayHandlerRef.current) {
      audio.removeEventListener("canplay", canPlayHandlerRef.current);
      canPlayHandlerRef.current = null;
    }

    const onCanPlay = () => {
      audio.removeEventListener("canplay", onCanPlay);
      canPlayHandlerRef.current = null;

      if (useAudioPlayStore.getState().playing && !endedRef.current) {
        audio.play().catch(() => useAudioPlayStore.getState().setPlaying(false));
      }
    };

    canPlayHandlerRef.current = onCanPlay;
    audio.addEventListener("canplay", onCanPlay);

    return () => {
      audio.removeEventListener("canplay", onCanPlay);
      canPlayHandlerRef.current = null;
    };
  }, [audioSrc]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !currentAudio) {
      setDuration(0);
      setCurrentTime(0);
      return;
    }

    endedRef.current = false;
    loadedRef.current = false;

    const onLoadedMetadata = () => {
      if (loadedRef.current) return;
      loadedRef.current = true;

      if (isFinite(audio.duration) && audio.duration > 0) {
        setDuration(audio.duration);
      }

      if (currentAudio.progress > 0) {
        audio.currentTime = currentAudio.progress;
        setCurrentTime(currentAudio.progress);
      } else {
        setCurrentTime(0);
      }
    };

    const onTimeUpdate = () => {
      setCurrentTime(audio.currentTime);
      if (isFinite(audio.duration) && audio.duration > 0) {
        setDuration(audio.duration);
      }
    };

    const onDurationChange = () => {
      if (isFinite(audio.duration) && audio.duration > 0) {
        setDuration(audio.duration);
      }
    };

    const onEnded = () => {
      endedRef.current = true;

      const currentPlaylist = useAudioPlayStore.getState().playlist;
      const isLast = currentPlaylist.length <= 1 || currentPlaylist[currentPlaylist.length - 1]?.id === currentAudio.id;

      if (isLast) {
        setPlaying(false);
      } else {
        playNext();
      }
    };

    const onStalled = () => {
      endedRef.current = false;
    };

    audio.addEventListener("loadedmetadata", onLoadedMetadata);
    audio.addEventListener("timeupdate", onTimeUpdate);
    audio.addEventListener("durationchange", onDurationChange);
    audio.addEventListener("ended", onEnded);
    audio.addEventListener("stalled", onStalled);

    if (audio.readyState >= 1) {
      onLoadedMetadata();
    }

    return () => {
      audio.removeEventListener("loadedmetadata", onLoadedMetadata);
      audio.removeEventListener("timeupdate", onTimeUpdate);
      audio.removeEventListener("durationchange", onDurationChange);
      audio.removeEventListener("ended", onEnded);
      audio.removeEventListener("stalled", onStalled);
    };
  }, [currentAudio?.id]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !currentAudio) return;
    if (endedRef.current || audio.ended) return;
    if (!audioSrc || audio.readyState < 2) return;

    const shouldPlay = playing && audio.paused;
    const shouldPause = !playing && !audio.paused;

    if (shouldPlay) {
      audio.play().catch(() => setPlaying(false));
    } else if (shouldPause) {
      audio.pause();
    }
  }, [playing, currentAudio?.id, audioSrc]);

  useEffect(() => {
    return () => {
      if (srcRef.current.startsWith("blob:")) {
        URL.revokeObjectURL(srcRef.current);
      }
    };
  }, []);

  const remainingTime = duration > currentTime ? duration - currentTime : 0;

  return {
    audioRef,
    audioSrc,
    playing,
    currentAudio,
    remainingTime,
    playlist,
    onPlay: () => {
      endedRef.current = false;
      play();
    },
    onPause: () => pause(audioRef.current?.currentTime || 0),
    onNext: () => playNext(),
  };
};

export function SidebarAudioPlayer({
  audioData,
  collapsed,
}: {
  audioData: ReturnType<typeof useAudioPlayer>;
  collapsed?: boolean;
}) {
  const { playing, currentAudio, remainingTime, playlist, onPlay, onPause, onNext } = audioData;
  const PlayIcon = playing ? IconPlayerPause : IconPlayerPlay;
  const isNextDisabled = playlist.length <= 1;

  if (collapsed) {
    return (
      <Stack gap="sm" align="center" className={styles.audioPlayer}>
        <Text className={styles.timeCollapsed}>{formatTime(remainingTime)}</Text>
        <ActionIcon size="lg" radius="xl" variant="filled" color="blue" onClick={playing ? onPause : onPlay}>
          <PlayIcon size={20} />
        </ActionIcon>
        <ActionIcon size="lg" radius="xl" variant="subtle" color="blue" onClick={onNext} disabled={isNextDisabled}>
          <IconPlayerSkipForward size={20} />
        </ActionIcon>
      </Stack>
    );
  }

  return (
    <>
      {currentAudio && (
        <Box className={styles.audioPlayer}>
          <Group justify="space-between" mb="xs">
            <Text className={styles.time}>{formatTime(remainingTime)}</Text>
            <Tooltip label={`${playlist.length} tracks`} position="top">
              <Text size="sm" c="dimmed">
                {playlist.length}
              </Text>
            </Tooltip>
          </Group>
          <Group justify="center" gap="sm">
            <ActionIcon
              size="xl"
              radius="xl"
              variant="filled"
              color="blue"
              onClick={playing ? onPause : onPlay}
              className={styles.playBtn}
            >
              <PlayIcon size={22} />
            </ActionIcon>
            <ActionIcon size="lg" radius="xl" variant="subtle" color="blue" onClick={onNext} disabled={isNextDisabled}>
              <IconPlayerSkipForward size={20} />
            </ActionIcon>
          </Group>
        </Box>
      )}
    </>
  );
}

export function HeaderAudioPlayer({ audioData }: { audioData: ReturnType<typeof useAudioPlayer> }) {
  const { playing, currentAudio, remainingTime, playlist, onPlay, onPause, onNext } = audioData;

  return (
    <>
      {currentAudio && (
        <Group gap="xs" ml="sm">
          <Text className={styles.time}>{formatTime(remainingTime)}</Text>
          <ActionIcon size="lg" radius="xl" variant="filled" color="blue" onClick={playing ? onPause : onPlay}>
            {playing ? <IconPlayerPause size={20} /> : <IconPlayerPlay size={20} />}
          </ActionIcon>
          <ActionIcon
            size="lg"
            radius="xl"
            variant="subtle"
            color="blue"
            onClick={onNext}
            disabled={playlist.length <= 1}
          >
            <IconPlayerSkipForward size={20} />
          </ActionIcon>
        </Group>
      )}
    </>
  );
}
