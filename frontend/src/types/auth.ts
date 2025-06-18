// 用户注册请求
export interface UserRegisterRequest {
  username: string;
  email: string;
  password: string;
  role?: 'student' | 'admin';
}

// 用户登录请求
export interface UserLoginRequest {
  email: string;
  password: string;
}

// 用户更新请求
export interface UserUpdateRequest {
  username?: string;
  phone?: string;
  student_id?: string;
  major?: string;
  grade?: string;
  avatar?: string;
}

// 用户响应
export interface UserResponse {
  id: number;
  username: string;
  email: string;
  role: string;
  avatar?: string;
  phone?: string;
  student_id?: string;
  major?: string;
  grade?: string;
  status: string;
  created_at: string;
}

// 登录响应
export interface LoginResponse {
  user: UserResponse;
  token: string;
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

// 刷新令牌响应
export interface RefreshTokenResponse {
  token: string;
} 