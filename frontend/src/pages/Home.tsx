import React from 'react';
import { Card, Row, Col, Button, Space, Statistic, message } from 'antd';
import { UserOutlined, TeamOutlined, FileTextOutlined, BellOutlined } from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../hooks/useAppSelector';
import { loginAsync, registerAsync } from '../store/slices/authSlice';

const Home: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);

  const features = [
    {
      icon: <TeamOutlined style={{ fontSize: '32px', color: '#1890ff' }} />,
      title: 'EPI实验室管理',
      description: '浏览和管理EPI实验室信息，了解招新要求',
    },
    {
      icon: <TeamOutlined style={{ fontSize: '32px', color: '#52c41a' }} />,
      title: '申请管理',
      description: '提交申请，跟踪申请状态，查看反馈',
    },
    {
      icon: <TeamOutlined style={{ fontSize: '32px', color: '#faad14' }} />,
      title: '成果展示',
      description: '展示EPI实验室研究成果和项目经验',
    },
  ];

  // 测试登录功能
  const handleTestLogin = async () => {
    try {
      await dispatch(loginAsync({
        email: 'admin@example.com',
        password: '123456'
      })).unwrap();
      message.success('测试登录成功！');
    } catch (error) {
      message.error('测试登录失败：' + (error as Error).message);
    }
  };

  // 测试注册功能
  const handleTestRegister = async () => {
    try {
      await dispatch(registerAsync({
        username: 'testuser',
        email: 'testuser@example.com',
        password: '123456',
        role: 'student'
      })).unwrap();
      message.success('测试注册成功！');
    } catch (error) {
      message.error('测试注册失败：' + (error as Error).message);
    }
  };

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ marginBottom: '24px' }}>
        <h1>欢迎使用EPI实验室招新平台</h1>
        <p>这是一个现代化的EPI实验室招新管理系统，提供完整的用户管理、EPI实验室管理、申请管理等功能。</p>
        
        {/* 测试功能区域 */}
        <Card title="功能测试" style={{ marginBottom: '24px' }}>
          <Space>
            <Button type="primary" onClick={handleTestLogin}>
              测试登录 (admin@example.com / 123456)
            </Button>
            <Button onClick={handleTestRegister}>
              测试注册新用户
            </Button>
          </Space>
          {isAuthenticated && user && (
            <div style={{ marginTop: '16px', padding: '12px', backgroundColor: '#f6ffed', borderRadius: '6px' }}>
              <p><strong>当前登录用户：</strong></p>
              <p>用户名：{user.username}</p>
              <p>邮箱：{user.email}</p>
              <p>角色：{user.role}</p>
            </div>
          )}
        </Card>
      </div>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={1128}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="EPI实验室数量"
              value={15}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="申请数量"
              value={93}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="未读通知"
              value={5}
              prefix={<BellOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
        <Col xs={24} lg={12}>
          <Card title="快速操作">
            <Space direction="vertical" style={{ width: '100%' }}>
              <Button type="primary" block>
                查看EPI实验室列表
              </Button>
              <Button block>
                提交申请
              </Button>
              <Button block>
                查看我的申请
              </Button>
              <Button block>
                查看通知
              </Button>
            </Space>
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="系统公告">
            <div>
              <p><strong>2024年春季招新开始</strong></p>
              <p>EPI实验室已开始接受2024年春季招新申请，请有意向的同学及时提交申请。</p>
              <p style={{ color: '#666', fontSize: '12px' }}>发布时间：2024-01-15</p>
            </div>
            <div style={{ marginTop: '16px' }}>
              <p><strong>系统维护通知</strong></p>
              <p>系统将于本周日凌晨2:00-4:00进行维护，期间可能无法正常访问。</p>
              <p style={{ color: '#666', fontSize: '12px' }}>发布时间：2024-01-10</p>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Home; 