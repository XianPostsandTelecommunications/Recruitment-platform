import React from 'react';
import { Provider } from 'react-redux';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { store } from './store';
import AppRouter from './router';
import './styles/global.css';

const App: React.FC = () => {
  console.log('App component rendering...');
  
  return (
    <Provider store={store}>
      <ConfigProvider locale={zhCN}>
        <AppRouter />
      </ConfigProvider>
    </Provider>
  );
};

export default App;
