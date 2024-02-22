#!/usr/bin/env bash
#########################################
#Copyright (c) 2024 ABLECLOUD Co. Ltd.
#
#설치되어 있는 smb를 구성하여 서비스 시작하는 스크립트
#
#최초작성자 : 정민철 주임(mcjeong@ablecloud.io)
#최초작성일 : 2024-02-20
#########################################
after_host=$(cat /etc/hosts | grep 'gwvm-mngt' | awk '{print $1}')
before_host=$(grep 'hosts' /usr/local/samba/etc/smb.conf | awk '{print $4}')
host_ip="${after_host:0:6}"
action=$1
user_id=$3
user_pw=$5
folder=$7
path=$9
if [ -n $action ]
then
        if [ $action == "create" ]
        then
                sed -i "s/$before_host/$host_ip/g" /usr/local/samba/etc/smb.conf
                if [ ${#user_id} -ne 0 ] && [ ${#user_pw} -ne 0 ] && [ ${#folder} -ne 0 ] && [ ${#path} -ne 0 ]
                then
                        echo -e "\n[$folder]" >> /usr/local/samba/etc/smb.conf
                        echo -e "\tpath = $path" >> /usr/local/samba/etc/smb.conf
                        echo -e "\twritable = yes" >> /usr/local/samba/etc/smb.conf
                        echo -e "\tpublic = yes" >> /usr/local/samba/etc/smb.conf
                        echo -e "\tcreate mask = 0777" >> /usr/local/samba/etc/smb.conf
                        echo -e "\tdirectory mask = 0777" >> /usr/local/samba/etc/smb.conf

                        # 사용자 추가를 위한 expect 스크립트
                        user_add=$(useradd -M $user_id > /dev/null 2>&1; echo $?)

                        if [ $user_add -ne 9 ]
                        then
                                expect -c "
                                spawn /usr/local/samba/bin/smbpasswd -a $user_id
                                expect "password:"
                                        send \"$user_pw\\r\"
                                        expect "password"
                                                send \"$user_pw\\r\"
                                expect eof
                                " > /dev/null
                                state=$(systemctl is-enabled smb)

                                if [ $state == "disabled" ]
                                then
                                        firewall-cmd --permanent --add-service=samba > /dev/null 2>&1
                                        firewall-cmd --reload > /dev/null 2>&1

                                        systemctl enable smb > /dev/null 2>&1
                                        systemctl start smb > /dev/null
                                fi
                        fi
                else
                        echo "ID, PW, Forder Name, PATH Check Please"
                fi
        elif [ $action == "user_create" ]
        then
                user=$(/usr/local/samba/bin/pdbedit -L | grep -v 'root' | cut -d ':' -f1 )
                for list in $user
                do
                        if [ $user_id == $list ]
                        then
                                echo "The same USER_ID exists."
                        else
                                useradd -M $user_id > /dev/null 2>&1

                                expect -c "
                                spawn /usr/local/samba/bin/smbpasswd -a $user_id
                                expect "password:"
                                        send \"$user_pw\\r\"
                                        expect "password"
                                                send \"$user_pw\\r\"
                                expect eof
                                "
                        fi
                done
        elif [ $action == "user_delete" ]
        then
                user_del=$(/usr/local/samba/bin/smbpasswd -x $user_id > /dev/null 2>&1; echo $?)
                if [ $user_del -eq 0 ]
                then
                        /usr/sbin/userdel $user_id > /dev/null 2>&1
                fi
        elif [ $action == "update" ]
        then
                user=$(/usr/local/samba/bin/pdbedit -L | grep -v 'root' | cut -d ':' -f1 )
                for list in $user
                do
                        if [ $user_id == $list ]
                        then
                                expect -c "
                                spawn /usr/local/samba/bin/smbpasswd -U $user_id
                                expect "password:"
                                        send \"$user_pw\\r\"
                                        expect "password"
                                                send \"$user_pw\\r\"
                                expect eof
                                " > /dev/null
                        else
                                echo "The same USER_ID exists."
                        fi
                done
        elif [ $action == "delete" ]
        then
                user=$(/usr/local/samba/bin/pdbedit -L | grep -v 'root' | cut -d ':' -f1 )
                for list in $user
                do
                        /usr/local/samba/bin/smbpasswd -x $list > /dev/null 2>&1
                        /usr/sbin/userdel $list > /dev/null 2>&1
                done
                        cat /dev/null > /usr/local/samba/etc/smb.conf
                        echo -e "\n[global]" >> /usr/local/samba/etc/smb.conf
                        echo -e "\tworkgroup = WORKGROUP" >> /usr/local/samba/etc/smb.conf
                        echo -e "\tsecurity = user" >> /usr/local/samba/etc/smb.conf
                        echo -e "\thosts allow = 100.100." >> /usr/local/samba/etc/smb.conf

                        systemctl stop smb > /dev/null 2>&1
                        systemctl disable smb > /dev/null 2>&1
                        firewall-cmd --permanent --remove-service=samba > /dev/null 2>&1
                        firewall-cmd --reload > /dev/null 2>&1
        fi
fi