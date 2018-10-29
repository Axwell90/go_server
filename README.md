Сервер, реализующий WEB API для сервиса отзывов автовладельцев. В данных сервиса есть три вида сущностей: User (автовладелец), Model (модель авто), Review (отзыв).

Описание API:

Получение данных сущности:

```GET /<entity>/<id>```

```entity``` принимает одно из следующих значение: ```user```, ```model``` или ```review```.

```id``` содержит строку, которая может принимать любые значения.

Получение списка отзывов автовладельца:

```GET /user/<id>/reviews```

Запрос может содержать следующие GET-параметры:

fromDate - отобразить отзывы с ```created``` >= ```fromDate```.
<br>
toDate - отобразить отзывы с ```created``` <= ```toDate```.

Получение средней оценки модели

```GET /model/<id>/mark```

Запрос может содержать следующие GET-параметры:

fromDate - учитывать отзывы с ```created``` >= ```fromDate```.
<br>
toDate - учитывать отзывы с ```created``` <= ```toDate```.
<br>
sex - учитывать отзывы автовладельцев с указанным полом.

Обновление данных сущности

```POST /<entity>/<id>```

```entity``` принимает одно из следующих значение: ```user```, ```model``` или ```review```.

```id``` содержит строку, которая может принимать любые значения.

Добавление новой сущности

```POST /<entity>```

```entity``` принимает одно из следующих значение: ```user```, ```model``` или ```review```.

Все данные, передаваемые в теле запроса/ответа, отформатированы с использованием json и имеют соответсвующий заголовок Content-Type: ```application/json```.

В случае успешного запроса, ожидается ответ с кодом 200 или ошибка с кодом 400/404.

