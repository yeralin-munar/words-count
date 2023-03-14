# Intro
Утилита подсчитывает количество уникальных запросов из входного файла.

В начале работы утилиты:
- создается папка, в которой будут хранится файлы для каждого уникального запроса
- удаляется файл output.txt, если он существует

Каждый уникальный запрос попадает в LRU кэш.
Если запрос удаляется из кэша, то он записывается в файл. 
Название файла строка состоит из Unicode кодов каждого символа запроса соединенных символом подчеркивания ("_").
Размер файла обозначает количество уникальных запросов. 
При каждом повторении запроса в файл записывается N байтов. N зависит от количества повторений запроса.

Далее происходит чтение из папки (с файлами уникальных запросов), и запись результатов в output.txt.
Читается N файлов и потом происходит их удаление.
В конце папка удаляется.

# How to run
```
> go run main.go -n=4 -f=input.txt
```