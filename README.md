### Запуск тестов

#### 1) Установка зависимостей
В корне проекта выполните:

```bash
go mod download
```

#### 2) Настройка переменных окружения
Проект использует файл `.env` в корне репозитория (рядом с `go.mod`).

Убедитесь, что в `.env` заданы переменные:

```env
API_URL=https://testboard.avito.com/api/v1
TEST_LOGIN=your_email@example.com
TEST_PASSWORD=your_password
DEBUG=false
```

#### 3) Запуск всех тестов
Из корня проекта:

```bash
go test -v -count=1 ./...
```

#### 4) Запуск конкретного пакета/сценария
Только сценарии объявлений:

```bash
go test -v ./tests/scenarios/advertisement
```

Только сценарии профиля:

```bash
go test -v ./tests/scenarios/myAdvertisement
```

#### 5) Windows (PowerShell)
Команды те же — запускать из PowerShell в корне проекта:

```powershell
go test -v ./...
```

Если переменные из `.env` не подхватываются, можно временно задать их прямо в PowerShell:

```powershell
$env:API_URL="https://testboard.avito.com/api/v1&quot;
$env:TEST_LOGIN="your_email@example.com"
$env:TEST_PASSWORD="your_password"
go test -v ./...
```