create_event_pb_files:
	protoc \
	--proto_path=internal/adapters/grpc/events_receiver/proto \
	internal/adapters/grpc/events_receiver/proto/*.proto \
	--go_out=internal/adapters/grpc/events_receiver

	protoc \
	--proto_path=internal/adapters/grpc/events_receiver/proto \
	internal/adapters/grpc/events_receiver/proto/*.proto \
	--go-grpc_out=internal/adapters/grpc/events_receiver

create_auth_pb_files:
	protoc \
	--proto_path=internal/adapters/grpc/auth/proto \
	internal/adapters/grpc/auth/proto/*.proto \
	--go_out=internal/adapters/grpc/auth

	protoc \
	--proto_path=internal/adapters/grpc/auth/proto \
	internal/adapters/grpc/auth/proto/*.proto \
	--go-grpc_out=internal/adapters/grpc/auth

clean_event_pb_files:
	rm internal/adapters/grpc/events_receiver/event_pb/*.go

clean_auth_pb_files:
	rm internal/adapters/grpc/auth/auth_pb/*.go

check_approved_tasks:
	curl \
	-v \
	--request GET \
	--header "Cookie: access=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODc5NzU0NTgsImlhdCI6MTQ4Nzk3MTg1OCwiaXNzIjoiYWNtZS5jb20iLCJzdWIiOiI4NThhNGIwMS02MmM4LTRjMmYtYmZhNy02ZDAxODgzM2JlYTciLCJhcHBsaWNhdGlvbklkIjoiM2MyMTllNTgtZWQwZS00YjE4LWFkNDgtZjRmOTI3OTNhZTMyIiwicm9sZXMiOlsiYWRtaW4iXX0.O29_m_NDa8Cj7kcpV7zw5BfFmVGsK1n3EolCj5u1M9hZ09EnkaOl5n68OLsIcpCrX0Ue58qsabag3MCNS6H4ldt6kMnH6k4bVg4TvIjoR8WE-yGcu_xDUObYKZYaHWiNeuDL1EuQQI_8HajQLND-c9juy5ILuz6Fhx8CLfHCziEHX_aQPt7jQ2IIasVzprKkgvWS07Hiv2Oskryx49wqCesl46b-30c6nfttHUDEQrVq9gaepca3Nhjj_cPtC400JgLCN9DOYIbtd69zvD8vDUOvVzMr2HGdWtKthqa35NF-3xMZKD8CShe8ZT74fNd9YZ0WRE-YeIf3T_Hv5p5V2w" \
	--url http://localhost:3000/analytics/v1/tasks/approved && echo "\n"

check_rejected_tasks:
	curl \
	-v \
	--request GET \
	--header "Cookie: access=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODc5NzU0NTgsImlhdCI6MTQ4Nzk3MTg1OCwiaXNzIjoiYWNtZS5jb20iLCJzdWIiOiI4NThhNGIwMS02MmM4LTRjMmYtYmZhNy02ZDAxODgzM2JlYTciLCJhcHBsaWNhdGlvbklkIjoiM2MyMTllNTgtZWQwZS00YjE4LWFkNDgtZjRmOTI3OTNhZTMyIiwicm9sZXMiOlsiYWRtaW4iXX0.O29_m_NDa8Cj7kcpV7zw5BfFmVGsK1n3EolCj5u1M9hZ09EnkaOl5n68OLsIcpCrX0Ue58qsabag3MCNS6H4ldt6kMnH6k4bVg4TvIjoR8WE-yGcu_xDUObYKZYaHWiNeuDL1EuQQI_8HajQLND-c9juy5ILuz6Fhx8CLfHCziEHX_aQPt7jQ2IIasVzprKkgvWS07Hiv2Oskryx49wqCesl46b-30c6nfttHUDEQrVq9gaepca3Nhjj_cPtC400JgLCN9DOYIbtd69zvD8vDUOvVzMr2HGdWtKthqa35NF-3xMZKD8CShe8ZT74fNd9YZ0WRE-YeIf3T_Hv5p5V2w" \
	--url http://localhost:3000/analytics/v1/tasks/rejected && echo "\n"

check_total_response_time:
	curl \
	-v \
	--request GET \
	--header "Cookie: access=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODc5NzU0NTgsImlhdCI6MTQ4Nzk3MTg1OCwiaXNzIjoiYWNtZS5jb20iLCJzdWIiOiI4NThhNGIwMS02MmM4LTRjMmYtYmZhNy02ZDAxODgzM2JlYTciLCJhcHBsaWNhdGlvbklkIjoiM2MyMTllNTgtZWQwZS00YjE4LWFkNDgtZjRmOTI3OTNhZTMyIiwicm9sZXMiOlsiYWRtaW4iXX0.O29_m_NDa8Cj7kcpV7zw5BfFmVGsK1n3EolCj5u1M9hZ09EnkaOl5n68OLsIcpCrX0Ue58qsabag3MCNS6H4ldt6kMnH6k4bVg4TvIjoR8WE-yGcu_xDUObYKZYaHWiNeuDL1EuQQI_8HajQLND-c9juy5ILuz6Fhx8CLfHCziEHX_aQPt7jQ2IIasVzprKkgvWS07Hiv2Oskryx49wqCesl46b-30c6nfttHUDEQrVq9gaepca3Nhjj_cPtC400JgLCN9DOYIbtd69zvD8vDUOvVzMr2HGdWtKthqa35NF-3xMZKD8CShe8ZT74fNd9YZ0WRE-YeIf3T_Hv5p5V2w" \
	--url "http://localhost:3000/analytics/v1/task/totalresponsetime?id=111" && echo "\n"

docs:
	docker build --tag swaggo/swag:1.8.1 . --file swaggo.Dockerfile && \
	docker run --rm --volume ${PWD}:/app --workdir /app swaggo/swag:1.8.1 /root/swag init \
		--parseDependency \
		--parseInternal \
		--dir ./internal/adapters/http \
		--generalInfo swagger.go \
		--output ./api/swagger/public \
		--parseDepth 1

tests/integration/analytics/approve:
	go test -v -tags=approve ./internal/tests/

tests/integration/analytics/reject:
	go test -v -tags=reject ./internal/tests/

up_local_environment:
	docker compose up -d

down_local_environment:
	docker compose down

start_analytics_microservice:
	go run cmd/main.go

send_create_task_event:
	docker container exec -it kafka \
	bash -c \
	'echo "{\"uuid\":\"7c866871-620d-4ee8-b665-e577e9bd62f3\",\"task_id\":1,\"time\":\"2022-07-26T20:13:09.260560304Z\",\"type\":\"create\",\"user\":\"author@mail.ru\",\"approvers_number\":2}" | /opt/bitnami/kafka/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test'