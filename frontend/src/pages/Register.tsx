import React from 'react';
import { Form, Input, Button, Card, Typography, message, Space, Select } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/useAppSelector';
import { registerAsync, clearError } from '../store/slices/authSlice';
import type { UserRegisterRequest } from '../types/auth';

const { Title, Text } = Typography;
const { Option } = Select;

const Register: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { loading, error } = useAppSelector((state) => state.auth);

  const onFinish = async (values: UserRegisterRequest) => {
    try {
      await dispatch(registerAsync(values)).unwrap();
      message.success('注册成功');
      navigate('/');
    } catch (err) {
      // 错误已在 Redux 中处理
    }
  };

  // 清除错误信息
  React.useEffect(() => {
    if (error) {
      message.error(error);
      dispatch(clearError());
    }
  }, [error, dispatch]);

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <Card
        style={{
          width: 400,
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
        }}
      >
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Title level={2} style={{ marginBottom: 8 }}>
            EPI实验室招新平台
          </Title>
          <Text type="secondary">创建您的账户</Text>
        </div>

        <Form
          name="register"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="username"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, max: 20, message: '用户名长度必须在3-20位之间' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名"
            />
          </Form.Item>

          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱地址' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input
              prefix={<MailOutlined />}
              placeholder="邮箱地址"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, max: 20, message: '密码长度必须在6-20位之间' },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
            />
          </Form.Item>

          <Form.Item
            name="role"
            initialValue="student"
          >
            <Select placeholder="选择角色">
              <Option value="student">学生</Option>
              <Option value="admin">管理员</Option>
            </Select>
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
            >
              注册
            </Button>
          </Form.Item>

          <div style={{ textAlign: 'center' }}>
            <Space>
              <Text type="secondary">已有账户？</Text>
              <Link to="/login">
                <Text>立即登录</Text>
              </Link>
            </Space>
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default Register; 