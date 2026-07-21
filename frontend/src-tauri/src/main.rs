#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use tauri::AppHandle;

#[derive(Serialize, Deserialize, Clone, Default)]
struct BgConfig {
    path: Option<String>,
    blur: u8,
    opacity: u8,
}

#[derive(Serialize, Deserialize)]
struct Config {
    openai_api_key: String,
    gemini_api_key: String,
    anthropic_api_key: String,
    ai_model: String,
    temperature: String,
}

#[derive(Serialize, Deserialize)]
struct ChatMessage {
    id: i64,
    sender: String,
    text: String,
    timestamp: i64,
}

#[derive(Serialize, Deserialize, Clone, Default)]
struct IconConfig {
    user_icon_path: Option<String>,
    user_icon_pos_x: i32,
    user_icon_pos_y: i32,
    ai_icon_path: Option<String>,
    ai_icon_pos_x: i32,
    ai_icon_pos_y: i32,
}

// Función auxiliar para obtener el directorio raíz de datos de la app
fn get_app_data_dir(app_handle: &AppHandle) -> Result<PathBuf, String> {
    tauri::api::path::app_local_data_dir(&app_handle.config())
        .ok_or_else(|| "No se pudo obtener el directorio de datos de la app".into())
}

#[tauri::command]
fn save_bg_config(app_handle: AppHandle, config: BgConfig) -> Result<(), String> {
    let mut config_path = get_app_data_dir(&app_handle)?;
    config_path.push("config");
    fs::create_dir_all(&config_path).map_err(|e| e.to_string())?;

    config_path.push("background_settings.json");
    let config_json = serde_json::to_string_pretty(&config).map_err(|e| e.to_string())?;
    fs::write(config_path, config_json).map_err(|e| e.to_string())?;

    Ok(())
}

#[tauri::command]
fn load_bg_config(app_handle: AppHandle) -> Result<BgConfig, String> {
    let mut config_path = get_app_data_dir(&app_handle)?;
    config_path.push("config");
    config_path.push("background_settings.json");

    if config_path.exists() {
        let content = fs::read_to_string(config_path).map_err(|e| e.to_string())?;
        let config: BgConfig = serde_json::from_str(&content).map_err(|e| e.to_string())?;
        Ok(config)
    } else {
        Ok(BgConfig {
            path: None,
            blur: 0,
            opacity: 100,
        })
    }
}

#[tauri::command]
fn save_background_image(
    app_handle: AppHandle,
    image_bytes: Vec<u8>,
    extension: String,
) -> Result<String, String> {
    if image_bytes.len() > 50 * 1024 * 1024 {
        return Err("La imagen excede el límite de 50MB".into());
    }

    let mut bg_dir = get_app_data_dir(&app_handle)?;
    bg_dir.push("assets");
    bg_dir.push("backgrounds");

    fs::create_dir_all(&bg_dir).map_err(|e| e.to_string())?;

    let filename = format!("current_bg.{}", extension);
    bg_dir.push(&filename);

    fs::write(&bg_dir, image_bytes).map_err(|e| e.to_string())?;
    Ok(bg_dir.to_string_lossy().to_string())
}

#[tauri::command]
fn save_icon_config(app_handle: AppHandle, config: IconConfig) -> Result<(), String> {
    let mut config_path = get_app_data_dir(&app_handle)?;
    config_path.push("config");
    fs::create_dir_all(&config_path).map_err(|e| e.to_string())?;

    config_path.push("icon_settings.json");
    let config_json = serde_json::to_string_pretty(&config).map_err(|e| e.to_string())?;
    fs::write(config_path, config_json).map_err(|e| e.to_string())?;

    Ok(())
}

#[tauri::command]
fn load_icon_config(app_handle: AppHandle) -> Result<IconConfig, String> {
    let mut config_path = get_app_data_dir(&app_handle)?;
    config_path.push("config");
    config_path.push("icon_settings.json");

    if config_path.exists() {
        let content = fs::read_to_string(config_path).map_err(|e| e.to_string())?;
        let config: IconConfig = serde_json::from_str(&content).map_err(|e| e.to_string())?;
        Ok(config)
    } else {
        Ok(IconConfig {
            user_icon_path: None,
            user_icon_pos_x: 50,
            user_icon_pos_y: 50,
            ai_icon_path: None,
            ai_icon_pos_x: 50,
            ai_icon_pos_y: 50,
        })
    }
}

