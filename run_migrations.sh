#!/bin/sh

echo "Ожидание готовности базы данных..."
until pg_isready -h db -U user -d calendar; do
  echo "Ждем..."
  sleep 2
done

echo "Выполнение миграций..."
goose -dir /migrations postgres "postgres://user:password@db:5432/calendar?sslmode=disable" up

echo "Миграции выполнены."