import { getTask, PodcastTaskResult, TaskStageName, TaskStageStatus } from "@/services";
import { create } from "zustand";
import { useAudioPlayStore } from "./audioPlayStore";

const maxRetries = 10;

interface QueueItem {
  batchNo: string;
  retryCount: number;
}

interface TaskPollingStore {
  queue: QueueItem[];
  timerId: NodeJS.Timeout | null;
  isPolling: boolean;

  addToQueue(batchNo: string): void;
  removeFromQueue: (batchNo: string) => void;
  startPolling(): void;
  stopPolling(): void;
  incrRetry(batchNo: string): void;
  checkQueue: () => Promise<void>;
  clearQueue(): void;
}

export const useTaskPollingStore = create<TaskPollingStore>((set, get) => ({
  queue: [],
  timerId: null,
  isPolling: false,

  addToQueue: (batchNo: string) => {
    const { queue, startPolling } = get();

    if (queue.some((item) => item.batchNo === batchNo)) return;

    set({ queue: [...queue, { batchNo, retryCount: 0 }] });

    if (!get().isPolling) startPolling();
  },
  removeFromQueue: (batchNo: string) => {
    const { queue, stopPolling } = get();

    const newQueue = queue.filter((item) => item.batchNo !== batchNo);

    set({ queue: newQueue });

    if (queue.length === 0) stopPolling();
  },
  startPolling: () => {
    const { timerId, checkQueue } = get();
    if (timerId) return;
    set({ isPolling: true });

    const id = setInterval(() => {
      checkQueue();
    }, 60000);

    set({ timerId: id });
  },
  stopPolling: () => {
    const { timerId } = get();

    if (timerId) {
      clearInterval(timerId);
      set({ timerId: null, isPolling: false });
    }
  },
  incrRetry: (batchNo: string) => {
    const { queue, removeFromQueue } = get();

    const item = queue.find((item) => item.batchNo === batchNo);

    if (!item) return;

    if (item.retryCount >= maxRetries) {
      removeFromQueue(batchNo);
      return;
    }

    const newQueue = queue.map((item) => {
      if (item.batchNo === batchNo) return { ...item, retryCount: item.retryCount + 1 };

      return item;
    });

    set({ queue: newQueue });
  },
  checkQueue: async () => {
    const { queue, removeFromQueue, incrRetry } = get();
    if (queue.length === 0) {
      return;
    }

    const checkPromises = queue.map(async (item) => {
      const task = await getTask(item.batchNo);
      if (!task) {
        incrRetry(item.batchNo);
        return;
      }

      if (task.result == PodcastTaskResult.Failed) {
        removeFromQueue(item.batchNo);
        return;
      }

      const stage = task.stages?.findLast((state) => state.stage === TaskStageName.TTS);
      if (!stage || stage.status === TaskStageStatus.Processing) {
        incrRetry(item.batchNo);
        return;
      }

      if (stage.status == TaskStageStatus.Failed) {
        removeFromQueue(item.batchNo);
        return;
      }

      if (stage.audio && stage.audio.data) {
        useAudioPlayStore.getState().addAudio(stage.id, stage.audio);
        removeFromQueue(item.batchNo);
      }

      incrRetry(item.batchNo);
    });

    await Promise.all(checkPromises);
  },
  clearQueue: () => {
    const { stopPolling } = get();
    stopPolling();
    set({ queue: [] });
  },
}));
