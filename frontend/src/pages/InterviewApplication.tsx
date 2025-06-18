import React, { useState } from 'react';
import {
  Form,
  Input,
  Button,
  Card,
  Typography,
  message,
  Space,
  Row,
  Col,
  Select,
  DatePicker,
  TimePicker,
  Divider,
  Alert,
} from 'antd';
import {
  UserOutlined,
  MailOutlined,
  PhoneOutlined,
  IdcardOutlined,
  BookOutlined,
  CalendarOutlined,
  SafetyOutlined,
  SendOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;

interface ApplicationForm {
  name: string;
  email: string;
  phone: string;
  student_id: string;
  major: string;
  grade: string;
  interview_date: Dayjs;
  interview_time: Dayjs;
  verification_code: string;
}

const InterviewApplication: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [codeLoading, setCodeLoading] = useState(false);
  const [codeSent, setCodeSent] = useState(false);
  const [countdown, setCountdown] = useState(0);

  // 发送验证码
  const handleSendCode = async () => {
    const email = form.getFieldValue('email');
    if (!email) {
      message.error('请先输入邮箱地址');
      return;
    }

    setCodeLoading(true);
    try {
      const response = await fetch('/api/v1/send-code', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();
      if (data.code === 200) {
        message.success({
          content: '📧 验证码已发送到邮箱，请注意查收',
          duration: 4,
          style: {
            fontSize: '14px',
          },
        });
        setCodeSent(true);
        setCountdown(60);
        const timer = setInterval(() => {
          setCountdown((prev) => {
            if (prev <= 1) {
              clearInterval(timer);
              return 0;
            }
            return prev - 1;
          });
        }, 1000);
      } else {
        message.error({
          content: data.message || '❌ 发送验证码失败，请稍后重试',
          duration: 4,
          style: {
            fontSize: '14px',
          },
        });
      }
    } catch (error) {
      message.error({
        content: '❌ 网络错误，请稍后重试',
        duration: 4,
        style: {
          fontSize: '14px',
        },
      });
    } finally {
      setCodeLoading(false);
    }
  };

  // 提交申请
  const handleSubmit = async (values: ApplicationForm) => {
    setLoading(true);
    try {
      const interviewDateTime = values.interview_date
        .hour(values.interview_time.hour())
        .minute(values.interview_time.minute())
        .format('YYYY-MM-DD HH:mm');

      const response = await fetch('/api/v1/apply', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: values.name,
          email: values.email,
          phone: values.phone,
          student_id: values.student_id,
          major: values.major,
          grade: values.grade,
          interview_time: interviewDateTime,
          verification_code: values.verification_code,
        }),
      });

      const data = await response.json();
      if (data.code === 200) {
        message.success({
          content: '🎉 申请提交成功！我们会尽快联系您安排面试。',
          duration: 5,
          style: {
            fontSize: '16px',
            fontWeight: 'bold',
          },
        });
        form.resetFields();
        setCodeSent(false);
        setCountdown(0);
      } else {
        message.error({
          content: data.message || '❌ 申请提交失败，请检查信息后重试',
          duration: 4,
          style: {
            fontSize: '14px',
          },
        });
      }
    } catch (error) {
      message.error({
        content: '❌ 网络错误，请稍后重试',
        duration: 4,
        style: {
          fontSize: '14px',
        },
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '20px',
      }}
    >
      <div style={{ maxWidth: '800px', margin: '0 auto' }}>
        <Card
          style={{
            borderRadius: '16px',
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
            backdropFilter: 'blur(10px)',
            backgroundColor: 'rgba(255, 255, 255, 0.95)',
          }}
        >
          {/* 头部信息 */}
          <div style={{ textAlign: 'center', marginBottom: '32px' }}>
            <Title level={2} style={{ color: '#1890ff', marginBottom: '8px' }}>
              🧪 实验室面试申请
            </Title>
            <Paragraph type="secondary" style={{ fontSize: '16px' }}>
              欢迎加入我们的实验室！请填写以下信息完成面试申请
            </Paragraph>
          </div>

          {/* 申请说明 */}
          <Alert
            message="申请说明"
            description="请确保填写的信息真实准确，我们会在收到申请后尽快安排面试时间。面试将通过邮件或电话通知。"
            type="info"
            showIcon
            style={{ marginBottom: '24px' }}
          />

          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            size="large"
            initialValues={{
              grade: '2024',
              interview_date: dayjs(),
              interview_time: dayjs().hour(14).minute(0),
            }}
          >
            <Row gutter={[16, 0]}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="name"
                  label="姓名"
                  rules={[{ required: true, message: '请输入姓名' }]}
                >
                  <Input
                    prefix={<UserOutlined />}
                    placeholder="请输入您的真实姓名"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="email"
                  label="邮箱"
                  rules={[
                    { required: true, message: '请输入邮箱' },
                    { type: 'email', message: '请输入有效的邮箱地址' },
                  ]}
                >
                  <Input
                    prefix={<MailOutlined />}
                    placeholder="请输入邮箱地址"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={[16, 0]}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="phone"
                  label="手机号码"
                  rules={[
                    { required: true, message: '请输入手机号码' },
                    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' },
                  ]}
                >
                  <Input
                    prefix={<PhoneOutlined />}
                    placeholder="请输入手机号码"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="student_id"
                  label="学号"
                  rules={[{ required: true, message: '请输入学号' }]}
                >
                  <Input
                    prefix={<IdcardOutlined />}
                    placeholder="请输入学号"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={[16, 0]}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="major"
                  label="专业"
                  rules={[{ required: true, message: '请选择专业' }]}
                >
                  <Select placeholder="请选择专业">
                    <Option value="计算机科学与技术">计算机科学与技术</Option>
                    <Option value="软件工程">软件工程</Option>
                    <Option value="人工智能">人工智能</Option>
                    <Option value="数据科学">数据科学</Option>
                    <Option value="网络工程">网络工程</Option>
                    <Option value="信息安全">信息安全</Option>
                    <Option value="其他">其他</Option>
                  </Select>
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="grade"
                  label="年级"
                  rules={[{ required: true, message: '请选择年级' }]}
                >
                  <Select placeholder="请选择年级">
                    <Option value="2021">2021级</Option>
                    <Option value="2022">2022级</Option>
                    <Option value="2023">2023级</Option>
                    <Option value="2024">2024级</Option>
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            <Divider orientation="left">面试时间安排</Divider>

            <Row gutter={[16, 0]}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="interview_date"
                  label="面试日期"
                  rules={[{ required: true, message: '请选择面试日期' }]}
                >
                  <DatePicker
                    style={{ width: '100%' }}
                    placeholder="选择面试日期"
                    disabledDate={(current) => current && current < dayjs().startOf('day')}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="interview_time"
                  label="面试时间"
                  rules={[{ required: true, message: '请选择面试时间' }]}
                >
                  <TimePicker
                    style={{ width: '100%' }}
                    placeholder="选择面试时间"
                    format="HH:mm"
                    minuteStep={30}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Divider orientation="left">邮箱验证</Divider>

            <Row gutter={[16, 0]}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="verification_code"
                  label="验证码"
                  rules={[{ required: true, message: '请输入验证码' }]}
                >
                  <Input
                    prefix={<SafetyOutlined />}
                    placeholder="请输入6位验证码"
                    maxLength={6}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item label=" " style={{ marginTop: '29px' }}>
                  <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSendCode}
                    loading={codeLoading}
                    disabled={countdown > 0}
                    style={{ width: '100%' }}
                  >
                    {countdown > 0 ? `${countdown}s后重发` : '发送验证码'}
                  </Button>
                </Form.Item>
              </Col>
            </Row>

            <Form.Item style={{ marginTop: '32px' }}>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                size="large"
                style={{
                  width: '100%',
                  height: '48px',
                  fontSize: '16px',
                  borderRadius: '8px',
                }}
              >
                提交面试申请
              </Button>
            </Form.Item>
          </Form>

          {/* 底部说明 */}
          <div style={{ textAlign: 'center', marginTop: '24px' }}>
            <Text type="secondary">
              如有疑问，请联系：lab@example.com | 电话：123-456-7890
            </Text>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default InterviewApplication; 