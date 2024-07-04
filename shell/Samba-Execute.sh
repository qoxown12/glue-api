#!/bin/sh
#########################################
#Copyright (c) 2024 ABLECLOUD Co. Ltd.
#
#설치되어 있는 smb를 구성하여 서비스 시작하는 스크립트
#
#최초작성자 : 정민철 주임(mcjeong@ablecloud.io)
#최초작성일 : 2024-02-20
#########################################
smb_conf="/etc/samba/smb.conf"
conf_json_file="/usr/local/glue-api/conf.json"

action=$1
sec_type=$2
user_id=$4
user_pw=$6
cache=$8
folder=${10}
path=${12}
fs_name=${14}
volume_path=${16}
realm=${18}
dns=${20}

if [ $action ]
then
        if [ $action == "create" ]
        then
                if [ $sec_type == "normal" ]
                then
                        user_check=$(useradd $user_id > /dev/null 2>&1 ; echo $?)
                        if [ $user_check -ne 9 ]
                        then
                                if [ ! -d $path ]
                                then
                                        mkdir -p $path
                                        chmod 777 $path
                                fi
                                if [ ${#user_id} -ne 0 ] && [ ${#user_pw} -ne 0 ] && [ ${#folder} -ne 0 ] && [ ${#path} -ne 0 ]
                                then
                                        mount -t ceph admin@.$fs_name=$volume_path $path
                                        fsid=$(cat /etc/ceph/ceph.conf | grep -m 1 'fsid' | awk '{print $3}')
                                        admin_key=$(cat /etc/ceph/ceph.client.admin.keyring | grep 'key' | awk '{print $3}')
                                        sed -i "$ a admin@$fsid.$fs_name=$volume_path $path ceph name=admin,secret=$admin_key,rw,relatime,seclabel,defaults 0 0" /etc/fstab

                                        cat /dev/null > $smb_conf

                                        echo -e "[global]" >> $smb_conf
                                        echo -e "workgroup = WORKGROUP" >> $smb_conf
                                        echo -e "hosts allow = 0.0.0.0/0.0.0.0" >> $smb_conf
                                        echo -e "security = user" >> $smb_conf
                                        echo -e "passdb backend = tdbsam" >> $smb_conf
                                        echo -e "usershare allow guests = yes" >> $smb_conf
                                        echo -e "guest account = root" >> $smb_conf
                                        echo -e "guest ok = yes" >>$smb_conf
                                        echo -e "\nforce user = root" >> $smb_conf

                                        echo -e "\nlog file = /var/log/samba/%m.log" >> $smb_conf
                                        echo -e "log level = 10" >> $smb_conf

                                        echo -e "\n[$folder]" >> $smb_conf
                                        echo -e "comment = Share Directories" >> $smb_conf
                                        echo -e "path = $path" >> $smb_conf
                                        echo -e "writable = yes" >> $smb_conf
                                        echo -e "public = yes" >> $smb_conf
                                        echo -e "create mask = 0777" >> $smb_conf
                                        echo -e "directory mask = 0777" >> $smb_conf

                                        if [ $cache == "true" ]
                                        then
                                                echo -e "csc policy = programs" >> $smb_conf
                                        fi

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

                                                systemctl enable --now smb.service > /dev/null 2>&1
                                                if [ $sec_type == "ads" ]
                                                then
                                                        systemctl enable --now winbind.service > /dev/null 2>&1
                                                fi
                                        fi
                                else
                                        echo "ID, PW, Forder Name, PATH Check Please"
                                fi
                        else
                                echo "The same user ID exists"
                        fi
                else
                        if [ ! -d $path ]
                        then
                                mkdir -p $path
                                chmod 777 $path
                        fi

                        if [ ${#user_id} -ne 0 ] && [ ${#user_pw} -ne 0 ] && [ ${#folder} -ne 0 ] && [ ${#path} -ne 0 ]
                        then
                                update-crypto-policies --set DEFAULT:AD-SUPPORT > /dev/null 2>&1

                                mount -t ceph admin@.$fs_name=$volume_path $path
                                fsid=$(cat /etc/ceph/ceph.conf | grep -m 1 'fsid' | awk '{print $3}')
                                admin_key=$(cat /etc/ceph/ceph.client.admin.keyring | grep 'key' | awk '{print $3}')
                                sed -i "$ a admin@$fsid.$fs_name=$volume_path $path ceph name=admin,secret=$admin_key,rw,relatime,seclabel,defaults 0 0" /etc/fstab

                                sed -i 's/normal/ads/g' $conf_json_file

                                cat /dev/null > $smb_conf

                                echo -e "[global]" >> $smb_conf
                                echo -e "hosts allow = 0.0.0.0/0.0.0.0" >> $smb_conf
                                echo -e "usershare allow guests = yes" >>$smb_conf
                                echo -e "usershare owner only = no" >> $smb_conf
                                echo -e "guest account = root" >> $smb_conf
                                echo -e "guest ok = yes" >>$smb_conf
                                echo -e "force user = root" >> $smb_conf

                                echo -e "\nsecurity = ads" >> $smb_conf
                                echo -e "winbind separator = +" >> $smb_conf
                                echo -e "idmap config * : unix_nss_info = yes" >> $smb_conf

                                echo -e "vfs objects = acl_xattr" >> $smb_conf
                                echo -e "map acl inherit = Yes" >> $smb_conf
                                echo -e "store dos attributes = Yes" >> $smb_conf

                                echo -e "\ndedicated keytab file = /etc/krb5.keytab" >> $smb_conf

                                echo -e "\nserver min protocol = SMB3" >> $smb_conf
                                echo -e "server max protocol = SMB3" >> $smb_conf

                                echo -e "\nlog file = /var/log/samba/%m.log" >> $smb_conf
                                echo -e "log level = 10" >> $smb_conf

                                echo -e "\n[$folder]" >> $smb_conf
                                echo -e "comment = Share Directories" >> $smb_conf
                                echo -e "path = $path" >> $smb_conf
                                echo -e "writable = yes" >> $smb_conf
                                echo -e "public = yes" >> $smb_conf
                                echo -e "create mask = 0777" >> $smb_conf
                                echo -e "directory mask = 0777" >> $smb_conf
                                echo -e "vfs objects = fake_compression" >> $smb_conf
                                if [ $cache == "true" ]
                                then
                                        echo -e "csc policy = programs" >> $smb_conf
                                fi

                                #DNS 설정
                                sed -i 's/DNS1/DNS2/g' /etc/sysconfig/network-scripts/ifcfg-enp0s20
                                echo -e "DNS1=$dns" >> /etc/sysconfig/network-scripts/ifcfg-enp0s20

                                systemctl restart NetworkManager

                                sleep 2

                                expect -c "
                                spawn realm join --membership-software=samba --client-software=winbind $realm -U $user_id
                                expect "password:"
                                        send \"$user_pw\\r\"
                                        expect "password"
                                                send \"$user_pw\\r\"
                                expect eof
                                "
                                state=$(systemctl is-enabled smb)

                                if [ $state == "disabled" ]
                                then
                                        firewall-cmd --permanent --add-service=samba > /dev/null 2>&1
                                        firewall-cmd --reload > /dev/null 2>&1

                                        systemctl enable --now smb.service > /dev/null 2>&1
                                        systemctl enable --now winbind.service > /dev/null 2>&1
                                fi
                        else
                                echo "ID, PW, Forder Name, PATH Check Please"
                        fi
                fi
        elif [ $action == "user_create" ]
        then
                user=$(pdbedit -L --debuglevel=1 | grep -v 'root' | grep -v 'ablecloud' | cut -d ':' -f1)
                if [ !$user]
                then
                        useradd $user_id > /dev/null 2>&1

                        expect -c "
                        spawn smbpasswd -a $user_id
                        expect "password:"
                                send \"$user_pw\\r\"
                                expect "password"
                                        send \"$user_pw\\r\"
                        expect eof
                        " > /dev/null
                else
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
                                        " > /dev/null
                                fi
                        done
                fi
        elif [ $action == "user_delete" ]
        then
                user_del=$(smbpasswd -x $user_id > /dev/null 2>&1; echo $?)
                if [ $user_del -eq 0 ]
                then
                        userdel -r $user_id > /dev/null 2>&1
                fi
        elif [ $action == "user_update" ]
        then
                user=$(pdbedit -L --debuglevel=1 | grep -v 'root' | grep -v 'ablecloud'| cut -d ':' -f1)
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
                path=$(cat $smb_conf | grep 'path' | awk '{print $3}')

                if [ $path ]
                then
                        umount -l -f $path > /dev/null 2>&1
                        fsid=$(cat /etc/ceph/ceph.conf | grep -m 1 'fsid' | awk '{print $3}')
                        sed -i "/$fsid/d" /etc/fstab

                        cat /dev/null > $smb_conf

                        echo -e "[global]" >> $smb_conf
                        echo -e "workgroup = WORKGROUP" >> $smb_conf
                        echo -e "hosts allow = 0.0.0.0/0.0.0.0" >> $smb_conf
                        echo -e "security = user" >> $smb_conf
                        echo -e "passdb backend = tdbsam" >> $smb_conf
                        echo -e "usershare allow guests = yes" >> $smb_conf
                        echo -e "guest account = root" >> $smb_conf
                        echo -e "guest ok = yes" >>$smb_conf
                        echo -e "force user = root" >> $smb_conf

                        echo -e "\nlog file = /var/log/samba/%m.log" >> $smb_conf
                        echo -e "log level = 10" >> $smb_conf

                        systemctl stop smb > /dev/null 2>&1
                        systemctl disable smb > /dev/null 2>&1
                        firewall-cmd --permanent --remove-service=samba > /dev/null 2>&1
                        firewall-cmd --reload > /dev/null 2>&1

                        if [[ "$(grep samba $conf_json_file | cut -d ':' -f2)" =~ ads ]]
                        then
                                sed -i 's/ads/normal/g' $conf_json_file
                                sed -i '/DNS1/d' /etc/sysconfig/network-scripts/ifcfg-enp0s20
                                sed -i 's/DNS2/DNS1/g' /etc/sysconfig/network-scripts/ifcfg-enp0s20

                                systemctl restart NetworkManager

                                systemctl stop winbind.service > /dev/null 2>&1
                                systemctl disable winbind.service > /dev/null 2>&1

                                user=$(pdbedit -L --debuglevel=1 | grep -v 'root' | grep -v 'ablecloud' | cut -d ':' -f1)

                                if [ $user ]
                                then
                                        for list in $user
                                        do
                                                smbpasswd -x $list > /dev/null 2>&1
                                                userdel -r $list > /dev/null 2>&1
                                        done
                                fi
                        else
                                user=$(pdbedit -L --debuglevel=1 | grep -v 'root' | grep -v 'ablecloud' | cut -d ':' -f1)

                                for list in $user
                                do
                                        smbpasswd -x $list > /dev/null 2>&1
                                        userdel -r $list > /dev/null 2>&1
                                done
                        fi
                else
                        echo "Already Deleted"
                fi
        elif [ $action == "select" ]
        then
                hostname=$(hostname | cut -d '.' -f1)
                ip_address=$(cat /etc/hosts | grep $hostname-mngt | awk '{print $1}')
                folder_name=$(grep -F '[' $smb_conf | grep -v 'global' | tr -d '[]')
                path=$(cat $smb_conf | grep path | awk '{print $3}')
                port_data=$(netstat -ltnp | grep  smb | grep -v tcp6 | awk '{print $4}' | cut -d ':' -f2 | tr "\n" ",")
                smb_names=$(systemctl show --no-pager smb | grep -w 'Names' | cut -d "=" -f2)
                smb_status=$(systemctl show --no-pager smb | grep -w 'ActiveState' | cut -d "=" -f2)
                smb_state=$(systemctl show --no-pager smb | grep -w 'UnitFileState' | cut -d "=" -f2)
                users_data=$(pdbedit -L --debuglevel=1 | grep -v 'root' | grep -v 'ablecloud'| cut -d ':' -f1)
                fs_name=$(mount | grep admin | grep smb | cut -d "." -f2 | cut -d "=" -f1)
                volume_path=$(mount | grep admin | grep smb | cut -d "=" -f2 | cut -d " " -f1)
                security_type=$(grep "samba" $conf_json_file | cut -d ':' -f2 | tr -d ' ' | sed 's/\"//g')

                if [[ "$(grep samba $conf_json_file | cut -d ':' -f2)" =~ ads ]]
                then
                winbind_names=$(systemctl show --no-pager winbind | grep -w 'Names' | cut -d "=" -f2)
                winbind_status=$(systemctl show --no-pager winbind | grep -w 'ActiveState' | cut -d "=" -f2)
                winbind_state=$(systemctl show --no-pager winbind | grep -w 'UnitFileState' | cut -d "=" -f2)
                realm=$(cat /etc/samba/smb.conf | grep 'realm' | awk '{print $3}')
                fi

                user=()
                for list in $users_data
                do
                       user+=\"$list\"\,
                done
                users=${user:0:${#user}-1}

                if [[ "$(grep samba $conf_json_file | cut -d ':' -f2)" =~ ads ]]
                then
                        if [ -z "$port_data" ]
                        then
                                printf '{"names":["%s","%s"],"status":["%s","%s"],"state":["%s","%s"],"hostname":"%s","security_type":"%s","ip_address":"%s","folder_name":"%s","path":"%s","port":[%s],"fs_name":"%s","volume_path":"%s","realm":"%s","users":[%s]}' "$smb_names" "$winbind_names" "$smb_status" "$winbind_status" "$smb_state" "$winbind_state" "$hostname" "$security_type" "$ip_address" "$folder_name" "$path" "$port_data" "$fs_name" "$volume_path" "$realm" "$users"
                        else
                                port=${port_data:0:${#port_data}-1}
                                printf '{"names":["%s","%s"],"status":["%s","%s"],"state":["%s","%s"],"hostname":"%s","security_type":"%s","ip_address":"%s","folder_name":"%s","path":"%s","port":[%s],"fs_name":"%s","volume_path":"%s","realm":"%s","users":[%s]}' "$smb_names" "$winbind_names" "$smb_status" "$winbind_status" "$smb_state" "$winbind_state" "$hostname" "$security_type" "$ip_address" "$folder_name" "$path" "$port" "$fs_name" "$volume_path" "$realm" "$users"
                        fi
                else
                        if [ -z "$port_data" ]
                        then
                                printf '{"names":"%s","status":"%s","state":"%s","hostname":"%s","ip_address":"%s","security_type":"%s","folder_name":"%s","path":"%s","port":[%s],"fs_name":"%s","volume_path":"%s","users":[%s]}' "$smb_names" "$smb_status" "$smb_state" "$hostname" "$ip_address" "$security_type" "$folder_name" "$path" "$port_data" "$fs_name" "$volume_path" "$users"
                        else
                                port=${port_data:0:${#port_data}-1}
                                printf '{"names":"%s","status":"%s","state":"%s","hostname":"%s","ip_address":"%s","security_type":"%s","folder_name":"%s","path":"%s","port":[%s],"fs_name":"%s","volume_path":"%s","users":[%s]}' "$smb_names" "$smb_status" "$smb_state" "$hostname" "$ip_address" "$security_type" "$folder_name" "$path" "$port" "$fs_name" "$volume_path" "$users"
                        fi
                fi
        else
                echo "Check the command"
        fi
else
        echo "A command is required"
fi