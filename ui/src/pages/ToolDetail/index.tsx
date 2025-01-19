import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Input, Button, Form, message } from 'antd';
import ReactMarkdown from 'react-markdown';

interface Tool {
  id: string;
  name: string;
  desc: string;
  content: string;
}

const ToolDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);
  const [tool, setTool] = useState<Tool | null>(null);
  const [loading, setLoading] = useState(true);
  const [form] = Form.useForm();

  useEffect(() => {
    loadToolDetail();
  }, [id]);

  const loadToolDetail = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/tools/${id}`);
      if (!response.ok) {
        throw new Error('加载失败');
      }
      const data = await response.json();
      setTool(data);
      form.setFieldsValue({
        content: data.content,
        description: data.desc
      });
    } catch (err) {
      message.error('加载失败');
      console.error('加载工具详情失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      const response = await fetch(`/api/tools/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('_token')}`,
        },
        body: JSON.stringify(values),
      });

      if (!response.ok) {
        throw new Error('更新失败');
      }

      message.success('更新成功');
      setIsEditing(false);
      loadToolDetail();
    } catch (err) {
      message.error('更新失败');
      console.error('更新工具详情失败:', err);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <Card loading={loading}>
        {!loading && tool && (
          <>
            <div className="flex justify-between items-center mb-4">
              <h1 className="text-2xl font-bold">{tool.name}</h1>
              <div>
                <Button
                  type={isEditing ? "primary" : "default"}
                  onClick={() => setIsEditing(!isEditing)}
                >
                  {isEditing ? '预览' : '编辑'}
                </Button>
                <Button
                  type="link"
                  className="ml-2"
                  onClick={() => navigate(`/admin/tools/${id}/post`)}
                >
                  编辑帖子
                </Button>
              </div>
            </div>
            {isEditing ? (
              <Form
                form={form}
                onFinish={handleSubmit}
                layout="vertical"
                initialValues={{
                  description: tool.desc,
                  content: tool.content
                }}
              >
                <Form.Item
                  name="description"
                  label="简短描述"
                  rules={[{ required: true, message: '请输入简短描述' }]}
                >
                  <Input placeholder="简短描述工具的主要功能" />
                </Form.Item>
                <Form.Item
                  name="content"
                  label="详细内容"
                  rules={[{ required: true, message: '请输入详细内容' }]}
                >
                  <Input.TextArea
                    rows={10}
                    placeholder="使用 Markdown 格式编写详细内容"
                  />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit">
                    保存
                  </Button>
                </Form.Item>
              </Form>
            ) : (
              <div>
                <div className="mb-4 text-gray-600">
                  {tool.desc}
                </div>
                <div className="prose max-w-none">
                  <ReactMarkdown>
                    {tool.content || '暂无详细内容'}
                  </ReactMarkdown>
                </div>
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
};

export default ToolDetail;