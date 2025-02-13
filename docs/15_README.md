сборка 

1) make up

2) make down

3) миграции применяются автоматически.
Удалить базу ранее созданных событий -  docker volume rm deployments_db_data

Интеграционные тесты запускаются командой make integration-tests
1) добавление события и выведение ошибки о попытке создания дублирующего события
2) получение листинга событий на день/неделю/месяц (тест сам создает три события, их же и выведет. Если создать другие
события с нужными характеристиками ранее, то выведет и их)
3) отправка уведомлений

После проведения успешных тестов окружение останавливается автоматически и выводится строка
Integration tests exited with code: 0


RabbitMQ:
http://localhost:15672
guest/guest

время событий используется в формате UNIX, для удобства использования команд есть декодер времени
из обычного формата в формат UNIX calendar/cmd/unix-transformer/main.go

Использование gRPC
1. Создание события
grpcurl -plaintext -d '{
"event": {
"title": "Название события 30",
"description": "Описание события 30",
"startTime": 1734006600,
"endTime": 1734010200,
"userId": "e9b1f4b2-dc3e-4ea0-a8f3-1234567890ab"
}
}' localhost:50051 api.EventService/CreateEvent

2. Обновление события
   grpcurl -plaintext -d '{
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab",
   "event": {
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab",
   "title": "Обновлённое название",
   "description": "Обновлённое описание",
   "startTime": 1609459200,
   "endTime": 1609462800,
   "userId": "e9b1f4b2-dc3e-4ea0-a8f3-1234567890ab"
   }
   }' localhost:50051 api.EventService/UpdateEvent

3. Удаление события
   grpcurl -plaintext -d '{
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab"
   }' localhost:50051 api.EventService/DeleteEvent

4. Получение события
   grpcurl -plaintext -d '{
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab"
   }' localhost:50051 api.EventService/GetEvent

5. Список всех событий
   grpcurl -plaintext -d '{}' localhost:50051 api.EventService/ListEvents

6. Список событий за день
   grpcurl -plaintext -d '{
   "date": 1731628800
   }' localhost:50051 api.EventService/ListEventsByDay

7. Список событий за неделю
   grpcurl -plaintext -d '{
   "start": 1609459200  // Время в формате Unix начала недели
   }' localhost:50051 api.EventService/ListEventsByWeek

8. Список событий за месяц // Время в формате Unix начала месяца
   grpcurl -plaintext -d '{
   "start": 1609459200  
   }' localhost:50051 api.EventService/ListEventsByMonth


