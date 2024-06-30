
## API Endpoints

### Registering New Account

```http
  POST /api/v1/user/register
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required**. your username |
| `password` | `string` | **Required**. your password |

### Loging In To An Account

```http
  POST /api/v1/user/login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Required**. your username |
| `password`      | `string` | **Required**. your password |

### Refreshing JWT Token

```http
  POST /api/v1/user/login/refresh
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `refresh`      | `string` | **Required**. your JWT refresh token |

```http
  PUT|PATCH /api/v1/users/{username}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Required**. your username |

#### Header Parameter:
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer 'your  JWT Access token' 

```http
  DELETE /api/v1/users/{username}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Required**. your username |

#### Header Parameter:
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer 'your  JWT Access token' 

### Getting A User

```http
  GET /api/v1/users/${username}/
```



### Generating New Short Link

```http
  POST /api/v1/url
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `url`      | `string` | **Required**. your long URL |

#### Header Parameter:
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer 'your  JWT Access token' 

### Updating Long URL of Short URL

```http
  PUT|PATCH /api/v1/url/{uuid}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `url`      | `string` | **Required**. your long URL |

#### Header Parameter:
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer 'your  JWT Access token' 

### Deleting a Short URL

```http
  DELETE /api/v1/url/{uuid} 
```
#### Header Parameter:
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer 'your  JWT Access token' 

#### Getting User's URL's

```http
  GET /api/v1/users/{username}/urls
```
