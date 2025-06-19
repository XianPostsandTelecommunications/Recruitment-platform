import instance from './instance';
import type { 
  UserLoginRequest, 
  UserUpdateRequest, 
  UserResponse, 
  LoginResponse, 
  ChangePasswordRequest, 
  RefreshTokenResponse 
} from '../types/auth';
import type { ApiResponse } from '../types/common';



// 用户登录
export const login = (data: UserLoginRequest) =>
  instance.post<ApiResponse<LoginResponse>>('/auth/login', data);

// 获取用户信息
export const getProfile = () =>
  instance.get<ApiResponse<UserResponse>>('/auth/profile');

// 更新用户信息
export const updateProfile = (data: UserUpdateRequest) =>
  instance.put<ApiResponse<UserResponse>>('/auth/profile', data);

// 修改密码
export const changePassword = (data: ChangePasswordRequest) =>
  instance.post<ApiResponse<null>>('/auth/change-password', data);

// 刷新令牌
export const refreshToken = () =>
  instance.post<ApiResponse<RefreshTokenResponse>>('/auth/refresh');

// 用户登出
export const logout = () =>
  instance.post<ApiResponse<null>>('/auth/logout'); 