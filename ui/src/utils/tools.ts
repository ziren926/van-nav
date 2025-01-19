import { useRef, useCallback } from "react";

export interface Tool {
  id: string;
  name: string;
  url: string;
  desc: string;
  logo: string;
  catelog: string;
  content?: string;  // 新增字段：详细内容
  updatedAt?: string; // 新增字段：更新时间
  createdAt?: string; // 新增字段：创建时间
}

// 如果需要，可以添加工具相关的辅助函数
export const formatToolDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

export const useDebounce = (fn: any, delay: number) => {
  const { current } = useRef<{ time: any }>({ time: null });
  return useCallback(
    (...args: any[]) => {
      if (current.time) {
        clearTimeout(current.time);
        current.time = null;
      }
      current.time = setTimeout(() => {
        fn.apply(this, args);
        clearTimeout(current.time);
        current.time = null;
      }, delay);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [current, delay]
  );
};