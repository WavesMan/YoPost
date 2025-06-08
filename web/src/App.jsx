import React, { useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import './App.css';
import Sidebar from './components/Sidebar';
import EmailList from './components/EmailList';
import EmailView from './components/EmailView';

function App() {
  const [activeFolder, setActiveFolder] = useState('Inbox');
  const [selectedEmail, setSelectedEmail] = useState(null);

  return (
    <Router>
      <div className="app">
        <div className="app-header">
          <h1>YoPost</h1>
        </div>
        <div className="app-body">
          <Sidebar 
            activeFolder={activeFolder} 
            setActiveFolder={setActiveFolder}
          />
          <EmailList 
            activeFolder={activeFolder} 
            setSelectedEmail={setSelectedEmail}
          />
          <EmailView 
            selectedEmail={selectedEmail}
          />
        </div>
      </div>
    </Router>
  );
}

export default App;
