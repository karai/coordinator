.PHONY: postgres adminer migrate migrate-down
postgres:
		docker run --rm -ti --network host 		\
		-e POSTGRES_PASSWORD=extrasupersecret \
		-e POSTGRES_DB=karai 									\
		-e POSTGRES_USER=postgres 						\
		postgres

adminer:
		docker run --rm -ti --network host adminer

migrate:
		migrate -source file://migrations 		\
						-database postgres://postgres:extrasupersecret@localhost/karai?sslmode=disable up
migrate-down:           
		migrate -source file://migrations 		\
						-database postgres://postgres:extrasupersecret@localhost/karai?sslmode=disable down
karai:
		go build 															\
		&& ./go-karai
