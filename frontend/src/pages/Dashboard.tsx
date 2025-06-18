import React from 'react';
import { Card, Row, Col, Statistic, Typography, Progress, List, Avatar } from 'antd';
import { 
  UserOutlined, 
  TeamOutlined, 
  BookOutlined, 
  TrophyOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined
} from '@ant-design/icons';
import { useAppSelector } from '../hooks/useAppSelector';

const { Title } = Typography;

const Dashboard: React.FC = () => {
  const { user } = useAppSelector((state) => state.auth);

  // 模拟数据
  const stats = [
    { title: 'EPI实验室总数', value: 12, icon: <TeamOutlined />, color: '#1890ff' },
    { title: '我的申请', value: 3, icon: <BookOutlined />, color: '#52c41a' },
    { title: '待审核', value: 1, icon: <ClockCircleOutlined />, color: '#faad14' },
    { title: '已通过', value: 2, icon: <CheckCircleOutlined />, color: '#52c41a' },
  ];

  const recentActivities = [
    {
      title: '申请了人工智能EPI实验室',
      description: '2024-01-15 14:30',
      avatar: <BookOutlined />,
      color: '#1890ff',
    },
    {
      title: '个人信息更新',
      description: '2024-01-14 10:20',
      avatar: <UserOutlined />,
      color: '#52c41a',
    },
    {
      title: '收到EPI实验室通知',
      description: '2024-01-13 16:45',
      avatar: <TrophyOutlined />,
      color: '#faad14',
    },
  ];

  const applicationStatus = [
    { status: '待审核', count: 1, color: '#faad14' },
    { status: '已通过', count: 2, color: '#52c41a' },
    { status: '已拒绝', count: 0, color: '#ff4d4f' },
  ];

  return (
    <div>
      <Title level={2}>仪表盘</Title>
      
      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {stats.map((stat, index) => (
          <Col xs={24} sm={12} lg={6} key={index}>
            <Card>
              <Statistic
                title={stat.title}
                value={stat.value}
                prefix={React.cloneElement(stat.icon, { style: { color: stat.color } })}
                valueStyle={{ color: stat.color }}
              />
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]}>
        {/* 申请状态 */}
        <Col xs={24} lg={12}>
          <Card title="申请状态" style={{ height: '100%' }}>
            {applicationStatus.map((item, index) => (
              <div key={index} style={{ marginBottom: 16 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <span>{item.status}</span>
                  <span style={{ color: item.color, fontWeight: 'bold' }}>
                    {item.count}
                  </span>
                </div>
                <Progress 
                  percent={item.count * 33.33} 
                  strokeColor={item.color}
                  showInfo={false}
                />
              </div>
            ))}
          </Card>
        </Col>

        {/* 最近活动 */}
        <Col xs={24} lg={12}>
          <Card title="最近活动" style={{ height: '100%' }}>
            <List
              itemLayout="horizontal"
              dataSource={recentActivities}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={
                      <Avatar 
                        icon={item.avatar} 
                        style={{ backgroundColor: item.color }}
                      />
                    }
                    title={item.title}
                    description={item.description}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>

      {/* 快速统计 */}
      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={8}>
          <Card title="申请成功率">
            <div style={{ textAlign: 'center' }}>
              <Progress
                type="circle"
                percent={66.7}
                format={(percent) => `${percent}%`}
                strokeColor="#52c41a"
              />
              <div style={{ marginTop: 16 }}>
                <Text type="secondary">2/3 申请已通过</Text>
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="活跃EPI实验室">
            <div style={{ textAlign: 'center' }}>
              <Statistic
                title="本月新增"
                value={3}
                suffix="个"
                valueStyle={{ color: '#1890ff' }}
              />
            </div>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="系统通知">
            <div style={{ textAlign: 'center' }}>
              <Statistic
                title="未读消息"
                value={2}
                suffix="条"
                valueStyle={{ color: '#faad14' }}
              />
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard; 