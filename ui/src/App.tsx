import React, { Suspense, useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { App as AntApp, Menu, Spin } from 'antd';
import { decodeTheme, initTheme } from './utils/theme';
import './App.css';

// 使用 React.lazy 懒加载组件
const Home = React.lazy(() => import('./pages/Home'));
const AdminPage = React.lazy(() => import('./pages/admin').then(module => ({ default: module.AdminPage })));
const Login = React.lazy(() => import('./pages/Login'));

// 懒加载管理后台的子页面
const Tools = React.lazy(() => import('./pages/admin/tabs/Tools').then(module => ({ default: module.Tools })));
const Catelog = React.lazy(() => import('./pages/admin/tabs/Catelog').then(module => ({ default: module.Catelog })));
const ApiToken = React.lazy(() => import('./pages/admin/tabs/ApiToken').then(module => ({ default: module.ApiToken })));
const Setting = React.lazy(() => import('./pages/admin/tabs/Setting').then(module => ({ default: module.Setting })));

// 导航栏组件
const Navigation = () => {
  const [current, setCurrent] = useState('home');
  const [isDarkMode, setIsDarkMode] = useState(false);

  useEffect(() => {
    const theme = initTheme();
    const decodedTheme = decodeTheme(theme);
    setIsDarkMode(decodedTheme.includes('dark'));
  }, []);

  const handleClick = (e: { key: string }) => {
    setCurrent(e.key);
  };

  return (
    <Menu
      mode="horizontal"
      className={`nav-menu ${isDarkMode ? 'dark' : 'light'}`}
      selectedKeys={[current]}
      onClick={handleClick}
    >
      <Menu.Item key="home">
        <Link to="/">首页</Link>
      </Menu.Item>
      <Menu.Item key="popular">
        <Link to="/popular">Chatgpt-Task介紹</Link>
      </Menu.Item>
      <Menu.Item key="new">
        <Link to="/new">我的收藏</Link>
      </Menu.Item>
    </Menu>
  );
};

// 加载中的占位组件
const LoadingFallback = () => {
  const [isDarkMode, setIsDarkMode] = useState(false);

  useEffect(() => {
    const theme = initTheme();
    const decodedTheme = decodeTheme(theme);
    setIsDarkMode(decodedTheme.includes('dark'));

    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.target instanceof HTMLElement) {
          setIsDarkMode(mutation.target.classList.contains('dark-mode'));
        }
      });
    });

    const body = document.querySelector('body');
    if (body) {
      observer.observe(body, {
        attributes: true,
        attributeFilter: ['class']
      });
    }

    return () => observer.disconnect();
  }, []);

  return (
    <div style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      backgroundColor: isDarkMode ? '#121212' : '#ffffff',
      color: isDarkMode ? 'rgba(255, 255, 255, 0.6)' : '#272e3b',
    }}>
      <Spin size="large" tip="加载中..." />
    </div>
  );
};

// 页面布局组件
const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="layout">
      <Navigation />
      {children}
    </div>
  );
};

function App() {
  return (
    <AntApp>
      <Router>
        <Suspense fallback={<LoadingFallback />}>
          <Routes>
            <Route path="/" element={<Layout><Home /></Layout>} />
            {/* 暂时注释掉未实现的路由 */}
            {/*
            <Route path="/popular" element={<Layout><Popular /></Layout>} />
            <Route path="/new" element={<Layout><New /></Layout>} />
            */}
            <Route path="/login" element={<Login />} />
            <Route path="/admin" element={<AdminPage />}>
              <Route index element={<Tools />} />
              <Route path="tools" element={<Tools />} />
              <Route path="categories" element={<Catelog />} />
              <Route path="api-token" element={<ApiToken />} />
              <Route path="settings" element={<Setting />} />
            </Route>
          </Routes>
        </Suspense>
      </Router>
    </AntApp>
  );
}

export default App;