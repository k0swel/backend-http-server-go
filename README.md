# RU
## Идея
В процессе выполнения лабораторной работы по дисциплине связанной с деплоем и архектурой обалчных приложений мною было дополнительно решено написать бекенд для своего сайта на языке GO. С этим компилируемым языком я знаком всего 2 недели на момент первого коммита в данном репозитории. Также я не разработчик, а инженер в облачном хостинг провайдере, так что не вручаюсь за корректность паттернов разработки. Буду рад замечаниям в разделе [**Discussions**](https://github.com/k0swel/backend-http-server-go/discussions) дабы повысить свою квалификацию в разработке.

## Сборка и запуск
Для сборки проекта необходимо перейти в директорию с проектом (```cd <dir> && go build .```), после чего скомпилируется ELF исполняемый бинарный файл, который нужно запустить на сервере.
Для запуска проекта важно определить следующие переменные окружения:

* ```"SQLITE3_FILE"``` - файл СУБД sqlite3

* ```"SQLITE3_TABLE"``` - единственная таблица нашего веб-приложения - таблица Users

* ```"MEMCACHED_ADDR"``` - адрес сервиса MEMCACHED "127.0.0.1:11211". Служит для хранения кук. Срок жизни куки = 24 часа

# ENG
## Idea
While working on a university assignment, i decided to develop ```http-server``` that handles questions from my-frontend using ```GOLANG```. I've only been familiar with this programming language for two weeks, since my first commit to the repository. Furthermore, I'm not a developer, but I'm an engineer at a hosting provider, so I'm not a development professional and not familiar with good development patterns. I'd appreciate your feedback in the [**Discussions**](https://github.com/k0swel/backend-http-server-go/discussions), as I want to improve my development skills.
## Compile and launch
For building the project you need to be in project's dir ```cd <dir> && go build .```, after you'll get executable ELF file that you can run at server. To launch project you need to define the following environment variables:
* ```"SQLITE3_FILE"``` - file SQLITE3.
* ```"SQLITE3_TABLE"``` - The only table our web-app - ```Users```
* ```"MEMCACHED_ADDR"``` - MEMCACHED is at "127.0.0.1:11211". It is storing cookies (life-time 24 hours)