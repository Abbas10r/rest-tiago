# rest-tiago
docker ps
docker exec -it containerId bash
psql -u postgres
select *

goose -s create add_some_column sql
goose postgres postgres://social:social@localhost/social?sslmode=disable status | up
swag init -g ./api/main.go -d cmd,internal