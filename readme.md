# ESTIMATE

## ABOUT

Создание ПО для управление бэклогом.

## INSTALL

* Установить [Go](https://golang.org/)
* Отредактировать файл _import.sh_ для Linux, MacOS, либо _import.cmd_ для Windows.
```Shell
#!/usr/bin/env bash
# Необходимо изменить переменную GOPATH - она должна указывать на корень нашего проекта
export GOPATH=~/projects/estima/

go get -u -v "github.com/dgrijalva/jwt-go"
go get -u -v "github.com/gorilla/context"
go get -u -v "github.com/gorilla/handlers"
go get -u -v "github.com/gorilla/mux"
go get -u -v "github.com/auth0/go-jwt-middleware"
go get -u -v "github.com/glycerine/zygomys/repl"
go get -u -v "gopkg.in/ldap.v2"
go get -u -v "github.com/go-errors/errors"

#
# Библиотека AranGO должна собираться из исходников, так как релиз слишком старый
# Поэтому используется загрузка из репозитория github 
#
mkdir ./src/github.com/diegogub/
cd ./src/github.com/diegogub/
# Удаление старой папки если есть
rm -rf aranGO

git init
git clone https://github.com/anatse/aranGO.git
git clone https://github.com/diegogub/napping
```
* Для сборки проекта используется команда go build, в результате появится исполняекмый файл **estima**
* Перед запуском программы необходимо также установить [arangodb](https://www.arangodb.com). Есть два способа установки:
  1. Через [docker](https://www.docker.com/)
        * Установить [docker](https://www.docker.com/)
        * Установить [arangodb](https://hub.docker.com/_/arangodb/)
        * Поправить переменные окружения через kitematic.
          * Пользователь: **root**
          * Пароль пользователя: **root**
          * Проброс порта - **8529**
        * Запустить arangodb
  2. Через [brew](https://brew.sh/)
        * Установить [brew](https://brew.sh/)
        * Установить, указать пароль для пользователя **root**, запустить [arangodb](https://www.arangodb.com)
        ```bash
        > brew install arangodb
        > /usr/local/opt/arangodb/sbin/arango-secure-installation
        > /usr/local/opt/arangodb/sbin/arangod
        ```
    * Через браузер войти в web интерфейс http://localhost:8529, создать там новую базу данных **estima**
* Проверить конфигурацию - файл config.json (комментарии в JSON не поддердиваются, здесь приведены для понимания)
    ```javascript
    {
      // Текущаий профиль 
      "active": "develop",
    
      // Список доступных профилей
      "profiles": [{
          "name": "develop",    // Имя профиля
          "secret": "secret",   // Ключ дл яшифрование куки 
          "Ldap": {             // Параметры подключения к LDAP
            "protocol": "fake", // Протокол, если установлен в fake, то проверка пользователя в LDAP не производится
            "host": "",
            "dn": "",
            "port": 389
          },
          "Database": {         // Параметры подключения к БД
            "url": "http://localhost:8529",
            "user": "root",
            "password": "root",
            "log": false,
            "name": "estima"
          },
          "Auth": {             // Параметры для формирование Auth куки
            "cookieName": "Estima",
            "maxAge": 10000
          }
        }, { // Следующий профиль
          "name": "test",
          "secret": "secret",
          "Ldap": {
            "protocol": "tcp",
            "host": "ldap.forumsys.com",
            "port": 389,
            "dn": "DC=example,DC=com"
          },
          "Database": {
            "url": "http://localhost:8529",
            "user": "root",
            "password": "123456",
            "log": false,
            "name": "estima"
          },
          "Auth": {
            "cookieName": "Estima",
            "maxAge": 10000
          }
      }]
    }
    ```

* Запустить приложение **go run**

----

### Запуска приложения в Intellij Idea.

 * Необходимо установить плагин [Go Lang Plugin](https://plugins.jetbrains.com/plugin/5047) для работы с Go.
 * Указать настройки:
    * Go > Go Libraries > Указать в global путь до проекта.
    * Run > Edit Config:
      * Создать Go Application.
      * В File указать путь до **estima.go**
      * В Working Directory указать путь до дириктории проекта.
      
* Теперт проект доступен по ссылке [localhost:9080](http://localhost:9080/)

## Unit тестирование

* Установить переменную окружения CONFIG_PATH = полный путь до файла config.json
* Запустить комманду 
```Shell 
go test ./src/ru/... -v
```