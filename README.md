<!-- ABOUT THE PROJECT -->
# Тестовое задание в Effective Mobile

Этот репозиторий содержит решение тестового задания на позицию Junior разработчика в компанию Effective Mobile.

## О проекте

Это REST API, состоящий из 4 эндпоинтов:

- `/get` — Возвращает данные людей с различными фильтрами и пагинацией.
- `/delete` — Удаляет человека по идентификатору
- `/update` — Изменяет сущность
- `/add` — Добавляет новых людей в формате:
```json
{
"name": "Dmitriy",
"surname": "Ushakov",
"patronymic": "Vasilevich" // необязательно
}

```



<!-- GETTING STARTED -->
## Запуск проекта

### Запустите с помощью Docker Compose
Вы можете легко запустить проект локально с помощью Docker Compose, выполнив следующие шаги:

#### 1. Скачайте `docker-compose.yaml` файл
В терминале на linux/macOS запустите следующую команду:

```sh
wget https://raw.githubusercontent.com/dafraer/effective-mobile-task/refs/heads/main/docker-compose.yaml
```  

#### 2. Настройка архитектуры и переменных окружения
- При необходимости вы можете изменить порт на тот, который вам подходит.
- Выберите корректный тег образа в соответствии с архитектурой вашей системы.
    - **Для x86_64 (AMD64):** используйте `3.0-amd64`
    - **Для ARM64 (e.g., Raspberry Pi):** используйте `3.0-arm64`

#### 3. Запуск
Выполните следующую команду в терминале для запуска проекта

```sh
sudo docker-compose up -d
```  

<br>

### Запустите с помощью Go

#### 1. Скачайте код проекта

```sh
git clone https://github.com/dafraer/effective-mobile-task.git
```  

#### 2. Откройте репозиторий проекта

```sh
cd ./effective-mobile-task
```  

#### 3. В терминале на linux/macOS запустите следующую команду указав подходящий вам порт:

```sh
go run ./cmd/main.go
```  

<!-- CONTACT -->
## Контакты

Камиль Нуриев - [telegram](https://t.me/dafraer) - kdnuriev@gmail.com