import React from 'react';
import { Provider } from 'react-redux';
import { ConfigProvider, App as AntdApp } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { store } from './store';
import AppRouter from './router';
import './styles/global.css';

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <ConfigProvider locale={zhCN}>
        <AntdApp>
          <AppRouter />
        </AntdApp>
      </ConfigProvider>
    </Provider>
  );
};

export default App;
