import React, { useState } from 'react';
import './EmailList.css';

function EmailList({ activeFolder, setSelectedEmail }) {
  const folderEmails = {
    Inbox: [
      { id: 1, sender: 'user1@example.com', subject: 'Meeting tomorrow', preview: 'About the project discussion...', time: '10:30 AM', read: false },
      { id: 2, sender: 'user2@example.com', subject: 'Weekly report', preview: 'Please find attached the weekly...', time: 'Yesterday', read: true },
      { id: 3, sender: 'user3@example.com', subject: 'Project update', preview: 'The development is going well...', time: 'Mar 15', read: true }
    ],
    Sent: [
      { id: 4, sender: 'me@example.com', subject: 'Project update', preview: 'Sent the latest project updates...', time: 'Mar 16', read: true }
    ],
    Drafts: [
      { id: 5, sender: 'me@example.com', subject: 'Project update', preview: 'Draft of the project updates...', time: 'Mar 16', read: true }
    ],
    Trash: [
      { id: 6, sender: 'user1@example.com', subject: 'Meeting tomorrow', preview: 'About the project discussion...', time: '10:30 AM', read: false }
    ]
  };

  const [selectedEmailId, setSelectedEmailId] = useState(null);
  const [selectedEmails, setSelectedEmails] = useState([]); // 新增批量选择状态
  const [isSelecting, setIsSelecting] = useState(false); // 批量选择模式

  const handleEmailClick = (email) => {
    if (isSelecting) {
      // 批量选择模式逻辑
      if (selectedEmails.includes(email.id)) {
        setSelectedEmails(selectedEmails.filter(id => id !== email.id));
      } else {
        setSelectedEmails([...selectedEmails, email.id]);
      }
    } else {
      // 原有单个选择逻辑
      if (selectedEmailId === email.id) {
        setSelectedEmailId(null);
        setSelectedEmail(null);
      } else {
        setSelectedEmailId(email.id);
        const updatedEmail = {
          ...email,
          read: true
        };
        setSelectedEmail(updatedEmail);
        folderEmails[activeFolder] = folderEmails[activeFolder].map(e => 
          e.id === email.id ? updatedEmail : e
        );
      }
    }
};

  const toggleSelectMode = () => {
    setIsSelecting(!isSelecting);
    setSelectedEmails([]);
  };

  const markAsImportant = () => {
    // 标记重要邮件逻辑
    folderEmails[activeFolder] = folderEmails[activeFolder].map(e => 
      selectedEmails.includes(e.id) ? {...e, important: true} : e
    );
  };

  const deleteSelected = () => {
    // 删除邮件逻辑
    folderEmails[activeFolder] = folderEmails[activeFolder].filter(
      e => !selectedEmails.includes(e.id)
    );
    setSelectedEmails([]);
  };

  // 在返回的JSX中修改controls部分
  return (
    <div className="email-list">
      <div className="email-list-header">
        <h2>{activeFolder}</h2>
        <div className="email-list-controls">
          <button onClick={toggleSelectMode}>
            {isSelecting ? 'Cancel' : 'Select'}
          </button>
          {isSelecting && (
            <>
              <button onClick={markAsImportant}>Mark Important</button>
              <button onClick={deleteSelected}>Delete</button>
            </>
          )}
          <input type="text" placeholder="Search emails" />
        </div>
      </div>
      <div className="email-list-items">
        {folderEmails[activeFolder] ? (
          folderEmails[activeFolder].map(email => (
            <div 
              key={email.id} 
              className={`email-list-item ${email.read ? 'read' : 'unread'} 
                ${selectedEmailId === email.id ? 'selected' : ''}
                ${selectedEmails.includes(email.id) ? 'selected-batch' : ''}`}
              onClick={() => handleEmailClick(email)}
            >
              <div className="email-sender">
                {!email.read && selectedEmailId !== email.id && <span className="unread-dot">•</span>}
                {email.sender}
              </div>
              <div className="email-subject">{email.subject}</div>
              <div className="email-preview">{email.preview}</div>
              <div className="email-time">{email.time}</div>
            </div>
          ))
        ) : (
          <div className="email-list-empty">
            No emails in this folder
          </div>
        )}
      </div>
    </div>
  );
}

export default EmailList;