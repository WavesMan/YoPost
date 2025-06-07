import React from 'react';
import './EmailView.css';

function EmailView({ selectedEmail, setSelectedEmail }) {
  const markAsImportant = () => {
    // 标记重要邮件逻辑
    const updatedEmail = {
      ...selectedEmail,
      important: true
    };
    setSelectedEmail(updatedEmail);
  };

  const deletePermanently = () => {
    // 永久删除逻辑
    setSelectedEmail(null);
  };

  return (
    <div className="email-view">
      {selectedEmail ? (
        <>
          <div className="email-view-header">
            <h3>{selectedEmail.subject}</h3>
            <div className="email-view-metadata">
              <span className="email-sender">{selectedEmail.sender}</span>
              <span className="email-time">{selectedEmail.time}</span>
            </div>
            <div className="email-view-actions">
              <button onClick={markAsImportant}>Mark Important</button>
              <button onClick={deletePermanently}>Delete Permanently</button>
            </div>
          </div>
          <div className="email-view-body">
            <p>{selectedEmail.preview}</p>
          </div>
        </>
      ) : (
        <div className="email-view-empty">
          <p>Select an email to view</p>
        </div>
      )}
    </div>
  );
}

export default EmailView;