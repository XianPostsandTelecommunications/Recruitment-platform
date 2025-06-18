// 通用响应格式
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  total?: number;
  page?: number;
  size?: number;
}

// 分页参数
export interface PaginationParams {
  page?: number;
  size?: number;
}

// 分页响应
export interface PaginatedResponse<T> {
  total: number;
  page: number;
  size: number;
  list: T[];
}

// 字符串数组类型（用于标签等）
export type StringSlice = string[];

// 文件上传响应
export interface UploadResponse {
  url: string;
  filename: string;
  size: number;
}

// 统计信息
export interface Stats {
  total: number;
  [key: string]: number;
} 