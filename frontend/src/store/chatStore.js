import { create } from 'zustand';

const useChatStore = create(
  (set) => ({
    chats: [],
    activeChatId: null,

    setChats: (historiales) => set({
      chats: historiales.map(h => ({
        id: h.code.toString(),
        dbId: h.id,
        title: h.nombre,
        titleGenerated: true,
        messages: [],
        personalityId: null,
        updatedAt: h.created_at
      }))
    }),

    setActiveChat: (id) => set({ activeChatId: id }),

    addMessage: (chatId, message) => set((state) => {
      const updatedChats = state.chats.map((chat) => {
        if (chat.id === chatId) {
          return {
            ...chat,
            messages: [...chat.messages, { ...message, id: Date.now(), timestamp: Date.now() }],
            updatedAt: Date.now()
          };
        }
        return chat;
      });
      return { chats: updatedChats };
    }),

    setMessages: (chatId, messages) => set((state) => {
      const updatedChats = state.chats.map((chat) => {
        if (chat.id === chatId) {
          return { ...chat, messages };
        }
        return chat;
      });
      return { chats: updatedChats };
    }),

    createNewChat: () => set((state) => {
      const newChat = {
        id: Date.now().toString(),
        title: 'Nuevo Chat',
        titleGenerated: false,
        messages: [],
        personalityId: null,
        updatedAt: Date.now()
      };
      return {
        chats: [newChat, ...state.chats],
        activeChatId: newChat.id
      };
    }),

    updateChatTitle: (chatId, title, dbId) => set((state) => {
      const updatedChats = state.chats.map((chat) => {
        if (chat.id === chatId) {
          return { ...chat, title, dbId, titleGenerated: true, updatedAt: Date.now() };
        }
        return chat;
      });
      return { chats: updatedChats };
    }),

    updateChatPersonality: (chatId, personalityId) => set((state) => {
      const updatedChats = state.chats.map((chat) => {
        if (chat.id === chatId) {
          return { ...chat, personalityId, updatedAt: Date.now() };
        }
        return chat;
      });
      return { chats: updatedChats };
    }),

    deleteChat: (id) => set((state) => {
      const filteredChats = state.chats.filter(chat => chat.id !== id);
      return {
        chats: filteredChats,
        activeChatId: state.activeChatId === id ? (filteredChats[0]?.id || null) : state.activeChatId
      };
    })
  })
);

export default useChatStore;
