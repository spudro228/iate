Для запуска.
1. Установить зависимости
2. сделать [.env](.env) файл и установить значения переменным
3. go run main.go или скомпилить и запустить
4. Чтобы вырубить приложение нужен SIGKILL ткт грасфул шатдаун не реализован
5. 
http://localhost:8080/test -  покажет ОК если все работает

TODO:
~~1. Доработать /products/getAll  чтобы передавался user_id и по нему фильтровались~~
~~2. Доработать сохранение позиций через бота (tryToParseAndSaveInfoFromUser), чтобы позиция была привязана к user_id~~
3. Внести варианты ответов для бота
4. докер компоуз?