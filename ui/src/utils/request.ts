// ui/src/utils/request.ts

interface RequestOptions extends RequestInit {
  data?: any;
}
export async function request(url: string, options: RequestOptions = {}) {
  // 确保 url 以 / 开头
  const fullUrl = url.startsWith('/') ? url : `/${url}`;

  const { data, ...rest } = options;

  const defaultOptions: RequestOptions = {
    headers: {
      'Content-Type': 'application/json',
    },
    ...rest,
  };

  if (data) {
    defaultOptions.body = JSON.stringify(data);
  }

  // 获取存储的 token
  const token = localStorage.getItem('_token');
  if (token) {
    defaultOptions.headers = {
      ...defaultOptions.headers,
      'Authorization': token,
    };
  }

  console.log('Request URL:', fullUrl);
  console.log('Request Options:', defaultOptions);

  try {
    const response = await fetch(fullUrl, defaultOptions);
    console.log('Response Status:', response.status);

    if (response.status === 401) {
      localStorage.removeItem('_token');
      window.location.href = '/login';
      return;
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      const jsonResponse = await response.json();
      console.log('Response Data:', jsonResponse);
      return jsonResponse;
    }

    const textResponse = await response.text();
    console.log('Response Text:', textResponse);
    return textResponse;

  } catch (error) {
    console.error('Request error:', error);
    throw error;
  }
}