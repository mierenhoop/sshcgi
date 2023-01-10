
all: sshcgi

.PHONY: sshcgi clean

sshcgi:
	go build sshcgi

run:
	go run sshcgi

connect:
	ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2222 0.0.0.0

clean:
	rm -f sshcgi
