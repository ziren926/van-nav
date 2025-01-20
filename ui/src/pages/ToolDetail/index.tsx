import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Button, message } from 'antd';
import ReactMarkdown from 'react-markdown';

interface Tool {
  id: string;
  name: string;
  desc: string;
  content: string;
  author?: string; // 添加作者字段
}

const ToolDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [tool, setTool] = useState<Tool | null>(null);
  const [loading, setLoading] = useState(true);
  const currentUser = 'ziren926'; // 当前用户

  const loadToolDetail = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/tools/${id}`);
      if (!response.ok) {
        throw new Error('加载失败');
      }
      const data = await response.json();

      // 验证返回的数据
      if (!data || typeof data !== 'object') {
        throw new Error('返回数据格式错误');
      }

      setTool(data);
    } catch (err) {
      message.error('加载失败: ' + (err instanceof Error ? err.message : '未知错误'));
      console.error('加载工具详情失败:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadToolDetail();
  }, [id]);

  // 检查是否是作者
  const isAuthor = tool?.author === currentUser;

  return (
    <div className="max-w-4xl mx-auto p-4">
      <Card loading={loading}>
        {!loading && tool && (
          <>
            <div className="flex justify-between items-center mb-4">
              <div>
                <h1 className="text-2xl font-bold">{tool.name}</h1>
                {tool.author && (
                  <div className="text-gray-500 text-sm mt-1">
                    作者：{tool.author}
                  </div>
                )}
              </div>
              {isAuthor && (
                <Button
                  type="link"
                  onClick={() => navigate(`/admin/tools/${id}/post`)}
                >
                  编辑帖子
                </Button>
              )}
            </div>
            <div>
              {/* 描述部分 */}
              {tool.desc && (
                <div className="mb-4 text-gray-600 p-4 bg-gray-50 rounded">
                  {tool.desc}
                </div>
              )}
              {/* 内容部分 */}
              <div className="prose max-w-none mt-6">
                {tool.content ? (
                  <ReactMarkdown>
                    {tool.content}
                  </ReactMarkdown>
                ) : (
                  <div className="text-gray-500">暂无详细内容</div>
                )}
              </div>
            </div>
          </>
        )}
      </Card>
    </div>
  );
};

export default ToolDetail;