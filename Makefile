docker-rm-volumes:
	docker volume rm $(docker volume ls -qf dangling=true)
run-generator:
	go run /home/pepachka/Develop/bachelor/internal/generator/cmd/main.go
