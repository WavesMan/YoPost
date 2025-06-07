import React, { useState } from 'react';
import './Sidebar.css';
import ComposeEmail from './ComposeEmail'; // 新增导入

function Sidebar({ activeFolder, setActiveFolder }) {
  const menuItems = ['Inbox', 'Starred', 'Sent', 'Drafts', 'Trash'];
  const [showCompose, setShowCompose] = useState(false); // 新增状态

  return (
    <div className="sidebar">
      <button 
        className="sidebar-compose" 
        onClick={() => setShowCompose(true)} // 新增点击事件
      >
        Compose
      </button>
      <div className="sidebar-menu">
        {menuItems.map(item => (
          <div 
            key={item}
            className={`sidebar-menu-item ${activeFolder === item ? 'active' : ''}`}
            onClick={() => setActiveFolder(item)}
          >
            {item}
          </div>
        ))}
      </div>
      {showCompose && <ComposeEmail onClose={() => setShowCompose(false)} />}
    </div>
  );
}

export default Sidebar;