main.go содержит весь код
go.mod содержит импортированные библиотеки
go.sum содержит консольный вывод команды go get ., выполненной в каталоге проекта

Framework:gin
БД:PostgreSQL
Для тестирования использовал Httpie
Запросы:
GET-запрос для показа полного списка клиентов микросервиса| http://localhost:8080/clients
POST-запрос для добавления средств на счет (если клиент новый, то создается новое поле в бд) http://localhost/clients/addfunds {"Id":"1","Balance":100} -H "Content-Type: application/json"
GET-запрос для показа баланса конкретного клиента по его ID | http://localhost:8080/clients/:id
POST-запрос для резервирования средств на отдельном счете | http://locallhosy:8080/clients/reserve {"Id_client":"1","usluga":"1","transaction":"2","Price":100} -H "Content-Type: application/json"
POST-запрос для снятия денег с отдельного счета | http://localhost:8080/clients/accept {"Id_client":"1","usluga":"1","transaction":"2","Price":100} -H "Content-Type: application/json"

Структура БД:
avito_users - таблица клиентов микросервиса. Поля: ИД клиента, Баланс
uslugi - справочная таблица услуг, предоставляемых компанией. Поля: ИД услуги, название услуги
reserved_accounts - таблица счетов для резервирования средств с основного счета. Поля: ИД клиента, баланс 
transactions - таблица операций с со спец счетами пользователей. Поля: ИД операции, ИД клиента, ИД услуги, Цена 
