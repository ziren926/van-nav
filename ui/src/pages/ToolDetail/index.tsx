// ui/src/pages/ToolDetail.tsx
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Input, Button, Form, message, Spin } from 'antd';
import ReactMarkdown from 'react-markdown';

const { TextArea } = Input;

const ToolDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [isEditing, setIsEditing] = useState(false);
  const [tool, setTool] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [form] = Form.useForm();

  useEffect(() => {
    loadToolDetail();
  }, [id]);

  const loadToolDetail = async () => {
    try {
      setLoading(true);
      // 使用现有的 api 方法
      const response = await fetch(`/api/tools/${id}`);
      const data = await response.json();
      setTool(data);
      form.setFieldsValue({
        content: data.content,
        description: data.desc
      });
    } catch (err) {
      message.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      await fetch(`/api/tools/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      });
      message.success('更新成功');
      setIsEditing(false);
      loadToolDetail();
    } catch (err) {
      message.error('更新失败');
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <Card loading={loading}>
        {!loading && (
          <>
            <div className="flex justify-between items-center mb-4">
              <h1 className="text-2xl font-bold">{tool?.name}</h1>
              <Button
                type={isEditing ? "primary" : "default"}
                onClick={() => setIsEditing(!isEditing)}
              >
                {isEditing ? '预览' : '编辑'}
              </Button>
            </div>
            {isEditing ? (
              <Form
                form={form}
                onFinish={handleSubmit}
                layout="vertical"
              >
                <Form.Item
                  name="description"
                  label="简短描述"
                >
                  <Input placeholder="简短描述工具的主要功能" />
                </Form.Item>
                <Form.Item
                  name="content"
                  label="详细内容"
                >
                  <TextArea
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
                  {tool?.desc}
                </div>
                <div className="prose max-w-none">
                  <ReactMarkdown>
                    {tool?.content || '暂无详细内容'}
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