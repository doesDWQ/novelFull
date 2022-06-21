.PHONY: admin-fronted
admin-fronted:
	cd admin-fronted && npm start

.PHONY: backend
backendDev:
	cd backend && air

.PHONY: relyon
relyon:
	docker-compose up -d