import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

import { TonConnectUIProvider} from '@tonconnect/ui-react';
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
        <TonConnectUIProvider restoreConnection={false} manifestUrl={`${window.location.origin}/ton.json`}>
    <App />
    </TonConnectUIProvider>
  </React.StrictMode>
);

