import { useState, useEffect, useCallback } from 'react';
import { Card, Table, Button, Space, Input, message, Popconfirm, Modal, Form, Spin } from 'antd';
import { useData } from '../hooks/useData';
import { fetchPosts, fetchAddPost, fetchDeletePost, fetchUpdatePost } from '../../../utils/api';

export const Posts = () => {
  const [showAddModel, setShowAddModel] = useState(false);
  const [addForm] = Form.useForm();
  const { store, loading, reload } = useData();
  const [searchString, setSearchString] = useState("");

  const handleSubmit = async (values: any) => {
      try {
          await fetchAddPost(values);
          message.success('发布成功');
          // 刷新帖子列表
          fetchPosts();
      } catch (error) {
          console.error('发布失败:', error);
          message.error('发布失败: ' + (error.message || '未知错误'));
      }
  };

  // 处理创建帖子
  const handleCreate = useCallback(
    async (record: any) => {
      try {
        await fetchAddPost(record);
        message.success("添加成功!");
        setShowAddModel(false);
        reload();
      } catch (err) {
        message.warning("添加失败!");
      }
    },
    [reload]
  );

  // 处理删除帖子
  const handleDelete = useCallback(
    async (id: number) => {
      try {
        await fetchDeletePost(id);
        message.success("删除成功!");
        reload();
      } catch (err) {
        message.warning("删除失败!");
      }
    },
    [reload]
  );

  // 表格列配置
  const columns = [
    {
      title: '序号',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '标题',
      dataIndex: 'title',
      width: 200,
    },
    {
      title: '内容',
      dataIndex: 'content',
      ellipsis: true,
    },
    {
      title: '发布时间',
      dataIndex: 'createTime',
      width: 180,
    },
    {
      title: '操作',
      width: 120,
      render: (_, record: any) => (
        <Space>
          <Button
            type="link"
            onClick={() => {
              // TODO: 实现编辑功能
            }}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这篇帖子吗？"
            onConfirm={() => handleDelete(record.id)}
          >
            <Button type="link">删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title={
        <Space>
          <span>{`当前共 ${store?.posts?.length ?? 0} 篇帖子`}</span>
        </Space>
      }
      extra={
        <Space>
          <Input.Search
            allowClear
            placeholder="搜索帖子"
            onSearch={(value: string) => {
              setSearchString(value.trim());
            }}
          />
          <Button
            type="primary"
            onClick={() => {
              setShowAddModel(true);
            }}
          >
            新建帖子
          </Button>
          <Button
            onClick={() => {
              reload();
            }}
          >
            刷新
          </Button>
        </Space>
      }
    >
      <Spin spinning={loading}>
        <Table
          dataSource={store?.posts || []}
          columns={columns}
          rowKey="id"
          size="small"
          pagination={{ defaultPageSize: 10 }}
        />
      </Spin>

      {/* 新建帖子弹窗 */}
      <Modal
        open={showAddModel}
        title="新建帖子"
        onCancel={() => {
          setShowAddModel(false);
          addForm.resetFields();
        }}
        onOk={() => {
          addForm.validateFields().then((values) => {
            handleCreate(values);
          });
        }}
      >
        <Form form={addForm}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入帖子标题' }]}
          >
            <Input placeholder="请输入帖子标题" />
          </Form.Item>
          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入帖子内容' }]}
          >
            <Input.TextArea rows={4} placeholder="请输入帖子内容" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default Posts;