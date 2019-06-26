#!/bin/sh

APP_NAME=${PWD##*/}

read -p '是否需要重新编译并运行? (y/n)：' rebuild

case $rebuild in
    y )
        case "$(pidof ${APP_NAME} | wc -w)" in
            0 )
                ;;
            1 )
                kill $(pidof ./${APP_NAME} | awk '{print $0}')
                ;;
            * )
                kill $(pidof ./${APP_NAME} | awk '{print $0}')
                ;;
        esac

        go build
        nohup ./${APP_NAME} &
        echo "Service ${APP_NAME} rebuild success!"
        ;;
esac
