import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import type { UserResponse, UserRegisterRequest } from '../../types/auth';
import * as authApi from '../../api/auth';

// 异步action：登录
export const loginAsync = createAsyncThunk(
  'auth/login',
  async (credentials: { email: string; password: string }) => {
    const response = await authApi.login(credentials);
    const data = response.data as any;
    if (data && data.user && data.token) {
      localStorage.setItem('token', data.token);
      return data.user;
    }
    throw new Error('登录失败');
  }
);

// 异步action：注册
export const registerAsync = createAsyncThunk(
  'auth/register',
  async (userData: UserRegisterRequest) => {
    const response = await authApi.register(userData);
    const data = response.data as any;
    if (data && data.user && data.token) {
      localStorage.setItem('token', data.token);
      return data.user;
    }
    throw new Error('注册失败');
  }
);

// 异步action：获取用户信息
export const getProfileAsync = createAsyncThunk(
  'auth/getProfile',
  async () => {
    const response = await authApi.getProfile();
    const data = response.data as any;
    if (data) {
      return data;
    }
    throw new Error('获取用户信息失败');
  }
);

// 异步action：更新用户信息
export const updateProfileAsync = createAsyncThunk(
  'auth/updateProfile',
  async (userData: { username?: string; phone?: string; student_id?: string; major?: string; grade?: string; avatar?: string }) => {
    const response = await authApi.updateProfile(userData);
    const data = response.data as any;
    if (data) {
      return data;
    }
    throw new Error('更新用户信息失败');
  }
);

// 异步action：登出
export const logoutAsync = createAsyncThunk(
  'auth/logout',
  async () => {
    await authApi.logout();
    localStorage.removeItem('token');
  }
);

interface AuthState {
  user: UserResponse | null;
  token: string | null;
  loading: boolean;
  error: string | null;
  isAuthenticated: boolean;
}

const initialState: AuthState = {
  user: null,
  token: localStorage.getItem('token'),
  loading: false,
  error: null,
  isAuthenticated: !!localStorage.getItem('token'),
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setUser: (state, action: PayloadAction<UserResponse>) => {
      state.user = action.payload;
      state.isAuthenticated = true;
    },
  },
  extraReducers: (builder) => {
    // 登录
    builder
      .addCase(loginAsync.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(loginAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
        state.token = localStorage.getItem('token');
        state.isAuthenticated = true;
      })
      .addCase(loginAsync.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || '登录失败';
        state.isAuthenticated = false;
      });

    // 注册
    builder
      .addCase(registerAsync.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(registerAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
        state.token = localStorage.getItem('token');
        state.isAuthenticated = true;
      })
      .addCase(registerAsync.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || '注册失败';
        state.isAuthenticated = false;
      });

    // 获取用户信息
    builder
      .addCase(getProfileAsync.pending, (state) => {
        state.loading = true;
      })
      .addCase(getProfileAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
        state.isAuthenticated = true;
      })
      .addCase(getProfileAsync.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || '获取用户信息失败';
        state.isAuthenticated = false;
      });

    // 更新用户信息
    builder
      .addCase(updateProfileAsync.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateProfileAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
      })
      .addCase(updateProfileAsync.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || '更新用户信息失败';
      });

    // 登出
    builder
      .addCase(logoutAsync.fulfilled, (state) => {
        state.user = null;
        state.token = null;
        state.error = null;
        state.isAuthenticated = false;
      });
  },
});

export const { clearError, setUser } = authSlice.actions;
export default authSlice.reducer; 