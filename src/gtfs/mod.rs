// GTFS-Realtime ingestion module

use crate::{datamodel::*, error::Result};
use chrono::{DateTime, Utc};
use reqwest::Client;
use std::time::Duration;
use tokio::time;
use tracing::{debug, error, info};

// Generated from gtfs-realtime.proto by prost
pub mod proto {
    include!(concat!(env!("OUT_DIR"), "/transit_realtime.rs"));
}

pub struct GtfsRtIngester {
    client: Client,
    vehicle_positions_url: String,
    trip_updates_url: String,
    service_alerts_url: String,
}

impl GtfsRtIngester {
    pub fn new(
        vehicle_positions_url: String,
        trip_updates_url: String,
        service_alerts_url: String,
    ) -> Self {
        let client = Client::builder()
            .timeout(Duration::from_secs(10))
            .build()
            .expect("Failed to create HTTP client");

        Self {
            client,
            vehicle_positions_url,
            trip_updates_url,
            service_alerts_url,
        }
    }

    /// Fetch vehicle positions from GTFS-RT feed
    pub async fn fetch_vehicle_positions(&self) -> Result<Vec<VehiclePosition>> {
        debug!("Fetching vehicle positions from {}", self.vehicle_positions_url);

        let bytes = self.client
            .get(&self.vehicle_positions_url)
            .send()
            .await?
            .bytes()
            .await?;

        let feed = prost::Message::decode(&bytes[..])
            .map_err(|e| crate::Error::GtfsRt(format!("Failed to decode protobuf: {}", e)))?;

        let positions = self.parse_vehicle_positions(feed)?;
        info!("Fetched {} vehicle positions", positions.len());

        Ok(positions)
    }

    /// Fetch trip updates from GTFS-RT feed
    pub async fn fetch_trip_updates(&self) -> Result<Vec<TripUpdate>> {
        debug!("Fetching trip updates from {}", self.trip_updates_url);

        let bytes = self.client
            .get(&self.trip_updates_url)
            .send()
            .await?
            .bytes()
            .await?;

        let feed: proto::FeedMessage = prost::Message::decode(&bytes[..])
            .map_err(|e| crate::Error::GtfsRt(format!("Failed to decode protobuf: {}", e)))?;

        let updates = self.parse_trip_updates(feed)?;
        info!("Fetched {} trip updates", updates.len());

        Ok(updates)
    }

    /// Fetch service alerts from GTFS-RT feed
    pub async fn fetch_service_alerts(&self) -> Result<Vec<ServiceAlert>> {
        debug!("Fetching service alerts from {}", self.service_alerts_url);

        let bytes = self.client
            .get(&self.service_alerts_url)
            .send()
            .await?
            .bytes()
            .await?;

        let feed: proto::FeedMessage = prost::Message::decode(&bytes[..])
            .map_err(|e| crate::Error::GtfsRt(format!("Failed to decode protobuf: {}", e)))?;

        let alerts = self.parse_service_alerts(feed)?;
        info!("Fetched {} service alerts", alerts.len());

        Ok(alerts)
    }

    fn parse_vehicle_positions(&self, feed: proto::FeedMessage) -> Result<Vec<VehiclePosition>> {
        let mut positions = Vec::new();

        for entity in feed.entity {
            if let Some(vehicle) = entity.vehicle {
                if let Some(position) = vehicle.position {
                    if let Some(trip) = vehicle.trip {
                        positions.push(VehiclePosition {
                            id: uuid::Uuid::new_v4(),
                            vehicle_id: vehicle.vehicle.map(|v| v.id).unwrap_or_default(),
                            trip_id: trip.trip_id.unwrap_or_default(),
                            route_id: trip.route_id.unwrap_or_default(),
                            timestamp: DateTime::from_timestamp(
                                vehicle.timestamp.unwrap_or(0) as i64,
                                0
                            ).unwrap_or_else(Utc::now),
                            latitude: position.latitude,
                            longitude: position.longitude,
                            speed: position.speed,
                            heading: position.bearing,
                            current_stop_sequence: vehicle.current_stop_sequence.map(|s| s as i32),
                            current_status: vehicle.current_status.map(|s| format!("{:?}", s)),
                        });
                    }
                }
            }
        }

        Ok(positions)
    }

