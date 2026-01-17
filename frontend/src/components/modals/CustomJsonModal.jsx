import React, { useState } from 'react';
import { X } from 'lucide-react';
import { useChat } from '../../context/ChatContext';

export default function CustomJsonModal({ isOpen, onClose }) {
  const { sendMessage } = useChat();
  const [customJson, setCustomJson] = useState('');

  if (!isOpen) return null;

  const handleSend = () => {
    try {
      const jsonObj = JSON.parse(customJson);
      sendMessage(jsonObj);
      setCustomJson('');
      onClose();
    } catch (e) {
      alert('Invalid JSON format: ' + e.message);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl">
        <div className="p-4 border-b border-gray-200 flex items-center justify-between">
          <h3 className="text-lg font-semibold">Send Custom JSON</h3>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            <X className="w-5 h-5" />
          </button>
        </div>
        <div className="p-4">
          <textarea
            value={customJson}
            onChange={(e) => setCustomJson(e.target.value)}
            className="w-full h-64 px-3 py-2 border border-gray-300 rounded-lg font-mono text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder={`{\n  "type": "messageRoom",\n  "user": "username",\n  "room": 1,\n  "content": "Hello"\n}`}
          />
          <div className="mt-4 flex justify-end gap-2">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSend}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}