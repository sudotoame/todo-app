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

- `GET /tasks` - Возвращает все **tasks**
  - `?completed=true`
  - `?completed=false`
    - **query** параметры для метода **GET**, возвращают **tasks** у которых поле **completed** соответствует query значении
- `GET /tasks/{id}` - Возвращает **task** по **ID**

- `PATCH /tasks/{id}` - Принимает **json** со статусом **task**, по переданному **id** меняет его статус либо на `true/false`
```json
{
	"created": true
}	
```

- `PUT /tasks/{id}` - Меняет указанный **task** по **id** используя поля из тела запроса

```json
{
	"title": "title",
	"description": "description"
}
```

- `DELETE /tasks/{id}` - Принимает `id` **task** по **path** параметру, чтобы удалить таску
