import { useState, useEffect } from 'react';
import { PlusCircle, Edit3, UserCircle, Save } from 'lucide-react';
import Modal from '../../components/common/Modal/Modal';
import { invoke, convertFileSrc } from '@tauri-apps/api/tauri';
import { ENDPOINTS } from '../../service/api';
import styles from './PersonalityPage.module.scss';

const isTauri = typeof window !== 'undefined' && window.__TAURI_IPC__ !== undefined;

const PersonalityPage = () => {
  const [personalities, setPersonalities] = useState([]);
  const [loading, setLoading] = useState(true);
  
  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState(null);
  
  // Form State
  const [formData, setFormData] = useState({ nombre: '', descripcion_corta: '', instrucciones: '', image: null });
  const [selectedFile, setSelectedFile] = useState(null);
  const [previewUrl, setPreviewUrl] = useState(null);

  const [alert, setAlert] = useState({ isOpen: false, title: '', message: '', isError: false });

  const fetchPersonalities = async () => {
    setLoading(true);
    try {
      const response = await fetch(ENDPOINTS.PERSONALITIES);
      if (response.ok) {
        const data = await response.json();
        
        if (isTauri) {
          const processedData = await Promise.all(data.map(async (p) => {
            if (p.image) {
              try {
                const fullPath = await invoke('get_personality_image_path', { filename: p.image });
                return { ...p, localImageUrl: convertFileSrc(fullPath) };
              } catch (e) {
                return p;
              }
            }
            return p;
          }));
          setPersonalities(processedData || []);
        } else {
          setPersonalities(data || []);
        }
      }
    } catch (e) {
      console.error("Error fetching personalities:", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPersonalities();
  }, []);

  const loadPersonalityImagePreview = async (imageFilename) => {
    if (!imageFilename || !isTauri) return null;
    try {
      const fullPath = await invoke('get_personality_image_path', { filename: imageFilename });
      return convertFileSrc(fullPath);
    } catch (e) {
      console.error("Error loading image path:", e);
      return null;
    }
  };

  const openNewModal = () => {
    setEditingId(null);
    setFormData({ nombre: '', descripcion_corta: '', instrucciones: '', image: null });
    setSelectedFile(null);
    setPreviewUrl(null);
    setIsModalOpen(true);
  };

  const openEditModal = async (p) => {
    setEditingId(p.id);
    setFormData({ 
      nombre: p.nombre || '', 
      descripcion_corta: p.descripcion_corta || '', 
      instrucciones: p.instrucciones || '',
      image: p.image || null
    });
    setSelectedFile(null);
    setPreviewUrl(null);
    setIsModalOpen(true);

    if (p.image) {
      const url = await loadPersonalityImagePreview(p.image);
      setPreviewUrl(url);
    }
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setSelectedFile(file);
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  const handleSave = async () => {
    if (!formData.nombre.trim() || !formData.instrucciones.trim()) {
      setAlert({ isOpen: true, title: 'Campos requeridos', message: 'El nombre y las instrucciones son obligatorios.', isError: true });
      return;
    }

    try {
      let finalFormData = { ...formData };
      
      if (selectedFile) {
        if (isTauri) {
           const arrayBuffer = await selectedFile.arrayBuffer();
           const bytes = Array.from(new Uint8Array(arrayBuffer));
           const extension = selectedFile.name.split('.').pop();
           const filename = `persona_${Date.now()}.${extension}`;
           
           const savedFilename = await invoke('save_personality_image', {
             imageBytes: bytes,
             filename: filename
           });
           finalFormData.image = savedFilename;
        } else {
           console.warn("Tauri no está disponible, no se guardará la imagen.");
        }
      }

      const isEdit = editingId !== null;
      const url = isEdit ? `${ENDPOINTS.PERSONALITIES}/${editingId}` : ENDPOINTS.PERSONALITIES;
      const method = isEdit ? 'PUT' : 'POST';

      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(finalFormData)
      });

      if (response.ok) {
        setIsModalOpen(false);
        fetchPersonalities();
      } else {
        const errData = await response.json();
        throw new Error(errData.error || 'Error al guardar');
      }
    } catch (e) {
      setAlert({ isOpen: true, title: 'Error', message: e.message, isError: true });
    }
  };

  return (
    <div className={styles.pageContainer}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Personalidades de IA</h1>
          <p className={styles.subtitle}>Define el comportamiento y el rol de tus agentes.</p>
        </div>
        <button className={styles.primaryBtn} onClick={openNewModal}>
          <PlusCircle size={20} />
          Crear Personalidad
        </button>
      </header>

      {loading ? (
        <div className={styles.loadingState}>Cargando personalidades...</div>
      ) : (
        <div className={styles.grid}>
          {personalities.map((p) => (
            <div key={p.id} className={styles.card}>
              <div className={styles.cardHeader}>
                <div className={styles.cardTitle}>
                  {p.localImageUrl ? (
                    <div className={styles.cardAvatar} style={{backgroundImage: `url(${p.localImageUrl})`}} />
                  ) : (
                    <UserCircle size={24} className={styles.cardIcon} />
                  )}
                  <h3>{p.nombre}</h3>
                </div>
                <button className={styles.editBtn} onClick={() => openEditModal(p)} title="Editar">
                  <Edit3 size={18} />
                </button>
              </div>
              <div className={styles.cardBody}>
                <p className={styles.description}>{p.descripcion_corta || 'Sin descripción.'}</p>
              </div>
            </div>
          ))}
          {personalities.length === 0 && (
            <div className={styles.emptyState}>No has creado ninguna personalidad.</div>
          )}
        </div>
      )}

      {/* Modal Creación/Edición */}
      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)}
        title={editingId ? 'Editar Personalidad' : 'Nueva Personalidad'}
        actions={
          <>
            <button className={styles.cancelBtn} onClick={() => setIsModalOpen(false)}>Cancelar</button>
            <button className={styles.saveBtn} onClick={handleSave}><Save size={16}/> Guardar</button>
          </>
        }
      >
        <div className={styles.formGroup}>
          <label>Imagen de Perfil</label>
          <div className={styles.imageUploadContainer}>
            <div 
              className={styles.imagePreview} 
              style={previewUrl ? { backgroundImage: `url(${previewUrl})` } : {}}
            >
              {!previewUrl && <UserCircle size={44} className={styles.cardIcon} style={{margin: '0'}}/>}
            </div>
            <input 
              type="file" 
              accept="image/*" 
              onChange={handleImageChange} 
              className={styles.fileInput}
            />
          </div>
        </div>
        <div className={styles.formGroup}>
          <label>Nombre de la Personalidad *</label>
          <input 
            className={styles.input} 
            value={formData.nombre} 
            onChange={e => setFormData({...formData, nombre: e.target.value})} 
            placeholder="Ej. Programador Experto" 
          />
        </div>
        <div className={styles.formGroup}>
          <label>Descripción corta</label>
          <input 
            className={styles.input} 
            value={formData.descripcion_corta} 
            onChange={e => setFormData({...formData, descripcion_corta: e.target.value})} 
            placeholder="Una breve descripción para identificarlo" 
          />
        </div>
        <div className={styles.formGroup}>
          <label>Instrucciones de Comportamiento (Prompt) *</label>
          <textarea 
            className={styles.textarea} 
            value={formData.instrucciones} 
            onChange={e => setFormData({...formData, instrucciones: e.target.value})} 
            placeholder="Ej. Eres un experto en Python. Debes responder con código limpio..." 
            rows={5}
          />
        </div>
      </Modal>

      {/* Modal Alertas */}
      <Modal 
        isOpen={alert.isOpen} 
        onClose={() => setAlert({ ...alert, isOpen: false })}
        title={alert.title}
        actions={
          <button 
            className={`${styles.saveBtn} ${alert.isError ? styles.errorBtn : ''}`}
            onClick={() => setAlert({ ...alert, isOpen: false })}
          >
            Entendido
          </button>
        }
      >
        <p>{alert.message}</p>
      </Modal>

    </div>
  );
};

export default PersonalityPage;
