create_event_pb_files:
	protoc \
	--proto_path=internal/adapters/grpc/events_receiver/proto \
	internal/adapters/grpc/events_receiver/proto/*.proto \
	--go_out=internal/adapters/grpc/events_receiver

	protoc \
	--proto_path=internal/adapters/grpc/events_receiver/proto \
	internal/adapters/grpc/events_receiver/proto/*.proto \
	--go-grpc_out=internal/adapters/grpc/events_receiver

clean_event_pb_files:
	rm internal/adapters/grpc/events_receiver/event_pb/*.go