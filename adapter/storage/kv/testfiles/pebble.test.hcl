
server {
    bind_addr = "0.0.0.0"

    ports {
        http = 8080
        grpc = 50051
    }
}

raft {
    bootstrap = true
    local_id = "1"
    base_dir = "./testfiles/raft"
}

logger {
    log_path = "./testfiles/log/test.log"
    log_level = "DEBUG"
    raft_log_path = "./testfiles/log"
}

chain {
    base_dir = "./testfiles/data"
}

message_broker {
    server_addr = "127.0.0.1:8333"
}