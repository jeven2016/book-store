import React from 'react';
import ReactDOM from 'react-dom/client';
import reportWebVitals from './reportWebVitals';
import 'react-windy-ui/dist/wui.css';
import '@/styles/default.scss';
import '@/common/config/i18n';

import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import dayjs from 'dayjs';
import Home2 from '@/pages/Home2';

//detect current timezone
dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.tz.guess();

const rootElem = document.getElementById('root');
ReactDOM.createRoot(rootElem as HTMLElement).render(<Home2 />);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals(console.log);
