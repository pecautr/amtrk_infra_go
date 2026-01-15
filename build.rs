// Build script to compile Protocol Buffers for GTFS-Realtime

use std::io::Result;

fn main() -> Result<()> {
    // Compile GTFS-Realtime .proto files
    prost_build::compile_protos(
        &[
            "proto/gtfs-realtime.proto",
        ],
        &["proto/"],
    )?;
    
    println!("cargo:rerun-if-changed=proto/gtfs-realtime.proto");
    Ok(())
}
