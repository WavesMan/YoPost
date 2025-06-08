import React, { useState } from 'react';
import './ComposeEmail.css';

function ComposeEmail({ onClose }) {
  const [isMaximized, setIsMaximized] = useState(false);

  const toggleMaximize = () => {
    setIsMaximized(!isMaximized);
  };

  return (
    <div className={`compose-modal ${isMaximized ? 'maximized' : ''}`}>
      <div className="compose-header">
        <h3>New Message</h3>
        <div className="compose-actions">
          <button onClick={toggleMaximize}>
            {isMaximized ? '↘' : '□'}
          </button>
          <button onClick={onClose}>×</button>
        </div>
      </div>
      <div className="compose-body">
        <input type="text" placeholder="To" />
        <input type="text" placeholder="Subject" />
        <textarea placeholder="Message"></textarea>
      </div>
      <div className="compose-footer">
        <button className="send-button">Send</button>
      </div>
    </div>
  );
}

export default ComposeEmail;