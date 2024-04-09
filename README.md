# Notification service
Сервис для доставки уведомлений пользователям. Подключается к кафке, получает новые события и, при включенных для пользователя уведомлениях, 
преобразовывает события в сообщения и отправляет их через бота в Телеграм.

Предполагается работа совместно с сервисом портфолио (tikkichest portfolio service).

## API
Сервис работает с форматом JSON.

Доступные методы:

    POST /notifications/{id} - сохраняет в базу данных айди пользователя в системе и его username в телеграме, работает как включение уведомлений для пользователя
    DELETE /notifications/{id} - удаляет из базы данные о пользователе по айди, работает как отключение уведомлений для выбранного пользователя
    PATCH /notifications/{id} - изменяет username в телеграме для пользователя с выбранным айди

Для включения уведомлений или изменении имени пользователя в телеграм в теле запроса ожидается объект вида:

    ID       int    `json:"id"`
    Username string `json:"telegram_username"`

## Kafka
Сервис ожидает получать из кафки события с айди пользователя в качестве ключа и объектом JSON в качестве значения:

    Object   Object `json:"object"` 
    ObjectID int    `json:"object_id"`
    Change   Change `json:"change"` 

Список доступных значений поля Object:

    Portfolio Object = "portfolio"
	Craft     Object = "craft"
	Content   Object = "content"

Список доступных значений поля Change:

    CreateObj Change = "created"
	UpdateObj Change = "changed"
	DeleteObj Change = "deleted"

## Telegram

Сервис отправляет уведомления в телеграм вида: "Your {Object} №{ObjectID} has been {Change}"

## Переменные окружения

Сервис умеет считывать переменные из файла .env в директории исполняемого файла (в корне проекта).

В примерах указаны дефолтные значения. Если программа не сможет считать пользовательские env, то возьмет их.

Переменные сервера:

    SERVER_LISTEN=:8080
    SERVER_READ_TIMEOUT=5s
    SERVER_WRITE_TIMEOUT=5s
    SERVER_IDLE_TIMEOUT=30s

Переменные Postgres:

    PG_USER=
	PG_PASSWORD=
	PG_HOST=localhost
	PG_PORT=5432
	PG_DATABASE=

Переменные Kafka:

    KAFKA_HOST=localhost
	KAFKA_PORT=9092
	KAFKA_TOPIC=
	KAFKA_CG_ID=             //Kafka consumer group id

Переменные Telegram:

	TG_TOKEN=                //token for telegram bot