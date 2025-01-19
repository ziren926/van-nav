import React, { useState, useEffect } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { Card, Input, Button, Form, message } from 'antd';
import { TextArea } from 'antd/lib/input';

const PostEditor: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const history = useHistory();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchPost();
  }, [id]);

  const fetchPost = async () => {
    try {
      const response = await fetch(`/api/admin/tool/${id}/post`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      form.setFieldsValue(data);
      setLoading(false);
    } catch (error) {
      message.error('加载帖子失败');
      setLoading(false);
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      await fetch(`/api/admin/tool/${id}/post`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(values),
      });
      message.success('保存成功');
      history.push(`/admin/tools/${id}`);
    } catch (error) {
      message.error('保存失败');
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <Card loading={loading}>
        <Form
          form={form}
          onFinish={handleSubmit}
          layout="vertical"
        >
          <Form.Item
            name="post_title"
            label="标题"
          >
            <Input placeholder="请输入标题" />
          </Form.Item>
          <Form.Item
            name="post_content"
            label="内容"
          >
            <TextArea rows={10} placeholder="请输入内容" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              保存
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default PostEditor;