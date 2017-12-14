#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
#set -x

YUCS_GITHUB_IO_DIR="/Users/yucs/work/yucs.github.io"
POST_DIR="${YUCS_GITHUB_IO_DIR}/source/_posts"


function hexo_deloy(){

     cd ${YUCS_GITHUB_IO_DIR}

     git checkout hexo-source 
    
     hexo g -d 
}


function cp_markdown(){

   for dir in `ls -F |grep "/$" | grep -v 'prepost'` 
   do  
        if [  `ls ${dir}` ]; then 
	       cp ${dir}/*.md  ${POST_DIR}
	    fi 
   done 
}




cp_markdown

hexo_deloy

