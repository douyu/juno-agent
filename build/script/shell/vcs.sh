#!/bin/bash
# get vcs revision number

basePath=$(dirname $(dirname $(dirname $(dirname $(readlink -f $0)))))
cd ${basePath}

svnInfo=$(svn info 2>&1)
if [[ ${svnInfo} =~ "svn: E155007" ]] || [[ ${svnInfo} =~ "svn: command not found" ]]; then
    # get git md5
    gitSha1=$(git rev-parse HEAD)
    echo $(expr substr ${gitSha1} 1 8)
else
    # get svn revision
    svnRevision=$(svn info |grep Revision: |cut -c11-)
    echo ${svnRevision}
fi
