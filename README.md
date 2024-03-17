## Randomgen

Этот сервис реализован на языке Go и служит инструментом для генерации случайных значений. Он предоставляет JSON API, доступный через HTTP. Каждой генерации присваивается уникальный идентификатор, что позволяет получить сгенерированное значение с помощью метода `retrieve`.

## Конечные точки

### POST /api/generate/

Этот конечный пункт генерирует случайное значение и его идентификатор.

#### Тело запроса

Тело запроса должно содержать JSON с необязательными параметрами:
- `type` (необязательно): Тип генерируемого случайного значения (`string`, `number`, `guid`, `alphanumeric` или пользовательские значения).
- `length` (необязательно): Длина генерируемого значения.

#### Ответ

После успешной генерации ответ будет содержать:
- `id`: Уникальный идентификатор для сгенерированного значения.

### GET /api/retrieve/

Этот конечный пункт извлекает ранее сгенерированное значение по его идентификатору.

#### Параметр запроса

- `id`: Идентификатор сгенерированного значения для извлечения.

#### Ответ

Ответ будет содержать сгенерированное значение.
