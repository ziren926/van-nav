// ui/src/utils/request.ts

interface RequestOptions extends RequestInit {
  data?: any;
}

export async function request(url: string, options: RequestOptions = {}) {
  const { data, ...rest } = options;

  const defaultOptions: RequestOptions = {
    headers: {
      'Content-Type': 'application/json',
    },
    ...rest,
  };

  // 如果有 data，添加到 body 中
  if (data) {
    defaultOptions.body = JSON.stringify(data);
  }

  // 获取存储的 token
  const token = localStorage.getItem('token');
  if (token) {
    defaultOptions.headers = {
      ...defaultOptions.headers,
      'Authorization': `Bearer ${token}`,
    };
  }

  try {
    const response = await fetch(url, defaultOptions);

    // 处理 401 未授权的情况
    if (response.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
      return;
    }

    // 如果响应不是 200 OK，抛出错误
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    // 尝试解析 JSON 响应
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    }

    // 如果不是 JSON，返回原始响应
    return await response.text();

  } catch (error) {
    console.error('Request error:', error);
    throw error;
  }
}