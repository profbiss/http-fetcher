# Запуск сервиса
```bash
go run main.go -workers=16 -syncmap -addr=:8080
```
### флаги:
- workers - задаёт колличество воркеров 
- addr - адрес http api
- syncmap - переключает реализацию хранилища тасков на sync.Map по умолчанию используется TaskStore представляющий из себя структуру

```go
type TaskStore struct {
	mx sync.RWMutex
	tasks map[interface{}]interface{}
}
```


# Роуты
```
POST /task - создаёт таск
GET /task - возвращает список тасков
GET /task/{id} - возвращает данные таска
DELETE /task/{id} - удаляет таск
```

### Пример тела запроса на создание таска:
```json
{
  "method":"POST",
  "url":"https://ya.ru",
  "headers":{"Accept":["*/*"]},
  "body": "test",
}
```
### Примеры ответа:
```json
{
  "id":"6f5f38eb-be18-4b8e-afd0-4803fa21cc42",
}
```
```json
{
  "id":"6f5f38eb-be18-4b8e-afd0-4803fa21cc42",
  "result":{
    "status":403,
    "headers":{
      "Content-Type":["text/html; charset=utf-8"],
    },
    "length":100
  }
}
```
```json
{
  "id":"6f5f38eb-be18-4b8e-afd0-4803fa21cc42",
  "internal_error":"Тескт ошибки при обработке таска"
}
```