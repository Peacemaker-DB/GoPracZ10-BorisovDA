# Практическое задание № 10 Борисов Денис Александрович ЭФМО-01-25

Цели:

1.	Понять устройство JWT и где его уместно применять в REST API.
2.	Сгенерировать и проверить JWT в Go (HS256), передавать его в Authorization: Bearer …. .
3.	Реализовать middleware-аутентификацию (достаёт токен, валидирует, кладёт клеймы в context).
4.	Добавить middleware-авторизацию (RBAC/права на эндпоинты).
5.	Встроить это в уже знакомую архитектуру HTTP-сервиса/роутера

# Выполнение практического задания.

1. Структура проекта

  Для выполнения практической работы была сделана следующая структура проекта

<img width="371" height="495" alt="image" src="https://github.com/user-attachments/assets/3e85a56a-fb5e-4615-90ef-620d17c0f341" />

   Так же были установлены все необходимые расширения для выполнения практики

<img width="767" height="108" alt="Снимок экрана 2025-11-17 000248" src="https://github.com/user-attachments/assets/97137ac0-0d86-4fe1-a864-d2db8aab780e" />


2.	Конфиг и запуск сервера.
  Для реализации JWT токинизации был создан файл config.go

<img width="661" height="712" alt="Снимок экрана 2025-11-16 235927" src="https://github.com/user-attachments/assets/b76f004f-4d1b-492f-9460-d2d1f5813ba2" />

3. Роутер и каркас.
 Для выполнения практики была написана файл router.go. В котором реализован машрутизация защищенных маршрутов.

<img width="561" height="728" alt="Снимок экрана 2025-11-17 005416" src="https://github.com/user-attachments/assets/1b153371-55cb-46e6-a893-8f4c15cf08c8" />



4. Пользователи и сервис.

   Затем был создан файл user.go, где описывается модель пользователя. А так же файл user_mem.go, где создается заготовленные пользователи для выполнения практики

  Файл user.go

<img width="372" height="181" alt="Снимок экрана 2025-11-17 005327" src="https://github.com/user-attachments/assets/511edb69-d4d9-4626-b346-52c50e232c2f" />

  Файл user_mem.go

<img width="953" height="937" alt="Снимок экрана 2025-11-17 005835" src="https://github.com/user-attachments/assets/bd386a8a-6946-4ccc-8648-77cab1012906" />

  После был написан jwt.go, где осуществляется работа с токеном

 <img width="766" height="957" alt="Снимок экрана 2025-11-17 005655" src="https://github.com/user-attachments/assets/9f83162a-463b-4ee1-9310-4424e72a2174" />
 
5. Точка входа и запуск сервера

   Для запуска сервера был написан main.go

<img width="461" height="535" alt="Снимок экрана 2025-11-16 213024" src="https://github.com/user-attachments/assets/8616d908-eb9f-4c38-8a2a-162b25ffeb28" />

  А так же был создан файл service.go

<img width="856" height="1059" alt="Снимок экрана 2025-11-17 005148" src="https://github.com/user-attachments/assets/59274768-840f-4c6a-8184-7d0b9ff472f8" />

6. Middleware: аутентификация и авторизация

  Для реалиизации аутентификация и авторизация, было написано два файла authn.go и authz.go

  Файл authn.go

<img width="732" height="670" alt="Снимок экрана 2025-11-17 005341" src="https://github.com/user-attachments/assets/5e570518-951b-4143-a90a-2b0bb59a5b3b" />

  Файл authz.go

  <img width="673" height="489" alt="Снимок экрана 2025-11-17 005350" src="https://github.com/user-attachments/assets/4bfcbad9-6e50-4b39-a59d-74bc3253c810" />


# Задания со «звёздочкой»

  Для перехода на RS256, был создан новый файл, где происходит работа с ключями

  <img width="501" height="787" alt="image" src="https://github.com/user-attachments/assets/244f287a-6bc9-4e4f-bcf7-a082761b6cf5" />

  Так же был изменен файл config.go, для правильной работы ключей по стандарту RS256

  <img width="331" height="324" alt="image" src="https://github.com/user-attachments/assets/812ff451-d473-4458-b76b-65c2a1f02012" />

  Для создания возможности проверки количество проверок за определенное время был создан файл ratelimit.go

  <img width="546" height="721" alt="image" src="https://github.com/user-attachments/assets/1e2e96d9-6cc1-41ab-a2c4-023a9e6bcb47" />

  Для вывода ошибки в едином формате был дополнен код в service.go

  <img width="491" height="167" alt="image" src="https://github.com/user-attachments/assets/26f73caa-0073-44fa-9c81-920289438e48" />


  
# Проверка работоспособности
  Для коррекной работы требуется впиисать следующие переменные окружения:
  
  $env:JWT_SECRET="dev-secret"; 
  $env:JWT_TTL="15min"; 
  $env:APP_PORT="8080";
  
  Для проверки работоспособности был запущен сервер, посел в Postman были проверено:

  Логин админом

<img width="686" height="432" alt="image" src="https://github.com/user-attachments/assets/014bdbe7-ab65-4041-9e90-17568eefe6ae" />

  Доступ к защищённым ручкам 

<img width="690" height="402" alt="image" src="https://github.com/user-attachments/assets/796347de-f751-4e41-a5b2-5825bbde6b01" />

  Логин админом пользователем

<img width="687" height="382" alt="image" src="https://github.com/user-attachments/assets/16859b18-6935-43b0-924d-45ccc2fedf4d" />

  Доступ к защищённым ручкам

<img width="687" height="337" alt="image" src="https://github.com/user-attachments/assets/ff7c39f3-ecfd-42a7-9ac1-4ca4deede77f" />

Данные пользователя в БД

<img width="1098" height="135" alt="image" src="https://github.com/user-attachments/assets/d9b4941a-4863-4db8-8211-18ceb7dd99d9" />

