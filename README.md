# website-check
Данная программа предназначена для мониторинга доступности списка веб-сайтов. 
Проверяет доступность веб-сайтов каждую минуту и записывает время ответа. 
Программа поддерживает нескольких пользователей, которые могут делать запросы 
для получения времени доступа к определенным веб-сайтам.

## Использование
Шаги, чтобы настроить и запустить приложение:
```
git clone https://github.com/iKayrat/website-check.git
```
```
go run main.go
```
*порт:"8080"

1. Пользователи могут делать запросы к следующим эндпоинтам:

- Получить время доступа к определенному сайту: отправьте GET-запрос на /access{url}, где {url} - имя веб-сайта, время доступа к которому требуется получить.
- Получить имя сайта с минимальным временем доступа: отправьте GET-запрос на /min.
- Получить имя сайта с максимальным временем доступа: отправьте GET-запрос на /max.
- Получить количества запросов пользователей по трем вышеперечисленным эндпойнтам: отправьте GET-запрос на /counts.
 
Ссылка для Postman:
- https://elements.getpostman.com/redirect?entityId=14424408-84a684b2-faa2-48f4-83dd-a1580f290ee3&entityType=collection