    fn parse_trip_updates(&self, feed: proto::FeedMessage) -> Result<Vec<TripUpdate>> {
        let mut updates = Vec::new();

        for entity in feed.entity {
            if let Some(trip_update) = entity.trip_update {
                if let Some(trip) = trip_update.trip {
                    let stop_time_updates: Vec<StopTimeUpdate> = trip_update.stop_time_update
                        .into_iter()
                        .map(|stu| StopTimeUpdate {
                            stop_id: stu.stop_id.unwrap_or_default(),
                            stop_sequence: stu.stop_sequence.unwrap_or(0) as i32,
                            arrival_delay: stu.arrival.as_ref().and_then(|a| a.delay).map(|d| d as i32),
                            departure_delay: stu.departure.as_ref().and_then(|d| d.delay).map(|d| d as i32),
                            arrival_time: stu.arrival.and_then(|a| a.time)
                                .and_then(|t| DateTime::from_timestamp(t, 0)),
                            departure_time: stu.departure.and_then(|d| d.time)
                                .and_then(|t| DateTime::from_timestamp(t, 0)),
                        })
                        .collect();

                    updates.push(TripUpdate {
                        trip_id: trip.trip_id.unwrap_or_default(),
                        route_id: trip.route_id.unwrap_or_default(),
                        timestamp: Utc::now(),
                        delay: trip_update.delay.unwrap_or(0) as i32,
                        stop_time_updates,
                    });
                }
            }
        }

        Ok(updates)
    }

    fn parse_service_alerts(&self, feed: proto::FeedMessage) -> Result<Vec<ServiceAlert>> {
        let mut alerts = Vec::new();

        for entity in feed.entity {
            if let Some(alert) = entity.alert {
                let active_period = alert.active_period.first();

                alerts.push(ServiceAlert {
                    id: entity.id,
                    cause: format!("{:?}", alert.cause.unwrap_or(1)),
                    effect: format!("{:?}", alert.effect.unwrap_or(1)),
                    header_text: alert.header_text
                        .and_then(|t| t.translation.first().map(|tr| tr.text.clone()))
                        .flatten()
                        .unwrap_or_default(),
                    description_text: alert.description_text
                        .and_then(|t| t.translation.first().map(|tr| tr.text.clone()))
                        .flatten()
                        .unwrap_or_default(),
                    affected_routes: alert.informed_entity
                        .into_iter()
                        .filter_map(|ie| ie.route_id)
                        .collect(),
                    active_period_start: active_period
                        .and_then(|ap| ap.start)
                        .and_then(|t| DateTime::from_timestamp(t as i64, 0))
                        .unwrap_or_else(Utc::now),
                    active_period_end: active_period
                        .and_then(|ap| ap.end)
                        .and_then(|t| DateTime::from_timestamp(t as i64, 0)),
                });
            }
        }

        Ok(alerts)
    }

    /// Start continuous ingestion loop
    pub async fn run(&self, interval_secs: u64) {
        let mut interval = time::interval(Duration::from_secs(interval_secs));

        loop {
            interval.tick().await;

            if let Err(e) = self.fetch_all().await {
                error!("Failed to fetch GTFS-RT data: {}", e);
            }
        }
    }

    async fn fetch_all(&self) -> Result<()> {
        let (positions, updates, alerts) = tokio::try_join!(
            self.fetch_vehicle_positions(),
            self.fetch_trip_updates(),
            self.fetch_service_alerts()
        )?;

        debug!(
            "Fetched {} positions, {} updates, {} alerts",
            positions.len(),
            updates.len(),
            alerts.len()
        );

        // TODO: Persist to database

        Ok(())
    }
}
