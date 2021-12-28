import React from 'react'
import ReactDOM from 'react-dom'
import { App } from './app.jsx';  // app.js の App をインポート
 
const root = document.createElement('div');
root.setAttribute('id', 'root');
document.body.appendChild(root);
 
ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
);