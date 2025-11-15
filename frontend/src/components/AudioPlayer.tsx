import { ActionIcon, Box, Group, Stack, Text, Tooltip } from "@mantine/core";
import { IconPlayerPause, IconPlayerPlay, IconPlayerSkipForward } from "@tabler/icons-react";
import { useEffect, useRef, useState } from "react";

import styles from "@/assets/styles/appLayout.module.css";
import { useAudioPlayStore } from "@/stores";

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

  const store = useAudioPlayStore();
  const currentAudio = store.getCurrentAudio();

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !currentAudio) {
      setDuration(0);
      setCurrentTime(0);
      return;
    }

    const onLoadedMetadata = () => {
      const audioDuration = audio.duration;
      if (isFinite(audioDuration)) {
        setDuration(audioDuration);
        if (currentAudio.progress > 0 && currentAudio.progress < audioDuration) {
          audio.currentTime = currentAudio.progress;
          setCurrentTime(currentAudio.progress);
        } else {
          setCurrentTime(0);
        }
      }
    };

    const onTimeUpdate = () => {
      setCurrentTime(audio.currentTime);
      if (isFinite(audio.duration)) {
        setDuration(audio.duration);
      }
    };

    const onEnded = () => {
      store.playNext();
    };

    audio.addEventListener("loadedmetadata", onLoadedMetadata);
    audio.addEventListener("timeupdate", onTimeUpdate);
    audio.addEventListener("ended", onEnded);

    if (audio.readyState >= 1) {
      onLoadedMetadata();
    }

    return () => {
      audio.removeEventListener("loadedmetadata", onLoadedMetadata);
      audio.removeEventListener("timeupdate", onTimeUpdate);
      audio.removeEventListener("ended", onEnded);
    };
  }, [currentAudio?.id, store]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !currentAudio) return;

    const shouldPlay = store.playing && audio.paused;
    const shouldPause = !store.playing && !audio.paused;

    if (shouldPlay) {
      audio.play().catch(() => {
        store.setPlaying(false);
      });
    } else if (shouldPause) {
      audio.pause();
    }
  }, [store.playing, currentAudio?.id, store]);

  const remainingTime = duration > currentTime ? duration - currentTime : 0;

  return {
    audioRef,
    playing: store.playing,
    currentAudio,
    remainingTime,
    playlist: store.playlist,
    onPlay: () => store.play(),
    onPause: () => store.pause(audioRef.current?.currentTime || 0),
    onNext: () => store.playNext(),
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
