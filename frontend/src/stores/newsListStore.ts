// stores/useNewsListStore.ts
import { create } from "zustand";
import { NewsDetail } from "@/services";
import { httpx } from "wailsjs/go/models";

export interface SearchFormValues {
  source: string;
  topic: string;
  publishDate: string;
}

interface NewsListState {
  currentRecordId: number;
  pagination: httpx.Pagination;
  searchForm: SearchFormValues;
  newsList: NewsDetail[];
  loading: boolean;
}

interface NewsListActions {
  setCurrentRecordId: (recordId: number) => void;
  updatePagination: (pagination: Partial<httpx.Pagination>) => void;
  updateSearchForm: (searchForm: Partial<SearchFormValues>) => void;
  setNewsList: (newsList: NewsDetail[]) => void;
  setLoading: (loading: boolean) => void;
  searchWithReset: (searchForm: Partial<SearchFormValues>) => void;
  resetState: () => void;
}

const initialState: NewsListState = {
  currentRecordId: 0,
  pagination: { page: 1, limit: 20, total: 0 },
  searchForm: { source: "", topic: "", publishDate: "" },
  newsList: [],
  loading: true,
};

export const useNewsListStore = create<NewsListState & NewsListActions>()((set, get) => ({
  ...initialState,

  setCurrentRecordId: (recordId: number) => {
    const currentState = get();
    if (currentState.currentRecordId !== recordId) {
      set({
        currentRecordId: recordId,
        pagination: { page: 1, limit: 20, total: 0 },
        searchForm: { source: "", topic: "", publishDate: "" },
        newsList: [],
        loading: true,
      });
    }
  },

  updatePagination: (pagination: Partial<httpx.Pagination>) => {
    set((state) => ({
      pagination: { ...state.pagination, ...pagination },
    }));
  },

  updateSearchForm: (searchForm: Partial<SearchFormValues>) => {
    set((state) => ({
      searchForm: { ...state.searchForm, ...searchForm },
    }));
  },

  setNewsList: (newsList: NewsDetail[]) => {
    set({ newsList });
  },

  setLoading: (loading: boolean) => {
    set({ loading });
  },

  searchWithReset: (searchForm: Partial<SearchFormValues>) => {
    set((state) => ({
      searchForm: { ...state.searchForm, ...searchForm },
      pagination: { ...state.pagination, page: 1 },
      loading: true,
    }));
  },

  resetState: () => {
    set(initialState);
  },
}));
