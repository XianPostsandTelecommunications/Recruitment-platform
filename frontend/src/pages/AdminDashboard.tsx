import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  message,
  Typography,
  Row,
  Col,
  Statistic,
  Tooltip,
  Popconfirm,
  DatePicker,
  TimePicker,
  Alert,
} from 'antd';
import {
  EyeOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  UserOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { Option } = Select;
const { TextArea } = Input;

interface Applicant {
  id: number;
  name: string;
  email: string;
  phone: string;
  student_id: string;
  major: string;
  grade: string;
  interview_time: string;
  status: string;
  first_remark: string;
  second_remark: string;
  third_remark: string;
  created_at: string;
  updated_at: string;
}

const AdminDashboard: React.FC = () => {
  const [applicants, setApplicants] = useState<Applicant[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedApplicant, setSelectedApplicant] = useState<Applicant | null>(null);
  const [form] = Form.useForm();
  const [addModalVisible, setAddModalVisible] = useState(false);
  const [addForm] = Form.useForm();
  const [addLoading, setAddLoading] = useState(false);

  // 获取申请者列表
  const fetchApplicants = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/v1/applicants');
      const data = await response.json();
      if (data.code === 200) {
        setApplicants(data.data || []);
      } else {
        message.error(data.message || '获取数据失败');
      }
    } catch (error) {
      message.error('网络错误，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplicants();
  }, []);

  // 更新申请状态
  const handleUpdateStatus = async (values: any) => {
    if (!selectedApplicant) return;

    try {
      const response = await fetch(`/api/v1/applicants/${selectedApplicant.id}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      });

      const data = await response.json();
      if (data.code === 200) {
        message.success('状态更新成功');
        setModalVisible(false);
        form.resetFields();
        fetchApplicants();
      } else {
        message.error(data.message || '更新失败');
      }
    } catch (error) {
      message.error('网络错误，请稍后重试');
    }
  };

  // 删除申请
  const handleDelete = async (id: number) => {
    try {
      const response = await fetch(`/api/v1/applicants/${id}`, {
        method: 'DELETE',
      });

      const data = await response.json();
      if (data.code === 200) {
        message.success('删除成功');
        fetchApplicants();
      } else {
        message.error(data.message || '删除失败');
      }
    } catch (error) {
      message.error('网络错误，请稍后重试');
    }
  };

  // 打开编辑模态框
  const showEditModal = (applicant: Applicant) => {
    setSelectedApplicant(applicant);
    form.setFieldsValue({
      status: applicant.status,
      first_remark: applicant.first_remark,
      second_remark: applicant.second_remark,
      third_remark: applicant.third_remark,
    });
    setModalVisible(true);
  };

  // 获取状态标签
  const getStatusTag = (status: string) => {
    const statusMap = {
      pending: { color: 'orange', text: '待面试', icon: <ClockCircleOutlined /> },
      first_pass: { color: 'blue', text: '一面通过', icon: <CheckCircleOutlined /> },
      second_pass: { color: 'purple', text: '二面通过', icon: <CheckCircleOutlined /> },
      passed: { color: 'green', text: '通过', icon: <CheckCircleOutlined /> },
      rejected: { color: 'red', text: '未通过', icon: <CloseCircleOutlined /> },
    };

    const config = statusMap[status as keyof typeof statusMap] || statusMap.pending;
    return (
      <Tag color={config.color} icon={config.icon}>
        {config.text}
      </Tag>
    );
  };

  // 统计信息
  const getStats = () => {
    const total = applicants.length;
    const pending = applicants.filter(a => a.status === 'pending').length;
    const passed = applicants.filter(a => a.status === 'passed').length;
    const rejected = applicants.filter(a => a.status === 'rejected').length;

    return { total, pending, passed, rejected };
  };

  const stats = getStats();

  // 添加面试申请
  const handleAddApplicant = async (values: any) => {
    console.log('handleAddApplicant called', values);
    setAddLoading(true);
    try {
      const interviewDateTime = values.interview_date
        .hour(values.interview_time.hour())
        .minute(values.interview_time.minute())
        .format('YYYY-MM-DD HH:mm');
      const response = await fetch('/api/v1/apply', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: values.name,
          email: values.email,
          phone: values.phone,
          student_id: values.student_id,
          major: values.major,
          grade: values.grade,
          interview_time: interviewDateTime,
          verification_code: 'admin', // 后端如必填可用占位
        }),
      });
      const data = await response.json();
      if (data.code === 200) {
        message.success('添加成功');
        setAddModalVisible(false);
        addForm.resetFields();
        fetchApplicants();
      } else {
        message.error(data.message || '添加失败');
      }
    } catch (error) {
      message.error('网络错误，请稍后重试');
    } finally {
      setAddLoading(false);
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '姓名',
      dataIndex: 'name',
      key: 'name',
      width: 100,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 180,
      ellipsis: true,
    },
    {
      title: '手机',
      dataIndex: 'phone',
      key: 'phone',
      width: 120,
    },
    {
      title: '学号',
      dataIndex: 'student_id',
      key: 'student_id',
      width: 120,
    },
    {
      title: '专业',
      dataIndex: 'major',
      key: 'major',
      width: 120,
    },
    {
      title: '年级',
      dataIndex: 'grade',
      key: 'grade',
      width: 80,
    },
    {
      title: '面试时间',
      dataIndex: 'interview_time',
      key: 'interview_time',
      width: 150,
      render: (text: string) => dayjs(text).format('MM-DD HH:mm'),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '申请时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 120,
      render: (text: string) => dayjs(text).format('MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_: any, record: Applicant) => (
        <Space size="small">
          <Tooltip title="编辑状态">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => showEditModal(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确定要删除这条申请记录吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="text"
                danger
                icon={<DeleteOutlined />}
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ marginBottom: '24px' }}>
        <Title level={2}>🧪 EPI实验室面试管理后台</Title>
        <Text type="secondary">管理所有面试申请和状态</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="总申请数"
              value={stats.total}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="待面试"
              value={stats.pending}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="已通过"
              value={stats.passed}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="未通过"
              value={stats.rejected}
              prefix={<CloseCircleOutlined />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 申请列表 */}
      <Card
        title="面试申请列表"
        extra={
          <Space>
            <Button icon={<ReloadOutlined />} onClick={fetchApplicants} loading={loading}>
              刷新
            </Button>
            <Button type="primary" onClick={() => setAddModalVisible(true)}>
              添加申请
            </Button>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={applicants}
          rowKey="id"
          loading={loading}
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
            pageSize: 10,
            pageSizeOptions: ['10', '20', '50'],
          }}
          scroll={{ x: 1200 }}
        />
      </Card>

      {/* 添加申请Modal */}
      <Modal
        title="添加面试申请"
        open={addModalVisible}
        onCancel={() => { setAddModalVisible(false); addForm.resetFields(); }}
        footer={null}
        width={600}
      >
        {/* 温馨提示 */}
        <Alert
          message="请填写完整且格式正确，否则无法提交。"
          type="info"
          showIcon
          style={{ marginBottom: '16px' }}
        />
        <Form
          form={addForm}
          layout="vertical"
          initialValues={{
            grade: '2024',
            interview_date: dayjs(),
            interview_time: dayjs().hour(14).minute(0),
          }}
          validateTrigger="onChange"
          onFinish={async (values) => {
            console.log('onFinish called', values);
            await handleAddApplicant(values);
          }}
          onFinishFailed={({ values, errorFields }) => {
            message.error('表单校验未通过，请检查红色项');
            console.log('onFinishFailed', values, errorFields);
          }}
        >
          <Row gutter={[16, 0]}>
            <Col xs={24} md={12}>
              <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]} validateTrigger="onChange">
                <Input placeholder="请输入姓名" allowClear onChange={e => addForm.setFieldValue('name', e.target.value)} />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="email" label="邮箱" rules={[{ required: true, message: '请输入邮箱' }, { type: 'email', message: '请输入有效邮箱' }]} validateTrigger="onChange">
                <Input placeholder="请输入邮箱" allowClear onChange={e => addForm.setFieldValue('email', e.target.value)} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 0]}>
            <Col xs={24} md={12}>
              <Form.Item name="phone" label="手机号码" rules={[{ required: true, message: '请输入手机号码' }, { len: 11, message: '请输入11位手机号码' }]} validateTrigger="onChange">
                <Input placeholder="请输入手机号码" maxLength={11} allowClear onChange={e => addForm.setFieldValue('phone', e.target.value)} />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="student_id" label="学号" rules={[{ required: true, message: '请输入学号' }]} validateTrigger="onChange">
                <Input placeholder="请输入学号" allowClear onChange={e => addForm.setFieldValue('student_id', e.target.value)} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 0]}>
            <Col xs={24} md={12}>
              <Form.Item name="major" label="专业" rules={[{ required: true, message: '请选择专业' }]} validateTrigger="onChange">
                <Select placeholder="请选择专业" allowClear onChange={v => addForm.setFieldValue('major', v)}>
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
              <Form.Item name="grade" label="年级" rules={[{ required: true, message: '请选择年级' }]} validateTrigger="onChange">
                <Select placeholder="请选择年级" allowClear onChange={v => addForm.setFieldValue('grade', v)}>
                  <Option value="2021">2021级</Option>
                  <Option value="2022">2022级</Option>
                  <Option value="2023">2023级</Option>
                  <Option value="2024">2024级</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 0]}>
            <Col xs={24} md={12}>
              <Form.Item name="interview_date" label="面试日期" rules={[{ required: true, message: '请选择面试日期' }]} validateTrigger="onChange">
                <DatePicker style={{ width: '100%' }} placeholder="选择面试日期" allowClear onChange={v => addForm.setFieldValue('interview_date', v)} />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="interview_time" label="面试时间" rules={[{ required: true, message: '请选择面试时间' }]} validateTrigger="onChange">
                <TimePicker style={{ width: '100%' }} placeholder="选择面试时间" format="HH:mm" minuteStep={30} allowClear onChange={v => addForm.setFieldValue('interview_time', v)} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item style={{ marginTop: '32px' }}>
            <Button type="primary" htmlType="submit" loading={addLoading} style={{ width: '100%' }}
              onClick={async () => {
                const values = addForm.getFieldsValue();
                const errors = await addForm.validateFields().catch(e => e);
                console.log('submit click', values, errors);
                if (!errors || (Array.isArray(errors) && errors.length === 0)) {
                  // 校验通过但onFinish未触发，手动触发
                  handleAddApplicant(values);
                }
              }}
            >提交</Button>
          </Form.Item>
        </Form>
        {/* 实时输出表单校验错误，便于调试 */}
        <pre style={{ color: 'red', fontSize: 12 }}>{JSON.stringify(addForm.getFieldsError && addForm.getFieldsError(), null, 2)}</pre>
      </Modal>

      {/* 编辑模态框 */}
      <Modal
        title="更新申请状态"
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
        }}
        footer={null}
        width={600}
      >
        {selectedApplicant && (
          <div style={{ marginBottom: '16px' }}>
            <Text strong>申请人：{selectedApplicant.name}</Text>
            <br />
            <Text type="secondary">邮箱：{selectedApplicant.email}</Text>
            <br />
            <Text type="secondary">面试时间：{dayjs(selectedApplicant.interview_time).format('YYYY-MM-DD HH:mm')}</Text>
          </div>
        )}
        
        <Form
          form={form}
          layout="vertical"
          onFinish={handleUpdateStatus}
        >
          <Form.Item
            name="status"
            label="面试状态"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select placeholder="请选择面试状态">
              <Option value="pending">待面试</Option>
              <Option value="first_pass">一面通过</Option>
              <Option value="second_pass">二面通过</Option>
              <Option value="passed">通过</Option>
              <Option value="rejected">未通过</Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="first_remark"
            label="一面备注"
          >
            <TextArea
              rows={3}
              placeholder="请输入一面面试备注"
            />
          </Form.Item>

          <Form.Item
            name="second_remark"
            label="二面备注"
          >
            <TextArea
              rows={3}
              placeholder="请输入二面面试备注"
            />
          </Form.Item>

          <Form.Item
            name="third_remark"
            label="三面备注"
          >
            <TextArea
              rows={3}
              placeholder="请输入三面面试备注"
            />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                更新状态
              </Button>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
              }}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AdminDashboard; 