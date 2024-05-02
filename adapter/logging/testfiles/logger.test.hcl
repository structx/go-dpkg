
server {
    bind_addr = "0.0.0.0"
    default_timeout = 15

    ports {
        http = 8080
        grpc = 50051
    }
}

logger {
    log_path = "./testfiles/log/test.log"
    log_level = "DEBUG"
    raft_log_path = "./testfiles/log"
}