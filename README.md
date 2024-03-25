Когда завершите задачу, в этом README опишите свой ход мыслей: как вы пришли к решению, какие были варианты и почему выбрали именно этот. 

# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.

# Решение 

Для удобства я решил написать простенький HTTP-сервер, который использует 1 ручку.

В качестве хранилища я решил использовать Redis, который инкрементирует проверяет, сколько в течение последних N секунд было от конкретного пользователя.
Подумал, что использовать SQL-базы данных слишком дорого с точки зрения ресурсов. Поскольку нам нет надобности отслеживать полную историю, а лишь количественный показатель запросов за последние N секунд, любого IN-MEMORY хранилища достаточно.
Я могу не беспокоиться об одновременном доступе, все операции в Redis атомарны. Для декрементации счетчика я запускаю горутину, которая ждет N секунд и отправляет запрос в Redis.
Для хендлинга ошибок, которые могут возникнуть в этих горутинах, использую errgroup.
Реализовал Graceful Shutdown, а также добавил конфигурирование.
