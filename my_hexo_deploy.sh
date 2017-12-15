#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
#set -x

YUCS_GITHUB_IO_DIR="/Users/yucs/work/yucs.github.io"
POST_DIR="${YUCS_GITHUB_IO_DIR}/source/_posts"
PWD_DIR=`pwd`

function chechout_branch(){
    
     cd ${YUCS_GITHUB_IO_DIR}
     git checkout hexo-source  
     cd ${PWD_DIR}
}

function cp_markdown(){

   for dir in `ls -F |grep "/$" | grep -v 'prepost'` 
   do  
        cd ${dir}

        if [  "`ls | grep '.md' `" != "" ]; then 
	       cp *.md  ${POST_DIR}
	    fi 

	    cd ..
   done 
}

function hexo_deloy(){

     cd ${YUCS_GITHUB_IO_DIR}    
      hexo g -d 
}


function hexo_local_test(){
	cd ${YUCS_GITHUB_IO_DIR} 
	hexo s
}





chechout_branch

cp_markdown

if [ $# != 0 ];then
   	hexo_deloy
else
	hexo_local_test

fi

