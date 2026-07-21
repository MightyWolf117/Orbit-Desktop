const BASE_URL = "http://localhost:8080/api";

const ENDPOINTS = {
  HEALTH: `${BASE_URL}/health`,
  MODELS: `${BASE_URL}/models`,
  CHAT: `${BASE_URL}/chat`,
  PERSONALITIES: `${BASE_URL}/personalities`,
  HISTORIAL: `${BASE_URL}/historial`
};

export { BASE_URL, ENDPOINTS };