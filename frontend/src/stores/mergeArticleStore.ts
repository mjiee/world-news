import { create } from "zustand";
import { TaskStage } from "@/services";

interface MergeArticleStage {
  stages: TaskStage[];
}

interface MergeArticleActions {
  addStage: (data: TaskStage) => void;
  removeStage: (stageId: number) => void;
  resetStage: () => void;
}

export const useMergeArticleStore = create<MergeArticleStage & MergeArticleActions>()((set, get) => ({
  stages: [],
  addStage: (data) => {
    set((state) => ({
      stages: [...state.stages.filter((stage) => stage.id !== data.id), data],
    }));
  },
  removeStage: (stageId) => {
    set((state) => ({
      stages: state.stages.filter((stage) => stage.id !== stageId),
    }));
  },
  resetStage: () => {
    set(() => ({ stages: [] }));
  },
}));
