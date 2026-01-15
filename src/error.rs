// Error handling

use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
};
use std::fmt;

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),

    #[error("GTFS-RT error: {0}")]
    GtfsRt(String),

    #[error("HTTP error: {0}")]
    Http(#[from] reqwest::Error),

    #[error("Optimization error: {0}")]
    Optimization(String),

    #[error("ML error: {0}")]
    ML(String),

    #[error("Configuration error: {0}")]
    Config(#[from] config::ConfigError),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Serialization error: {0}")]
    Serialization(#[from] serde_json::Error),

    #[error("Not found: {0}")]
    NotFound(String),

    #[error("Invalid input: {0}")]
    InvalidInput(String),

    #[error("{0}")]
    Other(String),
}

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let (status, error_message) = match self {
            Error::NotFound(msg) => (StatusCode::NOT_FOUND, msg),
            Error::InvalidInput(msg) => (StatusCode::BAD_REQUEST, msg),
            Error::Database(ref e) => (StatusCode::INTERNAL_SERVER_ERROR, format!("Database error: {}", e)),
            _ => (StatusCode::INTERNAL_SERVER_ERROR, self.to_string()),
        };

        let body = serde_json::json!({
            "error": error_message,
        });

        (status, axum::Json(body)).into_response()
    }
}