#[tauri::command]
fn save_icon_image(
    app_handle: AppHandle,
    image_bytes: Vec<u8>,
    extension: String,
    icon_type: String, // "user" o "ai"
) -> Result<String, String> {
    // 5 MB check para iconos
    if image_bytes.len() > 5 * 1024 * 1024 {
        return Err("La imagen del ícono excede el límite de 5MB".into());
    }

    let mut icon_dir = get_app_data_dir(&app_handle)?;
    icon_dir.push("assets");
    icon_dir.push("icons");

    fs::create_dir_all(&icon_dir).map_err(|e| e.to_string())?;

    let filename = format!("{}_icon.{}", icon_type, extension);
    icon_dir.push(&filename);

    fs::write(&icon_dir, image_bytes).map_err(|e| e.to_string())?;
    Ok(icon_dir.to_string_lossy().to_string())
}

#[tauri::command]
fn save_personality_image(
    app_handle: AppHandle,
    image_bytes: Vec<u8>,
    filename: String,
) -> Result<String, String> {
    if image_bytes.len() > 5 * 1024 * 1024 {
        return Err("La imagen excede el límite de 5MB".into());
    }

    let mut dir = get_app_data_dir(&app_handle)?;
    dir.push("assets");
    dir.push("personalities");

    fs::create_dir_all(&dir).map_err(|e| e.to_string())?;
    dir.push(&filename);

    fs::write(&dir, image_bytes).map_err(|e| e.to_string())?;
    Ok(filename)
}

#[tauri::command]
fn get_personality_image_path(
    app_handle: AppHandle,
    filename: String,
) -> Result<String, String> {
    let mut dir = get_app_data_dir(&app_handle)?;
    dir.push("assets");
    dir.push("personalities");
    dir.push(&filename);
    
    Ok(dir.to_string_lossy().to_string())
}

#[tauri::command]
fn save_chat_message(
    app_handle: AppHandle,
    chat_code: i64,
    message_order: i32,
    message: ChatMessage,
) -> Result<(), String> {
    let mut dir = get_app_data_dir(&app_handle)?;
    dir.push("assets");
    dir.push("chats");
    dir.push(chat_code.to_string());

    fs::create_dir_all(&dir).map_err(|e| e.to_string())?;

    let filename = format!("{}.json", message_order);
    dir.push(&filename);

    let json_data = serde_json::to_string(&message).map_err(|e| e.to_string())?;
    fs::write(&dir, json_data).map_err(|e| e.to_string())?;

    Ok(())
}

#[tauri::command]
fn get_chat_messages(
    app_handle: AppHandle,
    chat_code: i64,
) -> Result<Vec<ChatMessage>, String> {
    let mut dir = get_app_data_dir(&app_handle)?;
    dir.push("assets");
    dir.push("chats");
    dir.push(chat_code.to_string());

    let mut messages = Vec::new();

    if dir.exists() && dir.is_dir() {
        let entries = fs::read_dir(&dir).map_err(|e| e.to_string())?;
        
        let mut files: Vec<(i32, std::path::PathBuf)> = Vec::new();

        for entry in entries {
            if let Ok(entry) = entry {
                let path = entry.path();
                if path.extension().and_then(|e| e.to_str()) == Some("json") {
                    if let Some(stem) = path.file_stem().and_then(|s| s.to_str()) {
                        if let Ok(order) = stem.parse::<i32>() {
                            files.push((order, path));
                        }
                    }
                }
            }
        }

        files.sort_by_key(|k| k.0);

        for (_, path) in files {
            if let Ok(content) = fs::read_to_string(&path) {
                if let Ok(msg) = serde_json::from_str::<ChatMessage>(&content) {
                    messages.push(msg);
                }
            }
        }
    }

    Ok(messages)
}

#[tauri::command]
fn delete_local_chat(
    app_handle: AppHandle,
    chat_code: i64,
) -> Result<(), String> {
    let mut dir = get_app_data_dir(&app_handle)?;
    dir.push("assets");
    dir.push("chats");
    dir.push(chat_code.to_string());

    if dir.exists() {
        fs::remove_dir_all(&dir).map_err(|e| e.to_string())?;
    }
    
    Ok(())
}

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            let data_dir = tauri::api::path::app_local_data_dir(&app.config());
            if let Some(mut path) = data_dir {
                path.push("config");
                let _ = fs::create_dir_all(&path);

                path.pop();
                path.push("assets");
                path.push("backgrounds");
                let _ = fs::create_dir_all(&path);

                path.pop();
                path.push("icons");
                let _ = fs::create_dir_all(&path);

                path.pop();
                path.push("personalities");
                let _ = fs::create_dir_all(&path);

                path.pop();
                path.push("chats");
                let _ = fs::create_dir_all(&path);
            }
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            save_bg_config,
            load_bg_config,
            save_background_image,
            save_icon_config,
            load_icon_config,
            save_icon_image,
            save_personality_image,
            get_personality_image_path,
            save_chat_message,
            get_chat_messages,
            delete_local_chat
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
