# Принцип работы
![Throttler](https://user-images.githubusercontent.com/7625387/210429206-c806da7f-7f9e-4335-b448-c9fe4abdde96.png)

# Запуск
Make:
```
make gorun #через локальный go
make docker-run #через docker-compose
```
Вручную через docker-compose:
```
HTTP_PORT=9000 POSTGRESQL_PORT docker-compose up --build
```