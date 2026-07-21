import { Link, useLocation } from 'react-router-dom';
import { PlusCircle, Settings, MessageSquare, Trash2, Users } from 'lucide-react';
import { useEffect } from 'react';
import useChatStore from '../../store/chatStore';
import useSettingsStore from '../../store/settingsStore';
import { ENDPOINTS } from '../../service/api';
import styles from './Sidebar.module.scss';

const Sidebar = () => {
  const { chats, setChats, activeChatId, setActiveChat, createNewChat, deleteChat } = useChatStore();
  const { isOnline, checkHealth } = useSettingsStore();
  const location = useLocation();

  const handleNewChat = () => {
    createNewChat();
  };

  useEffect(() => {
    const fetchHistorial = async () => {
      try {
        const response = await fetch(`${ENDPOINTS.HISTORIAL}`);
        if (response.ok) {
          const historiales = await response.json();
          setChats(historiales || []);
        }
      } catch (e) {
        console.error("Error fetching historial", e);
      }
    };
    fetchHistorial();
  }, [setChats]);

  // Polling para el Health Check cada 5 minutos
  useEffect(() => {
    checkHealth(); // Ejecutar inmediatamente
    const interval = setInterval(() => {
      checkHealth();
    }, 300000);
    return () => clearInterval(interval);
  }, [checkHealth]);

  return (
    <div className={styles.sidebar}>
      <div className={styles.header}>
        <h2 className={styles.title}>Orbit</h2>
        <Link to="/" className={styles.newChatBtn} onClick={handleNewChat}>
          <PlusCircle size={20} />
          <span>Nuevo Chat</span>
        </Link>
      </div>

      <div className={styles.chatList}>
        <div className={styles.listTitle}>Historial</div>
        {chats.map((chat) => (
          <div 
            key={chat.id} 
            className={`${styles.chatItem} ${activeChatId === chat.id && location.pathname === '/' ? styles.active : ''}`}
            onClick={() => setActiveChat(chat.id)}
          >
            <Link to="/" className={styles.chatLink}>
              <MessageSquare size={18} />
              <span className={styles.chatTitle}>{chat.title || 'Chat sin título'}</span>
            </Link>
            <button 
              className={styles.deleteBtn} 
              onClick={async (e) => {
                e.stopPropagation();
                const deleteLocal = window.confirm("¿Deseas eliminar también los archivos locales de este chat?");
                
                if (chat.dbId) {
                  try {
                    await fetch(`${ENDPOINTS.HISTORIAL}/${chat.dbId}`, { method: 'DELETE' });
                  } catch (err) {
                    console.error("Error eliminando historial en nube:", err);
                  }
                }
                
                if (deleteLocal && typeof window !== 'undefined' && window.__TAURI_IPC__) {
                  const { invoke } = await import('@tauri-apps/api/tauri');
                  try {
                    await invoke('delete_local_chat', { chat_code: parseInt(chat.id) });
                  } catch (err) {
                    console.error("Error eliminando chat local:", err);
                  }
                }

                deleteChat(chat.id);
              }}
              title="Eliminar chat"
            >
              <Trash2 size={16} />
            </button>
          </div>
        ))}
        {chats.length === 0 && (
          <div className={styles.emptyState}>No hay chats recientes</div>
        )}
      </div>

      <div className={styles.footer}>
        <div 
          className={styles.healthIndicator} 
          title={isOnline ? "Backend Conectado (Click para refrescar)" : "Backend Desconectado (Click para reintentar)"}
          onClick={checkHealth}
        >
          <div className={`${styles.dot} ${isOnline ? styles.online : styles.offline}`}></div>
          <span>{isOnline ? 'Servicio en línea' : 'Servicio offline'}</span>
        </div>
        
        <Link 
          to="/personalities" 
          className={`${styles.footerBtn} ${location.pathname === '/personalities' ? styles.activeBtn : ''}`}
        >
          <Users size={20} />
          <span>Personalidades</span>
        </Link>

        <Link 
          to="/settings" 
          className={`${styles.footerBtn} ${location.pathname === '/settings' ? styles.activeBtn : ''}`}
        >
          <Settings size={20} />
          <span>Ajustes</span>
        </Link>
      </div>
    </div>
  );
};

export default Sidebar;
