#todo-app


## Запуск

- `make run`

## handlers

- `POST /tasks` - Принимает **json** запрос, и создает новый **task**

```json
{
	"title": "title",
	"description": "description"
}
```

- `GET /tasks` - Принимает **json** запрос и возвращает найденный **task**
  - `?completed=true`
  - `?completed=false`
    - **query** параметры для метода **GET**, возвращают **tasks** у которых поле **completed** соответствует query значении
  - Поле `"title"` можно сделать пустым `""`, чтобы получить все **tasks**

```json
{
	"title": "title"
}
```