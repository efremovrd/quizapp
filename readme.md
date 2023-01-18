# Анкета.ру

## Цель работы

Разработка web-приложения для проведения анкетирования.

## Решаемая проблема/предоставляемая возможность

Создание российской платформы для анкетирования в сети Интернет, что упростит процесс сбора информации.

## Краткий перечень функциональных требований

- создание новых анкет;
- управление уже созданными анкетами;
- просмотр ответов;
- заполнение анкет;
- авторизация в системе.

## Use-case диаграмма системы

![image](docs/images/usecase.png)

При наличии хотя бы одного ответа на анкету, заблокировать доступ к редактированию анкеты.

<details>
<summary>Исходный код PlantUML...</summary>
@startuml usecase

!theme plain

"Гость" as Guest
"Пользователь" as User

Guest --> (Зарегистрироваться)
Guest --> (Войти)

User --> (Создать анкету)
User --> (Заполнить анкету)
User --> (Просмотреть свои анкеты)

(Просмотреть свои анкеты) --> (Просмотреть анкету)
(Просмотреть свои анкеты) --> (Редактировать анкету)
(Просмотреть свои анкеты) --> (Удалить анкету)

(Просмотреть анкету) --> (Просмотреть ответы)

package НастройкаАнкеты {
    usecase "Добавить вопрос"
    usecase "Удалить вопрос"
    usecase "Редактировать вопрос"
}

(Создать анкету) --> НастройкаАнкеты
(Редактировать анкету) --> НастройкаАнкеты

@enduml
</details> 

## Экраны будущего приложения

[Miro Wireframes](https://miro.com/app/board/uXjVPPDIYuM=/ "На уровне черновых эскизов")

## Макет будущего приложения

[Figma](https://www.figma.com/file/YIrbTn95bmkG7yJeRqVVdD "Полноценный макет")

## ER-диаграмма сущностей системы

![image](docs/images/entity.png)

<details>
<summary>Исходный код PlantUML...</summary>
@startuml entity

!theme plain

entity User {
    id: string <<PK>>
    ---
    login: string
    password: string
}

entity Form {
    id: string <<PK>>
    ---
    user_id: string <<FK>>
    title: string
}

entity Question {
    id: string <<PK>>
    ---
    form_id: string <<FK>>
    header: string
}

entity Answer {
    question_id: string <<FK>>
    pool_answer_id: string <<FK>>
    ---
    value: string
}

entity PoolAnswer {
    id: string <<PK>>
    ---
    form_id: string <<FK>>
    user_id: string nullable <<FK>>
}

User ||--o{ Form

User ||--o{ PoolAnswer

Form ||--o{ Question

Form ||--o{ PoolAnswer

PoolAnswer ||--|{ Answer

Question ||--o{ Answer

@enduml
</details>

## Нагрузочное тестирование ApacheBenchmark

```
ab -n 50000 -c 500 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklkIjoiMyIsIkxvZ2luIjoiZ2Vmb3JjZXJ0eCIsIlBhc3N3b3JkIjoiMzA1MFRpIn19.RVjnA-_kDcdWJXlSrLdXMd2TSjTgG3sBgu3hHxJDJX8" 'http://localhost:9090/api/v1/forms/2/answers/14?offset=0&limit=10'
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Completed 40000 requests
Completed 45000 requests
Completed 50000 requests
Finished 50000 requests


Server Software:        quizapp
Server Hostname:        localhost
Server Port:            9090

Document Path:          /api/v1/forms/2/answers/14?offset=0&limit=10
Document Length:        150 bytes

Concurrency Level:      500
Time taken for tests:   71.348 seconds
Complete requests:      50000
Failed requests:        31830
   (Connect: 0, Receive: 0, Length: 31830, Exceptions: 0)
Total transferred:      5632700 bytes
HTML transferred:       2725500 bytes
Requests per second:    700.79 [#/sec] (mean)
Time per request:       713.482 [ms] (mean)
Time per request:       1.427 [ms] (mean, across all concurrent requests)
Transfer rate:          77.10 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   1.1      0      22
Processing:     0  708 1366.0      7   15879
Waiting:        0  706 1367.1      0   15879
Total:          0  708 1366.0      7   15879

Percentage of the requests served within a certain time (ms)
  50%      7
  66%    115
  75%   1047
  80%   1466
  90%   2509
  95%   3527
  98%   4971
  99%   6140
 100%  15879 (longest request)
```