#!/usr/bin/env bash

GLUEAPIFILE=Glue-API-Linux
SSH_TARGET_HOST="$1"
echo $SSH_TARGET_HOST
ssh $SSH_TARGET_HOST 'killall '$GLUEAPIFILE
scp /Users/ycyun/GolandProjects/Glue-API/$GLUEAPIFILE $SSH_TARGET_HOST:~/$GLUEAPIFILE
ssh $SSH_TARGET_HOST 'export GIN_MODE=release; ~/'$GLUEAPIFILE