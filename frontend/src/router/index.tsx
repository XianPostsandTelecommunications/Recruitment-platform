import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import { useAppSelector } from '../hooks/useAppSelector';
import MainLayout from '../layouts/MainLayout';
import Login from '../pages/Login';

import Home from '../pages/Home';
import Dashboard from '../pages/Dashboard';
import InterviewApplication from '../pages/InterviewApplication';
import AdminDashboard from '../pages/AdminDashboard';

// 懒加载组件
const LazyComponent = React.lazy(() => Promise.resolve({ default: () => <div>Loading...</div> }));

// 受保护的路由组件
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, loading } = useAppSelector(state => state.auth);
  
  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" />
      </div>
    );
  }
  
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

// 公共路由组件（已登录用户重定向到首页）
const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAppSelector(state => state.auth);
  
  return isAuthenticated ? <Navigate to="/" replace /> : <>{children}</>;
};

// 简化的测试组件
const TestHome: React.FC = () => {
  return (
    <div style={{ padding: '20px' }}>
      <h1>EPI实验室招新平台</h1>
      <p>欢迎使用EPI实验室招新平台！</p>
      <div style={{ 
        padding: '20px', 
        backgroundColor: '#f0f2f5', 
        borderRadius: '8px',
        marginTop: '20px'
      }}>
        <h3>功能特性：</h3>
        <ul>
          <li>面试申请提交</li>
          <li>申请状态管理</li>
          <li>面试时间安排</li>
          <li>数据统计分析</li>
        </ul>
      </div>
      <div style={{ marginTop: '20px' }}>
        <a href="/apply" style={{
          padding: '10px 20px',
          backgroundColor: '#1890ff',
          color: 'white',
          textDecoration: 'none',
          borderRadius: '6px',
          marginRight: '10px'
        }}>
          申请面试
        </a>
        <a href="/admin" style={{
          padding: '10px 20px',
          backgroundColor: '#52c41a',
          color: 'white',
          textDecoration: 'none',
          borderRadius: '6px'
        }}>
          管理后台
        </a>
      </div>
    </div>
  );
};

const AppRouter: React.FC = () => {
  return (
    <BrowserRouter>
      <React.Suspense fallback={
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
          <Spin size="large" />
        </div>
      }>
        <Routes>
          {/* 首页重定向到报名EPI实验室页面 */}
          <Route path="/" element={<Navigate to="/apply" replace />} />
          
          {/* 面试申请页面 */}
          <Route path="/apply" element={<InterviewApplication />} />
          
          {/* 管理员后台，登录守卫 */}
          <Route path="/admin" element={
            <ProtectedRoute>
              <AdminDashboard />
            </ProtectedRoute>
          } />
          
          {/* 公共路由 */}
          <Route path="/login" element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          } />

          
          {/* 受保护的路由 */}
          <Route path="/dashboard" element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          } />
          <Route path="/profile" element={
            <ProtectedRoute>
              <LazyComponent />
            </ProtectedRoute>
          } />
          <Route path="/labs" element={
            <ProtectedRoute>
              <LazyComponent />
            </ProtectedRoute>
          } />
          <Route path="/applications" element={
            <ProtectedRoute>
              <LazyComponent />
            </ProtectedRoute>
          } />
          <Route path="/notifications" element={
            <ProtectedRoute>
              <LazyComponent />
            </ProtectedRoute>
          } />
          {/* 404页面 */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </React.Suspense>
    </BrowserRouter>
  );
};

export default AppRouter; 