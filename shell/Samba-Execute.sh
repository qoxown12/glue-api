#!/bin/bash
#########################################
#Copyright (c) 2024 ABLECLOUD Co. Ltd.
#
#설치되어 있는 smb를 구성하여 서비스 시작하는 스크립트
#
#최초작성자 : 정민철 주임(mcjeong@ablecloud.io)
#최초작성일 : 2024-02-20
#########################################
smb_conf="/etc/samba/smb.conf"

after_host=$(cat /etc/hosts | grep 'gwvm-mngt' | awk '{print $1}')
before_host=$(grep 'hosts' $smb_conf | awk '{print $4}')
host_ip="${after_host:0:6}"
action=$1
user_id=$3
user_pw=$5
folder=$7
path=$9
fs_name=${11}
volume_path=${13}

if [ -n $action ]
then
        if [ $action == "create" ]
        then
                user_check=$(useradd $user_id > /dev/null 2>&1 ; echo $?)
                if [ $user_check -ne 9 ]
                then
                        if [ ! -d $path ]
                        then
                                mkdir -p $path
                                chmod 777 $path
                        fi
                                mount -t ceph admin@.$fs_name=$volume_path $path
                                fsid=$(cat /etc/ceph/ceph.conf | grep 'fsid' | awk '{print $3}')
                                admin_key=$(cat /etc/ceph/ceph.client.admin.keyring | grep 'key' | awk '{print $3}')
                                sed -i "$ a admin@$fsid.$fs_name=$volume_path $path ceph name=admin,secret=$admin_key,rw,relatime,seclabel,defaults 0 0" /etc/fstab

                        if [ ${#user_id} -ne 0 ] && [ ${#user_pw} -ne 0 ] && [ ${#folder} -ne 0 ] && [ ${#path} -ne 0 ]
                        then
                                cat /dev/null > $smb_conf

                                echo -e "[global]" >> $smb_conf
                                echo -e "workgroup = WORKGROUP" >> $smb_conf
                                echo -e "hosts allow = 0.0.0.0/0.0.0.0" >> $smb_conf
                                echo -e "security = user" >> $smb_conf
                                echo -e "passdb backend = tdbsam" >> $smb_conf
                                echo -e "usershare allow guests = yes" >> $smb_conf
                                echo -e "guest account = root" >> $smb_conf
                                echo -e "guest ok = yes" >>$smb_conf

                                echo -e "\nlog file = /var/log/samba/%m.log" >> $smb_conf
                                echo -e "log level = 10" >> $smb_conf

                                echo -e "\nforce user = root" >> $smb_conf
                                echo -e "csc policy = programs" >> $smb_conf

                                echo -e "\n[$folder]" >> $smb_conf
                                echo -e "comment = Share Directories" >> $smb_conf
                                echo -e "path = $path" >> $smb_conf
                                echo -e "writable = yes" >> $smb_conf
                                echo -e "public = yes" >> $smb_conf
                                echo -e "create mask = 0777" >> $smb_conf
                                echo -e "directory mask = 0777" >> $smb_conf

                                # 사용자 추가를 위한 expect 스크립트
                                useradd $user_id > /dev/null 2>&1

                                expect -c "
                                spawn smbpasswd -a $user_id
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

                        else
                                echo "ID, PW, Forder Name, PATH Check Please"
                        fi
                fi
        elif [ $action == "user_create" ]
        then
                user=$(pdbedit -L | grep -v 'root' | grep -v 'ablecloud' | cut -d ':' -f1 )
                for list in $user
                do
                        if [ $user_id == $list ]
                        then
                                echo "The same USER_ID exists."
                        else
                                useradd $user_id > /dev/null 2>&1

                                expect -c "
                                spawn smbpasswd -a $user_id
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
                user_del=$(smbpasswd -x $user_id > /dev/null 2>&1; echo $?)
                if [ $user_del -eq 0 ]
                then
                        userdel -r $user_id > /dev/null 2>&1
                fi
        elif [ $action == "update" ]
        then
                user=$(pdbedit -L | grep -v 'root' | grep -v 'ablecloud'| cut -d ':' -f1 )
                for list in $user
                do
                        if [ $user_id == $list ]
                        then
                                expect -c "
                                spawn smbpasswd -U $user_id
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
                path=$(cat $smb_conf | grep path | awk '{print $3}')
                umount -l -f $path
                sed '$ d' -i /etc/fstab

                user=$(pdbedit -L | grep -v 'root' | grep -v 'ablecloud' | cut -d ':' -f1)
                allow_ip=$(cat /etc/hosts | grep 'ccvm' | awk '{print $1}' | cut -d '.' -f1,2)

                for list in $user
                do
                        smbpasswd -x $list > /dev/null 2>&1
                        userdel -r $list > /dev/null 2>&1
                done
                        cat /dev/null > $smb_conf
                        echo -e "[global]" >> $smb_conf
                        echo -e "workgroup = WORKGROUP" >> $smb_conf
                        echo -e "hosts allow = 0.0.0.0/0.0.0.0" >> $smb_conf
                        echo -e "security = user" >> $smb_conf
                        echo -e "passdb backend = tdbsam" >> $smb_conf
                        echo -e "usershare allow guests = yes" >> $smb_conf
                        echo -e "guest account = root" >> $smb_conf
                        echo -e "guest ok = yes" >>$smb_conf

                        echo -e "\nlog file = /var/log/samba/%m.log" >> $smb_conf
                        echo -e "log level = 10" >> $smb_conf

                        echo -e "force user = root" >> $smb_conf
                        echo -e "csc policy = programs" >> $smb_conf

                        systemctl stop smb > /dev/null 2>&1
                        systemctl disable smb > /dev/null 2>&1
                        firewall-cmd --permanent --remove-service=samba > /dev/null 2>&1
                        firewall-cmd --reload > /dev/null 2>&1
        elif [ $action == "select" ]
        then
                hostname=$(hostname)
                ip_address=$(cat /etc/hosts | grep $hostname-mngt | awk '{print $1}')
                folder_name=$(grep -F '[' $smb_conf | grep -v 'global' | tr -d '[]')
                path=$(cat $smb_conf | grep path | awk '{print $3}')
                port_data=$(netstat -ltnp | grep  smb | grep -v tcp6 | awk '{print $4}' | cut -d ':' -f2 | tr "\n" ",")
                names=$(systemctl show --no-pager smb | grep -w 'Names' | cut -d "=" -f2)
                status=$(systemctl show --no-pager smb | grep -w 'ActiveState' | cut -d "=" -f2)
                state=$(systemctl show --no-pager smb | grep -w 'UnitFileState' | cut -d "=" -f2)
                users_data=$(pdbedit -L --debuglevel=1 | grep -v 'root' | grep -v 'ablecloud'| cut -d ':' -f1)
                fs_name=$(mount | grep admin | cut -d "." -f2 | cut -d "=" -f1)
                volume_path=$(mount | grep admin | cut -d "=" -f2 | cut -d " " -f1)
                user=()

                for list in $users_data
                do
                       user+=\"$list\"\,
                done
                users=${user:0:${#user}-1}
                if [ -z "$port_data" ]
                then
                        printf '{"Names":"%s","Status":"%s","State":"%s","hostname":"%s","ip_address":"%s","folder_name":"%s","path":"%s","port":[%s],"fs_name":"%s","volume_path":"%s","users":[%s]}' "$names" "$status" "$state" "$hostname" "$ip_address" "$folder_name" "$path" "$port_data" "$fs_name" "$volume_path" "$users"
                else
                        port=${port_data:0:${#port_data}-1}
                        printf '{"Names":"%s","Status":"%s","State":"%s","hostname":"%s","ip_address":"%s","folder_name":"%s","path":"%s","port":[%s],"fs_name":"%s","volume_path":"%s","users":[%s]}' "$names" "$status" "$state" "$hostname" "$ip_address" "$folder_name" "$path" "$port" "$fs_name" "$volume_path" "$users"
                fi
        fi
fi