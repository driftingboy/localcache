run:
	sh start.sh

test:
	@echo ">>> start test"
	curl "http://localhost:8001/cache/db1/1" && \
	curl "http://localhost:8002/cache/db1/2" && \
	curl "http://localhost:8003/cache/db1/3"
	@echo ">>> end test"

test-curr:
	for number in 1 2 ; do \
  	curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
		curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
  	curl "http://localhost:8001/cache/db1/1" && \
		curl "http://localhost:8001/cache/db1/1" ; \
  done

